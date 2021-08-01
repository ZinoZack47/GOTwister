package mem

import "unsafe"

const (
	dwLocalPlayer uint32 = 0xD8A2DC
	dwEntityList  uint32 = 0x4DA31DC
	m_iTeamNum    uint32 = 0xF4
	m_bSpotted    uint32 = 0x93D
	m_bDormant    uint32 = 0xED
	m_iHealth     uint32 = 0x100
)

func EntityFromId(id uint32) uint32 {
	var dwEnt uint32
	RPM(MemManager().Client+dwEntityList+(id-0x01)*0x10, uintptr(unsafe.Pointer(&dwEnt)), unsafe.Sizeof(dwEnt))
	return dwEnt
}

func LocalPlayer() uint32 {
	var dwLocal uint32
	RPM(MemManager().Client+dwLocalPlayer, uintptr(unsafe.Pointer(&dwLocal)), unsafe.Sizeof(dwLocal))
	return dwLocal
}

func Health(dwEnt uint32) int {
	var iHealth int
	RPM(dwEnt+m_iHealth, uintptr(unsafe.Pointer(&iHealth)), unsafe.Sizeof(iHealth))
	return iHealth
}

func Team(dwEnt uint32) int {
	var iTeam int
	RPM(dwEnt+m_iTeamNum, uintptr(unsafe.Pointer(&iTeam)), unsafe.Sizeof(iTeam))
	return iTeam
}

func IsDormant(dwEnt uint32) bool {
	var bIsDormant bool
	RPM(dwEnt+m_bDormant, uintptr(unsafe.Pointer(&bIsDormant)), unsafe.Sizeof(bIsDormant))
	return bIsDormant
}

func IsSpotted(dwEnt uint32) bool {
	var bIsSpotted bool
	RPM(dwEnt+m_bSpotted, uintptr(unsafe.Pointer(&bIsSpotted)), unsafe.Sizeof(bIsSpotted))
	return bIsSpotted
}

func SetSpotted(dwEnt uint32, bIsSpotted bool) {
	WPM(dwEnt+m_bSpotted, uintptr(unsafe.Pointer(&bIsSpotted)), unsafe.Sizeof(bIsSpotted))
}

func IsAlive(dwEnt uint32) bool {
	return Health(dwEnt) > 0
}
