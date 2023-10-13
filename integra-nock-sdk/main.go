package main

import (

	//admin route file path
	adminRouter "integra-nock-sdk/admin/adminRoutes"

	"log"

	database "integra-nock-sdk/database"
	"integra-nock-sdk/middlewares"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// var err error

	log.Println("inside main file")
	config :=
		database.Config{
			ServerName: "localhost:3306",
			User:       "riyaz",
			Password:   "RIYAZdebut@2021",
			DB:         "integradb",
		}

	connectionString := database.GetConnectionString(config)
	err := database.Connect(connectionString)
	if err != nil {
		panic(err.Error())
	}

	router := gin.Default()
	router.Use(CORSMiddleware())

	// =================================== User Authentication API'S ========================================

	router.POST("/user/register", adminRouter.Register)
	router.POST("/admin/login", adminRouter.AdminLogin)
	router.POST("/user/login", adminRouter.UserLogin)
	router.POST("/organization/join/status_update", adminRouter.JoinStatusUpdate)
	router.POST("/organization/sign/status_update", adminRouter.SignStatusUpdate)
	router.Use(middlewares.JwtAuthMiddleware())
	// router.Use(CORSMiddleware())
	router.GET("/user/home", adminRouter.GetCurrentUser)

	// user list
	router.GET("/user/list", adminRouter.GetUsersList)

	// Disable user
	router.POST("/disable/user/:user_id", adminRouter.DisableUserData)

	// Enable User
	router.POST("/enable/user/:user_id", adminRouter.EnableUserData)

	// =================================== ONLY ADMIN APIS ==================================================

	// ##################### CHAINCODE PART

	// router.GET("/chaincode/list", adminRouter.NetworkChaincodeList)
	router.POST("/chaincode/list", adminRouter.ChaincodeList)

	//route for check commitness
	router.POST("/chaincode/commit-readiness", adminRouter.CommitReadiness)

	//commit chaincode
	router.POST("/chaincode/commit", adminRouter.CommitChaincode)

	//create new chaincode/update
	router.POST("/chaincode/createupdates/:cc_id", adminRouter.CreateUpdates)

	// CC Relaeses(Updates) list
	router.GET("/releases/list", adminRouter.GetReleasesList)

	// view chaincode logs
	router.GET("/releases/logs/:cu_id", adminRouter.ViewLogs)

	// Delete chaincode new Release
	router.DELETE("/releases/delete/:cu_id", adminRouter.DeleteRelease)

	// ###################### ORGANIZATION PART

	//add organization post'
	router.POST("/organization", adminRouter.AddOrganizations)

	//add newly added org peers
	router.POST("/organization/peers", adminRouter.AddPeers)

	// ========================================= CLIENT CALLING ADMIN APIS ========================================

	// #################### CHAINCODE PART

	//api for storing chaincode data after approve and before commiting	``````````````````````````````````````````````````````````````````````
	router.POST("/chaincode/logs", adminRouter.ChaincodeLogs)

	//get route to get chaincode info by id
	router.POST("/chaincode/checkforupdates", adminRouter.ChaincodeUpdateCheck)

	//api to download and install update of chaincode

	router.POST("/chaincode/update", adminRouter.InstallUpdate)

	// #################### ORGANIZATIONS PART

	//fetching all orgs list except for logged in one for client
	router.POST("/organization/list", adminRouter.ListOrganizations)

	//fetch org
	router.POST("/single/organization", adminRouter.Organization)

	//fetching all orgs list except for logged in one
	router.GET("/organization", adminRouter.AdminOrganizationsList)

	//api to fetch modified config for particular org
	router.POST("/organization/modifiedconfig", adminRouter.GetModifiedConfig)

	//save signatures for newly added org
	router.POST("organization/savesign", adminRouter.SaveSignatures)

	// org3 join cahnnel
	router.POST("/organization/joinChannel/:org_id", adminRouter.JoinOrgToChannel)

	router.Run("localhost:5000")

}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
