package controllers

import (
	"../services"
	"../utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

const (
	paramUserId = "userId"
)

func GetUserFromApi(c *gin.Context) {

	userId, err := strconv.Atoi(c.Param(paramUserId))
	println(userId)
	if err != nil {
		apiErr := &utils.ApiError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
		c.JSON(apiErr.Status, apiErr)
		return
	}

	res, apiErr := services.GetUser(userId)

	if apiErr != nil {
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, res)

}
