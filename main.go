package main

import (
	"fmt"
)

func main() {
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
