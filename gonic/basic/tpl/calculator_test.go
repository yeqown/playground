package tpl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Calculate(t *testing.T) {
	got := Calculate("1 + ( ( 222 + 43 ) * 4 ) - 45")
	want := float32(1016)

	if got != want {
		t.Errorf("not Equal, got=%v, want=%v", got, want)
		t.FailNow()
	}
}

func Test_parseExprAsPolishV2(t *testing.T) {
	type args struct {
		expr string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		//{
		//	name: "case 0",
		//	args: args{
		//		expr: "1 + ( ( 222 + 43 ) * 4 ) - 45",
		//	},
		//	want: "",
		//},
		{
			name: "case 1",
			args: args{
				expr: "1 + 2*3.3 + (445 * 6 + 7) * 89.1",
			},
			// got:1,2,3.3,*,+,445,6,*,+,7,+,89.1,*
			want: "1,2,3.3,*,+,445,6,*,7,+,89.1,*,+",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseExprAsPolishV2(tt.args.expr)
			assert.Equal(t, tt.want, got)
		})
	}
}
