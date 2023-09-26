package document_history_log

import (
	"git.edenfarm.id/cuxs/cuxs/event"
)

func init() {
	listenDocumentHistoryUser()
}

func listenDocumentHistoryUser() {
	c := make(chan interface{})

	event.Listen("document_history_log::user", c)

	go func() {
		for {
			data := <-c
			userDataRaw := data.(*UserDocumentHistoryLog)

			makeUserDocumentHistoryLog(userDataRaw)
		}
	}()
}
