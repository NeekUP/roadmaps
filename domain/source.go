package domain

type Source struct {
	Id                   int64
	Title                string
	Identifier           string
	NormalizedIdentifier string // Identifier in uppercase
	Type                 SourceType
	Properties           string // json
	Img                  string
	Desc                 string
}
