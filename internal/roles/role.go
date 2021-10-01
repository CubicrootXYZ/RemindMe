package roles

// Role a matrix user has
type Role string

const (
	// RoleAdmin user needs to be set in the config file
	RoleAdmin = Role("admin")
	// RoleUser is the default role
	RoleUser = Role("user")
)
