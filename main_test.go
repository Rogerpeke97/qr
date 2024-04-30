package main

import (
	"fmt"
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

func TestGetAlphaVal(t *testing.T) {
	r := getAlphaVal(8)
	if r != 29 {
		fmt.Printf("\nWrong value. Got: %d. Expected %d\n", r, 29)
		t.Error("LOL")
	}

	r = getAlphaVal(9)
	if r != 58 {
		fmt.Printf("\nWrong value. Got: %d. Expected %d\n", r, 58)
		t.Error("LOL")
	}

	r = getAlphaVal(10)
	if r != 116 {
		fmt.Printf("\nWrong value. Got: %d. Expected %d\n", r, 116)
		t.Error("LOL")
	}

	r = getAlphaVal(11)
	if r != 232 {
		fmt.Printf("\nWrong value. Got: %d. Expected %d\n", r, 232)
		t.Error("LOL")
	}

	r = getAlphaVal(12)
	if r != 205 {
		fmt.Printf("\nWrong value. Got: %d. Expected %d\n", r, 205)
		t.Error("LOL")
	}
}

func TestGenGeneratorPolynomial(t *testing.T) {
	// r := genGeneratorPolynomial(3)
	// if len(r) > 0 {
	// 	fmt.Printf("\nGen polynomial is %+v\n", r)
	// 	t.Error("LOL")
	// }
}
