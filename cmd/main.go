package main

import (
	"fmt"

	"github.com/appened/note"
)

func main() {
	folios, err := note.LoadFolios()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(folios)
}
