package domain

type Rights int

const (
	All Rights = 1 << iota
	U   Rights = 2
	M   Rights = 4
	A   Rights = 8
	O   Rights = 16
)

func (this Rights) HasFlag(flag Rights) bool {
	return this|flag == this
}
