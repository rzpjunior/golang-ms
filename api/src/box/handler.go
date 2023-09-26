package box

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.read)
	r.GET("/box_fridge/:id", h.boxFridgeDetail)
	r.GET("/product_box", h.readProductBox)
	r.GET("/export_template", h.exportTemplate, auth.Authorized("filter_rdl"))
	r.POST("/product_box/upload", h.uploadTemplate, auth.Authorized("filter_rdl"))
	r.POST("/rfidLabel", h.postPrint, auth.Authorized("pco_prt"))
	r.PUT("/product_box_finish", h.updateBoxFridgeFinish, auth.Authorized("filter_rdl"))

}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.Box
	var total int64

	if data, total, e = repository.GetBoxes(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// read : function to get requested data based on parameters
func (h *Handler) readProductBox(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	// var data []*model.BoxItem
	// var total int64

	// if data, total, e = repository.GetBoxItems(rq); e == nil {
	// 	ctx.Data(data, total)
	// }
	statusFilter := ctx.QueryParam("status")
	var data []*model.ProductFridgeBoxListQuery
	var total int64

	if data, total, e = repository.GetBoxFridgeItems(rq, statusFilter); e == nil {
		ctx.Data(data, total)
	}
	return ctx.Serve(e)
}

// readFilter : function to get requested data based on parameters with filtered permission
func (h *Handler) readFilter(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.Box
	var total int64

	if data, total, e = repository.GetBoxes(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

type ProductFridgeBoxQuery struct {
	BoxId       string `orm:"column(box_id);null" json:"box_id,omitempty"`
	Name        string
	TotalPrice  float64 `orm:"column(total_price);null"`
	UnitPrice   float64 `orm:"column(unit_price);null"`
	TotalWeight float64
	ImageUrl    string
	Uom         string
	LastSeenAt  time.Time `orm:"column(last_seen_at);null"`
	Rfid        string    `orm:"column(rfid);null"`
	BoxFridgeId string    `orm:"column(box_fridge_id);null"`
}

type ProductFridgeBoxListQuery struct {
	ProductName     string    `orm:"column(product_name);null" json:"product_name,omitempty"`
	TotalWeight     string    `orm:"column(total_weight);null" json:"total_weight,omitempty"`
	ItemImage       string    `orm:"column(item_image);null" json:"item_image,omitempty"`
	Uom             string    `orm:"column(uom_name);null" json:"uom_name,omitempty"`
	ProcessedAt     time.Time `orm:"column(processed_at);null"  json:"processed_at,omitempty"`
	Rfid            string    `orm:"column(rfid);null"  json:"rfid,omitempty"`
	WarehouseId     int64     `orm:"column(warehouse_id);null"  json:"-"`
	WasteImage      string    `orm:"column(waste_image);null" json:"waste_image,omitempty"`
	FinishedAt      time.Time `orm:"column(finished_at);null"  json:"finished_at,omitempty"`
	BoxFridgeStatus int64     `orm:"column(box_fridge_status);null"  json:"box_fridge_status,omitempty"`
	BoxItemStatus   int64     `orm:"column(box_item_status);null"  json:"box_item_status,omitempty"`
	Status          string    `orm:"column(status);null"  json:"status,omitempty"`
	WarehouseName   string    `orm:"column(warehouse_name);null" json:"warehouse_name,omitempty"`
	BranchName      string    `orm:"column(branch_name);null" json:"branch_name,omitempty"`
}

type BoxFridgeQuery struct {
	ProductList  []*ProductFridgeBoxQuery
	ProductCart  []*ProductFridgeBoxQuery
	ProductSold  []*ProductFridgeBoxQuery
	ProductWaste []*ProductFridgeBoxQuery
}

// detail : function to get detailed data by id
func (h *Handler) productBoxFridge(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	var data []*ProductFridgeBoxQuery
	o := orm.NewOrm()
	o.Using("read_only")
	// if _, e := o.Raw("select p.name ,pb.total_price ,pb.total_weight,pi.image_url,u.name as uom,bf.last_seen_at,rfid,bf.id as box_fridge_id  from "+
	// 	"box_fridge bf join branch_fridge bf2 "+
	// 	"on bf.warehouse_id =bf2.warehouse_id "+
	// 	"join box_item pb "+
	// 	"on bf.box_id =pb.box_id "+
	// 	"join product p "+
	// 	"on pb.product_id =p.id "+
	// 	"join product_image pi on pi.product_id=p.id "+
	// 	"join uom u on p.uom_id=u.id "+
	// 	"join box b on pb.box_id=b.id "+
	// 	"where bf.status=1 and bf2.warehouse_id =? and bf.last_seen_at >= NOW()-INTERVAL ? SECOND", id, 10).QueryRows(&data); e != nil {
	// 	return ctx.Serve(e)
	// }
	if _, e := o.Raw("select  "+
		"p.id,p.name ,b.rfid ,bi.total_weight ,u.name ,pi2.image_url,bf.last_seen_at,bf.warehouse_id,bf.image_url as waste_image,bi.finished_at  "+
		"from box_item bi  "+
		"left join box_fridge bf on bi.id =bf.box_item_id  and bi.status =1 "+
		"join box b on b.id =bi.box_id  "+
		"join product p on p.id =bi.product_id "+
		"join product_image pi2 on pi2.product_id =p.id "+
		"join uom u on p.uom_id=u.id "+
		"where bi.status =1 ", id, 10).QueryRows(&data); e != nil {
		return ctx.Serve(e)
	}
	for _, dat := range data {
		dat.BoxFridgeId = common.Encrypt(dat.BoxFridgeId)
	}
	ctx.ResponseData = data

	return ctx.Serve(e)

}

// detail : function to get detailed data by id
func (h *Handler) productBoxFridgeAddList(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}
	var data []*ProductFridgeBoxQuery
	o := orm.NewOrm()
	o.Using("read_only")
	if _, e := o.Raw("select p.name ,pb.total_price ,pb.total_weight,pi.image_url,u.name as uom,bf.last_seen_at,rfid,bf.id as box_fridge_id  from "+
		"box_fridge bf join branch_fridge bf2 "+
		"on bf.warehouse_id =bf2.warehouse_id "+
		"join box_item pb "+
		"on bf.box_id =pb.box_id "+
		"join product p "+
		"on pb.product_id =p.id "+
		"join product_image pi on pi.product_id=p.id "+
		"join uom u on p.uom_id=u.id "+
		"join box b on pb.box_id=b.id "+
		"where bf.status=0 and bf2.warehouse_id =? 	and bf.last_seen_at >= NOW()-INTERVAL ? SECOND", id, 10).QueryRows(&data); e != nil {
		//"where bf.status=0 and bf2.warehouse_id =? 	and bf.last_seen_at >= NOW()-INTERVAL ? SECOND", id, 10).QueryRows(&data); e != nil {
		return ctx.Serve(e)
	}
	for _, dat := range data {
		dat.BoxFridgeId = common.Encrypt(dat.BoxFridgeId)
	}
	ctx.ResponseData = data

	return ctx.Serve(e)

}

// detail : function to get detailed data by id
func (h *Handler) productBoxFridgeOutList(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}
	var data []*ProductFridgeBoxQuery
	o := orm.NewOrm()
	o.Using("read_only")
	if _, e := o.Raw("select p.name ,pb.total_price ,pb.total_weight,pi.image_url,u.name as uom,bf.last_seen_at,rfid,bf.id as box_fridge_id  from "+
		"box_fridge bf join branch_fridge bf2 "+
		"on bf.warehouse_id =bf2.warehouse_id "+
		"join box_item pb "+
		"on bf.box_id =pb.box_id "+
		"join product p "+
		"on pb.product_id =p.id "+
		"join product_image pi on pi.product_id=p.id "+
		"join uom u on p.uom_id=u.id "+
		"join box b on pb.box_id=b.id "+
		"where bf.status=1 and bf2.warehouse_id =? 	and bf.last_seen_at <= NOW()-INTERVAL ? SECOND", id, 10).QueryRows(&data); e != nil {
		// "where bf.status=1 and bf2.warehouse_id =? 	and bf.last_seen_at <= NOW()-INTERVAL ? SECOND", id, 10).QueryRows(&data); e != nil {
		return ctx.Serve(e)
	}
	for _, dat := range data {
		dat.BoxFridgeId = common.Encrypt(dat.BoxFridgeId)
	}
	ctx.ResponseData = data

	return ctx.Serve(e)

}

// detail : function to get detailed data by id
func (h *Handler) productBoxFridgeSoldList(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}
	var data []*ProductFridgeBoxQuery
	o := orm.NewOrm()
	o.Using("read_only")
	if _, e := o.Raw("select p.name ,pb.total_price ,pb.total_weight,pi.image_url,u.name as uom,bf.last_seen_at,rfid,bf.id as box_fridge_id  from "+
		"box_fridge bf join branch_fridge bf2 "+
		"on bf.warehouse_id =bf2.warehouse_id "+
		"join box_item pb "+
		"on bf.box_id =pb.box_id "+
		"join product p "+
		"on pb.product_id =p.id "+
		"join product_image pi on pi.product_id=p.id "+
		"join uom u on p.uom_id=u.id "+
		"join box b on pb.box_id=b.id "+
		"where bf.status=2 and bf2.warehouse_id =? ", id).QueryRows(&data); e != nil {
		// "where bf.status=1 and bf2.warehouse_id =? 	and bf.last_seen_at <= NOW()-INTERVAL ? SECOND", id, 10).QueryRows(&data); e != nil {
		return ctx.Serve(e)
	}
	for _, dat := range data {
		dat.BoxFridgeId = common.Encrypt(dat.BoxFridgeId)
	}
	ctx.ResponseData = data

	return ctx.Serve(e)

}

// detail : function to get detailed data by id
func (h *Handler) productBoxDetail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	id := ctx.Param("id")
	var data *ProductFridgeBoxQuery
	o := orm.NewOrm()
	o.Using("read_only")
	e = o.Raw("select pb.box_id,p.name ,pb.total_price ,pb.total_weight,pi.image_url,u.name as uom from "+
		"product_box pb "+
		"join product p "+
		"on pb.product_id =p.id "+
		"join product_image pi on pi.product_id=p.id "+
		"join box b on pb.box_id=b.id "+
		"join uom u on p.uom_id=u.id "+
		"where pb.status=1 and b.rfid =?", id).QueryRow(&data)
	if e == nil {
		data.BoxId = common.Encrypt(data.BoxId)
		ctx.ResponseData = data
	}

	return ctx.Serve(e)
}

// detail : function to get detailed data by id
func (h *Handler) boxFridge(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	//rq := ctx.RequestQuery()

	// var data []*model.BoxFridge
	// var total int64

	// if data, total, e = repository.GetBoxFridges(rq); e == nil {
	// 	ctx.Data(data, total)
	// }
	//ctx := c.(*cuxs.Context)

	// var id int64
	// if id, e = ctx.Decrypt("id"); e != nil {
	// 	return ctx.Serve(e)
	// }

	var total int64

	var data []*ProductFridgeBoxListQuery
	o := orm.NewOrm()
	o.Using("read_only")
	if _, e := o.Raw("select  " +
		"p.name as product_name ,b.rfid ,bi.total_weight ,u.name as uom_name, " +
		"pi2.image_url as item_image,bf.last_seen_at as processed_at,bf.warehouse_id , " +
		"bf.image_url as waste_image,bi.finished_at,bf.status as box_fridge_status,bi.status as box_item_status    " +
		"from box_item bi  " +
		"left join box_fridge bf on bi.id =bf.box_item_id  and bi.status =1 " +
		"join box b on b.id =bi.box_id  " +
		"join product p on p.id =bi.product_id " +
		"join product_image pi2 on pi2.product_id =p.id " +
		"join uom u on p.uom_id=u.id ").QueryRows(&data); e != nil {
		return ctx.Serve(e)
	}
	if e := o.Raw("select count(*) " +
		"from box_item bi  " +
		"left join box_fridge bf on bi.id =bf.box_item_id  and bi.status =1 " +
		"join box b on b.id =bi.box_id  " +
		"join product p on p.id =bi.product_id " +
		"join product_image pi2 on pi2.product_id =p.id " +
		"join uom u on p.uom_id=u.id ").QueryRow(&total); e != nil {
		return ctx.Serve(e)
	}
	for _, dat := range data {
		if dat.WarehouseId != 0 {
			Warehouse := &model.Warehouse{ID: dat.WarehouseId}
			if e := Warehouse.Read("ID"); e != nil {
				//o.Failure("warehouse.id", e.Error())
				fmt.Println(e)
			}

			dat.WarehouseName = Warehouse.Name

			branchFridge := &model.BranchFridge{Warehouse: Warehouse}
			if e := branchFridge.Read("Warehouse"); e != nil {
				fmt.Println(e)
			}
			if branchFridge.ID != 0 {
				if e := branchFridge.Branch.Read("ID"); e != nil {
					fmt.Println(e)
				}
				dat.BranchName = branchFridge.Branch.Name
			}
			fmt.Println(dat)
			if dat.BoxFridgeStatus == 1 {
				dat.Status = "Active"
			}
		}
		if dat.BoxItemStatus == 3 {
			dat.Status = "finished"
		} else {
			if dat.BoxFridgeStatus == 1 {
				dat.Status = "active"
			} else if dat.BoxFridgeStatus == 2 {
				dat.Status = "sold"
			} else if dat.BoxFridgeStatus == 4 {
				dat.Status = "waste"
			}
			if dat.WarehouseId == 0 {
				dat.Status = "new"
			}
		}

	}
	//ctx.ResponseData = data
	ctx.Data(data, total)
	return ctx.Serve(e)

}

// detail : function to get detailed data by id
func (h *Handler) boxFridgeDetail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	rfid := ctx.Param("id")

	var data *ProductFridgeBoxListQuery
	o := orm.NewOrm()
	o.Using("read_only")
	if e := o.Raw("select  "+
		"p.name as product_name ,b.rfid ,bi.total_weight ,u.name as uom_name, "+
		"pi2.image_url as item_image,bf.last_seen_at as processed_at,bf.warehouse_id , "+
		"bf.image_url as waste_image,bi.finished_at,bf.status as box_fridge_status,bi.status as box_item_status    "+
		"from box_item bi  "+
		"left join box_fridge bf on bi.id =bf.box_item_id  and bi.status =1 "+
		"join box b on b.id =bi.box_id  "+
		"join product p on p.id =bi.product_id "+
		"left join product_image pi2 on pi2.product_id =p.id "+
		"join uom u on p.uom_id=u.id "+
		"where b.rfid = ? and bi.status=1", rfid).QueryRow(&data); e != nil {
		return ctx.Serve(e)
	}
	if data.BoxItemStatus == 3 {
		e = errors.New("data already finished")
		return ctx.Serve(e)
	}
	if data.WarehouseId != 0 {
		Warehouse := &model.Warehouse{ID: data.WarehouseId}
		if e := Warehouse.Read("ID"); e != nil {
			//o.Failure("warehouse.id", e.Error())
			fmt.Println(e)
		}

		data.WarehouseName = Warehouse.Name

		branchFridge := &model.BranchFridge{Warehouse: Warehouse}
		if e := branchFridge.Read("Warehouse"); e != nil {
			fmt.Println(e)
		}
		if branchFridge.ID != 0 {
			if e := branchFridge.Branch.Read("ID"); e != nil {
				fmt.Println(e)
			}
			data.BranchName = branchFridge.Branch.Name
		}
		if data.BoxFridgeStatus == 1 {
			data.Status = "active"
		}
	}
	if data.BoxItemStatus == 3 {
		data.Status = "finished"
	} else {
		if data.BoxFridgeStatus == 1 {
			data.Status = "active"
		} else if data.BoxFridgeStatus == 2 {
			data.Status = "sold"
		} else if data.BoxFridgeStatus == 4 {
			data.Status = "waste"
		}
		if data.WarehouseId == 0 {
			data.Status = "new"
		}
	}

	ctx.ResponseData = data
	//ctx.Data(data, total)
	return ctx.Serve(e)

}

// detail : function to get detailed data by id
func (h *Handler) resetStatusBox(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}
	o := orm.NewOrm()
	o.Using("read_only")
	_, e = o.Raw("update box_item set status=1 ").Exec()
	_, e = o.Raw("update box_fridge	set status=?,last_seen_at=?	where warehouse_id=? ", 1, time.Now(), id).Exec()

	return ctx.Serve(e)
}

// detail : function to get detailed data by id
func (h *Handler) exportTemplate(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var backdate time.Time
	backdate = time.Now()

	var file string
	if file, e = exportTemplateXls(backdate); e == nil {
		ctx.Files(file)
	}
	return ctx.Serve(e)
}

// detail : function to get detailed data by id
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	id := ctx.Param("id")
	var data *ProductFridgeBoxQuery
	o := orm.NewOrm()
	o.Using("read_only")
	e = o.Raw("select pb.box_id,p.name ,pb.total_price ,pb.total_weight,pi.image_url,u.name as uom from "+
		"product_box pb "+
		"join product p "+
		"on pb.product_id =p.id "+
		"join product_image pi on pi.product_id=p.id "+
		"join box b on pb.box_id=b.id "+
		"join uom u on p.uom_id=u.id "+
		"where pb.status=1 and b.code =?", id).QueryRow(&data)
	if e == nil {
		data.BoxId = common.Encrypt(data.BoxId)
		ctx.ResponseData = data
	}

	return ctx.Serve(e)
}

//create : function to create new data based on input
func (h *Handler) uploadTemplate(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r createRequestTemplate
	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}
	if e = ctx.Bind(&r); e == nil {
		ctx.ResponseData, e = SaveTemplate(r)
	} else {
		//post error
		errLog := util.ErrorLog{
			ErrorCode:    422,
			Name:         r.Session.Staff.Name,
			Email:        r.Session.Staff.User.Email,
			ErrorMessage: e.Error(),
			Function:     "create_box_fridge",
		}
		util.PostToServiceErrorLog(errLog)
	}
	return ctx.Serve(e)
}

//create : function to update new data based on input
func (h *Handler) updateBoxFridgeFinish(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updateFinishRequest
	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}
	r.MacAddress = c.Request().Header.Get("X-MAC")

	if e = ctx.Bind(&r); e == nil {
		ctx.ResponseData, e = UpdateFinish(r)
	} else {
		//post error
		errLog := util.ErrorLog{
			ErrorCode:    422,
			Name:         r.Session.Staff.Name,
			Email:        r.Session.Staff.User.Email,
			ErrorMessage: e.Error(),
			Function:     "update_box_fridge",
		}
		util.PostToServiceErrorLog(errLog)
	}
	return ctx.Serve(e)
}

//funcUploadImageFile : function to upload image to S3 storage
func (h *Handler) funcUploadImageFile(r echo.Context) (e error) {
	ctx := r.(*cuxs.Context)
	if err := r.Request().ParseMultipartForm(1024); err != nil {
		return err
	}

	handler, err := r.FormFile("file")
	typeRequest := r.FormValue("type")
	fileType := handler.Header.Get("Content-Type")

	if fileType != "image/jpeg" && fileType != "image/png" {
		err = errors.New("The provided file format is not allowed. Please upload a JPEG or PNG image")
		return err
	}

	uploadedFile, _ := handler.Open()
	if err != nil {
		return err
	}
	defer uploadedFile.Close()

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	fileName := handler.Filename
	fileExtension := strings.ToLower(filepath.Ext(fileName))
	validFileExtensions := [...]string{".jpg", ".jpeg", ".png"}
	isFileExtensionValid := util.ItemExists(validFileExtensions, fileExtension)

	if !isFileExtensionValid {
		err = errors.New("The provided file format is not allowed. Please upload a JPEG or PNG image")
		return err
	}

	fileLocation := filepath.Join(dir, "", fileName)
	targetFile, err := os.OpenFile(fileLocation, os.O_RDWR|os.O_CREATE, 0775)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, uploadedFile); err != nil {
		return err
	}

	fl, err := util.UploadImageToS3(fileName, fileLocation, typeRequest)
	if err != nil {
		return err
	}
	ctx.Data(fl)

	os.Remove(fileLocation)

	return ctx.Serve(err)
}

func (h *Handler) postPrint(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r printRequest
	var resp ResponsePrint
	req := make(map[string]interface{})

	// if _, e = auth.UserSession(ctx); e == nil {
	// 	if e = ctx.Bind(&r); e == nil {
	// 		req["qr"] = r

	// 		file := util.SendPrint(req, "read/qrcodeLabel")
	// 		resp.LinkPrint = file
	// 		ctx.ResponseData = resp

	// 	}
	// }
	if e = ctx.Bind(&r); e == nil {
		req["qr"] = r

		file := util.SendPrint(req, "read/qrcodeLabel")
		resp.LinkPrint = file
		ctx.ResponseData = resp

	}

	return ctx.Serve(e)
}

type ResponsePrint struct {
	LinkPrint  string  `json:"link_print"`
	TotalPrint float64 `json:"total_print"`
}
