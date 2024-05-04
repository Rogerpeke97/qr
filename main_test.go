package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestGetAlphaVal(t *testing.T) {
	r := getAlphaVal(8)
	if r != 29 {
		fmt.Printf("\nWrong value.\nGot: %d.\nExpected %d\n", r, 29)
		t.Error("Fail!")
	}

	r = getAlphaVal(9)
	if r != 58 {
		fmt.Printf("\nWrong value.\nGot: %d.\nExpected %d\n", r, 58)
		t.Error("Fail!")
	}

	r = getAlphaVal(10)
	if r != 116 {
		fmt.Printf("\nWrong value.\nGot: %d.\nExpected %d\n", r, 116)
		t.Error("Fail!")
	}

	r = getAlphaVal(11)
	if r != 232 {
		fmt.Printf("\nWrong value.\nGot: %d.\nExpected %d\n", r, 232)
		t.Error("Fail!")
	}

	r = getAlphaVal(12)
	if r != 205 {
		fmt.Printf("\nWrong value.\nGot: %d.\nExpected %d\n", r, 205)
		t.Error("Fail!")
	}

	r = getAlphaVal(245)
	if r != 233 {
		fmt.Printf("\nWrong value.\nGot: %d.\nExpected %d\n", r, 233)
		t.Error("Fail!")
	}
}

func TestGetExponentAndValFromCoefficient(t *testing.T) {
	exp, val := getExponentAndValFromCoefficient(7)
	if exp != 198 {
		fmt.Println(val)
		fmt.Printf("\nWrong value.\nGot: %d.\nExpected %d\n", exp, 198)
		t.Error("Fail!")
	}

	exp, val = getExponentAndValFromCoefficient(14)
	if exp != 199 {
		fmt.Println(val)
		fmt.Printf("\nWrong value.\nGot: %d.\nExpected %d\n", exp, 199)
		t.Error("Fail!")
	}

	exp, val = getExponentAndValFromCoefficient(127)
	if exp != 87 {
		fmt.Println(val)
		fmt.Printf("\nWrong value.\nGot: %d.\nExpected %d\n", exp, 87)
		t.Error("Fail!")
	}

	exp, val = getExponentAndValFromCoefficient(120)
	if exp != 78 {
		fmt.Println(val)
		fmt.Printf("\nWrong value.\nGot: %d.\nExpected %d\n", exp, 78)
		t.Error("Fail!")
	}

	exp, val = getExponentAndValFromCoefficient(0)
	if exp != 0 {
		fmt.Println(val)
		fmt.Printf("\nWrong value.\nGot: %d.\nExpected %d\n", exp, 0)
		t.Error("Fail!")
	}
}

func TestGetCoefficientIfAlphaBig(t *testing.T) {
	c := getCoefficientIfAlphaBig(250, 250, true)
	if c != 0 {
		fmt.Printf("\nWrong value.\nGot: %d.\nExpected %d\n", c, 0)
		t.Error("Fail!")
	}

	c = getCoefficientIfAlphaBig(250, 250, false)
	if c != 243 {
		fmt.Printf("\nWrong value.\nGot: %d.\nExpected %d\n", c, 243)
		t.Error("Fail!")
	}

	c = getCoefficientIfAlphaBig(1, 0, false)
	if c != 1 {
		fmt.Printf("\nWrong value.\nGot: %d.\nExpected %d\n", c, 1)
		t.Error("Fail!")
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

func genPolynomialString(
	polynomial []PolynomialMember,
) string {
	var r string
	for i, m := range polynomial {
		var operator string
		if i+1 < len(polynomial) {
			operator = " + "
		}
		exp := m.Exp
		if !m.IsX {
			exp = 0
		}
		r += fmt.Sprintf("%dx^%d%s", m.Coefficient, exp, operator)
	}

	return r
}

func TestGenGeneratorPolynomial(t *testing.T) {
	r := genGeneratorPolynomial(80)
	expect := "1x^80 + 84x^79 + 135x^78 + 16x^77 + 169x^76 + 62x^75 + 204x^74 + 151x^73 + 126x^72 + 108x^71 + 91x^70 + 227x^69 + 174x^68 + 59x^67 + 51x^66 + 79x^65 + 252x^64 + 110x^63 + 45x^62 + 78x^61 + 141x^60 + 107x^59 + 166x^58 + 132x^57 + 131x^56 + 154x^55 + 37x^54 + 63x^53 + 41x^52 + 169x^51 + 231x^50 + 153x^49 + 64x^48 + 117x^47 + 90x^46 + 183x^45 + 142x^44 + 188x^43 + 193x^42 + 173x^41 + 189x^40 + 30x^39 + 224x^38 + 40x^37 + 185x^36 + 119x^35 + 11x^34 + 95x^33 + 133x^32 + 19x^31 + 52x^30 + 22x^29 + 15x^28 + 246x^27 + 236x^26 + 93x^25 + 203x^24 + 81x^23 + 134x^22 + 160x^21 + 131x^20 + 99x^19 + 72x^18 + 43x^17 + 143x^16 + 188x^15 + 66x^14 + 242x^13 + 104x^12 + 123x^11 + 126x^10 + 164x^9 + 77x^8 + 49x^7 + 29x^6 + 137x^5 + 241x^4 + 236x^3 + 89x^2 + 198x^1 + 17x^0"
	r_s := genPolynomialString(r)
	if r_s != expect {
		fmt.Printf("\nWrong value.\nGot: %s.\nExpected %s\n", r_s, expect)
		t.Error("Fail!")
	}

	r = genGeneratorPolynomial(30)
	expect = "1x^30 + 212x^29 + 246x^28 + 77x^27 + 73x^26 + 195x^25 + 192x^24 + 75x^23 + 98x^22 + 5x^21 + 70x^20 + 103x^19 + 177x^18 + 22x^17 + 217x^16 + 138x^15 + 51x^14 + 181x^13 + 246x^12 + 72x^11 + 25x^10 + 18x^9 + 46x^8 + 228x^7 + 74x^6 + 216x^5 + 195x^4 + 11x^3 + 106x^2 + 130x^1 + 150x^0"
	r_s = genPolynomialString(r)
	if r_s != expect {
		fmt.Printf("\nWrong value.\nGot: %s.\nExpected %s\n", r_s, expect)
		t.Error("Fail!")
	}

	r = genGeneratorPolynomial(20)
	expect = "1x^20 + 152x^19 + 185x^18 + 240x^17 + 5x^16 + 111x^15 + 99x^14 + 6x^13 + 220x^12 + 112x^11 + 150x^10 + 69x^9 + 36x^8 + 187x^7 + 22x^6 + 228x^5 + 198x^4 + 121x^3 + 121x^2 + 165x^1 + 174x^0"
	r_s = genPolynomialString(r)
	if r_s != expect {
		fmt.Printf("\nWrong value.\nGot: %s.\nExpected %s\n", r_s, expect)
		t.Error("Fail!")
	}

}

func TestGenEcc(t *testing.T) {
	// 1-M QR
	msg_p := []PolynomialMember{
		{Exp: 25, Coefficient: 32, IsX: true},
		{Exp: 24, Coefficient: 91, IsX: true},
		{Exp: 23, Coefficient: 11, IsX: true},
		{Exp: 22, Coefficient: 120, IsX: true},
		{Exp: 21, Coefficient: 209, IsX: true},
		{Exp: 20, Coefficient: 114, IsX: true},
		{Exp: 19, Coefficient: 220, IsX: true},
		{Exp: 18, Coefficient: 77, IsX: true},
		{Exp: 17, Coefficient: 67, IsX: true},
		{Exp: 16, Coefficient: 64, IsX: true},
		{Exp: 15, Coefficient: 236, IsX: true},
		{Exp: 14, Coefficient: 17, IsX: true},
		{Exp: 13, Coefficient: 236, IsX: true},
		{Exp: 12, Coefficient: 17, IsX: true},
		{Exp: 11, Coefficient: 236, IsX: true},
		{Exp: 10, Coefficient: 17, IsX: true},
	}

	gen_p := genGeneratorPolynomial(10)
	r := getEcc(msg_p, gen_p)
	r_s := genPolynomialString(r)
	expect := "196x^9 + 35x^8 + 39x^7 + 119x^6 + 235x^5 + 215x^4 + 231x^3 + 226x^2 + 93x^1 + 23x^0"
	if r_s != expect {
		fmt.Printf("\nWrong value.\nGot: %s.\nExpected %s\n", r_s, expect)
		t.Error("Fail!")
	}
}

func TestGenMessagePolynomial(t *testing.T) {
	polynomial := []PolynomialMember{
		{Exp: 15, Coefficient: 32, IsX: true},
		{Exp: 14, Coefficient: 91, IsX: true},
		{Exp: 13, Coefficient: 11, IsX: true},
		{Exp: 12, Coefficient: 120, IsX: true},
		{Exp: 11, Coefficient: 209, IsX: true},
		{Exp: 10, Coefficient: 114, IsX: true},
		{Exp: 9, Coefficient: 220, IsX: true},
		{Exp: 8, Coefficient: 77, IsX: true},
		{Exp: 7, Coefficient: 67, IsX: true},
		{Exp: 6, Coefficient: 64, IsX: true},
		{Exp: 5, Coefficient: 236, IsX: true},
		{Exp: 4, Coefficient: 17, IsX: true},
		{Exp: 3, Coefficient: 236, IsX: true},
		{Exp: 2, Coefficient: 17, IsX: true},
		{Exp: 1, Coefficient: 236, IsX: true},
		{Exp: 0, Coefficient: 17, IsX: false},
	}

	binaries := []string{
		"00100000",
		"01011011",
		"00001011",
		"01111000",
		"11010001",
		"01110010",
		"11011100",
		"01001101",
		"01000011",
		"01000000",
		"11101100",
		"00010001",
		"11101100",
		"00010001",
		"11101100",
		"00010001",
	}
	msg_p, msg_p_s := genMessagePolynomial(binaries)
	expected := "32x^15+91x^14+11x^13+120x^12+209x^11+114x^10+220x^9+77x^8+67x^7+64x^6+236x^5+17x^4+236x^3+17x^2+236x^1+17x^0"
	if msg_p_s != expected {
		fmt.Printf("\nWrong value.\nGot: %s.\nExpected %s\n", msg_p_s, expected)
		t.Error("Fail!")
	}

	no_match := true
	for _, m := range msg_p {
		no_match = true
		for _, m2 := range polynomial {
			if m.Exp == m2.Exp && m.IsX == m2.IsX && m.Coefficient == m2.Coefficient {
				no_match = false
				break
			}
		}
		if no_match {
			fmt.Printf("\nNo match for: %+v\n", m)
		}

	}

	if no_match {
		fmt.Printf("\nWrong value.\nGot: %+v.\nExpected %+v\n", msg_p, polynomial)
		t.Error("Fail!")
	}

}
