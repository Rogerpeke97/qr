package main

import (
	"fmt"
	"strconv"
)

var power_of_code = byte('^')
var x_code = byte('x')
var ascii_integer_min_max = []int{48, 57}

// Encoded data starts with the mode
var BYTE_MODE_INDICATOR = "0100"

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

// 20 codewords per block. 1 group and 80 codwrods total Num of blocks 1
func genMessagePolynomial(
	codewords []string,
) string {
	polynomial := ""
	codewords_len := len(codewords)
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
			polynomial += fmt.Sprintf("%s%dx^%d", sign, d, (codewords_len-1)-idx)
		} else {
			polynomial += fmt.Sprintf("+%d", d)
		}
	}

	return polynomial
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

func main() {
	str := "HELLO WORLD"
	mode, char_count_indicator, data, total_bits := encode(str)
	fmt.Printf("\nMode is: %s\nChar count is: %s\nData is: %s\nTotal bits: %d\n", mode, char_count_indicator, data, total_bits)
	codewords := divideIntoCodeWords(mode + char_count_indicator + data)
	fmt.Printf("\nCodewords are: %+v\nAmount of codewords: %d\n", codewords, len(codewords))
	fmt.Println(genMessagePolynomial(codewords))
}
