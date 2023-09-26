package engine

import "git.edenfarm.id/project-version2/api/src/migration"

func init() {
	handlers["migration"] = &migration.Handler{}
}
