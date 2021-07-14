package main

import (
	"fmt"
	"math/rand"
	"time"
)

const DebugCasts = false

type Attributes struct {
	SP int

	HasteRating int
	CritRating  int
	HitRating   int
	Intellect   int
}

type Modifiers struct {
	// Talents
	BalanceOfPower   bool
	MoonkinForm      bool
	FocusedStarlight bool
	ImprovedMoonfire bool
	Moonfury         bool
	WrathOfCenarius  bool
	NaturesGrace     bool

	// Meta gems
	ChaoticSkyfireDiamond bool

	// Trinkets, used on cd.
	ScryersBloodgem      bool
	SilverCrescent       bool
	LightningCapacitator bool

	// Idols
	IdolOfTheMoongoddess bool
}

const HitRatingPerPercent = 12.6
const CritRatingPerPercent = 22.1
const IntPerCritPercent = 79.4
const ResistCoefficient = 0.5
const SharedTrinketCooldown = 30.0

type FightReport struct {
	TotalDamage int
	Duration    int
	DPS         float32
}

// At this point, all talents + passive gear effects are baked into stats
func runFight(
	fightDurationSeconds float32,
	baseSP int,
	baseCritChance float32,
	baseCritMultiplier float32,
	baseHitChance float32,
	spells []Spell,
	spellPriority []SpellID,
	hasNaturesGrace bool,
	procs []*Proc,
) float32 {
	naturesGraceUp := false
	totalDamage := float32(0)
	tick := float32(0)

	buffs := map[SpellID]ActiveBuff{}
	debuffs := map[SpellID]ActiveDebuff{}
	cooldowns := map[SpellID]float32{}

	for tick < fightDurationSeconds {

		if DebugCasts {
			fmt.Println(tick, "--------")
		}

		spellPower := baseSP
		critChance := baseCritChance
		critDamageMultiplier := baseCritMultiplier

		// Tick all dots / debuffs
		for i := range spellPriority {
			spellID := spellPriority[i]

			debuff, hasDebuff := debuffs[spellID]
			if hasDebuff {
				for len(debuff.Ticks) > 0 && debuff.Ticks[0] < tick {
					totalDamage += debuff.DamageTicks[0]

					if DebugCasts {
						fmt.Println(spells[spellID].Name, " ticked for ", debuff.DamageTicks[0])
					}

					debuff.DamageTicks = debuff.DamageTicks[1:]
					debuff.Ticks = debuff.Ticks[1:]

					debuffs[spellID] = debuff
				}

				if debuff.EndsAt < tick {
					if DebugCasts {
						fmt.Println(debuff.Name, " expired")
					}
					delete(debuffs, spellID)
				}
			}
		}

		// Figure out next spell to cast
		var toCast *SpellID
		for i := range spellPriority {
			spellID := spellPriority[i]

			cooldown, hasCooldown := cooldowns[spellID]
			if hasCooldown && cooldown > tick {
				// Can't cast, on CD.
				continue
			}

			_, hasDebuff := debuffs[spellID]
			if hasDebuff {
				// Don't case, already has the dot / debuff.
				continue
			}

			toCast = &spellID
			break
		}

		if toCast == nil {
			panic("invalid spellPriority")
		}

		spell := spells[*toCast]

		if DebugCasts {
			fmt.Println("Casting ", spell.Name)
		}

		// Figure out the cast time
		castTime := spell.CastTime
		if hasNaturesGrace && naturesGraceUp && castTime > 0.5 {
			castTime -= 0.5
			naturesGraceUp = false
		}

		// Fast forward through cast time.
		gcdCapped := false
		if !spell.OffGCD {
			if castTime > 1 {
				tick += castTime
			} else {
				gcdCapped = true
			}
		}

		// Put spell on cooldown
		cooldowns[*toCast] = tick + spell.Cooldown

		// Check resist
		if !spell.Unresistable {
			rng := rand.Float32() * 100
			if rng > baseHitChance {
				// Miss.
				if DebugCasts {
					fmt.Println("Miss!")
				}
				continue
			}
		}

		// Apply any buff modifiers
		for k := range buffs {
			buff := buffs[k]

			if buff.EndsAt > tick {
				spellPower += buff.SP
			} else {
				if DebugCasts {
					fmt.Println(buff.Name, " expired")
				}
				delete(buffs, buff.ID)
			}
		}

		// If not resisted, always apply the dot. Ticks can be resisted,
		// that logic exists in the section that applies the tick damage.
		if spell.Duration != 0 {

			if spell.Dot != nil {
				if DebugCasts {
					fmt.Println("Applied debuff")
				}
				debuff := ActiveDebuff{
					Name:   spell.Name,
					ID:     *toCast,
					EndsAt: tick + spell.Duration,
				}

				// Faerie fire is a debuff without damage, skip this
				if spell.Dot.AvgDamage > 0 {
					damageTicks := make([]float32, spell.Dot.NumTicks)
					tickTimings := make([]float32, spell.Dot.NumTicks)

					tickDelta := spell.Duration / float32(spell.Dot.NumTicks)
					for i := 0; i < spell.Dot.NumTicks; i++ {
						damageTicks[i] = (spell.Dot.AvgDamage + (float32(spellPower) * spell.Dot.Coefficient)) / float32(spell.Dot.NumTicks)
						tickTimings[i] = tick + spell.CastTime + (float32(i+1) * tickDelta)
					}

					debuff.DamageTicks = damageTicks
					debuff.Ticks = tickTimings

				}

				debuffs[*toCast] = debuff
			}

			if spell.Buff != nil {
				if DebugCasts {
					fmt.Println("Applied buff")
				}
				buff := ActiveBuff{
					Name:   spell.Name,
					ID:     *toCast,
					EndsAt: tick + spell.Duration,
					SP:     spell.Buff.SP,
				}

				buffs[*toCast] = buff
			}
		}

		// Hit. Check for partial resist
		// TODO: Partial resists are not implemented.
		if spell.Damage != nil {
			castDmg := float32(0)
			rng := rand.Float32() * 100
			if rng < critChance+spell.Damage.CritChanceModifier {
				// Crit!
				castDmg += float32(spell.Damage.AvgDamage) * float32(critDamageMultiplier)
				castDmg += float32(spell.Damage.Coefficient) * float32(spellPower) * float32(critDamageMultiplier)
				naturesGraceUp = true

				if DebugCasts {
					fmt.Println("Crit! ", castDmg)
				}

				// Chec procs on crits.
				for i := range procs {
					proc := procs[i]

					if !proc.ProcOnCrit {
						continue
					}

					// Can't proc due to proc cooldown
					if tick-proc.LastProc < proc.ProcCooldown {
						continue
					}

					proc.LastProc = tick
					proc.ProcCount++

					if proc.ProcCount < proc.ProcsRequired {
						continue
					}

					proc.ProcCount = 0

					if proc.Damage != nil {

						// Check for miss.
						procHitRng := rand.Float32() * 100
						if procHitRng > baseHitChance {
							// Proc missed
							continue
						}

						// Check for proc crit.
						procCritRng := rand.Float32() * 100
						if procCritRng < critChance {
							totalDamage += proc.Damage.AvgDamage * critDamageMultiplier
						} else {
							totalDamage += proc.Damage.AvgDamage
						}
					}

				}

			} else {
				// Regular damage
				castDmg += float32(spell.Damage.AvgDamage)
				castDmg += float32(spell.Damage.Coefficient) * float32(spellPower)

				if DebugCasts {
					fmt.Println("Hit ", castDmg)
				}
			}
			totalDamage += castDmg

		}

		if gcdCapped {
			tick = tick + 1
		}

	}

	return totalDamage
}

// Known issues :
// - Trinket cds aren't shared.
func Simulate(attrs Attributes, modifiers Modifiers, fightDurationSeconds float32) FightReport {
	rand.Seed(time.Now().UnixNano())

	// Add modifiers that apply to everything
	baseHitChance := 83.0 + float32(attrs.HitRating)/HitRatingPerPercent
	if modifiers.BalanceOfPower {
		baseHitChance += 4 // Balance of power
	}

	if baseHitChance > 99 {
		baseHitChance = 99
	}

	baseCritChance := float32(attrs.CritRating) / CritRatingPerPercent
	baseCritChance += float32(attrs.Intellect) / IntPerCritPercent
	if modifiers.MoonkinForm {
		baseCritChance += 5.0 // Moonkin aura
	}

	critDamageMultiplier := float32(2.0) // Vengeance
	if modifiers.ChaoticSkyfireDiamond {
		critDamageMultiplier = 2.09
	}

	spells := buildSpells(modifiers)

	trinkets := []SpellID{}
	if modifiers.SilverCrescent {
		trinkets = append(trinkets, SpellIdentSilverCrescent)
	}

	if modifiers.ScryersBloodgem {
		trinkets = append(trinkets, SpellIdentBloodgem)
	}

	spellPriority := []SpellID{
		SpellIdentFaerieFire,
	}

	spellPriority = append(spellPriority, trinkets...)
	spellPriority = append(spellPriority, SpellIdentMoonfire)
	spellPriority = append(spellPriority, SpellIdentStarfire)

	procs := []*Proc{}
	if modifiers.LightningCapacitator {
		procs = append(procs, &Proc{
			ProcChance:    100.0,
			ProcCooldown:  2.5,
			ProcsRequired: 3,

			ProcOnCrit: true,
			Damage: &Damage{
				AvgDamage: 750,
			},
		})
	}

	totalDamage := runFight(
		fightDurationSeconds,
		attrs.SP, baseCritChance,
		float32(critDamageMultiplier),
		baseHitChance,
		spells,
		spellPriority,
		// This could probably be modeled as a proc + buff instead.
		modifiers.NaturesGrace,
		procs,
	)

	// Done
	return FightReport{
		TotalDamage: int(totalDamage),
		Duration:    int(fightDurationSeconds),
		DPS:         totalDamage / fightDurationSeconds,
	}

}
