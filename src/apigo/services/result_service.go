package services

import (
	"../domains"
	"../utils"
	"sync"
)

func GetResult(userId int) (*domains.Result, *utils.ApiError) {
	user := domains.User{
		ID: userId,
	}

	err := user.Get()
	if err != nil {
		return nil, err
	}

	country := domains.Country{
		ID: user.CountryID,
	}

	site := domains.Site{
		ID: user.SiteID,
	}

	go country.Get()
	go site.Get()

	result := domains.Result{
		User:    &user,
		Site:    &site,
		Country: &country,
	}

	return &result, nil

}

func GetResultWg(userId int) (*domains.Result, *utils.ApiError) {
	user := domains.User{
		ID: userId,
	}

	err := user.Get()

	if err != nil {
		return nil, err
	}

	country := domains.Country{
		ID: user.CountryID,
	}

	site := domains.Site{
		ID: user.SiteID,
	}

	var waitGroup sync.WaitGroup
	apiError := utils.ApiError{}
	waitGroup.Add(2)
	go country.GetWg(&waitGroup,&apiError)
	go site.GetWg(&waitGroup,&apiError)
	waitGroup.Wait()
	if apiError.Status != 0 {

		return nil, &apiError
	}

	result := domains.Result{
		User:    &user,
		Site:    &site,
		Country: &country,
	}

	return &result, nil

}

func GetResultCh(userId int) (*domains.Result, *utils.ApiError) {
	user := domains.User{
		ID: userId,
	}

	err := user.Get()

	if err != nil {
		return nil, err
	}

	country := domains.Country{
		ID: user.CountryID,
	}

	site := domains.Site{
		ID: user.SiteID,
	}
	resultCh := make(chan domains.Result)
	go country.GetCh(resultCh)
	go site.GetCh(resultCh)

	finalResult := &domains.Result{}
	for i := 0; i < 2; i++ {
		result := <-resultCh
		if result.ApiError != nil {
			return nil, result.ApiError
			break
		}

		if result.Site != nil {
			finalResult.Site = result.Site
		}
		if result.Country != nil {
			finalResult.Country = result.Country
		}
	}

	finalResult.User = &user
	return finalResult, nil

}
