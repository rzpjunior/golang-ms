package engine

import "git.edenfarm.id/project-version2/api/src/jobs"

func init() {
	handlers["jobs"] = &jobs.Handler{}
}
