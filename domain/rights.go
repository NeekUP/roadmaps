package domain

type Rights int

const (
	All Rights = 0 << 0
	U   Rights = 1 << 1
	M   Rights = 1 << 2
	A   Rights = 1 << 3
	O   Rights = 1 << 4
)

func (right Rights) HasFlag(flag Rights) bool {
	return right|flag == right
}
