package repository

import (
	"git.edenfarm.id/project-version2/api/src/inventory/product/unit_test/entity"
	"github.com/stretchr/testify/mock"
)

type ProductRepositoryMock struct {
	Mock mock.Mock
}

func (repository *ProductRepositoryMock) CheckPrintProductLabel(id string) *entity.Product {
	arguments := repository.Mock.Called(id)

	if arguments.Get(0) == nil {
		return nil
	} else {
		product := arguments.Get(0).(entity.Product)
		return &product
	}
}
