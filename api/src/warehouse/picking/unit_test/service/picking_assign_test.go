package service

import (
	"git.edenfarm.id/project-version2/api/src/warehouse/picking/unit_test/entity"
	"git.edenfarm.id/project-version2/api/src/warehouse/picking/unit_test/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

var pickingAssignRepository = &repository.PickingAssignRepositoryMock{Mock: mock.Mock{}}
var pickingAssignService = PickingAssignService{Repository: pickingAssignRepository}

func TestPickingAssignApprove(t *testing.T) {
	pickingAssign := entity.PickingAssign{
		Id:     "2",
		Status: 2,
	}

	pickingAssignRepository.Mock.On("CheckStatusById", "2").Return(pickingAssign)

	result, err := pickingAssignService.CheckStatusApproveById("2")
	assert.Nil(t, err)
	assert.Equal(t, pickingAssign.Id, result.Id)
	assert.Equal(t, pickingAssign.Status, result.Status)
}

func TestPickingAssignReject(t *testing.T) {
	pickingAssign := entity.PickingAssign{
		Id:     "3",
		Status: 3,
	}

	pickingAssignRepository.Mock.On("CheckStatusById", "3").Return(pickingAssign)

	result, err := pickingAssignService.CheckStatusRejectById("3")
	assert.Nil(t, err)
	assert.Equal(t, pickingAssign.Id, result.Id)
	assert.Equal(t, pickingAssign.Status, result.Status)
}

func TestPickingPrintLabel(t *testing.T) {
	pickingAssign := entity.PickingAssign{
		Id:          "1",
		StatusPrint: 0,
	}

	pickingAssignRepository.Mock.On("CheckStatusPrintById", "1").Return(pickingAssign)

	result, err := pickingAssignService.CheckStatusPrintLabelById("1")
	assert.Nil(t, err)
	assert.Equal(t, pickingAssign.Id, result.Id)
	assert.Equal(t, pickingAssign.StatusPrint, result.StatusPrint)
}

func TestPickingPrintInvoice(t *testing.T) {
	pickingAssign := entity.PickingAssign{
		Id:          "4",
		StatusPrint: 1,
	}

	pickingAssignRepository.Mock.On("CheckStatusPrintById", "4").Return(pickingAssign)

	result, err := pickingAssignService.CheckStatusPrintInvoiceById("4")
	assert.Nil(t, err)
	assert.Equal(t, pickingAssign.Id, result.Id)
	assert.Equal(t, pickingAssign.StatusPrint, result.StatusPrint)
}

func TestPickingPrintState(t *testing.T) {
	pickingAssign := entity.PickingAssign{
		Id:         "5",
		PrintState: "label_picking",
	}

	pickingAssignRepository.Mock.On("CheckPrintState", "label_picking").Return(pickingAssign)

	result, err := pickingAssignService.CheckPickingPrintState("label_picking")
	assert.Nil(t, err)
	assert.Equal(t, pickingAssign.Id, result.Id)
	assert.Equal(t, pickingAssign.StatusPrint, result.StatusPrint)
}

func TestSJPrintState(t *testing.T) {
	pickingAssign := entity.PickingAssign{
		Id:         "6",
		PrintState: "sj",
	}

	pickingAssignRepository.Mock.On("CheckPrintState", "sj").Return(pickingAssign)

	result, err := pickingAssignService.CheckSJPrintState("sj")
	assert.Nil(t, err)
	assert.Equal(t, pickingAssign.Id, result.Id)
	assert.Equal(t, pickingAssign.StatusPrint, result.StatusPrint)
}

func TestPickingTolerable(t *testing.T) {
	pickingAssign := entity.PickingAssign{
		Id:        "7",
		Tolerable: 0.01,
	}

	pickingAssignRepository.Mock.On("CheckTolerable", "7").Return(pickingAssign)

	result, _ := pickingAssignService.CheckTolerableByID("7")
	assert.Equal(t, pickingAssign.Id, result.Id)
	assert.Equal(t, pickingAssign.StatusPrint, result.StatusPrint)
}

func TestPickingRefreshedSalesOrder(t *testing.T) {
	var items []entity.Item
	var item entity.Item

	item.ProductCode = "PRD0001"
	item.ProductCode = "PRD0002"
	items = append(items, item)

	pickingAssign := entity.PickingAssign{
		Id:               "8",
		LockedSalesOrder: 1,
		ProductItem:      items,
	}

	pickingAssignRepository.Mock.On("CheckRefreshedItemSalesOrder", "8").Return(pickingAssign)

	result, _ := pickingAssignService.CheckRefreshedSalesOrderByID("8")
	assert.Equal(t, pickingAssign.Id, result.Id)
	assert.Equal(t, pickingAssign.LockedSalesOrder, result.LockedSalesOrder)
}
