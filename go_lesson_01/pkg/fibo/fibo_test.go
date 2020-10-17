package fibo

import (
	"testing"
)

func TestNum(t *testing.T) {
	type args struct {
		pos uint
	}
	tests := []struct {
		name    string
		args    args
		wantNum uint
		wantErr string
	}{
		{
			name:    "Test_0",
			args:    args{pos: 0},
			wantErr: "",
			wantNum: 0,
		},
		{
			name:    "Test_1",
			args:    args{pos: 1},
			wantErr: "",
			wantNum: 1,
		},
		{
			name:    "Test_2",
			args:    args{pos: 2},
			wantErr: "",
			wantNum: 1,
		},
		{
			name:    "Test_8",
			args:    args{pos: 8},
			wantErr: "",
			wantNum: 21,
		},
		{
			name:    "Test_20",
			args:    args{pos: 20},
			wantErr: "",
			wantNum: 6765,
		},
		{
			name:    "Test_21",
			args:    args{pos: 21},
			wantErr: "getting numbers at pos > 20 is not supported",
			wantNum: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNum, err := Num(tt.args.pos)

			gotErr := ""
			if err != nil {
				gotErr = err.Error()
			}

			if gotErr != tt.wantErr {
				t.Errorf("Num() gotErr = %v, want %v", gotErr, tt.wantErr)
			}

			if gotNum != tt.wantNum {
				t.Errorf("Num() gotNum = %v, want %v", gotNum, tt.wantNum)
			}
		})
	}
}
