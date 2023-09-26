package main

import (
	"fmt"
	"sync"

	"github.com/your/package/bridgeService" // Import your package here
	"golang.org/x/net/context"
)

type SalesOrderDetail struct {
	ID     string
	Detail *bridgeService.GetSalesOrderGPListResponse
	Err    error
}

/*
  ini merupakan bentukan map yang perlu dikembalikan setelah proses go routine dilakukan.
  jadi maps ini merupakan representasi dari []*model.PackingOrderPack
*/

func getSalesOrderDetails(ctx context.Context, wg *sync.WaitGroup, idx int, sopnumbe string, client bridgeService.BridgeServiceGrpcClient, resultChan chan<- SalesOrderDetail) {
	defer wg.Done()

	salesOrderDetail, err := client.GetSalesOrderListGPByID(ctx, &bridgeService.GetSalesOrderGPListByIDRequest{
		Id: sopnumbe,
	})
	resultChan <- SalesOrderDetail{
		ID:     sopnumbe,
		Detail: salesOrderDetail,
		Err:    err,
	}
}

func main() {
	// Initialize your gRPC client and context here
	// client := initializeYourGrpcClient()
	// ctx := context.Background()

	// Example salesOrders.Data
	salesOrders := struct {
		Data []struct {
			Sopnumbe string
			// Other fields
		}
	}{ /* ... */ }

	var wg sync.WaitGroup
	resultChan := make(chan SalesOrderDetail, len(salesOrders.Data))

	for idx, salesOrder := range salesOrders.Data {
		wg.Add(1)
		go getSalesOrderDetails(ctx, &wg, idx, salesOrder.Sopnumbe, client, resultChan)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results in a map
	resultsMap := make(map[string]SalesOrderDetail)
	for result := range resultChan {
		resultsMap[result.ID] = result
	}

	// Print or use resultsMap as needed
	for id, result := range resultsMap {
		fmt.Printf("ID: %s, Detail: %v, Error: %v\n", id, result.Detail, result.Err)
	}
}

/*
dont forget to append array of slice into slice
especially
packItemMaps []*model.PackingOrderPack

snippet code below is using to replicate how to append array of slice into slice
*/

func main() {
	// Array of slices
	arrayOfSlices := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}

	// Target slice
	targetSlice := make([]int, 0)

	// Append slices to target slice
	for _, slice := range arrayOfSlices {
		targetSlice = append(targetSlice, slice...)
	}

	// Print target slice
	fmt.Println(targetSlice)
}
