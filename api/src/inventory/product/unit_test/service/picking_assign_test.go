package service

import (
	"git.edenfarm.id/project-version2/api/src/inventory/product/unit_test/entity"
	"git.edenfarm.id/project-version2/api/src/inventory/product/unit_test/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

var productRepository = &repository.ProductRepositoryMock{Mock: mock.Mock{}}
var productService = ProductService{Repository: productRepository}

func TestVendorCourierCombination(t *testing.T) {
	product := entity.Product{
		Id:          "1",
		TotalPrint:  10,
		ProductCode: "PRD000649",
		ProductName: "Gula Rose Brand Premium Hijau (1 kg/pack)",
	}

	productRepository.Mock.On("CheckPrintProductLabel", "1").Return(product)

	result, err := productService.CheckPrintProductLabelById("1")
	assert.Nil(t, err)
	assert.Equal(t, product.Id, result.Id)
	assert.Equal(t, product.ProductName, result.ProductName)
	assert.Equal(t, product.ProductCode, result.ProductCode)
}
