package dto

// Gmaps Get Geocode
type GetGeocodeRequest struct {
	Platform string     `json:"platform" valid:"required"`
	Data     GetGeocode `json:"data" valid:"required"`
}

type GetGeocode struct {
	PlaceID string `json:"place_id"`
	Latlng  string `json:"latlng"`
	Params  string `json:"-"`
}

type GetGeocodeResponse struct {
	Results []Geocode `json:"results"`
}

type Geocode struct {
	FormattedAddress  string              `json:"formatted_address"`
	Geometry          Location            `json:"geometry"`
	AddressComponents []AddressComponents `json:"address_components"`
	MainText          string              `json:"main_text"`
	SecondaryText     string              `json:"secondary_text"`
}

type Location struct {
	Location LatLong `json:"location"`
}

type LatLong struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type AddressComponents struct {
	LongName string   `json:"long_name"`
	Types    []string `json:"types"`
}

// Get Gmaps Autocomplete
type GetAutoCompleteRequest struct {
	Platform string          `json:"platform" valid:"required"`
	Data     GetAutoComplete `json:"data" valid:"required"`
}

type GetAutoComplete struct {
	Search string `json:"search"`
	Params string `json:"-"`
}

type GetAutoCompleteResponse struct {
	Predictions []Autocomplete `json:"predictions"`
}

type Autocomplete struct {
	Description          string              `json:"description"`
	PlaceID              string              `json:"place_id"`
	StructuredFormatting StructuredFormating `json:"structured_formatting"`
}

type StructuredFormating struct {
	MainText      string `json:"main_text"`
	SecondaryText string `json:"secondary_text"`
}
