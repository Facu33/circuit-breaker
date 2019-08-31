package main

import (
	"../controllers"
	"../domains"
	"github.com/gin-gonic/gin"
)

const (
	port = ":8084"
)

var (
	router = gin.Default()
)

func main() {
	domains.CircuitBreakerGlobal.NewCircuitBreaker("CLOSE", 0, 15000, 3 )

	router.GET("/users/:userId", controllers.GetUserFromApi)
	router.GET("/sites/:siteId", controllers.GetSiteFromApi)
	router.GET("/countries/:countryId", controllers.GetCountryFromApi)
	router.GET("/results/:userId", controllers.GetResultFromApi)
	router.GET("/resultsWg/:userId", controllers.GetResultWgFromApi)
	router.GET("/resultsCh/:userId", controllers.GetResultChFromApi)

	router.Run(port)

}
