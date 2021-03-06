package note

import (
	"fmt"
	"strconv"
	"time"
)

// Note is a single item appended to a Folio
type Note struct {
	index       int    // Note's index in the Folio
	Done        bool   // Is the note marked Done
	Text        string // Text of the note
	DateCreated int64  // Date of the note's creation
	DateDone    int64  // Date the note was marked done
	DateEdited  int64  // Date the note was last edited
}

// Index returns the note's index in the Folio
// Because 'Appened is append-only, this value is constant
func (n Note) Index() int {
	return n.index
}

// Edit will update the note's text and DateEdited
func (n *Note) Edit(text string) {
	n.Text = text
	n.DateEdited = time.Now().Unix()
}

// ToggleDone will toggle Done between True/False and if True, will update DateDone
func (n *Note) ToggleDone() {
	n.Done = !n.Done
	if n.Done {
		n.DateDone = time.Now().Unix()
	}
}

func (n Note) String() string {
	return fmt.Sprintf("%v,%v,%v,%v,%v", n.Text, n.Done, n.DateCreated, n.DateDone, n.DateEdited)
}

func (n Note) csvLine() []string {
	return []string{
		n.Text,
		strconv.FormatBool(n.Done),
		strconv.FormatInt(n.DateCreated, 10),
		strconv.FormatInt(n.DateDone, 10),
		strconv.FormatInt(n.DateEdited, 10),
	}
}

func (n Note) ListString() string {
	s := fmt.Sprintf("%v. %v", n.index+1, n.Text)
	if n.Done {
		s += " ✅"
	}
	return s
}
