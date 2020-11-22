package iface

import "testing"

func TestMaxRichPersonAge(t *testing.T) {
	tests := []struct {
		name    string
		persons []AgeReporter
		want    int
	}{
		{
			name:    "Employee is older",
			persons: []AgeReporter{NewRichEmployee(40), NewRichCustomer(20)},
			want:    40,
		},
		{
			name:    "Customer is older",
			persons: []AgeReporter{NewRichEmployee(30), NewRichCustomer(50)},
			want:    50,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaxRichPersonAge(tt.persons...); got != tt.want {
				t.Errorf("MaxPersonAge() = %v, want %v", got, tt.want)
			}
		})
	}
}
