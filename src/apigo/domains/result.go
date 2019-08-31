package domains
 import (
 	"../utils"
 )
type Result struct {
	User    *User
	Site    *Site
	Country *Country
	ApiError *utils.ApiError
}
