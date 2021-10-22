package note

import (
	"encoding/csv"
	"fmt"
	"os"
	"sync"
	"time"
)

type Folio struct {
	Name     string
	Notes    []Note
	filename string
	mu       *sync.RWMutex
}

func LoadFolios() (*[]Folio, error) {
	return nil, nil
}

func CreateFolio(name string) (*Folio, error) {
	fileName := fmt.Sprintf("../data/%s.csv", name)
	mu := &sync.RWMutex{}

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return nil, err
	}
	file.Close()

	f := &Folio{name, []Note{}, fileName, mu}

	return f, nil
}

func (f *Folio) Append(note string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	now := time.Now().Unix()
	n := Note{len(f.Notes), false, note, now, now, now}

	file, err := os.OpenFile(f.filename, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err = writer.Write(n.csvLine()); err != nil {
		return err
	}
	f.Notes = append(f.Notes, n)

	return nil
}
