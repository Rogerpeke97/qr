package main

import (
	"fmt"
	"strconv"
)

// x and degree
func findMultiplierOfHighestDegree(
	f string,
) (string, int) {
	power_of_code := byte('^')
	x_code := byte('x')
	highest_degree_f := 0
	highest_degree_f_x := ""
	ascii_integer_min_max := []int{48, 57}

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

func findGcd(
	f string,
	g string,
	fp uint,
) string {

}

func encoder() {
	// url_to_encode := "https://github.com/Rogerpeke97"
	url_to_encode := "AB"
	var bin_representation [][]int
	encoded := []byte(url_to_encode)
	for i := 0; i < len(encoded); i++ {
		bin_representation = append(bin_representation, []int{})
		// big endian so opposite
		for j := 7; j >= 0; j-- {
			mask := byte(1 << uint(j))
			bin := encoded[i] & mask
			var bin_to_add int
			if bin > 0 {
				bin_to_add = 1
			} else {
				bin_to_add = 0
			}

			bin_representation[i] = append(bin_representation[i], bin_to_add)
		}
	}

	fmt.Printf("\nBin representation is: %+v\n", bin_representation)
}

func main() {
	// f := "x^3 - 8"
	// g := "3x^2 - 6x"

}
