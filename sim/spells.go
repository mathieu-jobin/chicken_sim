package sim

type Spell struct {
	Name string

	CastTime float32

	Cooldown     float32
	OffGCD       bool
	Unresistable bool

	Duration float32
	SharedCD *SpellID

	Damage *Damage
	Dot    *Dot
	Buff   *Buff
}

type Damage struct {
	AvgDamage          float32
	CritChanceModifier float32
	Coefficient        float32
}

type Dot struct {
	AvgDamage   float32
	Coefficient float32
	NumTicks    int
}

type Buff struct {
	SP int
}

type SpellID int

const (
	SpellIdentMoonfire SpellID = iota
	SpellIdentStarfire
	SpellIdentFaerieFire

	// Trinkets
	SpellIdentBloodgem
	SpellIdentSilverCrescent
)

type ActiveDebuff struct {
	Name string
	ID   SpellID

	DamageTicks []float32
	Ticks       []float32
	EndsAt      float32
}

type ActiveBuff struct {
	Name string
	ID   SpellID

	SP     int
	EndsAt float32
}

func buildSpells(modifiers Modifiers) []Spell {
	moonfire := Spell{
		Name: "Moonfire",
		Damage: &Damage{
			AvgDamage:   331,
			Coefficient: 0.15,
		},
		Dot: &Dot{
			AvgDamage:   600,
			Coefficient: 0.52,
			NumTicks:    4,
		},
		Duration: 12.0,
	}

	starfire := Spell{
		Name:     "Starfire",
		CastTime: 3.0,
		Damage: &Damage{
			AvgDamage:   658,
			Coefficient: 1,
		},
	}

	ff := Spell{
		Name:     "Faerie Fire",
		Duration: 40.0,
		Dot:      &Dot{},
	}

	// Add modifiers that apply to specific spells
	if modifiers.FocusedStarlight {
		starfire.Damage.CritChanceModifier = 4.0
	}

	if modifiers.ImprovedMoonfire {
		moonfire.Damage.CritChanceModifier = 10.0 // Improved moonfire
	}

	if modifiers.WrathOfCenarius {
		starfire.Damage.Coefficient *= 1.2
	}

	// Modify only coefficients here, avg damage is multiplied lower
	if modifiers.Moonfury {
		starfire.Damage.Coefficient *= 1.1
		moonfire.Damage.Coefficient *= 1.1
		moonfire.Dot.Coefficient *= 1.1
	}

	if modifiers.IdolOfTheMoongoddess {
		starfire.Damage.AvgDamage += 55 * starfire.Damage.Coefficient
	}

	if modifiers.Moonfury {
		starfire.Damage.AvgDamage *= 1.1
		moonfire.Dot.AvgDamage *= 1.1
		moonfire.Damage.AvgDamage *= 1.1
	}

	// trinkets are modelled as spells.
	bloodgem := Spell{
		Name:         "Scryers Bloodgem",
		Cooldown:     90,
		Duration:     15,
		OffGCD:       true,
		Unresistable: true,
		Buff: &Buff{
			SP: 150,
		},
	}

	silverCrescent := Spell{
		Name:         "Silver Crescent",
		Cooldown:     90,
		Duration:     20,
		OffGCD:       true,
		Unresistable: true,
		Buff: &Buff{
			SP: 155,
		},
	}

	return []Spell{
		SpellIdentMoonfire:       moonfire,
		SpellIdentStarfire:       starfire,
		SpellIdentFaerieFire:     ff,
		SpellIdentBloodgem:       bloodgem,
		SpellIdentSilverCrescent: silverCrescent,
	}
}
