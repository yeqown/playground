package onoff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func genStateSeq(seqs ...int8) []Evt {
	evts := make([]Evt, 0, len(seqs))
	for _, seq := range seqs {
		var e Evt
		if seq > 0 {
			// Enter
			e = Evt{Typ: _ENTER, Seq: uint8(seq)}
		} else {
			// Leave
			e = Evt{Typ: _LEAVE, Seq: uint8(0 - seq)}
		}

		evts = append(evts, e)
	}

	return evts
}

func Test_genStateSeq(t *testing.T) {
	evts := genStateSeq([]int8{1, -1, 2, -2}...)
	t.Logf("%+v", evts)
}

func Test_onoff_stateMachine(t *testing.T) {
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
			m := newStateMachine()
			es := genStateSeq(tt.seqs...)
			m.Handle(es...)

			assert.Equal(t, tt.wantState, m.state())
		})
	}
}
