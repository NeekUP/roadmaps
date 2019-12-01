package domain

type Step struct {
	Id            int64
	PlanId        int
	ReferenceId   int64
	ReferenceType ReferenceType
	Position      int
	Source        *Source
}
