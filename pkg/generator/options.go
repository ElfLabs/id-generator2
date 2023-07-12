package generator

type Options struct {
	Sequencer
	Formatter

	Count int64
	Step  int64
}

type Option func(o *Options)

func NewOptions(opts ...Option) Options {
	var o Options
	o.Apply(opts)
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
