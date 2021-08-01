package features

import mem "gotwister/Mem"

func Radar() {
	dwLocalPlayer := mem.LocalPlayer()
	iLocalTeam := mem.Team(dwLocalPlayer)

	for id := uint32(1); id <= 32; id++ {
		dwEnt := mem.EntityFromId(id)
		bIsDormant := mem.IsDormant(dwEnt)

		if bIsDormant {
			continue
		}

		iTeam := mem.Team(dwEnt)

		if iTeam <= 1 || iTeam == iLocalTeam {
			continue
		}

		bIsAlive := mem.IsAlive(dwEnt)

		if !bIsAlive {
			continue
		}

		bIsSpotted := mem.IsSpotted(dwEnt)

		if bIsSpotted {
			continue
		}

		mem.SetSpotted(dwEnt, true)
	}
}
