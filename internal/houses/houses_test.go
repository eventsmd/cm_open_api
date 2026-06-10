package houses

import (
	"slices"
	"testing"
)

func TestExpand(t *testing.T) {
	tests := []struct {
		name         string
		numbers      string
		ranges       string
		wantHouses   []string
		wantUnresolv int
	}{
		{"numbers only", "1,2,3А", "", []string{"1", "2", "3А"}, 0},
		{"numbers with spaces", " 1 , 2 ", "", []string{"1", "2"}, 0},
		{"simple ranges", "", "1-3;20-22", []string{"1", "2", "3", "20", "21", "22"}, 0},
		{"numbers and ranges", "5", "1-2", []string{"5", "1", "2"}, 0},
		{"letter range unresolvable", "", "12а-12г", nil, 1},
		{"oversized range unresolvable", "", "1-500", nil, 1},
		{"inverted range unresolvable", "", "10-8", nil, 1},
		{"trailing dash trimmed", "", "1-3-", []string{"1", "2", "3"}, 0},
		{"single token in ranges", "", "7", []string{"7"}, 0},
		{"empty input", "", "", nil, 0},
		{"mixed ok and broken", "4", "1-2;12а-12г", []string{"4", "1", "2"}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			houses, unresolv := Expand(tt.numbers, tt.ranges)
			if !slices.Equal(houses, tt.wantHouses) {
				t.Errorf("houses = %v, want %v", houses, tt.wantHouses)
			}
			if unresolv != tt.wantUnresolv {
				t.Errorf("unresolvable = %d, want %d", unresolv, tt.wantUnresolv)
			}
		})
	}
}
