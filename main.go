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

// Error correction level L just to keep it simple for now
// I'll add more later if not lazy :)
// Version will be 4 = 33x33 pixels
// The amount of characters using byte encoding is 78
// char count must be 8 bits in byte mode for versions 1...9
// returns mode, char_count_indicator, encoded_data
func encode(
	str string,
) (string, string, string) {
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

	fmt.Printf("\nBin representation is: %s\n", data_bin_str)
	return BYTE_MODE_INDICATOR, char_count_bin, data_bin_str
}

func main() {
	str := "HELLO WORLD"
	mode, char_count_indicator, data := encode(str)
	fmt.Printf("\nMode is: %s\nChar count is: %s\nData is: %s\n", mode, char_count_indicator, data)
}
