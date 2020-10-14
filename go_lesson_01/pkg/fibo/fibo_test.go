package fibo_test

import (
	"go_core/go_lesson_01/fibo"
	"testing"
)

func TestNumberAtPos(t *testing.T) {
	type args struct {
		pos uint
	}
	tests := []struct {
		name         string
		args         args
		wantErrorMsg string
		wantNumber   uint
	}{
		{
			name:         "Test_0",
			args:         args{pos: 0},
			wantErrorMsg: "",
			wantNumber:   0,
		},
		{
			name:         "Test_1",
			args:         args{pos: 1},
			wantErrorMsg: "",
			wantNumber:   1,
		},
		{
			name:         "Test_2",
			args:         args{pos: 2},
			wantErrorMsg: "",
			wantNumber:   1,
		},
		{
			name:         "Test_8",
			args:         args{pos: 8},
			wantErrorMsg: "",
			wantNumber:   21,
		},
		{
			name:         "Test_20",
			args:         args{pos: 20},
			wantErrorMsg: "",
			wantNumber:   6765,
		},
		{
			name:         "Test_21",
			args:         args{pos: 21},
			wantErrorMsg: "Getting numbers at pos > 20 is not supported!",
			wantNumber:   0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErrorMsg, gotNumber := fibo.NumberAtPos(tt.args.pos)

			if gotErrorMsg != tt.wantErrorMsg {
				t.Errorf("NumberAtPos() gotErrorMsg = %v, want %v", gotErrorMsg, tt.wantErrorMsg)
			}

			if gotNumber != tt.wantNumber {
				t.Errorf("NumberAtPos() gotNumber = %v, want %v", gotNumber, tt.wantNumber)
			}
		})
	}
}
