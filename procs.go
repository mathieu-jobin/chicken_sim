package main

type Proc struct {
	ProcChance   float32
	ProcCooldown float32

	ProcsRequired int

	ProcOnCrit bool

	Buff   *Buff
	Damage *Damage

	LastProc  float32
	ProcCount int
}
