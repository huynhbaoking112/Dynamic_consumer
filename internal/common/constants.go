package common

const (
	// Queue Names
	ActivityLogQueue = "iam_activity_log_queue"

	// Exchange Names
	IAMEventsExchange = "iam_events_topic"

	// Binding Keys
	ActivityLogBindingKey = "#.log"

	// Event Topics
	WorkspaceCreatedLog  = "workspace.created.log"
	WorkspaceUpdatedLog  = "workspace.updated.log"
	WorkspaceDeletedLog  = "workspace.deleted.log"
	UserCreatedLog       = "user.created.log"
	UserUpdatedLog       = "user.updated.log"
	UserDeletedLog       = "user.deleted.log"
	MemberAddedLog       = "member.added.log"
	MemberRemovedLog     = "member.removed.log"
	MemberRoleChangedLog = "member.role_changed.log"

	// Collection Names
	ActivityLogCollection = "activity_logs"
)
