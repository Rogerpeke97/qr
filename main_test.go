package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
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

	r = getAlphaVal(245)
	if r != 233 {
		fmt.Printf("\nWrong value. Got: %d. Expected %d\n", r, 233)
		t.Error("LOL")
	}
}

func TestGetExponentAndValFromCoefficient(t *testing.T) {
	exp, val := getExponentAndValFromCoefficient(7)
	if exp != 198 {
		fmt.Println(val)
		fmt.Printf("\nWrong value. Got: %d. Expected %d\n", exp, 198)
		t.Error("LOL")
	}

	exp, val = getExponentAndValFromCoefficient(14)
	if exp != 199 {
		fmt.Println(val)
		fmt.Printf("\nWrong value. Got: %d. Expected %d\n", exp, 199)
		t.Error("LOL")
	}

	exp, val = getExponentAndValFromCoefficient(127)
	if exp != 87 {
		fmt.Println(val)
		fmt.Printf("\nWrong value. Got: %d. Expected %d\n", exp, 87)
		t.Error("LOL")
	}

	exp, val = getExponentAndValFromCoefficient(120)
	if exp != 78 {
		fmt.Println(val)
		fmt.Printf("\nWrong value. Got: %d. Expected %d\n", exp, 78)
		t.Error("LOL")
	}
}

func TestGetCoefficientIfAlphaBig(t *testing.T) {
	c := getCoefficientIfAlphaBig(250, 250, true)
	if c != 0 {
		fmt.Printf("\nWrong value. Got: %d. Expected %d\n", c, 0)
		t.Error("LOL")
	}

	c = getCoefficientIfAlphaBig(250, 250, false)
	if c != 243 {
		fmt.Printf("\nWrong value. Got: %d. Expected %d\n", c, 243)
		t.Error("LOL")
	}

}

/*
(7)
x^7 + 127x^6 + 122x^5 + 154x^4 + 164x^3 + 11x^2 + 68x + 117

(10)
α0x10 + α251x9 + α67x8 + α46x7 + α61x6 + α118x5 + α70x4 + α64x3 + α94x2 + α32x + α45

x^10 + 216x^9 + 194x^8 + 159x^7 + 111x^6 + 199x^5 + 94x^4 + 95x^3 + 113x^2 + 157x + 193
*/
func TestGenGeneratorPolynomial(t *testing.T) {
	r := genGeneratorPolynomial(30)
	if len(r) > 0 {
		fmt.Printf("\nGen polynomial is %+v\n", r)
		t.Error("LOL")
	}
}

// To parse and convert into x form
func parseEquation(equation string) string {
	re := regexp.MustCompile(`α(\d+)`)
	matches := re.FindAllStringSubmatch(equation, -1)

	var terms []string
	exponent := len(matches) - 1
	for _, match := range matches {
		alphaVal, _ := strconv.Atoi(match[1])
		coefficient := getAlphaVal(alphaVal)
		term := fmt.Sprintf("%dx^%d", coefficient, exponent)
		terms = append(terms, term)
		if exponent > 0 {
			exponent--
		}
	}

	return strings.Join(terms, " + ")
}
