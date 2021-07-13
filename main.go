package main

import (
	"fmt"
)

const NumFights = 50000 // How many fights to simulate
const FightLength = 420 // Fight is about 7 minutes

func main() {

	fmt.Println("Running sim for ", NumFights, " fights of ", FightLength, " seconds")

	results := []FightReport{}
	for i := 0; i < NumFights; i++ {
		result := Simulate(Attributes{
			SP:          1065,
			HasteRating: 0,
			CritRating:  228,
			HitRating:   130,
			Intellect:   347,
		}, Modifiers{
			// Talents
			BalanceOfPower:   true,
			MoonkinForm:      true,
			FocusedStarlight: true,
			ImprovedMoonfire: true,
			Moonfury:         true,
			WrathOfCenarius:  true,
			NaturesGrace:     true,

			// Trinkets
			ScryersBloodgem:      true,
			LightningCapacitator: true,

			// Meta gems
			ChaoticSkyfireDiamond: true,

			// Idols
			IdolOfTheMoongoddess: true,
		}, FightLength)

		results = append(results, result)
	}

	avg := float32(0)
	for i := range results {
		avg += results[i].DPS
	}
	avg /= float32(len(results))

	fmt.Println(avg)
}
