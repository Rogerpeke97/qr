package main

import (
	"testing"
)

func TestFindMultiplierOfHighestDegree(t *testing.T) {
	x, deg := findMultiplierOfHighestDegree("x^3 - 8")
	if x != "x^3" || deg != 3 {
		t.Errorf("\nWrong degree or x. Needed x^3 and deg 3. Got %s and %d\n", x, deg)
	}

	x, deg = findMultiplierOfHighestDegree("x^3 - x^4 - x + 3 - x^7")
	if x != "x^7" || deg != 7 {
		t.Errorf("\nWrong degree or x. Needed x^7 and deg 7. Got %s and %d\n", x, deg)
	}

	x, deg = findMultiplierOfHighestDegree("-8 + 7^2 - 32^^32 - 09 - x + x^32")
	if x != "x^32" || deg != 32 {
		t.Errorf("\nWrong degree or x. Needed x^32 and deg 32. Got %s and %d\n", x, deg)
	}
}
