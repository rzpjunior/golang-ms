package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/dto"
)

type ITalonService interface {
	UpdateCustomerProfileTalon(ctx context.Context, req *dto.TalonRequestUpdateCustomerProfile) (err error)
	UpdateCustomerSessionTalon(ctx context.Context, req *dto.TalonRequestUpdateCustomerSession) (responseData *dto.CustomerSessionReturn, err error)
	GetCustomerProfile(profileCode string) (responseData *dto.CustomerProfileData, err error)
}

type TalonService struct {
	opt opt.Options
}

func NewTalonService() ITalonService {
	return &TalonService{
		opt: global.Setup.Common,
	}
}

var (
	baseUrl                                                                                                                   string
	request                                                                                                                   *http.Request
	client                                                                                                                    = &http.Client{}
	response                                                                                                                  *http.Response
	errorResponse                                                                                                             *dto.ErrorResponse
	TalonHost, TalonApiKey, TalonApplicationID, TalonCampaignID, TalonLoyaltyID, TalonEmail, TalonPass, TalonToken, TalonFile string
)

// UpdateCustomerProfileTalon : func (s *TalonService) tion to insert or update customer profile in talon
func (s *TalonService) UpdateCustomerProfileTalon(ctx context.Context, req *dto.TalonRequestUpdateCustomerProfile) (err error) {
	var (
		responseData *dto.CustomerProfileReturn
	)

	TalonHost = s.opt.Env.GetString("talon_one.host")
	TalonApiKey = s.opt.Env.GetString("talon_one.api_key")

	baseUrl = TalonHost + "/v2/customer_profiles/" + req.ProfileCode

	attributes := map[string]interface{}{
		"region":            req.Region,
		"customer_type":     req.CustomerType,
		"registration_date": req.CreatedDate,
	}

	if len(req.ReferrerData) > 0 {
		attributes["advocate_id"] = req.ReferrerData[0]
		attributes["advocate_ref_code"] = req.ReferrerData[1]
	}

	m, b := map[string]interface{}{
		"attributes": attributes,
	}, new(bytes.Buffer)
	json.NewEncoder(b).Encode(m)
	if request, err = http.NewRequest("PUT", baseUrl, b); err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "ApiKey-v1 "+TalonApiKey)
	if response, err = client.Do(request); err != nil {
		err = errors.New("Invalid customer profile")
		defer response.Body.Close()
		return err
	}

	err = json.NewDecoder(response.Body).Decode(&responseData)

	defer response.Body.Close()

	return
}

// UpdateCustomerSessionTalon : func (s *TalonService) tion to insert or update customer session in talon
func (s *TalonService) UpdateCustomerSessionTalon(ctx context.Context, req *dto.TalonRequestUpdateCustomerSession) (responseData *dto.CustomerSessionReturn, err error) {
	var cartItems []interface{}

	TalonHost = s.opt.Env.GetString("talon_one.host")
	TalonApiKey = s.opt.Env.GetString("talon_one.api_key")

	baseUrl = TalonHost + "/v2/customer_sessions/" + req.IntegrationCode + "?dry=" + req.IsDry

	if req.Status != "cancelled" {
		cartItems = []interface{}{}
		for _, v := range req.ItemList {

			item := map[string]interface{}{
				"name":     v.ItemName,
				"sku":      v.ItemCode,
				"category": v.ClassName,
				"price":    v.UnitPrice,
				"quantity": v.OrderQty,
				"weight":   v.UnitWeight,
			}

			cartItems = append(cartItems, item)
		}
	}

	session := make(map[string]interface{})
	session = map[string]interface{}{
		"profileId": req.ProfileCode,
		"state":     req.Status,
		"cartItems": cartItems,
		"attributes": map[string]interface{}{
			"eden_point_earned":  0,
			"is_use_point":       req.IsUsePoint,
			"archetype":          req.Archetype,
			"price_set":          req.PriceSet,
			"succ_campaign_id":   0,
			"count_get_campaign": 0,
			"voucher_amount":     req.VouDiscAmount,
			"order_type":         req.OrderType,
		},
	}

	if req.ReferralCode != "" {
		session["referralCode"] = req.ReferralCode
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
	request.Header.Set("Authorization", "ApiKey-v1 "+TalonApiKey)
	if response, err = client.Do(request); err != nil {
		err = json.NewDecoder(response.Body).Decode(&errorResponse)
		if strings.Contains(errorResponse.Message, "Only open sessions can be closed") {
			if responseData, err = s.GetCustomerSession(req.IntegrationCode); err != nil {
				return nil, err
			}
			return responseData, nil
		}

		err = errors.New("Invalid customer session")
		defer response.Body.Close()
		return nil, err
	}
	body, _ := ioutil.ReadAll(response.Body)

	defer response.Body.Close()

	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return
	}

	return responseData, err
}

// SetUpCsvFileForReferral : func (s *TalonService) tion to set up csv file and import it into talon.one's referral api
func (s *TalonService) SetUpCsvFileForReferral(req *dto.TalonRequestSetUpCsvFileForReferral) (err error) {
	var (
		fileName   string
		csvFile    *os.File
		hitCounter int8
	)
	TalonHost = s.opt.Env.GetString("talon_one.host")
	TalonApiKey = s.opt.Env.GetString("talon_one.api_key")
	TalonApplicationID = s.opt.Env.GetString("talon_one.APPLICATION_ID")
	TalonCampaignID = s.opt.Env.GetString("talon_one.CAMPAIGN_ID")
	TalonFile = s.opt.Env.GetString("talon_one.FILE")

	data := [][]string{
		{"code", "advocateprofileintegrationid", "limitval"},
	}
	data = append(data, []string{req.ReferralCode, req.AdvocateID, "0"})

	fileName = TalonFile
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
	baseUrl = TalonHost + "/v1/applications/" + TalonApplicationID + "/campaigns/" + TalonCampaignID + "/import_referrals"

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
	request.Header.Set("Authorization", "Bearer "+TalonToken)

	if response, err = client.Do(request); err != nil || (err == nil && response.StatusCode != 200) {
		if response.StatusCode == 401 {
			if hitCounter == 3 {
				err = errors.New("Error token invalid")
				return err
			}
			hitCounter++
			s.CreateSession()
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

// GetCustomerSession : func (s *TalonService) tion to get customer session data of talon.one
func (s *TalonService) GetCustomerSession(integrationCode string) (responseData *dto.CustomerSessionReturn, err error) {

	TalonHost = s.opt.Env.GetString("talon_one.host")
	TalonApiKey = s.opt.Env.GetString("talon_one.api_key")

	baseUrl = TalonHost + "/v2/customer_sessions/" + integrationCode
	if request, err = http.NewRequest("GET", baseUrl, nil); err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "ApiKey-v1 "+TalonApiKey)

	if response, err = client.Do(request); err != nil || (err == nil && response.StatusCode != 200) {
		err = errors.New("Invalid get session")
		defer response.Body.Close()
		return nil, err
	}

	err = json.NewDecoder(response.Body).Decode(&responseData)

	defer response.Body.Close()

	return
}

// ChangeTalonPoints : func (s *TalonService) tion to add or reduce loyalty points of a customer in talon.one
func (s *TalonService) ChangeTalonPoints(changeType, reason, integrationCode string, points float64) (err error) {

	TalonHost = s.opt.Env.GetString("talon_one.host")
	TalonApiKey = s.opt.Env.GetString("talon_one.api_key")
	TalonLoyaltyID = s.opt.Env.GetString("talon_one.loyalty_id")

	baseUrl = TalonHost + "/v1/loyalty_programs/" + TalonLoyaltyID + "/profile/" + integrationCode + "/" + changeType

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

// GetCustomerProfile : func (s *TalonService) tion to get customer profile data in talon.one
func (s *TalonService) GetCustomerProfile(profileCode string) (responseData *dto.CustomerProfileData, err error) {
	TalonHost = s.opt.Env.GetString("talon_one.host")
	TalonApiKey = s.opt.Env.GetString("talon_one.api_key")

	baseUrl = TalonHost + "/v1/customer_profiles/" + profileCode + "/inventory?profile=true"
	if request, err = http.NewRequest("GET", baseUrl, nil); err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", "ApiKey-v1 "+TalonApiKey)

	if response, err = client.Do(request); err != nil || (err == nil && response.StatusCode != 200) {
		err = errors.New("Invalid get profile")
		defer response.Body.Close()
		return nil, err
	}

	err = json.NewDecoder(response.Body).Decode(&responseData)
	defer response.Body.Close()

	return
}

// CreateSession : function to create a session for management api talon.one
func (s *TalonService) CreateSession() (err error) {
	var responseData *dto.SessionResponse

	TalonHost = s.opt.Env.GetString("talon_one.host")
	TalonApiKey = s.opt.Env.GetString("talon_one.api_key")
	TalonEmail = s.opt.Env.GetString("talon_one.email")
	TalonPass = s.opt.Env.GetString("talon_one.pass")

	baseUrl = TalonHost + "/v1/sessions"
	m, b := map[string]interface{}{
		"email":    TalonEmail,
		"password": TalonPass,
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

	// os.SetEnv("TOKEN_TALON", responseData.Token)
	TalonToken = responseData.Token

	return nil
}
