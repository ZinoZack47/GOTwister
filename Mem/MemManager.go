package mem

import (
	"fmt"
	"reflect"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	errnoERROR_IO_PENDING = 997
)

var (
	errERROR_IO_PENDING error = syscall.Errno(errnoERROR_IO_PENDING)
	errERROR_EINVAL     error = syscall.EINVAL
)

func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return errERROR_EINVAL
	case errnoERROR_IO_PENDING:
		return errERROR_IO_PENDING
	}
	return e
}

type moduleEntry32W struct {
	DwSize        uint32
	Th32ModuleID  uint32
	Th32ProcessID uint32
	GlblcntUsage  uint32
	ProccntUsage  uint32
	ModBaseAddr   *uint8
	ModBaseSize   uint32
	HMODULE       windows.Handle
	SzModule      [256]uint16
	SzExePath     [260]uint16
}

var (
	modkernel32            = windows.NewLazySystemDLL("kernel32.dll")
	procModule32NextW      = modkernel32.NewProc("Module32NextW")
	procReadProcessMemory  = modkernel32.NewProc("ReadProcessMemory")
	procWriteProcessMemory = modkernel32.NewProc("WriteProcessMemory")
)

func RPM(dwAddress uint32, pRead uintptr, iSize uintptr) (err error) {
	r1, _, e1 := syscall.Syscall6(procReadProcessMemory.Addr(), 5, uintptr(MemManager().Proc), uintptr(dwAddress), pRead, iSize, 0, 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func WPM(dwAddress uint32, pWrite uintptr, iSize uintptr) (err error) {
	r1, _, e1 := syscall.Syscall6(procWriteProcessMemory.Addr(), 5, uintptr(MemManager().Proc), uintptr(dwAddress), pWrite, iSize, 0, 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func module32NextW(snapshot windows.Handle, mEntry *moduleEntry32W) (err error) {
	r1, _, e1 := syscall.Syscall(procModule32NextW.Addr(), 2, uintptr(snapshot), uintptr(unsafe.Pointer(mEntry)), 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

type MemManager_t struct {
	Proc   windows.Handle
	Client uint32
	Engine uint32
}

var instantiated *MemManager_t
var once sync.Once

func MemManager() *MemManager_t {
	once.Do(func() {
		hProc, iPid, err := catchGame()

		if err != nil {
			panic(err)
		}

		dwClient, dwEngine, err := catchModules(iPid)

		if err != nil {
			panic(err)
		}
		instantiated = &MemManager_t{hProc, uint32(dwClient), uint32(dwEngine)}
	})
	return instantiated
}

func catchGame() (windows.Handle, uint32, error) {
	var hProc windows.Handle
	var iPid uint32
	const PROCESS_ALL_ACCESS = 0x1F1FFB
	szTarget := [8]uint16{'c', 's', 'g', 'o', '.', 'e', 'x', 'e'}
	hHandle, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		fmt.Println("could not create snapshot")
		return hProc, iPid, err
	}
	var procEntry windows.ProcessEntry32
	procEntry.Size = uint32(unsafe.Sizeof(procEntry))
	err = windows.Process32First(hHandle, &procEntry)
	if err != nil {
		fmt.Println("could not get first process")
		return hProc, iPid, err
	}

	for {

		err = windows.Process32Next(hHandle, &procEntry)

		if err != nil {
			if err == syscall.ERROR_NO_MORE_FILES {
				time.Sleep(time.Second)
				return catchGame()
			}
			return hProc, iPid, err
		}

		if len(procEntry.ExeFile) < cap(szTarget) {
			continue
		}

		if reflect.DeepEqual(szTarget[:], procEntry.ExeFile[:len(szTarget)]) {
			iPid = procEntry.ProcessID
			hProc, err = windows.OpenProcess(PROCESS_ALL_ACCESS, false, iPid)
			if err != nil {
				fmt.Println("Could not open process")
				return hProc, iPid, err
			}
			break
		}
	}

	defer windows.CloseHandle(hHandle)
	return hProc, iPid, nil
}

func catchModules(iPid uint32) (uintptr, uintptr, error) {
	var dwClient uintptr
	var dwEngine uintptr
	hHandle, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPMODULE, iPid)
	szClient := [10]uint16{'c', 'l', 'i', 'e', 'n', 't', '.', 'd', 'l', 'l'}
	szEngine := [10]uint16{'e', 'n', 'g', 'i', 'n', 'e', '.', 'd', 'l', 'l'}
	if err != nil {
		fmt.Println("could not create module snapshot")
		return dwClient, dwEngine, err
	}

	var mEntry moduleEntry32W
	mEntry.DwSize = uint32(unsafe.Sizeof(mEntry))
	for {

		if dwClient != 0 && dwEngine != 0 {
			break
		}

		err = module32NextW(hHandle, &mEntry)

		if err != nil {
			if err == syscall.ERROR_NO_MORE_FILES {
				if dwClient == 0 || dwEngine == 0 {
					return dwClient, dwEngine, fmt.Errorf("client 0x%06x not found and/or engine 0x%06x not found", dwClient, dwEngine)
				}
			}
			return dwClient, dwEngine, err
		}

		if len(mEntry.SzModule) < cap(szClient) {
			continue
		}

		if reflect.DeepEqual(szClient[:], mEntry.SzModule[:len(szClient)]) {
			dwClient = uintptr(unsafe.Pointer(mEntry.ModBaseAddr))
		}

		if len(mEntry.SzModule) < cap(szEngine) {
			continue
		}

		if reflect.DeepEqual(szEngine[:], mEntry.SzModule[:len(szEngine)]) {
			dwEngine = uintptr(unsafe.Pointer(mEntry.ModBaseAddr))
		}
	}

	defer windows.CloseHandle(hHandle)

	return dwClient, dwEngine, nil
}
