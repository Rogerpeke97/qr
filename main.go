package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"strconv"
)

var power_of_code = byte('^')
var x_code = byte('x')
var ascii_integer_min_max = []int{48, 57}

// Encoded data starts with the mode
var BYTE_MODE_INDICATOR = "0100"
var WIDTH = 32
var HEIGHT = 32

// x and degree
func findMultiplierOfHighestDegree(
	f string,
) (string, int) {
	highest_degree_f := 0
	highest_degree_f_x := ""
	f_len := len(f)

	for i := 0; i < f_len; i++ {
		if f[i] == x_code {
			if i+1 < f_len && f[i+1] == power_of_code && i+2 < f_len {
				var deg_byte []byte
				stop := false
				idx_of_degree := i + 2
				for !stop {
					if idx_of_degree >= f_len {
						stop = true
						continue
					}

					not_an_integer := f[idx_of_degree] < byte(ascii_integer_min_max[0]) || f[idx_of_degree] > byte(ascii_integer_min_max[1])
					if not_an_integer {
						stop = true
						continue
					}

					deg_byte = append(deg_byte, f[idx_of_degree])
					idx_of_degree++
				}
				if len(deg_byte) == 0 {
					fmt.Printf("\nInvalid character after power of at char %d of f: %s\n", i+2, f)
					panic("Error parsing polynomial")
				}

				d, err := strconv.Atoi(string(deg_byte))
				if err != nil {
					fmt.Printf("\nInvalid character after power of at char %d of f: %s\n", i+2, f)
					panic("Error parsing polynomial")
				}

				if d > highest_degree_f {
					highest_degree_f = d
					highest_degree_f_x = f[i:i+2] + string(deg_byte)
					continue
				}
			}

			if 1 > highest_degree_f {
				highest_degree_f = 1
				highest_degree_f_x = string(f[i])
			}
		}

	}

	return highest_degree_f_x, highest_degree_f
}

func getAlphaVal(
	exp int,
) int {
	base := 2
	curr_exp := -1
	acc := 0

	for curr_exp < exp {
		curr_exp++
		r := acc * base
		if r >= 256 {
			acc = r ^ 285
		} else {
			if r == 0 {
				r++
			}

			acc = r
		}
	}

	return acc
}

// exp and val
func getExponentAndValFromCoefficient(
	c int,
) (int, int) {
	if c == 0 {
		return 0, 0
	}
	base := 2

	acc := 0

	exp := -1
	for {
		exp++
		r := acc * base
		if r >= 256 {
			acc = r ^ 285
		} else {
			if r == 0 {
				r++
			}

			acc = r
		}

		if acc == c {
			break
		}
	}

	return exp, acc
}

type PolynomialMember struct {
	Exp         int
	Coefficient int
	IsX         bool
}

func solveSameMembersAndUpdateAlpha(
	polynomial *[]PolynomialMember,
) {
	i := 0
	// make sure to solve all members of the same family
	for i < len(*polynomial) {
		m := (*polynomial)[i]

		for j, m2 := range *polynomial {
			if m.Exp == m2.Exp && m.IsX == m2.IsX && j != i {
				c := getCoefficientIfAlphaBig(m.Coefficient, m2.Coefficient, true)
				new_m := PolynomialMember{
					Exp:         m.Exp,
					Coefficient: c,
					IsX:         m.IsX,
				}
				(*polynomial)[i] = new_m
				*polynomial = append((*polynomial)[:j], (*polynomial)[j+1:]...)
				i = -1
				break
			}
		}

		i++
	}
}

func genGeneratorPolynomial(
	err_c_codewords int,
) []PolynomialMember {
	//Because ((a^0x - a^m)*(a^0x - a^m+1)) * (a^0x - a^m+1) *...
	acc := [][]PolynomialMember{
		{
			{
				Exp:         1,
				Coefficient: 1,
				IsX:         true,
			},
			{
				Exp:         1,
				Coefficient: 1,
				IsX:         false,
			},
		},
		nil,
	}
	m := 1

	for i := 1; i < err_c_codewords; i++ {
		if acc[1] == nil {
			acc[1] = []PolynomialMember{
				{
					Exp:         1,
					Coefficient: 1,
					IsX:         true,
				},
				{
					Exp:         1,
					Coefficient: getAlphaVal(m),
					IsX:         false,
				},
			}
			// fmt.Printf("\nAlpha VAL is: %d for m: %d\n", getAlphaVal(m), m)
			m++
		}

		var new_acc_zero []PolynomialMember
		for _, member := range acc[0] {
			for _, member2 := range acc[1] {
				is_x_mult := member.IsX && member2.IsX
				var exp int
				if is_x_mult {
					exp = member.Exp + member2.Exp
				} else {
					exp = member.Exp
				}
				c := getCoefficientIfAlphaBig(member.Coefficient, member2.Coefficient, false)

				new_acc_zero = append(
					new_acc_zero,
					PolynomialMember{
						Exp:         exp,
						Coefficient: c,
						IsX:         member.IsX || member2.IsX,
					},
				)
			}
		}

		// Check if coefficient in alpha exp the exp is not bigger than 255
		acc[0] = new_acc_zero
		acc[1] = nil
		solveSameMembersAndUpdateAlpha(&acc[0])
	}

	return acc[0]
}

func getCoefficientIfAlphaBig(
	c1 int,
	c2 int,
	is_sum bool,
) int {
	if is_sum {
		return c1 ^ c2
	}

	a_exp, _ := getExponentAndValFromCoefficient(c1)
	b_exp, _ := getExponentAndValFromCoefficient(c2)

	n := a_exp + b_exp
	if n > 255 {
		b := int(math.Floor(float64(n) / 256))
		n = (n % 256) + b
	}

	return getAlphaVal(n)
}

// 20 EC codewords per block. 1 group and 80 codwrods total Num of blocks 1
func genMessagePolynomial(
	codewords []string,
) ([]PolynomialMember, string) {
	polynomial_s := ""
	codewords_len := len(codewords)
	var polynomial []PolynomialMember
	for idx, codeword := range codewords {
		d, err := strconv.ParseInt(codeword, 2, 0)
		if err != nil {
			fmt.Printf("\nError, codeword %s failed with: %+v\n", codeword, err)
			panic("FAIL!!")
		}
		if idx < codewords_len-1 {
			var sign string
			if idx > 0 {
				sign = "+"
			}
			exp := (codewords_len - 1) - idx
			polynomial_s += fmt.Sprintf("%s%dx^%d", sign, d, exp)
			polynomial = append(polynomial, PolynomialMember{
				Coefficient: int(d),
				Exp:         exp,
				IsX:         true,
			})
		} else {
			polynomial_s += fmt.Sprintf("+%d", d)
			polynomial = append(polynomial, PolynomialMember{
				Coefficient: int(d),
				Exp:         1,
				IsX:         false,
			})
		}
	}

	return polynomial, polynomial_s
}

func divideIntoCodeWords(
	data string,
) []string {
	// 8 bits. Total = 80
	var codewords []string
	for i := 8; i <= len(data); i += 8 {
		var codeword string
		for j := i - 8; j < i; j++ {
			codeword += string(data[j])
		}

		codewords = append(codewords, codeword)
	}

	return codewords
}

// Error correction level L just to keep it simple for now
// I'll add more later if not lazy :)
// Version will be 4 = 33x33 pixels
// The amount of characters using byte encoding is 78
// char count must be 8 bits in byte mode for versions 1...9
// total data codewords   	EC codewords per block		n of blocks in group 1  	n of data codewords in each of group 1 blocks
//
//	80				20				1				80
//
// data pad bytes 11101100 00010001 if not enough
// returns mode, char_count_indicator, encoded_data, total_bits
func encode(
	str string,
) (string, string, string, int) {
	total_bits := 80 * 8
	char_count_max_bit_long := 8
	char_count := int64(len(str))
	// Goes after mode indicator
	char_count_bin := strconv.FormatInt(char_count, 2)
	char_count_bin_len := len(char_count_bin)
	if char_count_bin_len > char_count_max_bit_long {
		fmt.Printf("\nMax bits for char count exceeded. Need %d, got %d\n", char_count_max_bit_long, len(char_count_bin))
		panic("Failed!")
	}

	if char_count_bin_len < char_count_max_bit_long {
		padding := char_count_max_bit_long - char_count_bin_len
		var pad_zeros string
		for i := 0; i < padding; i++ {
			pad_zeros += "0"
		}

		char_count_bin = pad_zeros + char_count_bin
	}

	var data_bin_str string
	encoded := []byte(str)
	for i := 0; i < len(encoded); i++ {
		// assumes machine stores bin in little endian
		for j := 7; j >= 0; j-- {
			mask := byte(1 << uint(j))
			bin := encoded[i] & mask
			var bin_str_to_add string
			if bin > 0 {
				bin_str_to_add = "1"
			} else {
				bin_str_to_add = "0"
			}

			data_bin_str += bin_str_to_add
		}
	}
	padding := total_bits - (len(data_bin_str) + len(BYTE_MODE_INDICATOR) + len(char_count_bin))

	// add zero terminators first. Max of 4 zeros
	if padding > 0 {
		for i := 0; i < padding; i++ {
			if i > 3 && (len(data_bin_str)+len(BYTE_MODE_INDICATOR)+len(char_count_bin))%8 == 0 {
				break
			}

			data_bin_str += "0"
		}

		padding = total_bits - (len(data_bin_str) + len(BYTE_MODE_INDICATOR) + len(char_count_bin))
	}

	// if still too short add 236 and 17
	if padding > 0 {
		pad_bytes := []string{"11101100", "00010001"}
		idx := 0

		pad_bytes_to_add := padding / 8
		for i := 0; i < pad_bytes_to_add; i++ {
			data_bin_str += pad_bytes[idx]
			if idx > 0 {
				idx = 0
				continue
			}

			idx = 1
		}
	}

	encoded_data_total := len(BYTE_MODE_INDICATOR) + len(char_count_bin) + len(data_bin_str)
	if encoded_data_total > total_bits {
		fmt.Printf("\nTotal data codewords exceeds the permitted amount. Got %d, want %d\n", encoded_data_total, total_bits)
		panic("Failed!")
	}

	return BYTE_MODE_INDICATOR, char_count_bin, data_bin_str, encoded_data_total
}

func getEcc(
	msg_p []PolynomialMember,
	gen_p []PolynomialMember,
) []PolynomialMember {
	//step a || 0
	b := msg_p
	steps := len(msg_p) * 2

	var to_xor []PolynomialMember

	for i := 0; i <= steps; i++ {
		// fmt.Printf("\nSTILL RUNNING ON IDX: %d\n", i)
		// fmt.Printf("\nb IS: %+v\n", b)
		// fmt.Printf("\ngen_p IS: %+v\n", gen_p)
		//Multiply by lead
		if i%2 == 0 {
			for k, m := range gen_p {
				// fmt.Println("HERE")
				// fmt.Println(m.Coefficient, b[0].Coefficient)
				c := getCoefficientIfAlphaBig(m.Coefficient, b[0].Coefficient, false)
				// fmt.Printf("\nC IS %d\n", c)
				to_xor = append(to_xor, PolynomialMember{
					Exp:         b[0].Exp - k,
					Coefficient: c,
					IsX:         true,
				})
			}
			continue
		}

		// fmt.Printf("\nto_xor IS: %+v\n", to_xor)
		//Xor to_xor with b
		for j := range to_xor {
			if j < len(b) {
				c := getCoefficientIfAlphaBig(b[j].Coefficient, to_xor[j].Coefficient, true)
				b[j].Coefficient = c
				continue
			}

			b = append(b, PolynomialMember{
				Exp:         to_xor[j].Exp,
				Coefficient: to_xor[j].Coefficient,
				IsX:         to_xor[j].IsX,
			})
		}

		//Remove lead
		b = append(b[:0], b[1:]...)
		to_xor = nil
	}

	return b
}

func genQrImage() {
	up_left := image.Point{0, 0}
	low_right := image.Point{WIDTH, HEIGHT}
	img := image.NewRGBA(image.Rectangle{up_left, low_right})

	cyan := color.RGBA{100, 200, 200, 0xff}

	for x := 0; x < WIDTH; x++ {
		for y := 0; y < HEIGHT; y++ {
			switch {
			case x < WIDTH/2 && y < HEIGHT/2:
				img.Set(x, y, cyan)
			case x >= WIDTH/2 && y >= HEIGHT/2:
				img.Set(x, y, color.White)
			default:
			}
		}
	}

	f, _ := os.Create("image.png")
	png.Encode(f, img)
}

func main() {
	str := "HELLO WORLD"
	mode, char_count_indicator, data, total_bits := encode(str)
	fmt.Printf("\nMode is: %s\nChar count is: %s\nData is: %s\nTotal bits: %d\n", mode, char_count_indicator, data, total_bits)
	codewords := divideIntoCodeWords(mode + char_count_indicator + data)
	//4-L requires 20 EC codewords per block
	ec_codewords_needed := 20
	fmt.Printf("\nCodewords are: %+v\nAmount of codewords: %d\n", codewords, len(codewords))
	msg_p, msg_p_s := genMessagePolynomial(codewords)
	fmt.Printf("\nMsg Polynomial is %+v\n", msg_p)
	fmt.Printf("\nMsg Polynomial string is %+v\n", msg_p_s)

	gen_p := genGeneratorPolynomial(ec_codewords_needed)
	fmt.Printf("\nGen Polynomial is %+v\n", gen_p)

	// Find if p exp is < codewords | genPolyMaxExp
	if msg_p[0].Exp < ec_codewords_needed {
		increase_msg_by_exp := ec_codewords_needed - msg_p[0].Exp
		for i := range msg_p {
			msg_p[i].Exp += increase_msg_by_exp
		}

	}

	if gen_p[0].Exp < ec_codewords_needed {
		increase_gen_by_exp := ec_codewords_needed - gen_p[0].Exp
		for i := range gen_p {
			gen_p[i].Exp += increase_gen_by_exp
		}
	}

	fmt.Printf("\nECC POLYNOMIAL IS %+v\n", getEcc(msg_p, gen_p))

	genQrImage()
}
