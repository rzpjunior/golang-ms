package repository

import "git.edenfarm.id/project-version2/api/src/inventory/product/unit_test/entity"

type ProductRepository interface {
	CheckPrintProductLabel(id string) *entity.Product
}
