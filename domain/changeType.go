package domain

type ChangeType int

// Add new action only int the end
const (
	AddPlan    ChangeType = 1
	EditPlan   ChangeType = 2
	DeletePlan ChangeType = 3

	AddTopic    ChangeType = 4
	EditTopic   ChangeType = 5
	DeleteTopic ChangeType = 6

	AddProject    ChangeType = 7
	EditProject   ChangeType = 8
	DeleteProject ChangeType = 9

	AddResource    ChangeType = 10
	EditResource   ChangeType = 11
	DeleteResource ChangeType = 12

	AddComment    ChangeType = 13
	EditComment   ChangeType = 14
	DeleteComment ChangeType = 15

	AddUser    ChangeType = 16
	EditUser   ChangeType = 17
	DeleteUser ChangeType = 18
)
