package domain

type Step struct {
	Id            int
	PlanId        int
	ReferenceId   int
	ReferenceType ReferenceType
	Position      int
}
