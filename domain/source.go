package domain

type Source struct {
	Id                   string
	Title                string
	Identifier           string
	NormalizedIdentifier string // Identifier in uppercase
	Type                 ResourceType
	Properties           string // json
}
