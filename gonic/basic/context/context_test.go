package golang

import "testing"

func TestWorker(t *testing.T) {
	type args struct {
		secs int
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
	}{
		// TODO: Add test cases.
		{name: "case 1", args: args{secs: 3}, wantCode: -1},
		{name: "case 2", args: args{secs: 1}, wantCode: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCode := Worker(tt.args.secs); gotCode != tt.wantCode {
				t.Errorf("Worker() = %v, want %v", gotCode, tt.wantCode)
			}
		})
	}
}
