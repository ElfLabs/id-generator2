package generator

import (
	"fmt"
	"testing"
)

func count(list []ID, id ID) int {
	n := 0
	for _, item := range list {
		if item == id {
			n++
		}
	}
	return n
}

func TestGenerator(t *testing.T) {
	type args struct {
		region  int64
		node    int64
		options []Option
		times   int
	}
	tests := []struct {
		name  string
		args  args
		check func(args args, ids []ID) error
	}{
		{
			name: "snowflake",
			args: args{
				node:    1,
				options: []Option{},
				times:   4096 * 10,
			},
			check: func(args args, ids []ID) error {
				for idx, id := range ids {
					if count(ids, id) > 1 {
						return fmt.Errorf("%d is repetition", id)
					}
					if idx%(len(ids)/10) == 0 || idx < 100 || idx > len(ids)-100 {
						t.Logf("%02d %d(0x%x)", idx, id, id)
					}
				}
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := NewGenerator(tt.args.region, tt.args.node, tt.args.options...)
			if err != nil {
				t.Errorf("NewGenerator failed: %v", err)
				return
			}
			var ids = make([]ID, 0, tt.args.times)
			for i := 0; i < tt.args.times; i++ {
				id := g.Generate()
				ids = append(ids, id)
			}
			if tt.check != nil {
				err = tt.check(tt.args, ids)
				if err != nil {
					t.Errorf("Check Generate failed: %v", err)
				}
			}
		})
	}
}
