package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"
	"thatcodingguy/chicken_sim/sim"
)

func main() {
	fmt.Println("Hooking Simulate fucntion")
	js.Global().Set("MoonkinSim", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println("Simulating")
		numFights := args[0].Int()

		attrRaw := args[1].String()
		modRaw := args[2].String()

		var attrs sim.Attributes
		err := json.Unmarshal([]byte(attrRaw), &attrs)
		if err != nil {
			fmt.Println(err)
		}

		var mods sim.Modifiers
		err = json.Unmarshal([]byte(modRaw), &mods)
		if err != nil {
			fmt.Println(err)
		}

		results := []sim.FightReport{}
		for i := 0; i < numFights; i++ {
			result := sim.Simulate(
				attrs, mods, 420,
			)

			results = append(results, result)
		}

		avg := float32(0)
		for i := range results {
			avg += results[i].DPS
		}
		avg /= float32(len(results))

		fmt.Println(avg)
		return avg
	}))

	<-make(chan bool)
}
