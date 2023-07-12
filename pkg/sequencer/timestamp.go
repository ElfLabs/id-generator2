package sequencer

import (
	"errors"
	"time"
)

const (
	DefaultTimestampEpoch = 1689146115000 // 2023-07-12 15:15:15
	DefaultMaxBackwards   = 10
)

var ErrTimeRollback = errors.New("time rollback")

type TimeBackwardHandler func(offset int64)

type Timestamp struct {
	MaxBackwards    int64
	backwardHandler TimeBackwardHandler
	epoch           time.Time
	latest          int64
}

type TimestampOption func(t *Timestamp)

func NewTimestampSequencer(opts ...TimestampOption) *Timestamp {
	t := Timestamp{
		MaxBackwards: DefaultMaxBackwards,
	}

	t.setEpoch(DefaultTimestampEpoch)
	for _, opt := range opts {
		opt(&t)
	}
	if t.backwardHandler == nil {
		t.backwardHandler = t.defaultBackwardHandler
	}

	return &t
}

func (t *Timestamp) defaultBackwardHandler(offset int64) {
	panic(ErrTimeRollback)
}

func (t *Timestamp) setEpoch(epoch int64) {
	var now = time.Now()
	t.epoch = now.Add(time.Unix(epoch/1000, (epoch%1000)*1000000).Sub(now))
}

func (t *Timestamp) Current() int64 {
	return time.Since(t.epoch).Milliseconds()
}

func (t *Timestamp) Next() int64 {
	now := time.Since(t.epoch).Milliseconds()

	for now <= t.latest {
		offset := t.latest - now
		if offset > t.MaxBackwards {
			t.backwardHandler(offset)
		} else {
			time.Sleep(time.Microsecond * 100)
		}
		now = time.Since(t.epoch).Milliseconds() // waiting next millisecond
	}
	t.latest = now

	return now
}

func WithEpochTimestamp(epoch int64) TimestampOption {
	return func(t *Timestamp) {
		t.setEpoch(epoch)
	}
}

func WithLatestTimestamp(latest int64) TimestampOption {
	return func(t *Timestamp) {
		t.latest = latest
	}
}

func WithMaxBackwards(n int64) TimestampOption {
	return func(t *Timestamp) {
		t.MaxBackwards = n
	}
}

func WithTimeBackwardHandler(handle TimeBackwardHandler) TimestampOption {
	return func(t *Timestamp) {
		if handle == nil {
			return
		}
		t.backwardHandler = handle
	}
}
