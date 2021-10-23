package main

import (
	"fmt"

	"github.com/appened/note"
)

// TODO: Add marking done
// TODO: Add editing note
// TODO: Add surfacing a note

func main() {
	folios, err := note.LoadFolios()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(folios)
}
