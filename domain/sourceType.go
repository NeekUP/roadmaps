package domain

type SourceType string

const (
	Article SourceType = "Article"
	Book    SourceType = "Book"
	Video   SourceType = "Video"
	Audio   SourceType = "Audio"
)

func (r SourceType) IsValid() bool {
	return r == Article ||
		r == Book ||
		r == Video ||
		r == Audio
}
