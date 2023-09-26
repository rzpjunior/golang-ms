// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package talon

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

var (
	baseUrl       string
	request       *http.Request
	client        = &http.Client{}
	response      *http.Response
	errorResponse *model.ErrorResponse
)

// UpdateCustomerProfileTalon : function to insert or update customer profile in talon
func UpdateCustomerProfileTalon(profileCode, tagCustomer, area, businessType, paymentMethod, createdDate string, referrerData ...string) (err error) {
	var (
		customerTags []string
		responseData *model.CustomerProfileReturn
	)

	o := orm.NewOrm()
	o.Using("read_only")
	baseUrl = util.TalonHost + "/v2/customer_profiles/" + profileCode

	// start get customer tags name
	customerTags = []string{}
	if tagCustomer != "" {
		tagCustomerArr := strings.Split(tagCustomer, ",")
		for _, v := range tagCustomerArr {
			tagName := ""
			if err = o.Raw("select name from tag_customer where id = ?", v).QueryRow(&tagName); err != nil {
				continue
			}
			customerTags = append(customerTags, tagName)
		}
	}
	// end get customer tags name

	attributes := map[string]interface{}{
		"area":              area,
		"business_type":     businessType,
		"customer_tags":     customerTags,
		"payment_method":    paymentMethod,
		"registration_date": createdDate,
	}
	if len(referrerData) > 0 {
		attributes["advocate_id"] = referrerData[0]
		attributes["advocate_ref_code"] = referrerData[1]
	}
	m, b := map[string]interface{}{
		"attributes": attributes,
	}, new(bytes.Buffer)
	json.NewEncoder(b).Encode(m)
	if request, err = http.NewRequest("PUT", baseUrl, b); err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "ApiKey-v1 "+util.TalonApiKey)
	if response, err = client.Do(request); err != nil || (err == nil && response.StatusCode != 200) {
		err = errors.New("Invalid customer profile")
		defer response.Body.Close()
		return err
	}

	err = json.NewDecoder(response.Body).Decode(&responseData)

	defer response.Body.Close()

	return
}

// UpdateCustomerSessionTalon : function to insert or update customer session in talon
func UpdateCustomerSessionTalon(status, isDry, integrationCode, profileCode, archetype, priceSet, referralCode string, itemList []*model.SessionItemData, isUsePoint bool, vouDiscAmount float64, orderType string) (responseData *model.CustomerSessionReturn, err error) {
	var cartItems []interface{}

	o := orm.NewOrm()
	o.Using("read_only")
	baseUrl = util.TalonHost + "/v2/customer_sessions/" + integrationCode + "?dry=" + isDry

	if status != "cancelled" {
		cartItems = []interface{}{}
		for _, v := range itemList {
			item := map[string]interface{}{
				"name":       v.ProductName,
				"sku":        v.ProductCode,
				"category":   v.CategoryName,
				"price":      v.UnitPrice,
				"quantity":   v.OrderQty,
				"weight":     v.UnitWeight,
				"attributes": v.Attributes,
			}

			cartItems = append(cartItems, item)
		}
	}

	session := make(map[string]interface{})
	session = map[string]interface{}{
		"profileId": profileCode,
		"state":     status,
		"cartItems": cartItems,
		"attributes": map[string]interface{}{
			"eden_point_earned":  0,
			"is_use_point":       isUsePoint,
			"archetype":          archetype,
			"price_set":          priceSet,
			"succ_campaign_id":   0,
			"count_get_campaign": 0,
			"voucher_amount":     vouDiscAmount,
			"order_type":         orderType,
		},
	}

	if referralCode != "" {
		session["referralCode"] = referralCode
	}

	m, b := map[string]interface{}{
		"customerSession": session,
		"responseContent": []string{
			"customerProfile",
			"customerSession",
			"triggeredCampaigns",
			"ruleFailureReasons",
		},
	}, new(bytes.Buffer)
	json.NewEncoder(b).Encode(m)
	if request, err = http.NewRequest("PUT", baseUrl, b); err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "ApiKey-v1 "+util.TalonApiKey)
	if response, err = client.Do(request); err != nil || (err == nil && response.StatusCode != 200) {
		err = json.NewDecoder(response.Body).Decode(&errorResponse)
		if strings.Contains(errorResponse.Message, "Only open sessions can be closed") {
			if responseData, err = GetCustomerSession(integrationCode); err != nil {
				return nil, err
			}

			var customerProfile *model.CustomerProfileData
			if customerProfile, err = GetCustomerProfile(profileCode); err != nil {
				return nil, err
			}

			responseData.CustomerProfile = &model.Profile{
				ID:                customerProfile.Profile.ID,
				CreatedDate:       customerProfile.Profile.CreatedDate,
				IntegrationID:     customerProfile.Profile.IntegrationID,
				Attributes:        customerProfile.Profile.Attributes,
				AccountID:         customerProfile.Profile.AccountID,
				ClosedSessions:    customerProfile.Profile.ClosedSessions,
				TotalSales:        customerProfile.Profile.TotalSales,
				LoyaltyMembership: customerProfile.Profile.LoyaltyMembership,
				LastActivity:      customerProfile.Profile.LastActivity,
			}

			return responseData, nil
		}

		err = errors.New("Invalid customer session")
		defer response.Body.Close()
		return nil, err
	}

	err = json.NewDecoder(response.Body).Decode(&responseData)

	defer response.Body.Close()

	return responseData, err
}

// SetUpCsvFileForReferral : function to set up csv file and import it into talon.one's referral api
func SetUpCsvFileForReferral(refCode, advocateID string) (err error) {
	var (
		fileName   string
		csvFile    *os.File
		hitCounter int8
	)

	data := [][]string{
		{"code", "advocateprofileintegrationid", "limitval"},
	}
	data = append(data, []string{refCode, advocateID, "0"})

	fileName = util.TalonFile
	if csvFile, err = os.Create(fileName); err != nil {
		err = errors.New(fmt.Sprintf("Failed creating CSV file: %s", err))
		return err
	}

	csvwriter := csv.NewWriter(csvFile)
	for _, row := range data {
		_ = csvwriter.Write(row)
	}
	csvwriter.Flush()
	defer csvFile.Close()

HitEndpoint:
	baseUrl = util.TalonHost + "/v1/applications/" + util.TalonApplicationID + "/campaigns/" + util.TalonCampaignID + "/import_referrals"

	file, _ := os.Open(csvFile.Name())
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("upFile", filepath.Base(file.Name()))
	io.Copy(part, file)
	writer.Close()
	if request, err = http.NewRequest("POST", baseUrl, body); err != nil {
		return err
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Authorization", "Bearer "+util.TalonToken)

	if response, err = client.Do(request); err != nil || (err == nil && response.StatusCode != 200) {
		if response.StatusCode == 401 {
			if hitCounter == 3 {
				err = errors.New("Error token invalid")
				return err
			}
			hitCounter++
			CreateSession()
			goto HitEndpoint
		}

		err = json.NewDecoder(response.Body).Decode(&errorResponse)
		if errorResponse.Message != "Import Error: Duplicate code error" {
			defer response.Body.Close()
			return err
		}
	}

	os.Remove(csvFile.Name())

	return err
}

// GetCampaignList : function to get campaign list of talon.one
func GetCampaignList(state string) (campaignIDs []int64, err error) {
	var (
		responseData *model.CampaignList
		hitCounter   int8
		param        string
	)

HitEndpoint:
	if state != "" {
		param = "?campaignState=" + state
	}

	baseUrl = util.TalonHost + "/v1/applications/" + util.TalonApplicationID + "/campaigns" + param
	if request, err = http.NewRequest("GET", baseUrl, nil); err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+util.TalonToken)

	if response, err = client.Do(request); err != nil || (err == nil && response.StatusCode != 200) {
		if response.StatusCode == 401 {
			if hitCounter == 3 {
				err = errors.New("Error token invalid")
				return nil, err
			}
			hitCounter++
			CreateSession()
			goto HitEndpoint
		}

		err = errors.New("Fail get campaign IDs")
		defer response.Body.Close()
		return nil, err
	}

	if err = json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		defer response.Body.Close()
		return nil, err
	}

	defer response.Body.Close()

	for _, v := range responseData.Data {
		campaignIDs = append(campaignIDs, int64(v.ID))
	}

	return campaignIDs, nil
}

// GetCampaignDetail : function to get campaign detail of talon.one
func GetCampaignDetail(campaignID string) (responseData *model.CampaignDetail, err error) {
HitEndpoint:
	baseUrl = util.TalonHost + "/v1/applications/" + util.TalonApplicationID + "/campaigns/" + campaignID
	if request, err = http.NewRequest("GET", baseUrl, nil); err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", "Bearer "+util.TalonToken)

	if response, err = client.Do(request); err != nil || (err == nil && response.StatusCode != 200) {
		if response.StatusCode == 401 {
			CreateSession()
			goto HitEndpoint
		}

		err = errors.New("Fail get campaign")
		defer response.Body.Close()
		return nil, err
	}

	if err = json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		defer response.Body.Close()
		return nil, err
	}

	defer response.Body.Close()

	return
}

// GetCustomerSession : function to get customer session data of talon.one
func GetCustomerSession(integrationCode string) (responseData *model.CustomerSessionReturn, err error) {
	baseUrl = util.TalonHost + "/v2/customer_sessions/" + integrationCode
	if request, err = http.NewRequest("GET", baseUrl, nil); err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "ApiKey-v1 "+util.TalonApiKey)

	if response, err = client.Do(request); err != nil || (err == nil && response.StatusCode != 200) {
		err = errors.New("Invalid get session")
		defer response.Body.Close()
		return nil, err
	}

	err = json.NewDecoder(response.Body).Decode(&responseData)

	defer response.Body.Close()

	return
}

// CreateSession : function to create a session for management api talon.one
func CreateSession() (err error) {
	var responseData *model.SessionResponse

	baseUrl = util.TalonHost + "/v1/sessions"
	m, b := map[string]interface{}{
		"email":    util.TalonEmail,
		"password": util.TalonPass,
	}, new(bytes.Buffer)
	json.NewEncoder(b).Encode(m)
	if request, err = http.NewRequest("POST", baseUrl, b); err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	if response, err = client.Do(request); err != nil || (err == nil && response.StatusCode != 201) {
		return err
	}

	if err = json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		defer response.Body.Close()
		return err
	}

	defer response.Body.Close()

	os.Setenv("TOKEN_TALON", responseData.Token)
	util.TalonToken = responseData.Token

	return nil
}

// ChangeTalonPoints : function to add or reduce loyalty points of a customer in talon.one
func ChangeTalonPoints(changeType, reason, integrationCode string, points float64) (err error) {
	baseUrl = util.TalonHost + "/v1/loyalty_programs/" + util.TalonLoyaltyID + "/profile/" + integrationCode + "/" + changeType

	m, b := map[string]interface{}{
		"points": points,
		"name":   reason,
	}, new(bytes.Buffer)
	json.NewEncoder(b).Encode(m)
	if request, err = http.NewRequest("POST", baseUrl, b); err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	if response, err = client.Do(request); err != nil || (err == nil && response.StatusCode != 204) {
		err = errors.New("Failed to change points")
		defer response.Body.Close()
		return err
	}

	return
}

// GetCustomerProfile : function to get customer profile data in talon.one
func GetCustomerProfile(profileCode string) (responseData *model.CustomerProfileData, err error) {
	baseUrl = util.TalonHost + "/v1/customer_profiles/" + profileCode + "/inventory?profile=true"
	if request, err = http.NewRequest("GET", baseUrl, nil); err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", "ApiKey-v1 "+util.TalonApiKey)

	if response, err = client.Do(request); err != nil || (err == nil && response.StatusCode != 200) {
		err = errors.New("Invalid get profile")
		defer response.Body.Close()
		return nil, err
	}

	err = json.NewDecoder(response.Body).Decode(&responseData)

	defer response.Body.Close()

	return
}
