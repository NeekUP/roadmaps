package domain

type Source struct {
	Id                   int
	Title                string
	Identifier           string
	NormalizedIdentifier string // Identifier in uppercase
	Type                 SourceType
	Properties           string // json
	Img                  string
	Desc                 string
}
