package bin

import (
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// Save : function to save data requested into database
func Save(r createRequest) (bin *model.Bin, e error) {
	//generate codes for document
	o := orm.NewOrm()
	o.Begin()

	r.Code, e = util.GenerateCode(r.Code, "bin")

	bin = &model.Bin{
		Code:        r.Code,
		Name:        r.Name,
		Warehouse:   r.Warehouse,
		ServiceTime: r.ServiceTime,
		Latitude:    &r.Latitude,
		Longitude:   &r.Longitude,
		Note:        r.Note,
		CreatedAt:   time.Now(),
		CreatedBy:   r.Session.Staff,
		Status:      1,
	}
	if r.ContainProduct == true {
		bin.Product = r.Product
	} else {
		bin.Product = &model.Product{}
	}

	_, e = o.Insert(bin)
	if e != nil {
		o.Rollback()
		return nil, e
	}

	e = log.AuditLogByUser(r.Session.Staff, bin.ID, "bin", "create", "")
	if e != nil {
		o.Rollback()
		return nil, e
	}

	if r.ContainProduct {
		stock := &model.Stock{
			Warehouse: r.Warehouse,
			Product:   r.Product,
		}
		if e = o.Read(stock, "Product", "Warehouse"); e != nil {
			o.Rollback()
			return nil, e
		}

		stock.Bin = bin

		if _, e = o.Update(stock, "bin"); e != nil {
			o.Rollback()
			return nil, e
		}

		if r.BinAssociated == 0 {
			e = log.AuditLogByUser(r.Session.Staff, bin.ID, "bin", "insert product", r.Product.Name)
			if e != nil {
				o.Rollback()
				return nil, e
			}

		} else {
			previousBin := &model.Bin{ID: r.BinAssociated}
			if e = previousBin.Read("ID"); e != nil {
				o.Rollback()
				return nil, e
			}
			previousBin.Product = &model.Product{}
			previousBin.UpdatedAt = time.Now()
			previousBin.UpdatedBy = r.Session.Staff
			if _, e = o.Update(previousBin, "product", "updatedat", "updatedby"); e != nil {
				o.Rollback()
				return nil, e
			}

			note := "Moved product " + r.Product.Name + " from " + previousBin.Name + " to " + bin.Name

			e = log.AuditLogByUser(r.Session.Staff, r.BinAssociated, "bin", "moved product", note)
			if e != nil {
				o.Rollback()
				return nil, e
			}

			e = log.AuditLogByUser(r.Session.Staff, bin.ID, "bin", "moved product", note)
			if e != nil {
				o.Rollback()
				return nil, e
			}
		}
	}

	o.Commit()

	return bin, e
}

// Archive : function to archive data requested into database
func Archive(a archiveRequest) (u *model.Bin, e error) {
	o := orm.NewOrm()
	o.Begin()

	u = &model.Bin{
		ID:        a.ID,
		Product:   &model.Product{},
		Status:    2,
		UpdatedAt: time.Now(),
		UpdatedBy: a.Session.Staff,
	}

	if _, err := o.Update(u, "Product", "status", "updatedat", "updatedby"); err != nil {
		o.Rollback()
		return nil, err
	}

	e = log.AuditLogByUser(a.Session.Staff, u.ID, "bin", "archive", "")
	if e != nil {
		o.Rollback()
		return nil, e
	}

	if a.ContainStock == true {
		a.Stock.Bin = nil
		if _, e = o.Update(a.Stock, "bin"); e != nil {
			o.Rollback()
			return nil, e
		}

		product := &model.Product{ID: a.Stock.Product.ID}
		e = product.Read("id")
		if e != nil {
			o.Rollback()
			return nil, e
		}

		e = log.AuditLogByUser(a.Session.Staff, a.ID, "bin", "removed product", product.Name)
		if e != nil {
			o.Rollback()
			return nil, e
		}
	}

	o.Commit()
	return u, nil
}

// Update : function to update data requested into database
func Update(r updateRequest) (bin *model.Bin, e error) {
	o := orm.NewOrm()
	o.Begin()

	bin = &model.Bin{
		ID:          r.ID,
		ServiceTime: r.ServiceTime,
		Latitude:    &r.Latitude,
		Longitude:   &r.Longitude,
		Note:        r.Note,
		UpdatedAt:   time.Now(),
		UpdatedBy:   r.Session.Staff,
	}

	if r.ContainProduct == true {
		bin.Product = r.Product
	} else {
		bin.Product = &model.Product{}
	}

	_, e = o.Update(bin, "product", "servicetime", "latitude", "longitude", "note", "updatedat", "updatedby")
	if e != nil {
		o.Rollback()
		return nil, e
	}

	e = log.AuditLogByUser(r.Session.Staff, bin.ID, "bin", "update", "")
	if e != nil {
		o.Rollback()
		return nil, e
	}

	if r.ContainProduct == true {
		stock := &model.Stock{
			Warehouse: r.Warehouse,
			Product:   r.Product,
		}

		if e = o.Read(stock, "Product", "Warehouse"); e != nil {
			o.Rollback()
			return nil, e
		}

		stock.Bin = bin

		if _, e = o.Update(stock, "bin"); e != nil {
			o.Rollback()
			return nil, e
		}

		if r.BinAssociated == 0 {
			e = log.AuditLogByUser(r.Session.Staff, bin.ID, "bin", "insert product", r.Product.Name)
			if e != nil {
				o.Rollback()
				return nil, e
			}
		} else if bin.ID == r.BinAssociated {
			// if there's bin associated, delete the affected bin
		} else {
			e = bin.Read("ID")
			if e != nil {
				o.Rollback()
				return nil, e
			}

			binAssociated := &model.Bin{ID: r.BinAssociated}
			e = binAssociated.Read("ID")
			if e != nil {
				o.Rollback()
				return nil, e
			}
			binAssociated.Product = &model.Product{}
			binAssociated.UpdatedAt = time.Now()
			binAssociated.UpdatedBy = r.Session.Staff

			if _, e = o.Update(binAssociated, "product", "updatedat", "updatedby"); e != nil {
				o.Rollback()
				return nil, e
			}

			note := "Moved product " + r.Product.Name + " from " + binAssociated.Name + " to " + bin.Name

			e = log.AuditLogByUser(r.Session.Staff, r.BinAssociated, "bin", "moved product", note)
			if e != nil {
				o.Rollback()
				return nil, e
			}

			e = log.AuditLogByUser(r.Session.Staff, stock.Bin.ID, "bin", "moved product", note)
			if e != nil {
				o.Rollback()
				return nil, e
			}
		}
	}
	o.Commit()

	return bin, e
}
