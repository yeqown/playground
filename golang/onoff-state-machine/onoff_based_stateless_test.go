package onoff

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/qmuntal/stateless"
	"github.com/stretchr/testify/assert"
)

var (
	emptySeq = uint8(0)
)

func prepareStateMachine(gSeq uint8) *stateless.StateMachine {
	sm := stateless.NewStateMachine(_OFFLINE)

	sm.SetTriggerParameters(_ENTER, reflect.TypeOf(emptySeq))
	sm.SetTriggerParameters(_LEAVE, reflect.TypeOf(emptySeq))
	sm.Configure(_OFFLINE).
		PermitDynamic(_ENTER, func(ctx context.Context, i ...interface{}) (stateless.State, error) {
			seq := i[0].(uint8)
			fmt.Printf("[off] enter seq=%d\n", seq)
			if seq <= gSeq {
				return _OFFLINE, nil
			}

			gSeq = seq
			return _ONLINE, nil
		}).
		PermitDynamic(_LEAVE, func(ctx context.Context, i ...interface{}) (stateless.State, error) {
			seq := i[0].(uint8)
			fmt.Printf("[off] leave seq=%d\n", seq)
			if seq <= gSeq {
				return _OFFLINE, nil
			}

			gSeq = seq
			return _OFFLINE, nil
		})

	sm.Configure(_ONLINE).
		PermitDynamic(_ENTER, func(ctx context.Context, i ...interface{}) (stateless.State, error) {
			seq := i[0].(uint8)
			fmt.Printf("[on] enter seq=%d\n", seq)
			if seq <= gSeq {
				return _ONLINE, nil
			}

			gSeq = seq
			return _ONLINE, nil
		}).
		PermitDynamic(_LEAVE, func(ctx context.Context, i ...interface{}) (stateless.State, error) {
			seq := i[0].(uint8)
			fmt.Printf("[on] leave seq=%d\n", seq)
			if seq < gSeq {
				return _ONLINE, nil
			}
			gSeq = seq
			return _OFFLINE, nil
		})

	return sm
}

func Test_statelessMachine(t *testing.T) {
	tests := []struct {
		name      string
		seqs      []int8
		wantState state
	}{
		{
			name:      "case 0",
			seqs:      []int8{1, -1, 2, -2},
			wantState: _OFFLINE,
		},
		{
			name:      "case 1",
			seqs:      []int8{1, 2, -1, -2},
			wantState: _OFFLINE,
		},
		{
			name:      "case 2",
			seqs:      []int8{1, -1, -2, 2},
			wantState: _OFFLINE,
		},
		{
			name:      "case 3",
			seqs:      []int8{1, 2, 3, -1, -2},
			wantState: _ONLINE,
		},
		{
			name:      "case 4",
			seqs:      []int8{1, 2, -2, 3, -1},
			wantState: _ONLINE,
		},
		{
			name:      "case 5",
			seqs:      []int8{1, 2, -2, 3, -3, -1},
			wantState: _OFFLINE,
		},
		{
			name:      "case 6",
			seqs:      []int8{1, 2, -2, 3, -1, -3},
			wantState: _OFFLINE,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gSeq := uint8(0)
			m := prepareStateMachine(gSeq)
			es := genStateSeq(tt.seqs...)

			for _, e := range es {
				if err := m.Fire(e.Typ, e.Seq); err != nil {
					t.Error(err)
					t.FailNow()
				}
				//t.Logf("now state=%v\n", m.MustState())
			}

			assert.Equal(t, tt.wantState, m.MustState())
		})
	}
}

// It's BAD
func Test_graph_output(t *testing.T) {
	gSeq := uint8(0)
	m := prepareStateMachine(gSeq)

	file, err := os.OpenFile("./stateless.dot", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	_, err = file.WriteString(m.ToGraph())
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}
