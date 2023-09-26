package jobs

import (
	"encoding/json"
	"fmt"

	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/cuxs/mongodb"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.read)
	r.GET("/:id", h.detail)
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.Jobs
	var total int64

	if data, total, e = repository.GetJobs(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// detail : function to get detailed data by id
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var jobs model.Jobs

	a := ctx.Param("id")
	jobs.ID, e = primitive.ObjectIDFromHex(a)
	if e != nil {
		fmt.Println(e)
	}

	md := mongodb.NewMongo()
	ret, err := md.GetOneDataWithFilter("Jobs", jobs)
	if err != nil {
		fmt.Println(err)
		md.DisconnectMongoClient()
		return e
	}
	json.Unmarshal(ret, &jobs)

	ctx.ResponseData = jobs

	return ctx.Serve(e)
}
