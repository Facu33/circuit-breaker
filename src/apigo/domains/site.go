package domains

import (
	"../utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

type Site struct {
	ID                 string   `json:"id"`
	Name               string   `json:"name"`
	CountryID          string   `json:"country_id"`
	SaleFeesMode       string   `json:"sale_fees_mode"`
	MercadopagoVersion int      `json:"mercadopago_version"`
	DefaultCurrencyID  string   `json:"default_currency_id"`
	ImmediatePayment   string   `json:"immediate_payment"`
	PaymentMethodIds   []string `json:"payment_method_ids"`
	Settings           struct {
		IdentificationTypes      []string `json:"identification_types"`
		TaxpayerTypes            []string `json:"taxpayer_types"`
		IdentificationTypesRules []struct {
			IdentificationType string `json:"identification_type"`
			Rules              []struct {
				EnabledTaxpayerTypes []string `json:"enabled_taxpayer_types"`
				BeginsWith           string   `json:"begins_with"`
				Type                 string   `json:"type"`
				MinLength            int      `json:"min_length"`
				MaxLength            int      `json:"max_length"`
			} `json:"rules"`
		} `json:"identification_types_rules"`
	} `json:"settings"`
	Currencies []struct {
		ID     string `json:"id"`
		Symbol string `json:"symbol"`
	} `json:"currencies"`
	Categories []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"categories"`
}

//Receiver
func (site *Site) Get() *utils.ApiError {

	if site.ID == "" {
		return &utils.ApiError{
			Message: "Site ID is empty",
			Status:  http.StatusBadRequest,
		}
	}
	url := fmt.Sprintf("%s%s", utils.UrlSite, site.ID)

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

	if err := json.Unmarshal(data, &site); err != nil {
		return &utils.ApiError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return nil
}

func (site *Site) GetWg(group *sync.WaitGroup, apiError *utils.ApiError) {

	if site.ID == "" {
		group.Done()
		apiError.Message = "Site ID is empty"
		apiError.Status = http.StatusBadRequest
		return
	}
	url := fmt.Sprintf("%s%s", utils.UrlSite, site.ID)

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

	if err := json.Unmarshal(data, &site); err != nil {
		group.Done()
		apiError.Message = err.Error()
		apiError.Status = http.StatusInternalServerError
		return
	}
	print(apiError.Message)

	group.Done()
	return
}

func (site *Site) GetCh(results chan Result) {

	if site.ID == "" {
		results <- Result{
			ApiError: &utils.ApiError{
				Message: "Site ID is empty",
				Status:  http.StatusBadRequest,
			},
		}
		return
	}
	url := fmt.Sprintf("%s%s", utils.UrlSite, site.ID)

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

	if err := json.Unmarshal(data, &site); err != nil {
		results <- Result{
			ApiError: &utils.ApiError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			},
		}
		return
	}
	results <- Result{
		Site: site,
	}
	return
}
