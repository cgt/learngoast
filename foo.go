package learngoast

import "fmt"

func IsApplesauce(x string) bool {
	return x == "applesauce"
}

func PrintIfApplesauce(x string) {
	if IsApplesauce(x) {
		fmt.Println("It is applesauce!")
	}
}
