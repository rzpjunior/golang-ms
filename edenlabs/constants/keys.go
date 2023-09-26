package constants

type contextKey string

const (
	KeyToken     contextKey = "token"
	KeyUserID    contextKey = "user_id"
	KeyCourierID contextKey = "courier_id"
	KeySiteID    contextKey = "site_id"
	KeyPlatform  contextKey = "platform"
	KeyTimezone  contextKey = "timezone"
)
