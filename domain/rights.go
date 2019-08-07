package domain

type Rights int

const (
	U Rights = 1 << iota
	M Rights = 2
	A Rights = 4
	O Rights = 8
)

func (this Rights) HasFlag(flag Rights) bool {
	return this|flag == this
}
