package controllers

import (
	"../domains"
	"../services"
	"../utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetResultFromApi(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param(paramUserId))

	if err != nil {
		apiErr := &utils.ApiError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
		c.JSON(apiErr.Status, apiErr)
		return
	}

	res, apiErr := services.GetResult(userId)

	if apiErr != nil {
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, res)

}

func GetResultWgFromApi(c *gin.Context) {

	userId, err := strconv.Atoi(c.Param(paramUserId))

	if err != nil {
		apiErr := &utils.ApiError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
		c.JSON(apiErr.Status, apiErr)
		return
	}

	res, apiErr := services.GetResultWg(userId)

	if apiErr != nil {
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, res)

}

func GetResultChFromApi(c *gin.Context) {

	if domains.CircuitBreakerGlobal.State == "OPEN" {
		apiErr := domains.Circuit()
		c.JSON(apiErr.Status, apiErr)
		return
	}
	if domains.CircuitBreakerGlobal.State == "HALF" {
		apiErr := domains.Circuit()
		c.JSON(apiErr.Status, apiErr)
		return
	}

	userId, err := strconv.Atoi(c.Param(paramUserId))

	if err != nil {
		apiErr := &utils.ApiError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
		c.JSON(apiErr.Status, apiErr)
		return
	}

	res, apiErr := services.GetResultCh(userId)
	if apiErr != nil {
		if domains.CircuitBreakerGlobal.ErrorsLimit == domains.CircuitBreakerGlobal.Errors {
			err := domains.Circuit()
			c.JSON(err.Status, err)
			return
		}
		c.JSON(apiErr.Status, apiErr)
		domains.CircuitBreakerGlobal.Errors++
		return
	}

	c.JSON(http.StatusOK, res)
	return

}
