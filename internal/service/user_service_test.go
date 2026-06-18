package service

import (
	"testing"
	"time"
)

func TestCalculateAge(t *testing.T) {
	now := time.Date(2026, 6, 18, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name string
		dob  time.Time
		want int
	}{
		{
			name: "birthday already passed this year",
			dob:  time.Date(1990, 5, 10, 0, 0, 0, 0, time.UTC),
			want: 36,
		},
		{
			name: "birthday not yet this year",
			dob:  time.Date(1990, 12, 25, 0, 0, 0, 0, time.UTC),
			want: 35,
		},
		{
			name: "birthday is today",
			dob:  time.Date(1990, 6, 18, 0, 0, 0, 0, time.UTC),
			want: 36,
		},
		{
			name: "newborn (age 0)",
			dob:  time.Date(2026, 6, 18, 0, 0, 0, 0, time.UTC),
			want: 0,
		},
		{
			name: "one day before birthday",
			dob:  time.Date(1990, 6, 19, 0, 0, 0, 0, time.UTC),
			want: 35,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := CalculateAge(tc.dob, now)
			if got != tc.want {
				t.Errorf("CalculateAge(%v, %v) = %d; want %d", tc.dob, now, got, tc.want)
			}
		})
	}
}
