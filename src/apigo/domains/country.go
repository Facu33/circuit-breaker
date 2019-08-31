package domains

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"../utils"
	"sync"
)

type Country struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	Locale             string `json:"locale"`
	CurrencyID         string `json:"currency_id"`
	DecimalSeparator   string `json:"decimal_separator"`
	ThousandsSeparator string `json:"thousands_separator"`
	TimeZone           string `json:"time_zone"`
	GeoInformation     struct {
		Location struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"location"`
	} `json:"geo_information"`
	States []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"states"`
}

//Receiver
func (country *Country) Get() *utils.ApiError {

	if country.ID == "" {
		return &utils.ApiError{
			Message: "Country ID is empty",
			Status:  http.StatusBadRequest,
		}
	}
	url := fmt.Sprintf("%s%s", utils.UrlCountries, country.ID)

	res, err := http.Get(url)

	if err != nil {
		return &utils.ApiError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return &utils.ApiError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	if err := json.Unmarshal(data, &country); err != nil {
		return &utils.ApiError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return nil
}

func (country *Country) GetWg(group *sync.WaitGroup, apiError *utils.ApiError) {

	if country.ID == "" {
		group.Done()
		apiError.Message = "Country ID is empty"
		apiError.Status = http.StatusBadRequest
		return
	}
	url := fmt.Sprintf("%s%s", utils.UrlCountries, country.ID)

	res, err := http.Get(url)

	if err != nil {
		group.Done()
		apiError.Message = err.Error()
		apiError.Status = http.StatusInternalServerError
		return
	}

	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		group.Done()
		apiError.Message = err.Error()
		apiError.Status = http.StatusInternalServerError
		return
	}

	if err := json.Unmarshal(data, &country); err != nil {
		group.Done()
		apiError.Message = err.Error()
		apiError.Status = http.StatusInternalServerError
		return
	}
	group.Done()
	return
}

func (country *Country) GetCh(results chan Result) {

	if country.ID == "" {
		results <- Result{
			ApiError: &utils.ApiError{
				Message: "Country ID is empty",
				Status:  http.StatusBadRequest,
			},
		}
		return
	}
	url := fmt.Sprintf("%s%s", utils.UrlCountries, country.ID)

	res, err := http.Get(url)

	if err != nil {
		results <- Result{
			ApiError: &utils.ApiError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			},
		}
		return
	}

	data, err := ioutil.ReadAll(res.Body)

	if err != nil {

		results <- Result{
			ApiError: &utils.ApiError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			},
		}
		return
	}

	if err := json.Unmarshal(data, &country); err != nil {

		results <- Result{
			ApiError: &utils.ApiError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			},
		}
		return
	}
	results <- Result{
		Country: country,
	}
	return
}
