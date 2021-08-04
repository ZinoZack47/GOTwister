package mem

import "unsafe"

const (
	dwLocalPlayer                   uint32 = 0xD8A2DC
	dwEntityList                    uint32 = 0x4DA31DC
	m_iTeamNum                      uint32 = 0xF4
	m_bSpotted                      uint32 = 0x93D
	m_bDormant                      uint32 = 0xED
	m_iHealth                       uint32 = 0x100
	m_bGunGameImmunity              uint32 = 0x3944
	dwGameRulesProxy                uint32 = 0x52C02BC
	m_SurvivalGameRuleDecisionTypes uint32 = 0x1328
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

func HasGunGameImmunity(dwEnt uint32) bool {
	var bHasGunGameImmunity bool
	RPM(dwEnt+m_bGunGameImmunity, uintptr(unsafe.Pointer(&bHasGunGameImmunity)), unsafe.Sizeof(bHasGunGameImmunity))
	return bHasGunGameImmunity
}

func IsDangerZone() bool {
	var dwGameRules uint32
	RPM(MemManager().Client+dwGameRulesProxy, uintptr(unsafe.Pointer(&dwGameRules)), unsafe.Sizeof(dwGameRules))
	var iSurvivalGRDT int
	RPM(dwGameRules+m_SurvivalGameRuleDecisionTypes, uintptr(unsafe.Pointer(&iSurvivalGRDT)), unsafe.Sizeof(iSurvivalGRDT))
	return (iSurvivalGRDT != 0)
}

func IsValidTarget(dwEnt uint32, bGunGameCheck bool) bool {

	if IsDormant(dwEnt) {
		return false
	}

	if !IsAlive(dwEnt) {
		return false
	}

	if Team(dwEnt) <= 1 || Team(dwEnt) == Team(LocalPlayer()) && !IsDangerZone() {
		return false
	}

	if HasGunGameImmunity(dwEnt) && bGunGameCheck {
		return false
	}

	return true
}
