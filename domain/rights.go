package domain

type Rights int

const (
	U   Rights = 1 << 0
	M   Rights = 1 << 1
	A   Rights = 1 << 2
	O   Rights = 1 << 3
	God Rights = U | M | A | O
)

func (right Rights) HasFlag(flag Rights) bool {
	return right|flag == right
}
