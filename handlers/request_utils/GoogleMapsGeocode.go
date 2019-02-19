package request_utils

import (
	"fmt"
	"log"
	"os"
)

// Generated with https://mholt.github.io/json-to-go/
type GoogleResponse struct {
	PlusCode struct {
		CompoundCode string `json:"compound_code"`
		GlobalCode   string `json:"global_code"`
	} `json:"plus_code"`
	Results []struct {
		AddressComponents []struct {
			LongName  string   `json:"long_name"`
			ShortName string   `json:"short_name"`
			Types     []string `json:"types"`
		} `json:"address_components"`
		FormattedAddress string `json:"formatted_address"`
		Geometry         struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
			LocationType string `json:"location_type"`
			Viewport     struct {
				Northeast struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"northeast"`
				Southwest struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"southwest"`
			} `json:"viewport"`
		} `json:"geometry"`
		PlaceID  string `json:"place_id"`
		PlusCode struct {
			CompoundCode string `json:"compound_code"`
			GlobalCode   string `json:"global_code"`
		} `json:"plus_code,omitempty"`
		Types []string `json:"types"`
	} `json:"results"`
	Status string `json:"status"`
}

func GetAddressForLatLng(locLat string, locLng string) string {
	mapsKey := os.Getenv("GOOGLE_MAPS_KEY")

	request := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?latlng=%s,%s&key=%s",
		locLat, locLng, mapsKey)

	tmp := GoogleResponse{}

	response := GetJson(request, &tmp)

	log.Println(response)
	//result.data.results[0].formatted_address
	if tmp.Status != "OK" {
		log.Println(tmp.Status)
	}

	return tmp.Results[0].FormattedAddress
}
