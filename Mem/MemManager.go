package mem

import (
	"golang.org/x/sys/windows"
)

func FindProc() {
	const szTarget = "csgo.exe"
	const TH32CS_SNAPPROCESS = 0x00000002
	hHandle, err := windows.CreateToolhelp32Snapshot(TH32CS_SNAPPROCESS, 0)

}
