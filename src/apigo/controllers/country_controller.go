package controllers

import (
"github.com/gin-gonic/gin"
"../services"
"net/http"
)

const (
	paramCountryId = "countryId"
)

func GetCountryFromApi(c *gin.Context) {
	countryId := c.Param(paramCountryId)

	res, err := services.GetCountry(countryId)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, res)
}
