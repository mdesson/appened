package note

import "time"

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
