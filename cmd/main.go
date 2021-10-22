package main

import (
	"fmt"

	"github.com/appened/note"
)

func main() {
	f, err := note.CreateFolio("okay")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Okay!")
	}
	err = f.Append("wow gotta do it!")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Okay!")
	}
	err = f.Append("wow, gotta do it!")
}
