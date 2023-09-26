package service

import (
	"errors"
	"git.edenfarm.id/project-version2/api/src/inventory/product/unit_test/entity"
	"git.edenfarm.id/project-version2/api/src/inventory/product/unit_test/repository"
)

type ProductService struct {
	Repository repository.ProductRepository
}

func (service ProductService) CheckPrintProductLabelById(id string) (*entity.Product, error) {
	product := service.Repository.CheckPrintProductLabel(id)
	if product.ProductName != "Gula Rose Brand Premium Hijau (1 kg/pack)" {
		return nil, errors.New("Invalid product")
	} else if product.ProductCode != "PRD000649" {
		return nil, errors.New("Invalid product Code")
	} else if product.TotalPrint != 10 {
		return nil, errors.New("Invalid amount of print")
	} else {
		return product, nil
	}
}
