package shared

// const (
// 	UserTable           = "\"user\""
// 	RoleTable           = "role"
// 	UserRolesTable      = "user_role"
// 	TaskTable           = "task"
// 	AuthorTable         = "author"
// 	SolvedTasksTable    = "solved_task"
// 	TaskSubmissionTable = "task_submission"
// 	TeamTable           = "team"
// 	TeamMembersTable    = "team_member"
// 	SettingsTable       = "settings"
// )

// TODO: Replace constants defined above with Tables struct fields in repositories
var Tables = struct {
	User        string
	Role        string
	UserRole    string
	Task        string
	SolvedTask  string
	TaskHistory string
	TaskAuthor  string
	Team        string
	TeamMember  string
	Settings     string
}{
	User:        "\"user\"",
	Role:        "role",
	UserRole:    "user_role",
	Task:        "task",
	SolvedTask:  "solved_task",
	TaskHistory: "task_submission",
	TaskAuthor:  "author",
	Team:        "team",
	TeamMember:  "team_member",
	Settings:     "settings",
}
