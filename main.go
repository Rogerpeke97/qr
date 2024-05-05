package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
)

type PolynomialMember struct {
	Exp         int
	Coefficient int
	IsX         bool
}

type QrCoordinate struct {
	X               int    `json:"x"`
	Y               int    `json:"y"`
	Color           string `json:"color"`
	Reserved        bool   `json:"reserved"`
	IsTimingPattern bool   `json:"is_timing_pattern"`
	IsDarkModule    bool   `json:"is_dark_module"`
}

// Encoded data starts with the mode
var FINDER_PATTERN_W_H = 7
var SEPARATOR_W_H = FINDER_PATTERN_W_H + 1
var ALIGNMENT_PATTERN_W_H = 5
var VERTICAL_TIMING_PATTERN_X_COORD = SEPARATOR_W_H - 2

// TODO: CALCULATE IT RATHER THAN HARD CODE. IT'S EASY, JUST BORING
// IM ADDING ONE MORE REPEATED BIT SINCE THE IDX 7th REPEATS IN THE RESERVED SPACES
var FORMAT_AND_VERSION_INFO = "1110111111000100"

// ver 4L in byte mode
var MAX_CHAR_AMOUNT = 78
var VERSION = 4

// total data codewords
var TOTAL_BITS_REQUIRED = 8 * 80
var EC_CODEWORDS_NEEDED = 20

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
					if member.IsX || member2.IsX {
						if member.Exp > member2.Exp {
							exp = member.Exp
						} else {
							exp = member2.Exp
						}
					} else {
						exp = 0
					}
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
		var sign string
		if idx > 0 {
			sign = "+"
		}
		exp := (codewords_len - 1) - idx
		polynomial_s += fmt.Sprintf("%s%dx^%d", sign, d, exp)
		var is_x bool
		if exp > 0 {
			is_x = true
		} else {
			is_x = false
		}

		polynomial = append(polynomial, PolynomialMember{
			Coefficient: int(d),
			Exp:         exp,
			IsX:         is_x,
		})
	}

	return polynomial, polynomial_s
}

func divideIntoCodeWords(
	data string,
) []string {
	// 8 bits. Total = 80
	var codewords []string
	for i := 8; i <= len(data); i += 8 {
		codewords = append(codewords, data[i-8:i])
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
	encoded_msg string,
	decoded_msg_len int,
	byte_mode_indicator string,
	char_count_indicator_req_bits int,
	total_bits_required int,

) string {
	var all_bits string
	all_bits += byte_mode_indicator
	data_bin_char_count := strconv.FormatInt(int64(decoded_msg_len), 2)
	for len(data_bin_char_count) < char_count_indicator_req_bits {
		data_bin_char_count = "0" + data_bin_char_count
	}

	all_bits += data_bin_char_count
	all_bits += encoded_msg

	// add terminators
	if len(all_bits) < total_bits_required {
		for i := 0; i < 4; i++ {
			if len(all_bits) < total_bits_required {
				all_bits += "0"
			}
		}
	}

	for len(all_bits)%8 != 0 {
		all_bits += "0"
	}

	// if still too short add 236 and 17
	if len(all_bits) < total_bits_required {
		pad_bytes := []string{"11101100", "00010001"}
		idx := 0

		pad_bytes_to_add := (total_bits_required - len(all_bits)) / 8
		for i := 0; i < pad_bytes_to_add; i++ {
			all_bits += pad_bytes[idx]
			if idx > 0 {
				idx = 0
				continue
			}

			idx = 1
		}
	}

	return all_bits
}

func encodeMessage(
	message string,
) string {
	var encoded_msg string
	encoded := []byte(message)
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

			encoded_msg += bin_str_to_add
		}
	}

	return encoded_msg
}

func getEcc(
	msg_p []PolynomialMember,
	gen_p []PolynomialMember,
) []PolynomialMember {
	b := msg_p
	steps := len(msg_p) * 2

	var to_xor []PolynomialMember

	for i := 0; i <= steps; i++ {
		//Multiply by lead
		if i%2 == 0 {
			for k, m := range gen_p {
				c := getCoefficientIfAlphaBig(m.Coefficient, b[0].Coefficient, false)
				to_xor = append(to_xor, PolynomialMember{
					Exp:         b[0].Exp - k,
					Coefficient: c,
					IsX:         true,
				})
			}
			continue
		}

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

func addAlignmentPatterns(
	coordinates *[]QrCoordinate,
	x int,
	y int,
) {
	//ROW AND COL FOR VER 4 (6, 26)
	row_col := []int{26 - 2, 26 - 2}

	if x >= row_col[0] && x < row_col[0]+ALIGNMENT_PATTERN_W_H {
		if y >= row_col[1] && y < row_col[1]+ALIGNMENT_PATTERN_W_H {
			if x == row_col[0]+2 && y == row_col[1]+2 {
				*coordinates = append(*coordinates, QrCoordinate{X: x, Y: y, Color: "black", Reserved: true})
				return
			}

			if x >= row_col[0]+1 && x < row_col[0]+4 && y >= row_col[1]+1 && y < row_col[1]+4 {
				*coordinates = append(*coordinates, QrCoordinate{X: x, Y: y, Color: "white", Reserved: true})
				return
			}

			*coordinates = append(*coordinates, QrCoordinate{X: x, Y: y, Color: "black", Reserved: true})
		}
	}

}

func addSeparators(
	coordinates *[]QrCoordinate,
	x int,
	y int,
	width int,
	height int,
) {
	var sp_start_points = [][]int{{0, SEPARATOR_W_H - 1}, {0, height - SEPARATOR_W_H}}

	has_to_paint := false
	for _, sp := range sp_start_points {
		if x >= sp[0] && x < sp[0]+SEPARATOR_W_H {
			if (y >= sp[1] && y < sp[1]+1) || (x == SEPARATOR_W_H-1 && (y < SEPARATOR_W_H || y > height-SEPARATOR_W_H)) {
				has_to_paint = true
			}
		}
	}

	if has_to_paint {
		*coordinates = append(*coordinates, QrCoordinate{X: x, Y: y, Color: "white", Reserved: true})
		//the opposite of the first sp
		if y < SEPARATOR_W_H {
			*coordinates = append(*coordinates, QrCoordinate{X: width - (x + 1), Y: y, Color: "white", Reserved: true})
		}
	}
}

func addFinderPatterns(
	coordinates *[]QrCoordinate,
	x int,
	y int,
	width int,
	height int,
) {
	var curr_fp_color string
	var fp_start_points = [][]int{{0, 0}, {0, height - FINDER_PATTERN_W_H}, {width - FINDER_PATTERN_W_H, 0}}

	has_to_paint := false
	to_reduce := []int{0, 0}
	for _, fp := range fp_start_points {
		if x >= fp[0] && x < fp[0]+FINDER_PATTERN_W_H {
			if y >= fp[1] && y < fp[1]+FINDER_PATTERN_W_H {
				has_to_paint = true
				to_reduce[0] = fp[0]
				to_reduce[1] = fp[1]
			}
		}
	}
	if has_to_paint {
		reduced_x := x - to_reduce[0]
		reduced_y := y - to_reduce[1]
		if reduced_y%6 == 0 {
			curr_fp_color = "black"
		} else {
			is_inner_square := reduced_x > 1 && reduced_x < 5 && reduced_y > 1 && reduced_y < 5
			if reduced_x%6 == 0 || is_inner_square {
				curr_fp_color = "black"
			} else {
				curr_fp_color = "white"
			}
		}

		*coordinates = append(*coordinates, QrCoordinate{X: x, Y: y, Color: curr_fp_color, Reserved: true})
	}
}

func determinePixelsFromVersion(
	version int,
) int {
	px := 17
	for i := 0; i < version; i++ {
		px += 4
	}

	return px
}

func addTimingPatterns(
	coordinates *[]QrCoordinate,
	x int,
	y int,
	width int,
	height int,
) {
	left_strip_y_range := []int{SEPARATOR_W_H, width - SEPARATOR_W_H}
	left_strip_x_range := []int{SEPARATOR_W_H - 2, SEPARATOR_W_H - 1}

	top_strip_x_range := []int{SEPARATOR_W_H, height - SEPARATOR_W_H}
	top_strip_y_range := []int{SEPARATOR_W_H - 2, SEPARATOR_W_H - 1}

	relative_y_left_idx := y - left_strip_y_range[0]
	relative_x_top_idx := x - top_strip_x_range[0]
	if x >= left_strip_x_range[0] && x < left_strip_x_range[1] {
		if y >= left_strip_y_range[0] && y < left_strip_y_range[1] {
			var curr_s_color string
			if relative_y_left_idx%2 == 0 {
				curr_s_color = "black"
			} else {
				curr_s_color = "white"
			}

			*coordinates = append(*coordinates, QrCoordinate{
				X:               x,
				Y:               y,
				Color:           curr_s_color,
				Reserved:        true,
				IsTimingPattern: true,
			})
		}
	}

	if x >= top_strip_x_range[0] && x < top_strip_x_range[1] {
		if y >= top_strip_y_range[0] && y < top_strip_y_range[1] {
			var curr_s_color string
			if relative_x_top_idx%2 == 0 {
				curr_s_color = "black"
			} else {
				curr_s_color = "white"
			}

			*coordinates = append(*coordinates, QrCoordinate{
				X:               x,
				Y:               y,
				Color:           curr_s_color,
				Reserved:        true,
				IsTimingPattern: true,
			})
		}
	}
}

func addDarkModuleAndReservedSpaces(
	coordinates *[]QrCoordinate,
	x int,
	y int,
	width int,
	height int,
) {
	dark_module_coordinate := []int{(VERSION * 4) + 9, 8}
	if x == dark_module_coordinate[0] && y == dark_module_coordinate[1] {
		*coordinates = append(*coordinates, QrCoordinate{
			X:            x,
			Y:            y,
			Color:        "black",
			Reserved:     true,
			IsDarkModule: true,
		})
	}

	var rs_range = [][]int{{0, 0}, {SEPARATOR_W_H, SEPARATOR_W_H}}
	has_to_paint := false

	if x >= rs_range[0][0] && x <= rs_range[1][0] {
		if x < rs_range[1][0] {
			if y == rs_range[1][1] {
				has_to_paint = true
			}
		} else {
			if y <= rs_range[1][1] {
				has_to_paint = true
			}
		}
	}

	if has_to_paint {
		*coordinates = append(*coordinates, QrCoordinate{X: x, Y: y, Color: "blue", Reserved: true})
		//the opposite of the first sp

		if y == rs_range[1][1] && width-x < width {
			*coordinates = append(*coordinates, QrCoordinate{X: width - x, Y: y, Color: "blue", Reserved: true})
		}
		if x == rs_range[1][0] && y <= rs_range[1][1] && height-y < height {
			*coordinates = append(*coordinates, QrCoordinate{X: x, Y: height - y, Color: "blue", Reserved: true})
		}
	}

}

func addDataBits(
	data string,
	coordinates *[]QrCoordinate,
	width int,
	height int,
) {
	// TODO: EVALUATE THE MASK NUMBER WITH THE PENALTY SCORES
	direction := "up"
	y_start_idx := (height - 1) * -1
	y_less_than := 1
	var can_paint_at_x []int
	var data_color string
	no_more_data := false

	for x := width - 1; x >= 0; x -= 2 {
		if no_more_data {
			break
		}

		if x == VERTICAL_TIMING_PATTERN_X_COORD {
			x--
		}

		for y := y_start_idx; y < y_less_than; y++ {
			if len(data) == 0 {
				no_more_data = true
				break
			}

			can_paint_at_x = []int{x, x - 1}
			y_abs := y
			if y < 0 {
				y_abs = y * -1
			}

			for _, coord := range *coordinates {
				if y_abs == coord.Y {
					if can_paint_at_x[0] == coord.X {
						can_paint_at_x[0] = -1
					}
					if can_paint_at_x[1] == coord.X {
						can_paint_at_x[1] = -1
					}

				}
			}

			if can_paint_at_x[0] != -1 {
				// fmt.Printf("ADDING DATA. y is: %d, y_start_idx: %d, y_less_than: %d\n", y, y_start_idx, y_less_than)
				//ascii 48 == 0 and 49 == 1
				if data[0] == 48 {
					// Add mask number 3
					// (row + column) mod 3 == 0 => flip bit
					if (can_paint_at_x[0]+y)%2 == 0 {
						data_color = "black"
					} else {
						data_color = "white"
					}
				} else {
					if (can_paint_at_x[0]+y)%2 == 0 {
						data_color = "white"
					} else {
						data_color = "black"
					}
				}

				*coordinates = append(*coordinates, QrCoordinate{X: x, Y: y_abs, Color: data_color})
				data = data[1:]
			}

			if can_paint_at_x[1] != -1 {
				if can_paint_at_x[1] >= 0 && len(data) >= 1 {
					if data[0] == 48 {
						if (can_paint_at_x[1]+y)%2 == 0 {
							data_color = "black"
						} else {
							data_color = "white"
						}
					} else {

						if (can_paint_at_x[1]+y)%2 == 0 {
							data_color = "white"
						} else {
							data_color = "black"
						}
					}

					*coordinates = append(*coordinates, QrCoordinate{X: can_paint_at_x[1], Y: y_abs, Color: data_color})
					data = data[1:]
				}
			}

		}

		if direction == "up" {
			direction = "down"
			y_start_idx = 0
			y_less_than = height
		} else {
			direction = "up"
			y_start_idx = (height - 1) * -1
			y_less_than = 1
		}
	}
}

func addPatterns(
	coordinates *[]QrCoordinate,
	width int,
	height int,
) {
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			addFinderPatterns(coordinates, x, y, width, height)
			addSeparators(coordinates, x, y, width, height)
			addAlignmentPatterns(coordinates, x, y)
			addDarkModuleAndReservedSpaces(coordinates, x, y, width, height)
			addTimingPatterns(coordinates, x, y, width, height)
		}
	}
}

func genCoordinatesForStaticPatterns(
	width int,
	height int,
) []QrCoordinate {
	var coordinates []QrCoordinate
	addPatterns(&coordinates, width, height)

	return coordinates
}

func fromPolynomialToBits(
	polynomial *[]PolynomialMember,
) string {
	var bin string

	for _, m := range *polynomial {
		c_bin := strconv.FormatInt(int64(m.Coefficient), 2)
		if len(c_bin) < 8 {
			var pad string
			for i := 0; i < 8-len(c_bin); i++ {
				pad += "0"
			}

			bin += pad + c_bin
			continue
		}

		bin += c_bin
	}

	return bin
}

func addFormatAndVersionInfo(
	coordinates *[]QrCoordinate,
	width int,
	height int,
) {
	horizontal_reserved_space_y := SEPARATOR_W_H
	vertical_reserved_space_x := SEPARATOR_W_H
	horizontal_bits_to_add := FORMAT_AND_VERSION_INFO
	vertical_bits_to_add := FORMAT_AND_VERSION_INFO

	for x := 0; x < width; x++ {
		add_bit_at_coords_idx := -1
		for i, c := range *coordinates {
			if c.X == x && c.Y == horizontal_reserved_space_y {
				if c.Reserved {
					add_bit_at_coords_idx = i
				}
				if c.IsTimingPattern || c.IsDarkModule {
					add_bit_at_coords_idx = -1
					break
				}
			}
		}

		if add_bit_at_coords_idx > -1 && len(horizontal_bits_to_add) > 0 {
			if horizontal_bits_to_add[0] == 48 {
				(*coordinates)[add_bit_at_coords_idx].Color = "white"
			} else {
				(*coordinates)[add_bit_at_coords_idx].Color = "black"
			}

			horizontal_bits_to_add = horizontal_bits_to_add[1:]
		}
	}

	for y := height - 1; y >= 0; y-- {
		add_bit_at_coords_idx := -1
		for i, c := range *coordinates {
			if c.X == vertical_reserved_space_x && c.Y == y {
				if c.Reserved {
					add_bit_at_coords_idx = i
				}
				if c.IsTimingPattern || c.IsDarkModule {
					add_bit_at_coords_idx = -1
					break
				}
			}
		}

		if add_bit_at_coords_idx > -1 && len(vertical_bits_to_add) > 0 {
			if vertical_bits_to_add[0] == 48 {
				(*coordinates)[add_bit_at_coords_idx].Color = "white"
			} else {
				(*coordinates)[add_bit_at_coords_idx].Color = "black"
			}

			vertical_bits_to_add = vertical_bits_to_add[1:]
		}
	}
}

func getCharCountIndicatorBitsLenAndModeIndicatorBits(
	indicator_mode string,
) (int, string) {
	if indicator_mode == "byte_mode" {
		return 8, "0100"
	}

	// for alphanumeric
	return 9, "0010"
}

// return coords and pixels
func GenQrCode(
	message string,
	version int,
	indicator_mode string,
	total_bits_required int,
	ec_codewords_needed int,
) (*[]QrCoordinate, int) {
	pixels := determinePixelsFromVersion(version)
	encoded_msg := encodeMessage(message)
	// encoded_data := encode(
	// 	encoded_msg,
	// 	len(data),
	// 	"0010",
	// 	9,
	// 	104,
	// )
	indicator_bits_len, indicator_bits := getCharCountIndicatorBitsLenAndModeIndicatorBits(indicator_mode)
	encoded_data := encode(
		encoded_msg,
		len(message),
		indicator_bits,
		indicator_bits_len,
		total_bits_required,
	)

	fmt.Printf("\nEncoded DATA IS: %+v\n", encoded_data)
	codewords := divideIntoCodeWords(encoded_data)
	msg_p, _ := genMessagePolynomial(codewords)
	gen_p := genGeneratorPolynomial(ec_codewords_needed)

	// multiply msg polynomial by ec needed
	for i := range msg_p {
		msg_p[i].Exp += ec_codewords_needed
	}

	msg_bits := fromPolynomialToBits(&msg_p)
	ecc := getEcc(msg_p, gen_p)

	fmt.Printf("\nMSG IN BITS: %s\n", msg_bits)
	fmt.Printf("\nECC IN BITS: %s\n", fromPolynomialToBits(&ecc))
	//add msg_p bits and ecc
	final_bin := msg_bits + fromPolynomialToBits(&ecc)
	//Remainder bits for ver 4 are 7
	// final_bin += "0000000"
	fmt.Printf("\nFINAL MESSAGE: %s\nLen of it: %d\n", final_bin, len(final_bin))

	coordinates := genCoordinatesForStaticPatterns(
		pixels,
		pixels,
	)
	addDataBits(final_bin, &coordinates, pixels, pixels)
	addFormatAndVersionInfo(&coordinates, pixels, pixels)

	return &coordinates, pixels
}

func main() {
	coordinates, pixels := GenQrCode(
		"https://www.youtube.com/watch?v=dQw4w9WgXcQ",
		VERSION,
		"byte_mode",
		TOTAL_BITS_REQUIRED,
		EC_CODEWORDS_NEEDED,
	)

	url := "http://localhost:8080"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/coordinates", func(w http.ResponseWriter, r *http.Request) {
		var coords *[]QrCoordinate
		var pxs int
		if r.Method == http.MethodPost {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer r.Body.Close()

			var new_msg struct{ Message string }
			err = json.Unmarshal(body, &new_msg)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			coords, pxs = GenQrCode(
				new_msg.Message,
				VERSION,
				"byte_mode",
				TOTAL_BITS_REQUIRED,
				EC_CODEWORDS_NEEDED,
			)
		} else {
			coords = coordinates
			pxs = pixels
		}

		json_coordinates, err := json.Marshal(map[string]interface{}{
			"coordinates": coords,
			"pixels":      pxs,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(json_coordinates)
	})

	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		fmt.Printf("Error opening URL: %v", err)
		return
	}

	fmt.Println("Opening URL:", url)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
