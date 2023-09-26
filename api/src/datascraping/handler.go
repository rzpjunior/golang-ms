// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package datascraping

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("/public_price_set", h.getPublicPriceSet, auth.Authorized("prd_mtc_upl"))
	r.GET("/dashboard", h.dashboard)
	r.GET("/public_price_1", h.getPublicPrice1)
	r.GET("/public_price_2", h.getPublicPrice2)
	r.GET("/public_product", h.getPublicProduct)
	r.GET("/product_matching", h.getProductMatchingTemplate, auth.Authorized("prd_mtc_upl"))
	// r.GET("/public_price_3", h.public_price_3)
	r.POST("/product_matching", h.updateProductMatching, auth.Authorized("prd_mtc_upl"))
}

// getPublicPriceSet : handler to get list of public price set
func (h *Handler) getPublicPriceSet(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.PublicPriceSet
	var total int64

	if data, total, e = repository.GetPublicPriceSets(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// dashboard : handler to get data from dashboard
func (h *Handler) dashboard(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var matchedArea []*model.MatchedArea

	if matchedArea, _, e = repository.GetMatchedAreas(rq); e == nil {
		e = SaveDashboardProduct(matchedArea)
	}

	return ctx.Serve(e)
}

// getPublicPrice1 : handler to get data from public data 1
func (h *Handler) getPublicPrice1(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var baseURL string
	var request *http.Request
	var client = &http.Client{}
	var response *http.Response
	var resData *model.PublicData1
	var startNum, currentPage, totalPages int64
	var areas []*model.MatchedArea

	if areas, _, e = repository.GetMatchedAreas(rq, "1"); e == nil {
		for _, v := range areas {
			baseURL = "https://tanihub.com/call"
			for currentPage, totalPages = 1, 1; currentPage <= totalPages; {
				m, b := map[string]interface{}{
					"url":    "v2/product-search",
					"method": "post",
					"body": map[string]interface{}{
						"query": "\nquery sellings(\n  $from: Int,\n  $size: Int,\n  $regionId: Int,\n  $regionSlug: String,\n  $groupId: Int,\n  $nowDate: String,\n  $sortField: String,\n  $sortOrder: String,\n  $searchKey: String,\n  $randomizerSeed: Float,\n  $relatedName: String,\n  $relatedId: Int,\n  $isPromoted: Boolean\n) {\n  sellings(\n    from: $from,\n    size: $size,\n    regionId: $regionId,\n    regionSlug: $regionSlug,\n    groupId: $groupId,\n    nowDate: $nowDate,\n    sortField: $sortField,\n    sortOrder: $sortOrder,\n    searchKey: $searchKey,\n    randomizerSeed: $randomizerSeed,\n    relatedName: $relatedName,\n    relatedId: $relatedId,\n    isPromoted: $isPromoted\n  ) {\n    count\n    totalPages\n    currentPage\n    params {\n      from\n      size\n      sortField\n      sortOrder\n      regionId\n      searchKey\n    }\n    items {\n      id\n      product {\n        id\n        name\n        slug\n        commercialSkuContent\n        brand {\n          name\n        }\n        productImages {\n          isDefault\n          imageURL\n        }\n        unit {\n          description\n        }\n        productPackaging {\n          id\n          name\n          multiplier\n        }\n        groups {\n          id\n          name\n          imageUrl\n        }\n        grade {\n          id\n        }\n      }\n      productPrices {\n        id\n        discount\n        discountType\n        minQty\n        maxQty\n        price\n      }\n      discount\n      discountInRupiah\n      discountType\n      minOrder\n      maxOrder\n      isActive\n      showedPrice\n      showedPriceAfterDiscount\n      stockLowerLimit\n      stockQty\n    }\n  }\n}\n",
						"variables": map[string]interface{}{
							"from":      startNum,
							"groupId":   0,
							"regionId":  v.PublicDataArea1.ID,
							"size":      100,
							"sortField": "_score",
							"sortOrder": "desc",
						},
					},
				}, new(bytes.Buffer)
				json.NewEncoder(b).Encode(m)

				if request, e = http.NewRequest("POST", baseURL, b); e == nil {
					request.Header.Set("Content-Type", "application/json")

					if response, e = client.Do(request); e == nil && response.StatusCode == 200 {
						e = json.NewDecoder(response.Body).Decode(&resData)

						e = SavePublicProduct1(resData.Data.Sellings.Items, v.PublicDataArea1)
					}

					defer response.Body.Close()
				}

				startNum = 100 * currentPage
				currentPage++
				totalPages = resData.Data.Sellings.TotalPages
			}
		}
	}

	return ctx.Serve(e)
}

// getPublicPrice2 : handler to get data from public data 2
func (h *Handler) getPublicPrice2(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var baseURL string
	var request *http.Request
	var client = &http.Client{}
	var response *http.Response
	var token *model.TokenData
	var resData *model.PublicData2
	var areas []*model.MatchedArea

	baseURL = "https://www.googleapis.com/identitytoolkit/v3/relyingparty/signupNewUser?key=AIzaSyCB5tP70uz9FZa-lZOqPWOUfe73nBYqa0s"
	m, b := map[string]bool{"returnSecureToken": true}, new(bytes.Buffer)
	json.NewEncoder(b).Encode(m)

	if request, e = http.NewRequest("POST", baseURL, b); e == nil {
		request.Header.Set("Content-Type", "application/json")

		response, e = client.Do(request)
		if e == nil && response.StatusCode == 200 {
			e = json.NewDecoder(response.Body).Decode(&token)

			if areas, _, e = repository.GetMatchedAreas(rq, "2"); e == nil {
				tomorrowDate := time.Now().Add(time.Hour*7).AddDate(0, 0, 1).Format("Monday, 02 January 2006")
				for _, v := range areas {
					page := 1
					for hasNextPage := true; hasNextPage == true; {
						baseURL = "https://api.sayurbox.io/graphql"

						variables := map[string]interface{}{
							"deliveryArea": v.PublicDataArea2.Name,
							"deliveryDate": tomorrowDate,
							"deliveryCode": v.PublicDataArea2.Code,
							"limit":        100,
							"page":         page,
							"type":         "CATEGORY",
							"value":        "SemuaProduk",
						}

						test := map[string]interface{}{
							"operationName": "getCatalogVariant",
							"variables":     variables,
							"query":         "query getCatalogVariant($deliveryDate: String!, $deliveryArea: String!, $deliveryCode: String, $limit: Int!, $page: Int!, $type: CatalogType, $value: String) {\n  catalogVariantList(deliveryDate: $deliveryDate, deliveryArea: $deliveryArea, deliveryCode: $deliveryCode, limit: $limit, page: $page, type: $type, value: $value) {\n    limit\n    page\n    size\n    hasNextPage\n    category {\n      displayName\n    }\n    list {\n      key\n      availability\n      categories\n      farmers {\n        image\n        name\n      }\n      image {\n        md\n        sm\n        lg\n      }\n      isDiscount\n      discount\n      labelDesc\n      labelName\n      maxQty\n      name\n      displayName\n      nextAvailableDates\n      packDesc\n      packNote\n      price\n      priceFormatted\n      actualPrice\n      actualPriceFormatted\n      shortDesc\n      stockAvailable\n      type\n      emptyMessageHtml\n      promoMessageHtml\n    }\n  }\n}\n",
						}

						m, b := test, new(bytes.Buffer)
						json.NewEncoder(b).Encode(m)
						if request, e = http.NewRequest("POST", baseURL, b); e == nil {
							request.Header.Set("Authorization", token.IDToken)
							request.Header.Set("Content-Type", "application/json")
							request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36")
							request.Header.Set("version", "1.41.0")

							if response, e = client.Do(request); e == nil && response.StatusCode == 200 {
								if e = json.NewDecoder(response.Body).Decode(&resData); e == nil {
									e = SavePublicProduct2(resData.Data.CatalogVariantList.List, v.PublicDataArea2)
								}
							}
						}

						defer response.Body.Close()

						hasNextPage = resData.Data.CatalogVariantList.HasNextPage
						page++
					}
				}
			}
		}
	}

	return ctx.Serve(e)
}

// getPublicProduct : handler to download product of public data
func (h *Handler) getPublicProduct(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var total int64
	var file string

	isExport := ctx.QueryParam("export") == "1"
	idPriceSet, _ := common.Decrypt(ctx.QueryParam("public_price_set"))
	if isExport {
		var data []*model.PublicProductForXls
		if idPriceSet == 1 {
			data, e = repository.GetPublicProduct1ForXls()
		} else if idPriceSet == 2 {
			data, e = repository.GetPublicProduct2ForXls()
		}

		if e == nil {
			if file, e = DownloadPublicProductXls(time.Now(), data, idPriceSet); e == nil {
				ctx.Files(file)
			}
		}
	} else {
		var data []*model.PublicProduct1
		if data, total, e = repository.GetPublicProduct1s(rq); e == nil {
			ctx.Data(data, total)
		}
	}

	return ctx.Serve(e)
}

// getProductMatchingTemplate : handler to download product matching template
func (h *Handler) getProductMatchingTemplate(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var data []*model.ProductMatchingTemplate
	var file string

	if data, e = repository.GetMatchingProduct(); e == nil {
		if file, e = DownloadProductMatchingXls(time.Now(), data); e == nil {
			ctx.Files(file)
		}
	}

	return ctx.Serve(e)
}

// updateProductMatching : handler to update product matching
func (h *Handler) updateProductMatching(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var r matchingRequest
	var successCount int64

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if e = ctx.Bind(&r); e == nil {
			successCount, e = SaveUpdateMatching(r)

			totalCount := len(r.Data)
			ctx.ResponseData = strconv.Itoa(int(successCount)) + " of " + strconv.Itoa(totalCount) + " data has been saved successfully"
		}
	}

	return ctx.Serve(e)
}

// func (h *Handler) getPublicPrice3(c echo.Context) (e error) {
// 	ctx := c.(*cuxs.Context)

// 	var baseURL string
// 	var request *http.Request
// 	var client = &http.Client{}
// 	var response *http.Response
// 	var resData *model.DataPasarnow
// 	var resultData []string

// 	baseURL = "https://api.pasarnow.com/api/appProductListNotLoggedIn?page=0"

// 	if request, e = http.NewRequest("GET", baseURL, nil); e == nil {
// 		response, e = client.Do(request)
// 		if e == nil && response.StatusCode == 200 {
// 			if e = json.NewDecoder(response.Body).Decode(&resData); e == nil {
// 				for i := 1; i < resData.Pages; i++ {
// 					baseURL = "https://api.pasarnow.com/api/appProductListNotLoggedIn?page=" + strconv.Itoa(i)
// 					if request, e := http.NewRequest("GET", baseURL, nil); e == nil {
// 						request.Header.Set("Location", "Jabodetabek")

// 						if response, e = client.Do(request); e == nil && response.StatusCode == 200 {
// 							if e = json.NewDecoder(response.Body).Decode(&resData); e == nil {
// 								for _, v := range resData.Products {
// 									resultData = append(resultData, v.ProductName)
// 								}
// 								ctx.ResponseData = resultData
// 							}
// 						}
// 						defer response.Body.Close()
// 					}
// 				}
// 			}
// 		}
// 		defer response.Body.Close()
// 	}

// 	return ctx.Serve(e)
// }
