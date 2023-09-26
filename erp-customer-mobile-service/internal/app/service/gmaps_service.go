package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/dto"
)

type IGmapsService interface {
	GetAutoComplete(ctx context.Context, req *dto.GetAutoCompleteRequest) (predictions *dto.GetAutoCompleteResponse, err error)
	GetGeocode(ctx context.Context, req *dto.GetGeocodeRequest) (res *dto.GetGeocodeResponse, err error)
}

type GmapsService struct {
	opt opt.Options
}

func NewGmapsService() IGmapsService {
	return &GmapsService{
		opt: global.Setup.Common,
	}
}

func (s *GmapsService) GetAutoComplete(ctx context.Context, req *dto.GetAutoCompleteRequest) (predictions *dto.GetAutoCompleteResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "GmapsService.Get")
	defer span.End()

	if req.Data.Search == "" {
		err = edenlabs.ErrorValidation("search", "kata kunci harus diisi")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	req.Data.Search = strings.ReplaceAll(req.Data.Search, " ", "%20")

	autoCompleteHost := s.opt.Env.GetString("gmaps.auto_complete_host")
	apiKey := s.opt.Env.GetString("gmaps.api_key")

	// Set Up Params
	params := fmt.Sprintf("?input=%s&language=id&components=country:id&key=%s", req.Data.Search, apiKey)

	client := &http.Client{}
	request, err := http.NewRequest("GET", autoCompleteHost+params, nil)
	if err != nil {
		return predictions, err
	}
	response, err := client.Do(request)
	if err != nil {
		return predictions, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return predictions, err
	}

	err = json.Unmarshal(body, &predictions)
	if err != nil {
		return predictions, err
	}

	return

}

func (s *GmapsService) GetGeocode(ctx context.Context, req *dto.GetGeocodeRequest) (data *dto.GetGeocodeResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "GmapsService.GetDetail")
	defer span.End()

	var (
		params, mainText, secondaryText, typeAddress string
		isHaveMainText                               bool
	)

	// both fields cannot be empty
	if req.Data.Latlng == "" && req.Data.PlaceID == "" {
		err = edenlabs.ErrorValidation("data.invalid", "place id atau lat long harus diisi")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// in gmaps only allowed one field, cannot both
	if req.Data.Latlng != "" && req.Data.PlaceID != "" {
		err = edenlabs.ErrorValidation("data.invalid", "pilih salah satu, place id atau lat long")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
	}

	if req.Data.PlaceID != "" {
		req.Data.Params = "place_id=" + req.Data.PlaceID
	}

	if req.Data.Latlng != "" {
		req.Data.Params = "latlng=" + req.Data.Latlng + "&result_type=premise|point_of_interest|route|establishment"
	}

	geocodeHost := s.opt.Env.GetString("gmaps.geocode_host")
	apiKey := s.opt.Env.GetString("gmaps.api_key")

	// Set Up Params
	params = fmt.Sprintf("?%s&language=id&key=%s", req.Data.Params, apiKey)

	client := &http.Client{}
	request, e := http.NewRequest("GET", geocodeHost+params, nil)
	if e != nil {
		return data, e
	}
	response, e := client.Do(request)
	if e != nil {
		return data, e
	}

	defer response.Body.Close()

	body, e := ioutil.ReadAll(response.Body)
	if e != nil {
		return data, e
	}
	var responseData dto.GetGeocodeResponse
	e = json.Unmarshal(body, &responseData)
	if e != nil {
		return data, e
	}
	// Handling for null data response
	if len(responseData.Results) == 0 {
		return data, e
	}

	// Only one data needed
	var results []dto.Geocode
	results = append(results, responseData.Results[0])
	data = &dto.GetGeocodeResponse{
		Results: results,
	}

	// filter type of address component to set on main text and secondary text, it's needed to show on mobile
	for _, v := range data.Results[0].AddressComponents {

		typeAddress = strings.Join(v.Types, ",")

		// Set up for main text
		if strings.Contains(typeAddress, "premise") {
			isHaveMainText = true
			mainText = v.LongName
		}

		if strings.Contains(typeAddress, "point_of_interest") && !isHaveMainText {
			isHaveMainText = true
			mainText = v.LongName
		}

		if strings.Contains(typeAddress, "route") && !isHaveMainText {
			isHaveMainText = true
			mainText = v.LongName
		}

		if strings.Contains(typeAddress, "administrative_area_level_4") && !isHaveMainText {
			isHaveMainText = true
			mainText = v.LongName
		}

		if strings.Contains(typeAddress, "administrative_area_level_3") && !isHaveMainText {
			isHaveMainText = true
			mainText = v.LongName
		}

		// Set up for secondary text
		if strings.Contains(typeAddress, "administrative_area_level_3") {
			secondaryText = v.LongName + ", "
		}

		if strings.Contains(typeAddress, "administrative_area_level_2") {
			secondaryText += v.LongName
		}
	}
	// set Main Text and Secondary Text
	data.Results[0].MainText = mainText
	data.Results[0].SecondaryText = secondaryText
	// Didn't need address components on response
	data.Results[0].AddressComponents = nil

	return
}
