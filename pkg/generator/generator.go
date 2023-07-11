package generator

import (
	"errors"
	"strconv"
	"sync"
)

var ErrCountOverflow = errors.New("count overflow")

type generator struct {
	mu sync.Mutex

	region int64
	node   int64
	count  int64
	step   int64

	countMax int64
	stepMax  int64

	sequencer Sequencer
	formatter Formatter
}

func NewGenerator(region, node int64, opts ...Option) (Generator, error) {
	return NewGeneratorWithOptions(region, node, NewOptions(opts...))
}

func NewGeneratorWithOptions(region, node int64, options Options) (Generator, error) {
	g := &generator{
		region:    region,
		node:      node,
		count:     options.GetCount(),
		step:      options.GetStep(),
		sequencer: options.Sequencer,
		formatter: options.Formatter,
	}

	if g.sequencer == nil {
		return nil, errors.New("generator sequencer is nil")
	}
	if g.formatter == nil {
		return nil, errors.New("generator formatter is nil")
	}

	regionMax := GetRegionMax(g.formatter)
	if g.region < 0 || g.region > regionMax {
		return nil, errors.New("Region number must be between 0 and " + strconv.FormatInt(regionMax, 10))
	}

	nodeMax := GetNodeMax(g.formatter)
	if g.node < 0 || g.node > nodeMax {
		return nil, errors.New("generator number must be between 0 and " + strconv.FormatInt(nodeMax, 10))
	}

	g.countMax = GetCountMax(g.formatter)
	g.stepMax = GetStepMax(g.formatter)

	return g, nil
}

func (g *generator) GetSequencer() Sequencer {
	return g.sequencer
}

func (g *generator) GetFormatter() Formatter {
	return g.formatter
}

func (g *generator) shift() ID {
	id := (g.region << g.formatter.RegionShift()) |
		(g.node << g.formatter.NodeShift()) |
		(g.count << g.formatter.CountShift()) |
		(g.step << g.formatter.StepShift())
	return ID(id)
}

func (g *generator) Generate() ID {
	g.mu.Lock()
	defer g.mu.Unlock()

	current := g.sequencer.Current()

	if current == g.count {
		g.step = (g.step + 1) & g.stepMax
		if g.step == 0 {
			for current <= g.count {
				current = g.sequencer.Next()
			}
		}
	} else {
		g.step = 0
	}

	if current > g.countMax {
		panic(ErrCountOverflow)
	}

	g.count = current

	return g.shift()
}
