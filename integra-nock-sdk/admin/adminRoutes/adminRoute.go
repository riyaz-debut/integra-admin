package adminRoute

import (

	// chaincodeController "restapi-gonic/src/chaincode-controller"

	adminController "integra-nock-sdk/admin/adminController"
	"log"

	reg "integra-nock-sdk/regsdk"

	configImpl "github.com/hyperledger/fabric-sdk-go/pkg/core/config"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	_ "github.com/spacemonkeygo/openssl"
)

var configProvider = configImpl.FromFile("./connection-org1.yaml")
var sdk *fabsdk.FabricSDK
var mspClient *msp.Client

// var channelID = "mychannel"

var resmgmtClient, sdkreg, clCtx = reg.RegSdk()

// user register route
func Register(c *gin.Context) {
	response := adminController.RegisterUser(c)
	c.IndentedJSON(response.Status, response)
}

// admin login route
func AdminLogin(c *gin.Context) {
	response := adminController.AdminLogin(c)
	c.IndentedJSON(response.Status, response)
}

// user login route
func UserLogin(c *gin.Context) {
	response := adminController.UserLogin(c)
	c.IndentedJSON(response.Status, response)
}

//register new user route
func GetCurrentUser(c *gin.Context) {
	response := adminController.GetUserData(c)
	c.IndentedJSON(response.Status, response)
}

// get users list UserList DisableUser
func GetUsersList(c *gin.Context) {
	response := adminController.GetAllUsersData(c)
	c.IndentedJSON(response.Status, response)
}

// DisableUser
func DisableUserData(c *gin.Context) {
	id := c.Param("user_id")
	response := adminController.DisableUser(c, id)
	c.IndentedJSON(response.Status, response)
}

// DisableUser
func EnableUserData(c *gin.Context) {
	id := c.Param("user_id")
	response := adminController.EnableUser(c, id)
	c.IndentedJSON(response.Status, response)
}

func ChaincodeList(c *gin.Context) {

	log.Println("@@@@@@@@@@@@@@@ inside list")
	response := adminController.ChaincodesList(c)
	c.IndentedJSON(response.Status, response)
}

// chaincode commit status route
func ChaincodeLogs(c *gin.Context) {
	response := adminController.CcLogs(c)
	c.IndentedJSON(response.Status, response)
}

//commit readiness check
func CommitReadiness(c *gin.Context) {
	log.Println("############################ sdkreg", sdkreg)
	response := adminController.CommitReadiness(resmgmtClient, c)
	c.IndentedJSON(response.Status, response)

}

// calling commitchaincode from controller file
func CommitChaincode(c *gin.Context) {
	log.Println("############################ sdkreg", sdkreg)
	response := adminController.CcCommit(resmgmtClient, c)
	c.IndentedJSON(response.Status, response)
}

//create chaincode updates by admin
func CreateUpdates(c *gin.Context) {
	cc_id := c.Param("cc_id")
	response := adminController.CreateUpdates(c, cc_id)
	c.IndentedJSON(response.Status, response)
}

//get chaincode checkupdate by id
func ChaincodeUpdateCheck(c *gin.Context) {
	response := adminController.CcUpdateCheck(c)
	c.IndentedJSON(response.Status, response)

}

// get all chaincode releases(updates) list
func GetReleasesList(c *gin.Context) {
	response := adminController.GetAllCCReleases(c)
	c.IndentedJSON(response.Status, response)
}

// get all chaincode releases(updates) list
func ViewLogs(c *gin.Context) {
	cu_id := c.Param("cu_id")
	response := adminController.Viewupdatelogs(c, cu_id)
	c.IndentedJSON(response.Status, response)
}

// DeleteCCRelease
func DeleteRelease(c *gin.Context) {
	cu_id := c.Param("cu_id")
	response := adminController.DeleteCCRelease(c, cu_id)
	c.IndentedJSON(response.Status, response)
}

//to get new update values from db
func InstallUpdate(c *gin.Context) {
	response := adminController.CcUpdate(c)
	c.IndentedJSON(response.Status, response)

}

// @@@@@@@@@@@@@@@@@@@@@@@@ Organization @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

//add organization route
func AddOrganizations(c *gin.Context) {
	response := adminController.AddOrg(c, sdk, resmgmtClient, clCtx)
	c.IndentedJSON(response.Status, response)
}

//add organization peers route
func AddPeers(c *gin.Context) {
	response := adminController.AddPeers(c)
	c.IndentedJSON(response.Status, response)
}

//func get list of Org except the login one
func ListOrganizations(c *gin.Context) {
	response := adminController.GetOrgs(c)
	c.IndentedJSON(response.Status, response)

}

//func get organization from a particular given org id
func Organization(c *gin.Context) {
	response := adminController.GetSingleOrgs(c)
	c.IndentedJSON(response.Status, response)

}

// GetAdminOrgs
//func get list of Org except the login one
func AdminOrganizationsList(c *gin.Context) {
	log.Println("i am in admin route")
	response := adminController.GetAdminOrgs(c)
	c.IndentedJSON(response.Status, response)

}

//func to sign newly added org
func GetModifiedConfig(c *gin.Context) {
	response := adminController.MConfig(c)
	c.IndentedJSON(response.Status, response)

}

//func to save signature for added org
func SaveSignatures(c *gin.Context) {
	response := adminController.SaveSign(c)
	c.IndentedJSON(response.Status, response)

}

// join new org to channel
func JoinOrgToChannel(c *gin.Context) {
	log.Println(" i am in join channel route ")
	org_id := c.Param("org_id")
	log.Println(" id is :", org_id)
	response := adminController.SaveChannelConfig(resmgmtClient, org_id)
	c.IndentedJSON(response.Status, response)
}

// admin login route
func JoinStatusUpdate(c *gin.Context) {
	response := adminController.OrgJoinStatus(c)
	c.IndentedJSON(response.Status, response)
}

// admin login route
func SignStatusUpdate(c *gin.Context) {
	response := adminController.OrgSignStatus(c)
	c.IndentedJSON(response.Status, response)
}
