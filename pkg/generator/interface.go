package generator

type Generator interface {
	Generate() ID
}

type Sequencer interface {
	Current() int64
	Next() int64
}

type Formatter interface {
	RegionBits() uint8
	NodeBits() uint8
	CountBits() uint8
	StepBits() uint8

	RegionShift() uint8
	NodeShift() uint8
	CountShift() uint8
	StepShift() uint8
}

type Planner interface {
	GetEpoch() int64
	Sequencer
	Formatter
}
