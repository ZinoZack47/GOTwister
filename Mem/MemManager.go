package mem

import (
	"fmt"
	"reflect"
	"syscall"
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

type moduleEntry32 struct {
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
	modkernel32       = windows.NewLazySystemDLL("kernel32.dll")
	procModule32NextW = modkernel32.NewProc("Module32NextW")
)

func module32Next(snapshot windows.Handle, mEntry *moduleEntry32) (err error) {
	r1, _, e1 := syscall.Syscall(procModule32NextW.Addr(), 2, uintptr(snapshot), uintptr(unsafe.Pointer(mEntry)), 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func Init() bool {
	hProc, iPid, err := catchGame()

	if err != nil {
		panic(err)
	}

	dwClient, dwEngine, err := catchModules(iPid)

	if err != nil {
		panic(err)
	}

	return true
}

func catchGame() (windows.Handle, uint32, error) {
	var hProc windows.Handle
	var iPid uint32
	const PROCESS_ALL_ACCESS = 0x1F1FFB
	szTarget := [8]uint16{'c', 's', 'g', 'o', '.', 'e', 'x', 'e'}
	var procEntry windows.ProcessEntry32
	hHandle, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		fmt.Println("could not create snapshot")
		return hProc, iPid, err
	}
	err = windows.Process32First(hHandle, &procEntry)
	if err != nil {
		fmt.Println("could not get first process")
		return hProc, iPid, err
	}

	for windows.Process32Next(hHandle, &procEntry) == nil {

		if len(procEntry.ExeFile) < 8 {
			continue
		}

		if reflect.DeepEqual(szTarget[:8], procEntry.ExeFile[:8]) {
			iPid = procEntry.ProcessID
			hProc, err = windows.OpenProcess(PROCESS_ALL_ACCESS, false, iPid)
			if err != nil {
				fmt.Println("Could not open process")
				return hProc, iPid, err
			}
			break
		}
	}

	windows.CloseHandle(hHandle)
	return hProc, iPid, nil
}

func catchModules(iPid uint32) (uint32, uint32, error) {
	var dwClient uint32
	var dwEngine uint32
	hHandle, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPMODULE32, iPid)
	szClient := []uint16{'c', 'l', 'i', 'e', 't', '.', 'd', 'l', 'l'}
	szEngine := []uint16{'e', 'n', 'g', 'i', 'n', 'e', '.', 'd', 'l', 'l'}
	if err != nil {
		fmt.Println("could not create module snapshot")
		return dwClient, dwEngine, err
	}

	var mEntry moduleEntry32
	for module32Next(hHandle, &mEntry) == nil {
		if len(mEntry.SzModule) < 9 {
			continue
		}

		if reflect.DeepEqual(szClient[:9], mEntry.SzModule[:9]) {
			dwClient = uint32(*mEntry.ModBaseAddr)
		}

		if len(mEntry.SzModule) < 10 {
			continue
		}

		if reflect.DeepEqual(szEngine[:10], mEntry.SzModule[:10]) {
			dwEngine = uint32(*mEntry.ModBaseAddr)
		}
	}

	windows.CloseHandle(hHandle)

	return dwClient, dwEngine, nil
}
