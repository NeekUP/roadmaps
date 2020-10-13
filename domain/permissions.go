package domain

type Permissions uint64
type GlobalPermissions uint64

// GLOBAL
const (
	G_CREATE_TOPICS   GlobalPermissions = 1 << 0
	G_REMOVE_TOPICS   GlobalPermissions = 1 << 1
	G_EDIT_TOPICS     GlobalPermissions = 1 << 2
	G_CREATE_PLANS    GlobalPermissions = 1 << 3
	G_REMOVE_PLANS    GlobalPermissions = 1 << 4
	G_EDIT_PLANS      GlobalPermissions = 1 << 5
	G_CREATE_SOURCES  GlobalPermissions = 1 << 6
	G_EDIT_SOURCES    GlobalPermissions = 1 << 7
	G_REMOVE_SOURCES  GlobalPermissions = 1 << 8
	G_CREATE_COMMENTS GlobalPermissions = 1 << 9
	G_CREATE_GROUPS   GlobalPermissions = 1 << 10
	G_EDIT_GROUPS     GlobalPermissions = 1 << 11
	G_REMOVE_GROUPS   GlobalPermissions = 1 << 12

	G_MANAGER     GlobalPermissions = 1 << 15
	G_TOP_MANAGER GlobalPermissions = 1 << 16

	G_MANAGE_RIGHTS        GlobalPermissions = 1 << 61
	G_MANAGE_G_RIGHTS      GlobalPermissions = 1 << 62
	G_MANAGE_HIGH_G_RIGHTS GlobalPermissions = 1 << 64
)

// TOPICS
const (
	TOPIC_EDIT          = 1 << 0
	TOPIC_ADD_PLANS     = 1 << 1
	TOPIC_CONFIRM_PLANS = 1 << 2
)

// PLANS
const (
	PLAN_EDIT              = 1 << 0
	PLAN_ADD_COMMENTS      = 1 << 1
	PLAN_MODERATE_COMMENTS = 1 << 2
)

// SOURCES
const (
	SOURCE_EDIT         = 1 << 0
	SOURCE_ADD_COMMENTS = 1 << 0
)

const (
	COMMENT_EDIT   = 1 << 0
	COMMENT_REMOVE = 1 << 1
)

func (right Permissions) HasFlag(flag Permissions) bool {
	return right|flag == right
}
