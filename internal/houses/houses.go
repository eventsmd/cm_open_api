// Package houses expands incident house lists serialized by events-workflow:
// house_numbers — CSV ("1,2,3А"), house_ranges — semicolon-separated dash
// ranges ("1-10;15-20").
package houses

import (
	"strconv"
	"strings"
)

// MaxRangeSpan caps range expansion: anything wider is counted as
// unresolvable (callers may fall back to street-level). This protects
// against pathological "1-9999" inputs.
const MaxRangeSpan = 200

// Expand returns individual house labels and the count of range groups that
// could not be expanded (non-numeric bounds, inverted or oversized ranges).
// Output may contain duplicates; callers should dedup if needed.
func Expand(numbers string, ranges string) (houses []string, unresolvable int) {
	// Keep delimiters in sync with internal/postgres.GetOutages.
	for _, n := range strings.Split(numbers, ",") {
		if n = strings.TrimSpace(n); n != "" {
			houses = append(houses, n)
		}
	}
	for _, rg := range strings.Split(ranges, ";") {
		rg = strings.TrimSuffix(strings.TrimSpace(rg), "-")
		if rg == "" {
			continue
		}
		parts := strings.Split(rg, "-")
		if len(parts) == 1 {
			houses = append(houses, parts[0])
			continue
		}
		if len(parts) != 2 {
			unresolvable++
			continue
		}
		a, errA := strconv.Atoi(strings.TrimSpace(parts[0]))
		b, errB := strconv.Atoi(strings.TrimSpace(parts[1]))
		if errA != nil || errB != nil || b < a || b-a > MaxRangeSpan {
			unresolvable++
			continue
		}
		for i := a; i <= b; i++ {
			houses = append(houses, strconv.Itoa(i))
		}
	}
	return houses, unresolvable
}
