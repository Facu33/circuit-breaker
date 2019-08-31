package controllers

import (
	"../services"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	paramSiteId = "siteId"
)

func GetSiteFromApi(c *gin.Context) {
	siteId := c.Param(paramSiteId)

	res, err := services.GetSite(siteId)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, res)
}
