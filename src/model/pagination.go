package model

type Page struct {
	Page int
	Size int
}

type Sort struct {
	// Name of the record that used for sorting
	Key string
	// Direction of sort. Any negative number means descending, any positive number or zero means ascending.
	Dir int8
}

type Pagination[T any] struct {
	Data  T
	Count int
}

func (s *Sort) DirString() string {
	if s.Dir < 0 {
		return "DESC"
	}

	return "ASC"
}
