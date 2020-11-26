package iface

import (
	"reflect"
	"testing"
)

func TestMaxPersonAge(t *testing.T) {
	tests := []struct {
		name    string
		persons []interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name:    "Employee is older",
			persons: []interface{}{NewEmployee(40), NewCustomer(20)},
			want:    NewEmployee(40),
			wantErr: false,
		},
		{
			name:    "Customer is older",
			persons: []interface{}{NewEmployee(30), NewCustomer(50)},
			want:    NewCustomer(50),
			wantErr: false,
		},
		{
			name:    "Wrong arg type",
			persons: []interface{}{NewEmployee(30), "customer"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MaxPersonAge(tt.persons...)
			if (err != nil) != tt.wantErr {
				t.Errorf("MaxPersonAge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MaxPersonAge() = %v, want %v", got, tt.want)
			}
		})
	}
}
