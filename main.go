// +---------------------------------------------------+
// | Author: Niklas Eichkorn
// | Date: 07.08.25
// | Version: 1.0
// |---------------------------------------------------+
// | Notes: merke https://gist.github.com/jordansissel/1e08b1c65157bde0f30a87c4fb569237
// +---------------------------------------------------+

package main

import (
	"fmt"

	"github.com/bit101/go-ansi"
)

func main() {
	print_start_screen()

}

func print_start_screen() {
	ansi.ClearScreen()
	ansi.SetBg(ansi.Blue)
	fmt.Println("Hello World")
}
