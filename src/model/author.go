package model

type Author struct {
	Id   int64
	Name string
}

type OptionalAuthor struct {
	Id   *int64
	Name *string
}
