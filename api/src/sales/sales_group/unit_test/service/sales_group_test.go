package service

import (
	"git.edenfarm.id/project-version2/api/src/sales/sales_group/unit_test/entity"
	"git.edenfarm.id/project-version2/api/src/sales/sales_group/unit_test/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

var salesGroupRepository = &repository.SalesGroupRepositoryMock{Mock: mock.Mock{}}
var salesGroupService = SalesGroupService{Repository: salesGroupRepository}

func TestSalesGroupStatus(t *testing.T) {
	salesGroup := entity.SalesGroup{
		Id:               "1",
		SalesGroupName:   "Eden JKT 1",
		SalesGroupStatus: 2,
	}

	salesGroupRepository.Mock.On("CheckSalesGroup", "1").Return(salesGroup)

	result, err := salesGroupService.CheckSalesGroupById("1")
	assert.Nil(t, err)
	assert.Equal(t, salesGroup.Id, result.Id)
	assert.Equal(t, salesGroup.SalesGroupName, result.SalesGroupName)
	assert.Equal(t, salesGroup.SalesGroupStatus, result.SalesGroupStatus)
}

func TestSalesGroupDuplicateName(t *testing.T) {
	salesGroup := entity.SalesGroup{
		Id:               "1",
		SalesGroupName:   "jakarta2",
		SalesGroupStatus: 2,
	}

	salesGroupRepository.Mock.On("CheckDuplicateNameSalesGroup", "jakarta2").Return(salesGroup)

	result, err := salesGroupService.CheckDuplicateNameSalesGroupByName("jakarta2")
	assert.Nil(t, err)
	assert.Equal(t, salesGroup.Id, result.Id)
	assert.Equal(t, salesGroup.SalesGroupName, result.SalesGroupName)
	assert.Equal(t, salesGroup.SalesGroupStatus, result.SalesGroupStatus)
}
