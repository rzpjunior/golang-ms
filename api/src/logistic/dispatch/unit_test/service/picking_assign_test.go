package service

import (
	"git.edenfarm.id/project-version2/api/src/logistic/dispatch/unit_test/entity"
	"git.edenfarm.id/project-version2/api/src/logistic/dispatch/unit_test/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

var dispatchRepository = &repository.DispatchRepositoryMock{Mock: mock.Mock{}}
var dispatchService = DispatchService{Repository: dispatchRepository}

func TestVendorCourierCombination(t *testing.T) {
	dispatch := entity.Dispacth{
		Id:          "1",
		VendorName:  "Eden Farm",
		CourierName: "Doni",
	}

	dispatchRepository.Mock.On("CheckCourierVendorCombination", "1").Return(dispatch)

	result, err := dispatchService.CheckVendorCourierCombinationById("1")
	assert.Nil(t, err)
	assert.Equal(t, dispatch.Id, result.Id)
	assert.Equal(t, dispatch.VendorName, result.VendorName)
	assert.Equal(t, dispatch.CourierName, result.CourierName)
}

func TestListDispatch(t *testing.T) {
	dispatch := entity.Dispacth{
		Id:             "2",
		SalesOrderCode: "SO-BASG9867-12210019",
		PickingStatus:  2,
		DispatchStatus: 1,
	}

	dispatchRepository.Mock.On("CheckListDispatch", "2").Return(dispatch)

	result, err := dispatchService.CheckListDispatchById("2")
	assert.Nil(t, err)
	assert.Equal(t, dispatch.Id, result.Id)
	assert.Equal(t, dispatch.PickingStatus, result.PickingStatus)
	assert.Equal(t, dispatch.DispatchStatus, result.DispatchStatus)
}
