package note

type Folio struct {
	Name     string
	Notes    []Note
	filename string
	// TODO: Add a channel to push updates into or use a mutex, or one shared channel for all folios
}

func CreateFolio(name string) (*Folio, error) {
	return nil, nil
}
