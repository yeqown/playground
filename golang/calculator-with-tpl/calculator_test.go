package calculator_with_tpl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Calculate(t *testing.T) {
	type args struct {
		expr string
		fn   ParseExprAsPolishFunc
	}

	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "case 0",
			args: args{
				expr: "1 + ( ( 222 + 43 ) * 4 ) - 45",
				fn:   defaultParseExprAsPolish,
			},
			want: 1016,
		},
		{
			name: "case 1",
			args: args{
				expr: "1 + 2*3.3 + (445 * 6 + 7) * 89.1",
				fn:   defaultParseExprAsPolish,
			},
			want: 238528.3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Calculate(tt.args.expr, tt.args.fn)
			assert.Equal(t, nil, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_calcWithPolish(t *testing.T) {
	type args struct {
		polishExpr string
	}

	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "case 0",
			args: args{
				polishExpr: "5,1,2,+,4,*,+,3,-",
			},
			want: 14,
		},
		{
			name: "case 1",
			args: args{
				polishExpr: "1,222,43,+,4,*,+,45,-",
			},
			want: 1016,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calcWithPolish(tt.args.polishExpr)
			assert.Equal(t, tt.want, got)
		})
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
		{
			name: "case 0",
			args: args{
				expr: "1 + ( ( 222 + 43 ) * 4 ) - 45",
			},
			// got: 1,222,+,43,*,4,+,45,-
			want: "1,222,43,+,4,*,+,45,-",
		},
		{
			name: "case 1",
			args: args{
				expr: "1 + 2*3.3 + (445 * 6 + 7) * 89.1",
			},
			want: "1,2,3.3,*,+,445,6,*,7,+,89.1,*,+",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := defaultParseExprAsPolish(tt.args.expr)
			assert.Equal(t, nil, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_isOpLower(t *testing.T) {
	type args struct {
		op1 rune
		op2 rune
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "case 0",
			args: args{
				op1: '+',
				op2: '*',
			},
			want: true,
		},
		{
			name: "case 1",
			args: args{
				op1: '(',
				op2: '*',
			},
			want: false,
		},
		{
			name: "case 2",
			args: args{
				op1: '(',
				op2: ')',
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isOpLower(tt.args.op1, tt.args.op2)
			assert.Equal(t, tt.want, got)
		})
	}
}
