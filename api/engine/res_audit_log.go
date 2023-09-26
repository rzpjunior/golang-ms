package engine

import "git.edenfarm.id/project-version2/api/src/audit_log"

func init() {
	handlers["audit_log"] = &audit_log.Handler{}
}
