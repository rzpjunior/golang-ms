package stock

import (
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/cuxs/cuxs/event"
)

func init() {
	listenItemDeliveryOutInsert()
	listenStockopnameCommited()
	listenWasteEntryCommited()
}

func listenItemDeliveryOutInsert() {
	c := make(chan interface{})

	event.Listen("delivery::delivery", c)
	go func() {
		for {
			data := <-c
			doi := data.(*model.DeliveryOrderItem)

			makeItemLogDeliveryOutInsert(doi)
		}
	}()
}

func listenStockopnameCommited() {
	c := make(chan interface{})
	event.Listen("stockopname::commited", c)

	go func() {
		for {
			data := <-c
			soi := data.(*model.StockOpnameItem)

			makeLogStockOpnameCommitted(soi)
		}
	}()
}
//
func listenWasteEntryCommited() {
	c := make(chan interface{})
	event.Listen("wasteentry::commited", c)

	go func() {
		for {
			data := <-c
			wei := data.(*model.WasteEntryItem)

			makeLogWasteEntryCommitted(wei)
		}
	}()
}
