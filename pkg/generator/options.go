package generator

import (
	"github.com/ElfLabs/id-generator/pkg/format"
	"github.com/ElfLabs/id-generator/pkg/planner/snowflake"
	"github.com/ElfLabs/id-generator/pkg/sequencer"
)

type Options struct {
	Sequencer
	Formatter

	Count int64
	Step  int64

	ErrCh chan<- error
}

type Option func(o *Options)

func NewOptions(opts ...Option) Options {
	var o Options
	o.Apply(opts)
	o.Init()
	return o
}

func (o *Options) Init() *Options {
	switch {
	case o.Sequencer == nil && o.Formatter == nil:
		planner := snowflake.NewSnowflakePlanner()
		o.Sequencer = planner
		o.Formatter = planner
	case o.Sequencer == nil:
		o.Sequencer = sequencer.NewTimestampSequencer()
	case o.Formatter == nil:
		o.Formatter = format.SnowflakeFormat{}
	}
	return o
}

func (o *Options) Apply(opts []Option) *Options {
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func (o *Options) GetCount() int64 {
	return o.Count
}

func (o *Options) GetStep() int64 {
	return o.Step
}

func WithCount(count int64) Option {
	return func(o *Options) {
		o.Count = count
	}
}

func WithStep(step int64) Option {
	return func(o *Options) {
		o.Step = step
	}
}

func WithRecover(count, step int64) Option {
	return func(o *Options) {
		o.Count = count
		o.Step = step
	}
}

func WithSequencer(sequencer Sequencer) Option {
	return func(o *Options) {
		o.Sequencer = sequencer
	}
}

func WithFormatter(formatter Formatter) Option {
	return func(o *Options) {
		o.Formatter = formatter
	}
}

func WithPlanner(planner Planner) Option {
	return func(o *Options) {
		o.Sequencer = planner
		o.Formatter = planner
	}
}

func WithErrorChan(ch chan<- error) Option {
	return func(o *Options) {
		o.ErrCh = ch
	}
}
