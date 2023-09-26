// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"encoding/json"
	"strconv"
	"strings"

	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetAdmDivisions : function to get data from database based on parameters
func GetAdmDivisions(rq *orm.RequestQuery, gb string) (m []*model.AdmDivision, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.AdmDivision))

	if gb != "" {
		sliceGb := strings.Split(gb, ",")
		q = q.GroupBy(sliceGb...)
	}

	if total, err = q.Filter("area_status", 1).Filter("country_status", 1).Filter("province_status", 1).Filter("city_status", 1).Filter("district_status", 1).Filter("sub_district_status", 1).All(&m, rq.Fields...); err == nil {
		return m, total, nil
	}

	return nil, total, err
}

// GetFilterAdmDivisions : function to get data from database based on parameters with filtered permission
func GetFilterAdmDivisions(rq *orm.RequestQuery, gb, polygon string) (m []*model.AdmDivision, total int64, err error) {
	o := orm.NewOrm()
	o.Using("read_only")
	q, _ := rq.QueryReadOnly(new(model.AdmDivision))

	if gb != "" {
		sliceGb := strings.Split(gb, ",")
		q = q.GroupBy(sliceGb...)
	}

	if total, err = q.Filter("area_status", 1).Filter("country_status", 1).Filter("province_status", 1).Filter("city_status", 1).Filter("district_status", 1).Filter("sub_district_status", 1).All(&m, rq.Fields...); err != nil {
		return nil, total, err
	}

	if polygon == "1" {
		for _, v := range m {
			//Centroid
			var centroidString string
			var centroidArr []string
			var latitude float64
			var longitude float64

			o.Raw("SELECT ST_AsText(ST_Centroid(`polygon`)) as centroid from `adm_division_geometry` `adg`  where `sub_district_id` = ?", v.SubDistrictId).QueryRow(&centroidString)
			centroidString = strings.Trim(centroidString, "POINT()")

			centroidArr = strings.Split(centroidString, " ")
			if len(centroidArr) == 2 {
				latitude, _ = strconv.ParseFloat(centroidArr[1], 64)
				longitude, _ = strconv.ParseFloat(centroidArr[0], 64)
			}

			v.Centroid = model.Centroid{
				Latitude:  latitude,
				Longitude: longitude,
			}

			// Polygon
			var polygonString string
			var polygon model.PolygonJson
			o.Raw("SELECT ST_AsGeoJSON(`polygon`) as geojson from `adm_division_geometry` `adg`  where `sub_district_id` = ?", v.SubDistrictId).QueryRow(&polygonString)
			json.Unmarshal([]byte(polygonString), &polygon)

			if len(polygon.Coordinates) > 0 {
				if len(polygon.Coordinates[0]) > 0 {
					for _, v1 := range polygon.Coordinates[0] {
						var polygonPart []float64
						// append latitude first then longitude
						polygonPart = append(polygonPart, v1[1])
						polygonPart = append(polygonPart, v1[0])

						v.Polygon = append(v.Polygon, polygonPart)
					}
				}
			}
		}
	}

	return m, total, nil
}
