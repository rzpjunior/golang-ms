package engine

import "git.edenfarm.id/project-version2/api/src/thirdparty/sms_viro"

func init() {
	handlers["sms_viro"] = &sms_viro.Handler{}
}
