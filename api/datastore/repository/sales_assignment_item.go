// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// GetSalesAssignmentItems : function to get data from database based on parameters
func GetSalesAssignmentItems(rq *orm.RequestQuery) (m []*model.SalesAssignmentItem, total int64, err error) {
	o := orm.NewOrm()
	o.Using("read_only")

	q, _ := rq.QueryReadOnly(new(model.SalesAssignmentItem))

	if total, err = q.Exclude("status", 4).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.SalesAssignmentItem
	if _, err = q.Exclude("status", 4).All(&mx, rq.Fields...); err == nil {
		for _, v := range mx {
			o.Raw("SELECT g.value_name FROM glossary g WHERE g.table = 'sales_assignment_item' AND g.attribute = 'task' AND g.value_int = ?", v.Task).QueryRow(&v.TaskStr)
			if v.Status == 14 {
				o.Raw("SELECT * FROM sales_failed_visit sfv WHERE sfv.sales_assignment_item_id = ?", v.ID).QueryRow(&v.SalesFailedVisit)
				if v.SalesFailedVisit != nil {
					v.SalesFailedVisit.FailedImageList = strings.Split(v.SalesFailedVisit.FailedImage, ",")
				}
			}
		}
		return mx, total, nil
	}

	return nil, total, err
}

// ValidSalesAssignmentItem : function to check if id is valid in database
func ValidSalesAssignmentItem(id int64) (SalesAssignmentItem *model.SalesAssignmentItem, e error) {
	SalesAssignmentItem = &model.SalesAssignmentItem{ID: id}
	e = SalesAssignmentItem.Read("ID")

	return
}

// GetSubmissionSA : function to get data from database based on parameters
func GetSubmissionSA(rq *orm.RequestQuery, task string) (m []*model.SalesAssignmentItem, total int64, err error) {
	o := orm.NewOrm()
	o.Using("read_only")

	q, _ := rq.QueryReadOnly(new(model.SalesAssignmentItem))
	q = q.Exclude("status", 4).RelatedSel("Branch", "CustomerAcquisition").Filter("task", task)

	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.SalesAssignmentItem

	branchFinishedVisit := make(map[time.Time]int64)
	branchFinishedFU := make(map[time.Time]int64)
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		for _, v := range mx {
			o.Raw("SELECT g.value_name FROM glossary g WHERE g.table = 'sales_assignment_item' AND g.attribute = 'task' AND g.value_int = ?", v.Task).QueryRow(&v.TaskStr)
			var salesGroupID int

			if v.OutofRoute == 1 {
				if _, err = o.LoadRelated(v, "SalesPerson", 0); err != nil {
					return nil, 0, err
				}
				if v.SalesPerson.SalesGroupID != 0 {
					salesGroupID = int(v.SalesPerson.SalesGroupID)
					var sa *model.SalesGroup
					if err = o.Raw("SELECT * FROM sales_group WHERE id = ?", v.SalesPerson.SalesGroupID).QueryRow(&sa); err != nil {
						return nil, 0, err
					}
					v.SalesAssignment = &model.SalesAssignment{
						SalesGroup: sa,
					}
				}
			}

			if v.Status == 2 {
				// get effective call
				if v.OutofRoute != 1 {
					if err = o.Raw("SELECT sales_group_id FROM sales_assignment sa WHERE sa.id = ?", v.SalesAssignment.ID).QueryRow(&salesGroupID); err != nil {
						return nil, 0, err
					}
				}
				var so []*model.SalesOrder
				if v.CustomerType == 1 && salesGroupID != 0 {
					if _, err = o.Raw("SELECT * FROM sales_order so WHERE so.branch_id = ? AND so.salesperson_id = ? AND so.sales_group_id = ? AND so.recognition_date = ? AND so.status NOT IN (3,4)",
						v.Branch.ID, v.SalesPerson.ID, salesGroupID, v.FinishDate.Format("2006-01-02")).QueryRows(&so); err != nil && !errors.Is(err, orm.ErrNoRows) {
						return nil, 0, err
					}
				}

				if v.Task == 1 && v.CustomerType == 1 {
					branchFinishedVisit[v.StartDate] = v.Branch.ID
				} else if v.Task == 2 && v.CustomerType == 1 {
					branchFinishedFU[v.StartDate] = v.Branch.ID
				}

				if so != nil {
					v.EffectiveCall = true
					if branchFinishedVisit[v.StartDate] != 0 {
						if branchFinishedVisit[v.StartDate] != v.Branch.ID {
							for _, rec := range so {
								v.RevenueEffectiveCall += rec.TotalCharge
							}
						}
					}
					if branchFinishedFU[v.StartDate] != 0 {
						if branchFinishedFU[v.StartDate] != v.Branch.ID {
							for _, rec := range so {
								v.RevenueEffectiveCall += rec.TotalCharge
							}
						}
					}
				}
			}
		}
		return mx, total, nil
	}

	return nil, total, err
}

// GetSubmissionSA : function to get data from database based on parameters
func GetSubmissionVisitAndFollowUp(rq *orm.RequestQuery) (m []*model.SalesAssignmentItem, total int64, err error) {
	o := orm.NewOrm()
	o.Using("read_only")

	q, _ := rq.QueryReadOnly(new(model.SalesAssignmentItem))

	cond := q.GetCond()

	cond1 := orm.NewCondition()
	cond1 = cond1.And("task", 1)

	cond2 := orm.NewCondition()
	cond2 = cond2.And("task", 2)

	cond3 := cond1.OrCond(cond2)
	cond = cond.AndCond(cond3)

	q = q.SetCond(cond).Exclude("status", 4)

	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.SalesAssignmentItem

	if _, err = q.RelatedSel("Branch", "CustomerAcquisition").All(&mx, rq.Fields...); err == nil {
		for _, v := range mx {
			o.Raw("SELECT g.value_name FROM glossary g WHERE g.table = 'sales_assignment_item' AND g.attribute = 'task' AND g.value_int = ?", v.Task).QueryRow(&v.TaskStr)

			var salesGroupID int
			if v.OutofRoute == 1 {
				if _, err = o.LoadRelated(v, "SalesPerson", 0); err != nil {
					return nil, 0, err
				}
				if v.SalesPerson.SalesGroupID != 0 {
					salesGroupID = int(v.SalesPerson.SalesGroupID)
					var sa *model.SalesGroup
					if err = o.Raw("SELECT * FROM sales_group WHERE id = ?", v.SalesPerson.SalesGroupID).QueryRow(&sa); err != nil {
						return nil, 0, err
					}
					v.SalesAssignment = &model.SalesAssignment{
						SalesGroup: sa,
					}
				}
			}
			if v.Status == 2 {
				// get effective call
				if v.OutofRoute != 1 {
					if err = o.Raw("SELECT sales_group_id FROM sales_assignment sa WHERE sa.id = ?", v.SalesAssignment.ID).QueryRow(&salesGroupID); err != nil {
						return nil, 0, err
					}
				}
				var so []*model.SalesOrder
				if v.CustomerType == 1 {

					if _, err = o.Raw("SELECT * FROM sales_order so WHERE so.branch_id = ? AND so.salesperson_id = ? AND so.sales_group_id = ? AND so.recognition_date = ? AND so.status NOT IN (3,4)",
						v.Branch.ID, v.SalesPerson.ID, salesGroupID, v.FinishDate.Format("2006-01-02")).QueryRows(&so); err != nil && !errors.Is(err, orm.ErrNoRows) {
						return nil, 0, err
					}
				}

				if so != nil {
					v.EffectiveCall = true
					for _, rec := range so {
						v.RevenueEffectiveCall += rec.TotalCharge
					}
				}
			}
		}
		return mx, total, nil
	}

	return nil, total, err
}

// GetSubmissionSalesAssignmentDetail : function to get submitted sales assignment item by id
func GetSubmissionSalesAssignmentDetail(field string, values ...interface{}) (*model.SalesAssignmentItem, error) {
	m := new(model.SalesAssignmentItem)
	o := orm.NewOrm()
	o.Using("read_only")
	if err := o.QueryTable(m).RelatedSel("Branch", "CustomerAcquisition", "Salesperson").Filter(field, values...).Limit(1).One(m); err != nil {
		return nil, err
	}

	//// Initialize minio client object.
	minioClient, _ := minio.New(util.S3endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(util.S3accessKeyID, util.S3secretAccessKey, ""),
		Secure: true,
	})
	reqParams := make(url.Values)
	//// Retrieve URL valid for 60 second

	if m.Status == 14 {
		o.Raw("SELECT * FROM sales_failed_visit WHERE sales_assignment_item_id = ?", m.ID).QueryRow(&m.SalesFailedVisit)
		if m.SalesFailedVisit != nil {
			if m.SalesFailedVisit.FailedImage != "" {
				failedImageArr := strings.Split(m.SalesFailedVisit.FailedImage, ",")
				for _, tp := range failedImageArr {
					tempImage := strings.Split(tp, "/")
					preSignedURLImage, _ := minioClient.PresignedGetObject(context.Background(), util.S3bucketNameImage, tempImage[4]+"/"+tempImage[5], time.Second*60, reqParams)
					m.SalesFailedVisit.FailedImageList = append(m.SalesFailedVisit.FailedImageList, preSignedURLImage.String())
				}
			}
		}
	} else {
		if m.TaskPhoto != "" {
			m.TaskPhotoArr = strings.Split(m.TaskPhoto, ",")
			for _, tp := range m.TaskPhotoArr {
				tempImage := strings.Split(tp, "/")
				preSignedURLImage, _ := minioClient.PresignedGetObject(context.Background(), util.S3bucketNameImage, tempImage[4]+"/"+tempImage[5], time.Second*60, reqParams)
				m.TaskPhotoList = append(m.TaskPhotoList, preSignedURLImage.String())
			}
		}
	}

	o.Raw("SELECT g.value_name FROM glossary g WHERE g.table = 'sales_assignment_item' AND g.attribute = 'task' AND g.value_int = ?", m.Task).QueryRow(&m.TaskStr)
	o.Raw("SELECT g.value_name FROM glossary g WHERE g.table = 'sales_assignment_item' AND g.attribute = 'task_answer' AND g.value_int = ?", m.AnswerOption).QueryRow(&m.AsnsweOptionStr)
	o.Raw("SELECT g.value_name FROM glossary g WHERE g.table = 'sales_assignment_item' AND g.attribute = 'customer_type' AND g.value_int = ?", m.CustomerType).QueryRow(&m.CustomerTypeStr)

	return m, nil
}

// GetSalesAssignmentItemsGroup : function to get data from database based on parameters
func GetSalesAssignmentItemsGroup(rq *orm.RequestQuery, fromDate, toDate time.Time, salesGroupID, salesPersonID int64) (m []*model.SalesAssignmentItem, total int64, err error) {
	o := orm.NewOrm()
	o.Using("read_only")
	var dateLabel string
	var staffIds []int

	q, _ := rq.QueryReadOnly(new(model.SalesAssignmentItem))

	if total, err = q.Exclude("status", 4).GroupBy("salesperson_id").Count(); err != nil {
		return nil, total, err
	}

	var mx []*model.SalesAssignmentItem
	if _, err = q.Exclude("status", 4).GroupBy("salesperson_id").RelatedSel().All(&mx, rq.Fields...); err == nil {
		for _, v := range mx {
			staffIds = append(staffIds, int(v.SalesPerson.ID))
			o.Raw("SELECT g.value_name FROM glossary g WHERE g.table = 'sales_assignment_item' AND g.attribute = 'task' AND g.value_int = ?", v.Task).QueryRow(&v.TaskStr)

			var mx1 []*model.SalesAssignmentItem
			var allSO []*model.SalesOrder
			if _, err = o.Raw("SELECT * FROM sales_assignment_item sai WHERE sai.salesperson_id = ? AND sai.start_date BETWEEN ? AND ? AND sai.status != 4",
				v.SalesPerson.ID, fromDate, toDate).QueryRows(&mx1); err != nil {
				return nil, 0, err
			}
			if _, err = o.Raw("SELECT * FROM sales_order so WHERE so.salesperson_id = ? AND so.sales_group_id = ? AND so.recognition_date BETWEEN ? AND ? AND so.status NOT IN (3,4)",
				v.SalesPerson.ID, v.SalesPerson.SalesGroupID, fromDate, toDate).QueryRows(&allSO); err != nil {
				return nil, 0, err
			}

			var totalEC, totalTaskFinished, revenueEffectiveCall, revenueTotal float64
			var totalCA int
			sumPlanVisit := make(map[time.Time]int64)
			sumPlanFU := make(map[time.Time]int64)
			sumPlanVisitFU := make(map[time.Time]int64)
			sumFinishedVisit := make(map[time.Time]int64)
			branchFinishedVisit := make(map[string]int64)
			sumFinishedFU := make(map[time.Time]int64)
			branchFinishedFU := make(map[string]int64)
			sumFinishedVisitFU := make(map[time.Time]int64)
			sumPerformanceVisit := make(map[time.Time]float64)
			sumPerformanceFU := make(map[time.Time]float64)
			sumEffectiveCall := make(map[time.Time]int64)
			sumEffectiveCallPercentage := make(map[time.Time]float64)
			layout := "02/01/06"

			for _, data := range mx1 {
				if data.Status == 1 {
					if data.Task == 1 {
						sumPlanVisit[data.StartDate] += 1
					} else if data.Task == 2 {
						sumPlanFU[data.StartDate] += 1
					}
					sumPlanVisitFU[data.StartDate] += 1
				}
				if data.Status == 2 {
					if data.Task == 1 {
						sumPlanVisit[data.StartDate] += 1
						sumFinishedVisit[data.StartDate] += 1
						if v.CustomerType == 1 {
							branchFinishedVisit[data.StartDate.Format(layout)] = data.Branch.ID
						} else if v.CustomerType == 2 {
							branchFinishedVisit[data.StartDate.Format(layout)] = data.CustomerAcquisition.ID
						}
					} else if data.Task == 2 {
						sumPlanFU[data.StartDate] += 1
						sumFinishedFU[data.StartDate] += 1
						if v.CustomerType == 1 {
							branchFinishedFU[data.StartDate.Format(layout)] = data.Branch.ID
						} else if v.CustomerType == 2 {
							branchFinishedFU[data.StartDate.Format(layout)] = data.CustomerAcquisition.ID
						}
					}
					sumFinishedVisitFU[data.StartDate] += 1

					// get effective call
					if data.CustomerType == 1 {
						var salesGroupID int
						if data.OutofRoute != 1 {
							if err = o.Raw("SELECT sales_group_id FROM sales_assignment sa WHERE sa.id = ?", data.SalesAssignment.ID).QueryRow(&salesGroupID); err != nil {
								return nil, 0, err
							}
						} else {
							if err = o.Raw("SELECT sales_group_id FROM staff WHERE id = ?", data.SalesPerson.ID).QueryRow(&salesGroupID); err != nil {
								return nil, 0, err
							}
						}
						var so []*model.SalesOrder
						if _, err = o.Raw("SELECT * FROM sales_order so WHERE so.branch_id = ? AND so.salesperson_id = ? AND so.sales_group_id = ? AND so.recognition_date = ? AND so.status NOT IN (3,4)",
							data.Branch.ID, data.SalesPerson.ID, salesGroupID, data.FinishDate.Format("2006-01-02")).QueryRows(&so); err != nil && !errors.Is(err, orm.ErrNoRows) {
							return nil, 0, err
						}
						if so != nil {
							sumEffectiveCall[data.StartDate] += 1
							totalEC += 1
							for _, rec := range so {
								revenueEffectiveCall += rec.TotalCharge
							}
						}
					}
					totalTaskFinished += 1
				}
				if data.Status == 14 {
					if data.Task == 1 {
						sumPlanVisit[data.StartDate] += 1
					} else if data.Task == 2 {
						sumPlanFU[data.StartDate] += 1
					}
					sumPlanVisitFU[data.StartDate] += 1
				}
			}

			if err = o.Raw("SELECT count(id) as total FROM customer_acquisition ca WHERE ca.sales_group_id IS NOT NULL AND ca.salesperson_id = ? AND ca.sales_group_id = ? AND DATE(ca.submit_date) BETWEEN ? AND ? AND ca.status != 3",
				v.SalesPerson.ID, v.SalesPerson.SalesGroupID, fromDate, toDate).QueryRow(&totalCA); err != nil {
				return nil, 0, err
			}

			for k, v1 := range sumPlanVisit {
				// visit performance calculation
				performance := (float64(sumFinishedVisit[k]) / float64(v1)) * 100
				if performance > 100 {
					performance = 100
				}
				sumPerformanceVisit[k] = performance
			}
			for k, v1 := range sumPlanFU {
				// FU performance calculation
				performance := (float64(sumFinishedFU[k]) / float64(v1)) * 100
				if performance > 100 {
					performance = 100
				}
				sumPerformanceFU[k] = performance
			}

			for k, v1 := range sumFinishedVisitFU {
				// call effective percentage calculation
				performanceEffectiveCall := (float64(sumEffectiveCall[k]) / float64(v1)) * 100
				if performanceEffectiveCall > 100 {
					performanceEffectiveCall = 100
				}
				sumEffectiveCallPercentage[k] = performanceEffectiveCall
			}

			dateMap := make(map[string]int64)
			for i := fromDate; !i.After(toDate); i = i.AddDate(0, 0, 1) {
				dateLabel = i.Format(layout)
				dateMap[dateLabel] = 0
			}

			sumPlanVisitFUStr := make(map[string]int64)
			for key, value := range sumPlanVisitFU {
				dateLabel = key.Format(layout)
				sumPlanVisitFUStr[dateLabel] = value
			}
			sumPlanVisitStr := make(map[string]int64)
			for key, value := range sumPlanVisit {
				dateLabel = key.Format(layout)
				sumPlanVisitStr[dateLabel] = value
			}
			sumPlanFUStr := make(map[string]int64)
			for key, value := range sumPlanFU {
				dateLabel = key.Format(layout)
				sumPlanFUStr[dateLabel] = value
			}
			sumFinishedVisitStr := make(map[string]int64)
			for key, value := range sumFinishedVisit {
				dateLabel = key.Format(layout)
				sumFinishedVisitStr[dateLabel] = value
			}
			sumFinishedFUStr := make(map[string]int64)
			for key, value := range sumFinishedFU {
				dateLabel = key.Format(layout)
				sumFinishedFUStr[dateLabel] = value
			}
			sumFinishedVisitFUStr := make(map[string]int64)
			for key, value := range sumFinishedVisitFU {
				dateLabel = key.Format(layout)
				sumFinishedVisitFUStr[dateLabel] = value
			}

			sumPerformanceVisitStr := make(map[string]float64)
			for key, value := range sumPerformanceVisit {
				dateLabel = key.Format(layout)
				sumPerformanceVisitStr[dateLabel] = value
			}
			sumPerformanceFUStr := make(map[string]float64)
			for key, value := range sumPerformanceFU {
				dateLabel = key.Format(layout)
				sumPerformanceFUStr[dateLabel] = value
			}
			sumEffectiveCallPercStr := make(map[string]float64)
			for key, value := range sumEffectiveCallPercentage {
				dateLabel = key.Format(layout)
				sumEffectiveCallPercStr[dateLabel] = value
			}

			for key := range dateMap {
				if value, ok := sumPlanVisitFUStr[key]; ok {
					dateMap[key] = value
				}
			}

			keys := make([]string, 0, len(dateMap))
			for k := range dateMap {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			for _, k := range keys {
				v.PlanVisit += sumPlanVisitStr[k]
				v.PlanFollowUp += sumPlanFUStr[k]
				v.VisitActual += sumFinishedVisitStr[k]
				v.FollowUpActual += sumFinishedFUStr[k]
			}
			for _, r := range allSO {
				revenueTotal += r.TotalCharge
			}

			if len(sumPlanVisitStr) > 0 {
				v.VisitPercentage = (float64(v.VisitActual) / float64(v.PlanVisit)) * 100
			}
			if len(sumPlanFUStr) > 0 {
				v.FollowUpPercentage = (float64(v.FollowUpActual) / float64(v.PlanFollowUp)) * 100
			}
			if len(sumEffectiveCallPercStr) > 0 {
				v.EffectiveCallPercentage = (totalEC / totalTaskFinished) * 100
			}
			v.RevenueEffectiveCall = revenueEffectiveCall
			v.RevenueTotal = revenueTotal
			v.TotalCA = totalCA

		}
	}

	if len(mx) > 0 {
		var allSO []*model.SalesOrder
		query := "SELECT salesperson_id, sales_group_id, total_charge FROM sales_order so WHERE so.recognition_date BETWEEN ? AND ? AND so.status NOT IN (3,4) AND (so.sales_group_id != 0 AND so.sales_group_id IS NOT NULL) "
		if salesGroupID != 0 {
			query += fmt.Sprintf("AND sales_group_id = %d ", salesGroupID)
		}
		if salesPersonID != 0 {
			query += fmt.Sprintf("AND salesperson_id = %d ", salesPersonID)
		}
		qMark := ""
		for _, _ = range staffIds {
			qMark = qMark + "?,"
		}
		qMark = strings.TrimSuffix(qMark, ",")
		query += "AND salesperson_id NOT IN (" + qMark + ") GROUP BY salesperson_id"
		if _, err = o.Raw(query, fromDate, toDate, staffIds).QueryRows(&allSO); err != nil {
			return nil, 0, err
		}
		for _, v := range allSO {
			totalRev := float64(0)
			var so []*model.SalesOrder
			var staff model.Staff
			var totalCA int
			staffIds = append(staffIds, int(v.Salesperson.ID))
			if _, err = o.Raw("SELECT total_charge FROM sales_order so WHERE so.salesperson_id = ? AND so.sales_group_id = ? AND so.recognition_date BETWEEN ? AND ? AND so.status NOT IN (3,4)",
				v.Salesperson.ID, v.SalesGroupID, fromDate.Format("2006-01-02"), toDate.Format("2006-01-02")).QueryRows(&so); err != nil && !errors.Is(err, orm.ErrNoRows) {
				return nil, 0, err
			}
			if err = o.Raw("SELECT * FROM staff WHERE id = ? AND sales_group_id = ? AND status = 1",
				v.Salesperson.ID, v.SalesGroupID).QueryRow(&staff); err != nil && !errors.Is(err, orm.ErrNoRows) {
				return nil, 0, err
			}
			if err = o.Raw("SELECT name FROM sales_group WHERE id = ? AND status = 1",
				v.SalesGroupID).QueryRow(&staff.SalesGroupName); err != nil && !errors.Is(err, orm.ErrNoRows) {
				return nil, 0, err
			}

			for _, order := range so {
				totalRev += order.TotalCharge
			}

			if err = o.Raw("SELECT count(id) as total FROM customer_acquisition ca WHERE ca.sales_group_id IS NOT NULL AND ca.salesperson_id = ? AND ca.sales_group_id = ? AND DATE(ca.submit_date) BETWEEN ? AND ? AND ca.status != 3",
				v.Salesperson.ID, v.SalesGroupID, fromDate, toDate).QueryRow(&totalCA); err != nil {
				return nil, 0, err
			}

			if staff.Name != "" {
				sai := &model.SalesAssignmentItem{
					SalesPerson:  &staff,
					RevenueTotal: totalRev,
					TotalCA:      totalCA,
				}
				mx = append(mx, sai)
			}
		}
		if len(allSO) > 0 {
			var allCA []*model.CustomerAcquisition
			query := "SELECT id, salesperson_id, sales_group_id FROM customer_acquisition ca WHERE ca.sales_group_id IS NOT NULL AND DATE(ca.submit_date) BETWEEN ? AND ? "
			if salesGroupID != 0 {
				query += fmt.Sprintf("AND ca.sales_group_id = %d ", salesGroupID)
			}
			if salesPersonID != 0 {
				query += fmt.Sprintf("AND ca.salesperson_id = %d ", salesPersonID)
			}
			qMark := ""
			for _, _ = range staffIds {
				qMark = qMark + "?,"
			}
			qMark = strings.TrimSuffix(qMark, ",")
			query += "AND ca.status != 3 AND ca.salesperson_id NOT IN (" + qMark + ") GROUP BY salesperson_id"
			if _, err = o.Raw(query, fromDate, toDate, staffIds).QueryRows(&allCA); err != nil {
				return nil, 0, err
			}

			for _, v := range allCA {
				var totalCA int
				var staff model.Staff
				if err = o.Raw("SELECT count(id) as total FROM customer_acquisition ca WHERE ca.sales_group_id IS NOT NULL AND ca.salesperson_id = ? AND ca.sales_group_id = ? AND DATE(ca.submit_date) BETWEEN ? AND ? AND ca.status != 3",
					v.Salesperson.ID, v.Salesgroup.ID, fromDate, toDate).QueryRow(&totalCA); err != nil {
					return nil, 0, err
				}
				if err = o.Raw("SELECT * FROM staff WHERE id = ? AND sales_group_id = ? AND status = 1",
					v.Salesperson.ID, v.Salesgroup.ID).QueryRow(&staff); err != nil && !errors.Is(err, orm.ErrNoRows) {
					return nil, 0, err
				}
				if err = o.Raw("SELECT name FROM sales_group WHERE id = ? AND status = 1",
					v.Salesgroup.ID).QueryRow(&staff.SalesGroupName); err != nil && !errors.Is(err, orm.ErrNoRows) {
					return nil, 0, err
				}

				if staff.Name != "" {
					sai := &model.SalesAssignmentItem{
						SalesPerson: &staff,
						TotalCA:     totalCA,
					}

					mx = append(mx, sai)
				}
			}
		} else {
			var allCA []*model.CustomerAcquisition
			query := "SELECT id, salesperson_id, sales_group_id FROM customer_acquisition ca WHERE ca.sales_group_id IS NOT NULL AND DATE(ca.submit_date) BETWEEN ? AND ? "
			if salesGroupID != 0 {
				query += fmt.Sprintf("AND ca.sales_group_id = %d ", salesGroupID)
			}
			if salesPersonID != 0 {
				query += fmt.Sprintf("AND ca.salesperson_id = %d ", salesPersonID)
			}
			qMark := ""
			for _, _ = range staffIds {
				qMark = qMark + "?,"
			}
			qMark = strings.TrimSuffix(qMark, ",")
			query += "AND ca.status != 3 AND ca.salesperson_id NOT IN (" + qMark + ") GROUP BY salesperson_id"
			if _, err = o.Raw(query, fromDate, toDate, staffIds).QueryRows(&allCA); err != nil {
				return nil, 0, err
			}

			for _, v := range allCA {
				var totalCA int
				var staff model.Staff
				if err = o.Raw("SELECT count(id) as total FROM customer_acquisition ca WHERE ca.sales_group_id IS NOT NULL AND ca.salesperson_id = ? AND ca.sales_group_id = ? AND DATE(ca.submit_date) BETWEEN ? AND ? AND ca.status != 3",
					v.Salesperson.ID, v.Salesgroup.ID, fromDate, toDate).QueryRow(&totalCA); err != nil {
					return nil, 0, err
				}
				if err = o.Raw("SELECT * FROM staff WHERE id = ? AND sales_group_id = ? AND status = 1",
					v.Salesperson.ID, v.Salesgroup.ID).QueryRow(&staff); err != nil && !errors.Is(err, orm.ErrNoRows) {
					return nil, 0, err
				}
				if err = o.Raw("SELECT name FROM sales_group WHERE id = ? AND status = 1",
					v.Salesgroup.ID).QueryRow(&staff.SalesGroupName); err != nil && !errors.Is(err, orm.ErrNoRows) {
					return nil, 0, err
				}

				if staff.Name != "" {
					sai := &model.SalesAssignmentItem{
						SalesPerson: &staff,
						TotalCA:     totalCA,
					}

					mx = append(mx, sai)
				}
			}
		}
	} else {
		var allSO []*model.SalesOrder
		query := "SELECT * FROM sales_order so WHERE so.recognition_date BETWEEN ? AND ? AND so.status NOT IN (3,4) AND (so.sales_group_id != 0 AND so.sales_group_id IS NOT NULL) "
		if salesGroupID != 0 {
			query += fmt.Sprintf("AND sales_group_id = %d ", salesGroupID)
		}
		if salesPersonID != 0 {
			query += fmt.Sprintf("AND salesperson_id = %d ", salesPersonID)
		}
		query += "GROUP BY salesperson_id"
		if _, err = o.Raw(query, fromDate, toDate).QueryRows(&allSO); err != nil {
			return nil, 0, err
		}
		for _, v := range allSO {
			totalRev := float64(0)
			var so []*model.SalesOrder
			var staff model.Staff
			var totalCA int
			staffIds = append(staffIds, int(v.Salesperson.ID))
			if _, err = o.Raw("SELECT * FROM sales_order so WHERE so.salesperson_id = ? AND so.sales_group_id = ? AND so.recognition_date BETWEEN ? AND ? AND so.status NOT IN (3,4)",
				v.Salesperson.ID, v.SalesGroupID, fromDate.Format("2006-01-02"), toDate.Format("2006-01-02")).QueryRows(&so); err != nil && !errors.Is(err, orm.ErrNoRows) {
				return nil, 0, err
			}
			if err = o.Raw("SELECT * FROM staff WHERE id = ? AND sales_group_id = ? AND status = 1",
				v.Salesperson.ID, v.SalesGroupID).QueryRow(&staff); err != nil && !errors.Is(err, orm.ErrNoRows) {
				return nil, 0, err
			}
			if err = o.Raw("SELECT name FROM sales_group WHERE id = ? AND status = 1",
				v.SalesGroupID).QueryRow(&staff.SalesGroupName); err != nil && !errors.Is(err, orm.ErrNoRows) {
				return nil, 0, err
			}

			for _, order := range so {
				totalRev += order.TotalCharge
			}

			if err = o.Raw("SELECT count(id) as total FROM customer_acquisition ca WHERE ca.sales_group_id IS NOT NULL AND ca.salesperson_id = ? AND ca.sales_group_id = ? AND DATE(ca.submit_date) BETWEEN ? AND ? AND ca.status != 3",
				v.Salesperson.ID, v.SalesGroupID, fromDate, toDate).QueryRow(&totalCA); err != nil {
				return nil, 0, err
			}

			if staff.Name != "" {
				sai := &model.SalesAssignmentItem{
					SalesPerson:  &staff,
					RevenueTotal: totalRev,
					TotalCA:      totalCA,
				}

				mx = append(mx, sai)
			}
		}
		if len(allSO) > 0 {
			var allCA []*model.CustomerAcquisition
			query := "SELECT id, salesperson_id, sales_group_id FROM customer_acquisition ca WHERE ca.sales_group_id IS NOT NULL AND DATE(ca.submit_date) BETWEEN ? AND ? "
			if salesGroupID != 0 {
				query += fmt.Sprintf("AND ca.sales_group_id = %d ", salesGroupID)
			}
			if salesPersonID != 0 {
				query += fmt.Sprintf("AND ca.salesperson_id = %d ", salesPersonID)
			}
			qMark := ""
			for _, _ = range staffIds {
				qMark = qMark + "?,"
			}
			qMark = strings.TrimSuffix(qMark, ",")
			query += "AND ca.status != 3 AND ca.salesperson_id NOT IN (" + qMark + ") GROUP BY salesperson_id"
			if _, err = o.Raw(query, fromDate, toDate, staffIds).QueryRows(&allCA); err != nil {
				return nil, 0, err
			}

			for _, v := range allCA {
				var totalCA int
				var staff model.Staff
				if err = o.Raw("SELECT count(id) as total FROM customer_acquisition ca WHERE ca.sales_group_id IS NOT NULL AND ca.salesperson_id = ? AND ca.sales_group_id = ? AND DATE(ca.submit_date) BETWEEN ? AND ? AND ca.status != 3",
					v.Salesperson.ID, v.Salesgroup.ID, fromDate, toDate).QueryRow(&totalCA); err != nil {
					return nil, 0, err
				}
				if err = o.Raw("SELECT * FROM staff WHERE id = ? AND sales_group_id = ? AND status = 1",
					v.Salesperson.ID, v.Salesgroup.ID).QueryRow(&staff); err != nil && !errors.Is(err, orm.ErrNoRows) {
					return nil, 0, err
				}
				if err = o.Raw("SELECT name FROM sales_group WHERE id = ? AND status = 1",
					v.Salesgroup.ID).QueryRow(&staff.SalesGroupName); err != nil && !errors.Is(err, orm.ErrNoRows) {
					return nil, 0, err
				}
				if staff.Name != "" {
					sai := &model.SalesAssignmentItem{
						SalesPerson: &staff,
						TotalCA:     totalCA,
					}

					mx = append(mx, sai)
				}
			}
		}
		if len(allSO) == 0 {
			var allCA []*model.CustomerAcquisition
			query := "SELECT id, salesperson_id, sales_group_id FROM customer_acquisition ca WHERE ca.sales_group_id IS NOT NULL AND DATE(ca.submit_date) BETWEEN ? AND ? "
			if salesGroupID != 0 {
				query += fmt.Sprintf("AND ca.sales_group_id = %d ", salesGroupID)
			}
			if salesPersonID != 0 {
				query += fmt.Sprintf("AND ca.salesperson_id = %d ", salesPersonID)
			}
			query += "AND ca.status != 3 GROUP BY salesperson_id"
			if _, err = o.Raw(query, fromDate, toDate).QueryRows(&allCA); err != nil {
				return nil, 0, err
			}

			for _, v := range allCA {
				var totalCA int
				var staff model.Staff
				if err = o.Raw("SELECT count(id) as total FROM customer_acquisition ca WHERE ca.sales_group_id IS NOT NULL AND ca.salesperson_id = ? AND ca.sales_group_id = ? AND DATE(ca.submit_date) BETWEEN ? AND ? AND ca.status != 3",
					v.Salesperson.ID, v.Salesgroup.ID, fromDate, toDate).QueryRow(&totalCA); err != nil {
					return nil, 0, err
				}
				if err = o.Raw("SELECT * FROM staff WHERE id = ? AND sales_group_id = ? AND status = 1",
					v.Salesperson.ID, v.Salesgroup.ID).QueryRow(&staff); err != nil && !errors.Is(err, orm.ErrNoRows) {
					return nil, 0, err
				}
				if err = o.Raw("SELECT name FROM sales_group WHERE id = ? AND status = 1",
					v.Salesgroup.ID).QueryRow(&staff.SalesGroupName); err != nil && !errors.Is(err, orm.ErrNoRows) {
					return nil, 0, err
				}

				if staff.Name != "" {
					sai := &model.SalesAssignmentItem{
						SalesPerson: &staff,
						TotalCA:     totalCA,
					}

					mx = append(mx, sai)
				}
			}
		}
	}

	return mx, total, err
}

// GetSalesAssignmentItemsGroup : function to get data from database based on parameters
func GetSalesAssignmentItemsTracker(salespersonId int64, fromDate, toDate time.Time) (m *model.TrackerPerformance, err error) {
	o := orm.NewOrm()
	o.Using("read_only")

	var mx []*model.SalesAssignmentItem
	var totalEffectiveCall, countFinished int64
	if _, err := o.Raw("SELECT * FROM sales_assignment_item sai WHERE sai.salesperson_id = ? AND sai.start_date BETWEEN ? AND ? AND sai.status != 4",
		salespersonId, fromDate, toDate).QueryRows(&mx); err != nil {
		return nil, err
	}

	visit := model.VisitFUTracker{}
	fu := model.VisitFUTracker{}
	for _, v := range mx {
		if v.Task == 1 {
			if v.Status == 1 {
				visit.TotalPlan += 1
			} else if v.Status == 2 {
				visit.TotalPlan += 1
				visit.TotalFinished += 1
			} else if v.Status == 3 {
				visit.TotalCancelled += 1
			} else if v.Status == 14 {
				visit.TotalPlan += 1
				visit.TotalFailed += 1
			}

			if v.OutofRoute == 1 {
				visit.TotalOutOfRoute += 1
			}
		} else if v.Task == 2 {
			if v.Status == 1 {
				fu.TotalPlan += 1
			} else if v.Status == 2 {
				fu.TotalPlan += 1
				fu.TotalFinished += 1
			} else if v.Status == 3 {
				fu.TotalCancelled += 1
			} else if v.Status == 14 {
				fu.TotalPlan += 1
				fu.TotalFailed += 1
			}

			if v.OutofRoute == 1 {
				fu.TotalOutOfRoute += 1
			}
		}
		if v.Status == 2 {
			// get effective call
			var salesGroupID int
			if v.OutofRoute != 1 {
				if err = o.Raw("SELECT sales_group_id FROM sales_assignment sa WHERE sa.id = ?", v.SalesAssignment.ID).QueryRow(&salesGroupID); err != nil {
					return nil, err
				}
			} else {
				if err = o.Raw("SELECT sales_group_id FROM staff WHERE id = ?", v.SalesPerson.ID).QueryRow(&salesGroupID); err != nil {
					return nil, err
				}
			}
			var so *model.SalesOrder
			if err = o.Raw("SELECT * FROM sales_order so WHERE so.branch_id = ? AND so.salesperson_id = ? AND so.sales_group_id = ? AND so.recognition_date = ? AND so.status NOT IN (3,4)",
				v.Branch.ID, v.SalesPerson.ID, salesGroupID, v.FinishDate.Format("2006-01-02")).QueryRow(&so); err != nil && !errors.Is(err, orm.ErrNoRows) {
				return nil, err
			}
			if so != nil {
				totalEffectiveCall += 1
			}
			countFinished += 1
		}
	}
	effectiveCallPercentage := float64(0)
	if countFinished != 0 {
		effectiveCallPercentage = (float64(totalEffectiveCall) / float64(countFinished)) * 100
	}
	resp := model.TrackerPerformance{
		VisitTracker:            &visit,
		FollowUpTracker:         &fu,
		EffectiveCallPercentage: effectiveCallPercentage,
	}

	return &resp, nil
}
