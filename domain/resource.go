package domain

type Resource struct {
	Id                   string
	Title                string
	Identifier           string
	NormalizedIdentifier string // Identifier in uppercase
	Type                 string // ResourceType
	Properties           string // json
}
