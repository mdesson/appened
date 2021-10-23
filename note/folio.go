package note

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

type Folio struct {
	Name     string
	Notes    []Note
	filename string
	mu       *sync.RWMutex
}

func LoadFolios() ([]*Folio, error) {
	// Get file names
	files, err := os.ReadDir("../data/")
	if err != nil {
		return nil, err
	}

	// Init folios
	folios := []*Folio{}

	// Fetch each folio
	for _, file := range files {
		filename := file.Name()
		folio, err := parseFolioCSV(filename)
		if err != nil {
			return nil, err
		}
		folios = append(folios, folio)
	}

	return folios, nil
}

func parseFolioCSV(filename string) (*Folio, error) {
	// Open file for reading
	filePath := fmt.Sprintf("../data/%s", filename)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Get folio's name
	name := filename[:len(filename)-4]

	// Get all records as strings
	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Create Notes from csv string records
	notes := []Note{}
	for i, record := range records {
		note := Note{}
		note.index = i
		note.Text = record[0]
		note.Done, err = strconv.ParseBool(record[1])
		if err != nil {
			return nil, err
		}
		note.DateCreated, err = strconv.ParseInt(record[2], 10, 64)
		if err != nil {
			return nil, err
		}
		note.DateDone, err = strconv.ParseInt(record[3], 10, 64)
		if err != nil {
			return nil, err
		}
		note.DateEdited, err = strconv.ParseInt(record[4], 10, 64)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	// Create folio
	folio := Folio{name, notes, filePath, &sync.RWMutex{}}

	return &folio, nil
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
