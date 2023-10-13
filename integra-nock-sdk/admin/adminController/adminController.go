package adminController

import (
	"bytes"
	"fmt"
	"html"
	config "integra-nock-sdk/config"
	"integra-nock-sdk/database"
	"integra-nock-sdk/utils"
	"io/ioutil"
	"log"
	"os"

	// "encoding/json"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"integra-nock-sdk/models"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-config/protolator"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	contextApi "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"

	// configImpl "github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/resource"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"golang.org/x/crypto/bcrypt"
)

var userId uint
var orgId int

type User struct {
	Id        uint      `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	UserName  string    `json:"user_name" gorm:"unique"`
	Password  string    `json:"password"`
	OrgId     int       `json:"org_id"`
	OrgName   string    `json:"org_name"  binding:"required"`
	OrgMsp    string    `json:"org_msp"`
	Role      string    `json:"role"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" gorm:"autoUpdateTime"`
}

// Register or add new user
func RegisterUser(body *gin.Context) utils.Response {

	type Users struct {
		UserName string `json:"username" binding:"required" `
		Password string `json:"password" binding:"required"`
		OrgId    int    `json:"org_id"  binding:"required"`
		OrgName  string `json:"org_name"  binding:"required"`
		OrgMsp   string `json:"org_msp"  binding:"required"`
		Role     string `json:"role" binding:"required"`
		Status   int    `json:"status"`
	}

	var data Users
	log.Println("data in admin controller: ", data)

	if err := body.ShouldBindJSON(&data); err != nil {
		log.Println("getting error in mapping json data to input variable", err)
		response := utils.Response{
			Status:  422,
			Message: "api error",
			Err:     err,
		}
		return response
	}

	user := User{}

	user.UserName = data.UserName
	user.Password = data.Password
	user.OrgId = data.OrgId
	user.OrgName = data.OrgName
	user.OrgMsp = data.OrgMsp
	user.Role = data.Role
	user.Status = data.Status

	//turn password into hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("err encrypting password", err)
	}
	user.Password = string(hashedPassword)

	//remove spaces in username
	user.UserName = html.EscapeString(strings.TrimSpace(user.UserName))

	tableExists := database.Connector.HasTable(User{})

	// Create table for chaincode_creation
	if !tableExists {
		if err := database.Connector.CreateTable(User{}).Error; err != nil {
			log.Println("error creating user table ", err)
			response := utils.Response{
				Status:  500,
				Message: "users table does not created successfully",
				Err:     err,
			}

			return response
		}
	}

	if err = database.Connector.Create(&user).Error; err != nil {
		log.Println("error registering user", err)
		log.Println("err post route @@@@@@@@@@ ", err)
		response := utils.Response{
			Status:  500,
			Message: "user not registered",
			Err:     err,
		}
		return response
	}
	response := utils.Response{
		Status:  200,
		Message: "user registered successfully",
		Data:    user,
	}
	return response
}

func AdminLogin(body *gin.Context) utils.Response {

	type Users struct {
		UserName string `json:"user_name" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var data Users

	if err := body.ShouldBindJSON(&data); err != nil {
		log.Println("getting error in mapping json data to input variable")
		response := utils.Response{
			Status:  500,
			Message: "api error",
			// Data: ,
			Err: err,
		}
		return response
	}

	user := models.User{}

	user.Username = data.UserName
	user.Password = data.Password

	log.Println("user in admin controller: ", user)

	token, err := models.AdminCredentialsCheck(user.Username, user.Password)
	if err != nil {
		log.Println("error checking credentials", err)
		response := utils.Response{
			Status:  500,
			Message: "User token did not receive successfully",
			Err:     err,
		}
		return response
	}
	log.Println("token is :", token)
	response := utils.Response{
		Status:  200,
		Message: "User token receive successfully",
		Data:    token,
	}
	return response
}

//user login controller of admin
func UserLogin(body *gin.Context) utils.Response {

	type UsersData struct {
		UserName string `json:"user_name" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var data UsersData

	if err := body.ShouldBindJSON(&data); err != nil {
		log.Println("getting error in mapping json data to input variable")
		response := utils.Response{
			Status:  500,
			Message: "api error",
			// Data: ,
			Err: err,
		}
		return response
	}

	user := models.User{}

	user.Username = data.UserName
	user.Password = data.Password

	// logic for check user disable or not
	// struct type variable
	var userInfo Users

	log.Println("username in userlogin fx at admin side :", user.Username)
	if err := database.Connector.Table("users").Find(&userInfo, "user_name = ?", user.Username).Error; err != nil {

		response := utils.Response{
			Status:  500,
			Message: "User data not found",
			// Data: ,
			Err: err,
		}
		return response

	}

	status := userInfo.Status
	role := userInfo.Role
	log.Println("user id is :", userInfo.Id)

	type Organizations struct {
		Id             int    `gorm:"primaryKey;autoIncrement"`
		Name           string `json:"name"`
		MspId          string `json:"msp_id"`
		PeersCount     int    `json:"peers_count"`
		Config         string `json:"file" gorm:"type:text"`
		ModifiedConfig string `json:"modified_config" gorm:"type:text"`
		Join_Status    int    `json:"join_status"`
		CreatedAt      string `json:"created_at"`
		UpdatedAt      string `json:"updated_at" gorm:"autoUpdateTime"`
	}

	var organizations Organizations
	if err := database.Connector.Table("organizations").Find(&organizations, "name = ?", userInfo.OrgName).Error; err != nil {
		log.Println("error in org query :", err)
	}
	log.Println("organization data in user login fx :", organizations.Join_Status)
	// check for status for disabling the user
	if status == 0 {
		response := utils.Response{
			Status:  500,
			Message: "User does not have permission to login",
			// Data:    userInfo,
		}
		return response
	}

	token, err := models.UserCredentialsCheck(user.Username, user.Password)
	if err != nil {
		log.Println("error checking credentials", err)
		response := utils.Response{
			Status:  200,
			Message: "User token did not receive successfully",
			Err:     err,
		}
		return response
	}
	log.Println("token is :", token)

	// userRole := "admin"
	// if role == userRole{
	// 	response := utils.Response{
	// 		Status:  200,
	// 		Message: "User token receive successfully",
	// 		Data:    token,
	// 		Role:	 role,
	// 	}
	// 	return response
	// }
	response := utils.Response{
		Status:  200,
		Message: "User token receive successfully",
		Data:    token,
		Role:    role,
		OrgData: organizations,
	}
	return response
}

type Users struct {
	Id       uint   `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	UserName string `json:"user_name" gorm:"unique"`
	OrgId    int    `json:"org_id"`
	OrgName  string `json:"org_name"  binding:"required"`
	OrgMsp   string `json:"org_msp"`
	Role     string `json:"role"`
	Status   int    `json:"status"`
}

//user home dashboard
func GetUserData(body *gin.Context) utils.Response {
	log.Println("data in admin controller of currentuser data: ", body)

	//get user_id from context
	user_id := body.MustGet("user_id")

	//change user_id that is of interface type to uint type
	userId := user_id.(uint)

	//Extract user data from user table using user id
	var user Users

	if err := database.Connector.Table("users").Find(&user, "id = ?", userId).Error; err != nil {

		response := utils.Response{
			Status:  500,
			Message: "User data not found",
			// Data: ,
			Err: err,
		}
		return response

	}

	userId = user.Id
	orgId = user.OrgId

	response := utils.Response{
		Status:  200,
		Message: "User data found",
		Data:    user,
	}
	return response

}

// get users list
func GetAllUsersData(body *gin.Context) utils.Response {
	log.Println("data in admin controller of GetAllUsersData data: ", body)

	var user []Users
	if err := database.Connector.Table("users").Find(&user).Error; err != nil {

		response := utils.Response{
			Status:  500,
			Message: "User data not found",
			// Data: ,
			Err: err,
		}
		return response

	}
	response := utils.Response{
		Status:  200,
		Message: "User data found",
		Data:    user,
	}
	return response

}

// disabling the user
func DisableUser(body *gin.Context, id string) utils.Response {
	log.Println("data in admin controller of disable data: ", body)

	//converting user id into int
	userId, err := strconv.Atoi(id)
	if err != nil {
		response := utils.Response{Status: 500, Message: "error converting org id to int", Err: err}
		return response
	}

	// convert int id to uint type
	user_id := uint(userId)

	// struct type variable
	var user Users

	if err := database.Connector.Table("users").Find(&user, "id = ?", user_id).Error; err != nil {

		response := utils.Response{
			Status:  500,
			Message: "User data not found",
			// Data: ,utils.Response
			Err: err,
		}
		return response

	}

	user_id = user.Id
	status := user.Status
	role := user.Role
	// check for status for disabling the user
	if status == 1 && role != "admin" {
		database.Connector.Table("users").Where("id = ?", user_id).Update("status", 0)
		response := utils.Response{
			Status:  200,
			Message: "Disable user successfully",
			Data:    user,
		}
		return response
	}

	response := utils.Response{
		Status:  200,
		Message: "Disabling not alloed",
		Data:    user,
	}
	return response

}

// enable user
func EnableUser(body *gin.Context, id string) utils.Response {
	log.Println("data in admin controller of enable data: ", body)

	//converting user id into int
	userId, err := strconv.Atoi(id)
	if err != nil {
		response := utils.Response{Status: 500, Message: "error in enable fx converting org id to int", Err: err}
		return response
	}

	// convert int id to uint type
	user_id := uint(userId)

	// struct type variable
	var user Users

	if err := database.Connector.Table("users").Find(&user, "id = ?", user_id).Error; err != nil {
		response := utils.Response{
			Status:  500,
			Message: "User data not found",
			// Data: ,
			Err: err,
		}
		return response

	}

	user_id = user.Id
	status := user.Status
	// check for status for disabling the user
	if status == 0 {
		database.Connector.Table("users").Where("id = ?", user_id).Update("status", 1)
		response := utils.Response{
			Status:  200,
			Message: "enable user successfully",
			Data:    user,
		}
		return response
	}
	response := utils.Response{
		Status:  200,
		Message: "Enabling not allowed",
		Data:    user,
	}
	return response

}

type ChaincodeLists struct {
	Id int `json:"id" gorm:"primary key;autoincrement"`
	// CC_ID     int    `json:"cc_id"`
	Name      string `json:"name"`
	Label     string `json:"label"`
	Version   string `json:"version"`
	Sequence  int    `json:"sequence"`
	Status    int    `json:"status"`
	Url       string `json:"url"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func ChaincodesList(body *gin.Context) utils.Response {

	var chaincode_lists []ChaincodeLists

	if err := database.Connector.Table("chaincode_lists").Find(&chaincode_lists).Error; err != nil {
		log.Println("error in organizations querying ", err)
		response := utils.Response{
			Status:  404,
			Message: "No chaincode found",
			Err:     err,
		}
		return response
	}

	log.Println("chaincodes found is ", chaincode_lists)

	response := utils.Response{
		Status:  200,
		Message: "chaincodes found successfully",
		Data:    chaincode_lists,
	}
	log.Println("response", response)
	return response
}

type ChaincodeUpdates struct {
	Id        int       `json:"id" gorm:"primary key;autoincrement"`
	CC_ID     int       `json:"cc_id"` // gorm:"foreignKey:CC_ID"
	Name      string    `json:"name"`
	Label     string    `json:"label"`
	Version   string    `json:"version"`
	Sequence  int       `json:"sequence"`
	Status    int       `json:"status"`
	Url       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"constraint:OnUpdate:CASCADE"`
}

//create update package by admin
func CreateUpdates(body *gin.Context, cc_id string) utils.Response {
	log.Println("create update fx")
	//converting user id into int
	id, err := strconv.Atoi(cc_id)
	if err != nil {
		response := utils.Response{Status: 500, Message: "error in enable fx converting org id to int", Err: err}
		return response
	}
	log.Println("Id in create update fx :", id)
	var ccList ChaincodeLists
	if err := database.Connector.Table("chaincode_lists").Find(&ccList, "id = ?", id).Error; err != nil {
		response := utils.Response{
			Status:  500,
			Message: "chaincode data not found",
			// Data: ,
			Err: err,
		}
		return response

	}

	id = ccList.Id
	var data ChaincodeUpdates
	if err := body.BindJSON(&data); err != nil {
		log.Println("err post route @@@@@@@@@@ ", err)
		response := utils.Response{
			Status:  0,
			Message: "payload error in data binding",
			Err:     err,
		}
		return response

	} else {
		log.Println("payload from backend  @@@@@@@@@@@@@@@", data)
	}

	tableExists := database.Connector.HasTable(&ChaincodeUpdates{})
	if !tableExists {
		if err := database.Connector.CreateTable(ChaincodeUpdates{}).Error; err != nil {
			log.Println("error creating table ", err)
			response := utils.Response{
				Status:  500,
				Message: "chaincode updates table does not created successfully",
				Err:     err,
			}
			return response
		}
	}

	static_status := 0

	chaincode_updates := ChaincodeUpdates{CC_ID: ccList.Id, Name: data.Name, Label: data.Name, Version: data.Version, Sequence: data.Sequence, Status: static_status, Url: data.Url, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := database.Connector.Create(&chaincode_updates).Error; err != nil {
		response := utils.Response{
			Status:  500,
			Message: "Cc Update Table Creation Failed",
			Err:     err,
		}
		return response
	}
	log.Println("data found is ", chaincode_updates)
	response := utils.Response{
		Status:  200,
		Message: "chaincode update created successfully",
		Data:    chaincode_updates,
	}
	return response
}

//struct to insert chaincode entry and commit status
type ChaincodeCommit struct {
	Id int `json:"id" gorm:"primary key;autoincrement"`
	// ccID int       `json:"cc_id"`
	Name      string    `json:"name"`
	Label     string    `json:"label"`
	Version   string    `json:"version"`
	Sequence  int       `json:"sequence"`
	Status    int       `json:"status,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" gorm:"autoUpdateTime"`
}

type ChaincodeLog struct {
	Id        int       `json:"id" gorm:"primary key;autoincrement"`
	CuId      int       `json:"cu_id" gorm:"foreignKey:CuId"`
	Name      string    `json:"name"`
	Label     string    `json:"label"`
	Version   string    `json:"version"`
	Sequence  int       `json:"sequence"`
	OrgId     int       `json:"org_id"`
	OrgName   string    `json:"org_name"`
	OrgMsp    string    `json:"msp_id"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func CcLogs(body *gin.Context) utils.Response {
	//struct to insert chaincode entry and commit status
	type ChaincodeDataReceive struct {
		Id       int    `json:"id,omitempty"`
		Name     string `json:"name,omitempty"`
		Label    string `json:"label,omitempty"`
		Version  string `json:"version,omitempty"`
		Sequence int    `json:"sequence,omitempty"`
		OrgName  string `json:"org_name,omitempty"`
		OrgId    int    `json:"org_id,omitempty"`
		OrgMsp   string `json:"msp_id,omitempty"`
		Status   int    `json:"status,omitempty"`
	}

	var data ChaincodeDataReceive
	if err := body.BindJSON(&data); err != nil {
		response := utils.Response{
			Status:  422,
			Message: "payload error",
			// Data: ,
			Err: err,
		}
		return response

	} else {
		log.Println("payload from backend  @@@@@@@@@@@@@@@", data)
		fmt.Println("req.body = ", reflect.TypeOf(data))
	}

	tableExists := database.Connector.HasTable(&ChaincodeLog{})
	// Create table for chaincode-logs
	if !tableExists {
		if err := database.Connector.CreateTable(&ChaincodeLog{}).Error; err != nil {
			log.Println("error creating table ", err)
			response := utils.Response{
				Status:  500,
				Message: "chaincode logs table does not created successfully",

				Err: err,
			}
			return response
		}
	}

	chaincode_logs := &ChaincodeLog{CuId: data.Id, Name: data.Name, Label: data.Label, Version: data.Version, Sequence: data.Sequence, OrgName: data.OrgName, OrgId: data.OrgId, OrgMsp: data.OrgMsp, Status: data.Status}

	if err := database.Connector.Create(&chaincode_logs).Error; err != nil {
		response := utils.Response{
			Status:  500,
			Message: "Data not entered successfully in chaincode_logs",
			Err:     err,
		}
		return response
	}
	response := utils.Response{
		Status:  200,
		Message: "chaincodes approval done for " + data.OrgName,
		Data:    chaincode_logs,
	}
	return response
}

//struct for checking chaincode commitreadiness
type Commitreadiness struct {
	ChannelID string `json:"channel_name"`
	CcId      string `json:"cc_id"`
	Version   string `json:"cc_version"`
	Sequence  int    `json:"sequence"`
}

// checking commitness
func CommitReadiness(orgResMgmt *resmgmt.Client, body *gin.Context) utils.Response {
	var check Commitreadiness
	log.Println("entering commit-readiness chaincode controller", body)
	// Call BindJSON to bind the received JSON to
	if err := body.BindJSON(&check); err != nil {
		response := utils.Response{
			Status:  422,
			Message: "error approving installed package",
			Err:     err,
		}
		return response

	} else {
		log.Println("payload from backend  @@@@@@@@@@@@@@@", check)

	}

	req := resmgmt.LifecycleCheckCCCommitReadinessRequest{
		Name:              check.CcId,
		Version:           check.Version,
		EndorsementPlugin: "escc",
		ValidationPlugin:  "vscc",
		Sequence:          int64(check.Sequence),
		InitRequired:      true,
	}
	resp, err := orgResMgmt.LifecycleCheckCCCommitReadiness(check.ChannelID, req, resmgmt.WithTargetEndpoints(config.PEER2), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		response := utils.Response{
			Status:  500,
			Message: "chaincode commit-readiness failed",
			Err:     err,
		}
		log.Fatalln(err)
		return response
	}
	var data interface {
	} = resp
	response := utils.Response{
		Status:  200,
		Message: "chaincode ready to commit",
		Data:    data,
	}
	return response
}

//function for chaincode commiting
func CcCommit(orgResMgmt *resmgmt.Client, body *gin.Context) utils.Response {
	log.Println("in commit controller fx")
	// struct to send data as payload in another api
	type ChaincodeData struct {
		ReleaseId int `json:"release_id"`
	}

	var commit ChaincodeData
	log.Println("Body commit controller fx", body)
	//binding payload data
	if err := body.BindJSON(&commit); err != nil {
		log.Println("2nd commit controller fx")
		response := utils.Response{
			Status:  422,
			Message: "error in payload data ",
			// Data:    nil,
			Err: err,
		}
		return response

	} else {
		log.Println("payload from backend  @@@@@@@@@@@@@@@", commit)

	}
	log.Println("3rd commit controller fx")
	type ChaincodeUpdates struct {
		Id        int       `json:"id"`
		CcId      string    `json:"cc_id"`
		Name      string    `json:"name"`
		Label     string    `json:"label"`
		Version   string    `json:"version"`
		Sequence  int       `json:"sequence"`
		Status    int       `json:"status"`
		Url       string    `json:"url"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	}

	var chaincode_updates ChaincodeUpdates
	if err := database.Connector.Table("chaincode_updates").Find(&chaincode_updates, "id = ?", commit.ReleaseId).Error; err != nil {
		response := utils.Response{
			Status:  500,
			Message: "User data not found",
			// Data: ,
			Err: err,
		}
		return response

	}
	log.Println("chaincode update data in commit fx :", chaincode_updates)
	req := resmgmt.LifecycleCommitCCRequest{
		Name:              chaincode_updates.Name,
		Version:           chaincode_updates.Version,
		Sequence:          int64(chaincode_updates.Sequence),
		EndorsementPlugin: "escc",
		ValidationPlugin:  "vscc",
		// SignaturePolicy:   ccPolicy,
		InitRequired: true,
	}
	txnID, err := orgResMgmt.LifecycleCommitCC(config.CHANNEL_ID, req, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithTargetEndpoints(config.PEER1, config.PEER2), resmgmt.WithOrdererEndpoint(config.ORDERER_ENDPOINT))
	if err != nil {
		log.Println("commitCC Fatal", err)
		response := utils.Response{
			Status:  500,
			Message: "Commiting chaincode failed",
			Err:     err,
		}
		return response
	}

	log.Println("commited chaincode transaction id @@@@@@@@@@@@@@@@@@@@", txnID)
	log.Println("######################################")
	log.Println("query commitedness")
	log.Println("######################################")

	reqQuery := resmgmt.LifecycleQueryCommittedCCRequest{
		Name: chaincode_updates.Name,
	}
	resp, err := orgResMgmt.LifecycleQueryCommittedCC(config.CHANNEL_ID, reqQuery, resmgmt.WithTargetEndpoints(config.PEER1), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		response := utils.Response{
			Status:  404,
			Message: "querying chaincode commitness failed",
			Err:     err,
		}

		return response
	}

	if err := database.Connector.Table("chaincode_updates").Where("id = ?", commit.ReleaseId).Update("status", 1).Error; err != nil {
		response := utils.Response{
			Status:  500,
			Message: "orgnizations table does not created successfully",
			Err:     err,
		}
		return response
	}

	if err := database.Connector.Table("chaincode_updates").Find(&chaincode_updates, "id = ?", commit.ReleaseId).Error; err != nil {
		response := utils.Response{
			Status:  500,
			Message: "User data not found",
			// Data: "",
			Err: err,
		}
		return response

	}

	// chaincod list struct
	type Chaincode_lists struct {
		Id        int       `json:"id" gorm:"primaryKey;autoIncrement"`
		Name      string    `json:"name"`
		Label     string    `json:"label"`
		Version   string    `json:"version"`
		Sequence  int       `json:"sequence"`
		Status    int       `json:"status"`
		Url       string    `json:"url"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	}

	tableExists := database.Connector.HasTable(&Chaincode_lists{})
	// Create table for chaincode_creation
	if !tableExists {
		if err := database.Connector.CreateTable(Chaincode_lists{}).Error; err != nil {
			log.Println("error creating table ", err)
			response := utils.Response{
				Status:  500,
				Message: "orgnizations table does not created successfully",
				Err:     err,
			}

			return response
		}
	}

	// var organizations Organizations
	chaincode_lists := &Chaincode_lists{Name: chaincode_updates.Name, Label: chaincode_updates.Label, Version: chaincode_updates.Version, Sequence: chaincode_updates.Sequence, Status: chaincode_updates.Status, Url: chaincode_updates.Url, CreatedAt: time.Now(), UpdatedAt: time.Now()}

	insert := database.Connector.Create(&chaincode_lists)
	if insert == nil {
		log.Println("error inserting ")
		response := utils.Response{
			Status:  500,
			Message: "chaincode list data did not successfully add in db",
			Err:     err,
		}
		return response

	}
	response := utils.Response{
		Status:  200,
		Message: "chaincode successfully commited",
		Data:    resp,
	}
	return response
}

type OrgInfo struct {
	Id int `json:"id"`
}

func InstallChaincode(body *gin.Context) utils.Response {
	log.Println("body is ", body)
	var orgId OrgInfo

	if err := body.BindJSON(&orgId); err != nil {
		response := utils.Response{
			Status:  0,
			Message: "payload error",
			// Data: ,
			Err: err,
		}
		return response

	} else {
		log.Println("payload from backend  @@@@@@@@@@@@@@@", orgId)
		fmt.Println("req.body = ", reflect.TypeOf(orgId))
	}

	//struct for sending chaincode info as response
	type ChaincodeList struct {
		// Id       int    `json:"id"`
		Name      string    `json:"name"`
		Label     string    `json:"label"`
		Version   string    `json:"version"`
		Sequence  int       `json:"sequence"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	//fetching organization installed chaincode from chaincode_logs
	var chaincode_logs ChaincodeLog
	if err := database.Connector.Table("chaincode_logs").Last(&chaincode_logs, "org_id = ? AND status = ?", orgId.Id, 1).Error; err != nil {
		response := utils.Response{
			Status:  500,
			Message: "No chaincode update found",
			// Data: ,
			Err: err,
		}
		return response

	}

	// chaincode insatlled on user/org
	chaincodeinfo := &ChaincodeList{
		Name:      chaincode_logs.Name,
		Label:     chaincode_logs.Label,
		Version:   chaincode_logs.Version,
		Sequence:  chaincode_logs.Sequence,
		CreatedAt: chaincode_logs.CreatedAt,
		UpdatedAt: chaincode_logs.UpdatedAt,
	}

	var data interface {
	} = chaincodeinfo

	response := utils.Response{
		Status:  1,
		Message: "chaincode successfully found",
		Data:    data,
		// Err:     nil,
	}
	return response
}

// @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ CHAINCODE update blOCK @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

func CcUpdateCheck(body *gin.Context) utils.Response {

	type ChaincodeUpdates struct {
		Id        int    `json:"id"`
		CC_ID     int    `json:"cc_id"`
		Name      string `json:"name"`
		Label     string `json:"label"`
		Version   string `json:"version"`
		Sequence  int    `json:"sequence"`
		Url       string `json:"url"`
		Status    int    `json:"status"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at" gorm:"constraint:OnUpdate:CASCADE"`
	}

	//fetching payload from calling api from client side
	type GetData struct {
		OrgId int `json:"org_id"`
	}

	var getData GetData
	if err := body.BindJSON(&getData); err != nil {
		response := utils.Response{
			Status:  422,
			Message: "payload error",
			Err:     err,
		}
		return response

	} else {
		log.Println("payload from backend  @@@@@@@@@@@@@@@", getData)

	}

	//fetching all chaincodes installed on particular client
	var chaincode_logs []ChaincodeLog
	if err := database.Connector.Find(&chaincode_logs, "org_id = ? ", getData.OrgId).Error; err != nil {
		response := utils.Response{
			Status:  500,
			Message: "No chaincode update  found or given chaincode_id",
			// Data:    nil,
			Err: err,
		}
		return response
	}

	//comparing each chaincode id intsalled received from logs for client with update available ids
	var chaincode_updates []ChaincodeUpdates

	cuId := []int{}

	for _, chaincodes := range chaincode_logs {
		cuId = append(cuId, chaincodes.CuId)
		log.Println("array", cuId)

	}

	//db.Where("name IN (?)", []string{"jinzhu", "jinzhu 2"}).Find(&users)
	if err := database.Connector.Not(cuId).Find(&chaincode_updates).Error; err != nil {
		// if err := database.Connector.Where("cc_id IN (?)", cuId).Find(&chaincode_updates).Error; err != nil {
		response := utils.Response{
			Status:  404,
			Message: "No chaincode update  found or given client",
			Err:     err,
		}
		return response
	}
	response := utils.Response{
		Status:  200,
		Message: "chaincode updates found",
		Data:    chaincode_updates,
	}
	return response
}

//
// get all releases list Viewupdatelogs
func GetAllCCReleases(body *gin.Context) utils.Response {
	log.Println("data in admin controller of GetAllReleases data: ", body)

	type ChaincodeUpdates struct {
		Id        int       `json:"id" gorm:"primary key;autoincrement"`
		CC_ID     int       `json:"cc_id"`
		Name      string    `json:"name"`
		Label     string    `json:"label"`
		Version   string    `json:"version"`
		Sequence  int       `json:"sequence"`
		Status    int       `json:"status"`
		Url       string    `json:"url"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	}

	// release array of struct
	var releases []ChaincodeUpdates
	if err := database.Connector.Table("chaincode_updates").Find(&releases).Error; err != nil {
		response := utils.Response{
			Status:  500,
			Message: "releases data not found",
			// Data: ,
			Err: err,
		}
		return response
	}

	log.Println("releases data in admin controller: ", releases)

	response := utils.Response{
		Status:  200,
		Message: "releases data found",
		Data:    releases,
	}
	return response

}

// releases logs
func Viewupdatelogs(body *gin.Context, id string) utils.Response {
	log.Println("data in admin controller of Viewupdatelogs data: ", body)
	//converting user id into int
	cuId, err := strconv.Atoi(id)
	if err != nil {
		response := utils.Response{Status: 500, Message: "error in enable fx converting org id to int", Err: err}
		return response
	}

	log.Println("ccID in view logs fx starting :", cuId)

	type ChaincodeUpdate struct {
		Status int `json:"status"`
	}

	var chaincode_updates ChaincodeUpdate

	if err := database.Connector.Table("chaincode_updates").Find(&chaincode_updates, "id = ?", cuId).Error; err != nil {

		response := utils.Response{
			Status:  500,
			Message: "chaincode data not found",
			// Data: ,
			Err: err,
		}
		return response

	}
	log.Println("status before : ", chaincode_updates)
	status := chaincode_updates.Status
	log.Println("status after: ", status)

	type ChaincodeLog struct {
		Id        int    `json:"id" gorm:"primary key;autoincrement"`
		OrgId     int    `json:"org_id"`
		OrgName   string `json:"org_name"`
		CreatedAt string `json:"created_at"`
		// CommitStatus   string    `json:"commit_status"`
		// UpdatedAt string `json:"updated_at" gorm:"autoUpdateTime"`
	}

	var chaincode_logs = []ChaincodeLog{}
	if err := database.Connector.Table("chaincode_logs").Find(&chaincode_logs, "cu_id = ?", cuId).Error; err != nil {

		response := utils.Response{
			Status:  500,
			Message: "chaincode data not found",
			// Data: ,
			Err: err,
		}
		return response

	}

	ccLogsLength := len(chaincode_logs)
	var ccLogsLengthData float64 = float64(ccLogsLength)
	log.Println("ccLogsLengthData length is :", ccLogsLengthData)

	type Organizations struct {
		Id             int    `gorm:"primaryKey;autoIncrement"`
		Name           string `json:"name"`
		MspId          string `json:"msp_id"`
		PeersCount     int    `json:"peers_count"`
		Config         string `json:"file" gorm:"type:text"`
		ModifiedConfig string `json:"modified_config" gorm:"type:text"`
		Join_Status    int    `json:"join_status"`
		CreatedAt      string `json:"created_at"`
		UpdatedAt      string `json:"updated_at" gorm:"autoUpdateTime"`
	}

	var organizations []Organizations
	if err := database.Connector.Table("organizations").Find(&organizations).Error; err != nil {
		response := utils.Response{
			Status:  500,
			Message: "Organization data not found",
			// Data: ,
			Err: err,
		}
		return response

	}
	// log.Println("Organization list at admin side :", organizations)
	orgLength := len(organizations)
	log.Println("Organization length is :", orgLength)

	var orgLengthData float64 = float64(orgLength)

	var midOfOrg float64 = orgLengthData / 2

	log.Println("midOfOrg is :", midOfOrg)

	// check for signing majority
	if status == 1 {
		response := utils.Response{
			Status:       200,
			Message:      "Already committed",
			Data:         chaincode_logs,
			CommitStatus: "false",
		}
		return response
	} else if status == 0 && ccLogsLengthData > midOfOrg {
		response := utils.Response{
			Status:       200,
			Message:      "Enable to commit",
			Data:         chaincode_logs,
			CommitStatus: "true",
		}
		return response
	}
	response := utils.Response{
		Status:       200,
		Message:      "Disable to commit",
		Data:         chaincode_logs,
		CommitStatus: "false",
	}
	return response

}

////////////////////////
// releases logs
func DeleteCCRelease(body *gin.Context, id string) utils.Response {
	log.Println("data in admin controller of delete release api: ", body)
	//converting user id into int
	cuId, err := strconv.Atoi(id)
	if err != nil {
		response := utils.Response{Status: 500, Message: "error in enable fx converting org id to int", Err: err}
		return response
	}

	log.Println("ccID in view logs fx starting :", cuId)

	type ChaincodeUpdate struct {
		Id        int       `json:"id" gorm:"primary key;autoincrement"`
		CC_ID     int       `json:"cc_id"`
		Name      string    `json:"name"`
		Label     string    `json:"label"`
		Version   string    `json:"version"`
		Sequence  int       `json:"sequence"`
		Status    int       `json:"status"`
		Url       string    `json:"url"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	}

	var chaincode_updates ChaincodeUpdate
	//if err := database.Connector.Table("users").Find(&user, "id = ?", userId).Error; err != nil {
	if err := database.Connector.Table("chaincode_updates").Delete(&chaincode_updates, "id = ?", cuId).Error; err != nil {

		response := utils.Response{
			Status:  500,
			Message: "Unable to delete cc release",
			// Data: ,
			Err: err,
		}
		return response

	}
	log.Println("deleted release in delete fx : ", chaincode_updates)
	response := utils.Response{
		Status:  200,
		Message: "Release deleted successfully",
		// Data:    chaincode_updates,
	}
	return response

}

//fetch new update values and send
func CcUpdate(body *gin.Context) utils.Response {
	type GetData struct {
		ChaincodeId int `json:"cc_id"`
	}
	var getData GetData
	if err := body.BindJSON(&getData); err != nil {
		log.Println("err post route @@@@@@@@@@ ", err)
		response := utils.Response{
			Status:  422,
			Message: "payload error",
			Err:     err,
		}
		return response
	} else {
		log.Println("payload from backend  @@@@@@@@@@@@@@@", getData)
	}

	var chaincode_updates ChaincodeUpdates
	if err := database.Connector.Find(&chaincode_updates, "id = ?", getData.ChaincodeId).Error; err != nil {
		log.Println("error finding chaincode updates")
		response := utils.Response{
			Status:  404,
			Message: "No update data found",
			// Data: ,
			Err: err,
		}
		return response

	}

	// database.Connector.Table("organizations").Find(&orginfo)
	log.Println("chaincode update found with given id  ", getData.ChaincodeId)

	response := utils.Response{
		Status:  200,
		Message: "chaincode update sent",
		Data:    chaincode_updates,
		// Err:     nil,
	}
	return response

}

//@@@@@@@@@@@@@@@@@@@@@@@ ORGANIZATION FUNCTIONS @@@@@@@@@@@@@@@@@@@@@@@@@@@@
//add org function

type Organizationinfo struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	MspId      string `json:"msp_id"`
	PeersCount int    `json:"peers_count"`
	Config     string `json:"file"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// Add new organization
func AddOrg(body *gin.Context, sdk *fabsdk.FabricSDK, rsmgmnt *resmgmt.Client, clCtx contextApi.ClientProvider) utils.Response {

	type Organizations struct {
		Id             int       `json:"id" gorm:"primaryKey;autoIncrement"`
		Name           string    `json:"name"`
		MspId          string    `json:"msp_id"`
		PeersCount     int       `json:"peers_count"`
		Config         string    `json:"file" gorm:"type:text"`
		ModifiedConfig string    `json:"modified_config" gorm:"type:text"`
		Join_Status    int       `json:"join_status"`
		CreatedAt      time.Time `json:"created_at"`
		UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	}
	log.Println("enter add org function")

	name := body.Request.PostFormValue("name")
	log.Println("org name ", name)
	msp_id := body.Request.PostFormValue("msp_id")
	peers_count := body.Request.PostFormValue("peers_count")
	var count_peers int
	if i, err := strconv.Atoi(peers_count); err == nil {
		fmt.Printf("i=%d, type: %T\n", i, i)
		count_peers = i
	}

	join_status := 0

	file, handler, err := body.Request.FormFile("file")
	if err != nil {

		log.Println("err in uploading", err, file)
		response := utils.Response{
			Status:  500,
			Message: "error getting config file",
			Err:     err,
		}
		return response

	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		response := utils.Response{
			Status:  500,
			Message: "error reading org config file",
			Err:     err,
		}
		return response

	}

	newOrgJsonPath := "public/jsonFile/" + handler.Filename
	err = ioutil.WriteFile(newOrgJsonPath, fileBytes, 0777)
	if err != nil {
		log.Println("errrrrr ", err)
		response := utils.Response{
			Status:  500,
			Message: "error writing file to folder",
			Err:     err,
		}
		return response

	}

	fmt.Print("Successfully Uploaded File\n")

	tableExists := database.Connector.HasTable(&Organizations{})

	// Create table for chaincode_creation
	if !tableExists {
		if err := database.Connector.CreateTable(Organizations{}).Error; err != nil {
			log.Println("error creating table ", err)
			response := utils.Response{
				Status:  500,
				Message: "orgnizations table does not created successfully",
				Err:     err,
			}

			return response
		}
	}

	// var organizations Organizations
	organizations := &Organizations{Name: name, MspId: msp_id, PeersCount: count_peers, Config: string(fileBytes), Join_Status: join_status, CreatedAt: time.Now(), UpdatedAt: time.Now()}

	insert := database.Connector.Create(&organizations)
	if insert == nil {
		log.Println("error inserting ")
		response := utils.Response{
			Status:  500,
			Message: "Org data did not successfully add in db",
			Err:     err,
		}
		return response

	}

	type OrgData struct {
		Id    int    `json:"id" gorm:"primaryKey;autoIncrement"`
		Name  string `json:"name"`
		MspId string `json:"msp_id"`
	}

	orgInfo := OrgData{}

	// user.UserName = data.UserName
	orgInfo.Id = organizations.Id
	log.Println("orgniazation id :", orgInfo.Id)
	orgInfo.Name = organizations.Name
	log.Println("orgniazation name :", orgInfo.Name)
	orgInfo.MspId = organizations.MspId
	log.Println("orgniazation mspid :", orgInfo.MspId)

	_, err = GetConfigBlock(organizations.Id, name, msp_id, rsmgmnt, sdk, clCtx, newOrgJsonPath)
	if err != nil {
		log.Println("error getting config block", err)
		response := utils.Response{
			Status:  500,
			Message: "error getting config block",
			// Data:  ,
			Err: err,
		}
		return response
	}

	// create envelope file from modified config
	configUpdatedJson, err := os.ReadFile("public/jsonFile/config_update_new.json")
	if err != nil {
		response := utils.Response{
			Status:  500,
			Message: "error getting config file",
			Err:     err,
		}
		return response
	}

	log.Println("######### 2")
	ModifiedConfig := &common.Config{}

	err = protolator.DeepUnmarshalJSON(bytes.NewReader(configUpdatedJson), ModifiedConfig)
	if err != nil {
		response := utils.Response{
			Status:  500,
			Message: "error getting config file",
			Err:     err,
		}
		return response
	}
	log.Println("######### 3")
	var bufModifiedConfig bytes.Buffer
	if err := protolator.DeepMarshalJSON(&bufModifiedConfig, ModifiedConfig); err != nil {
		log.Fatalf("DeepMarshalJSON returned error: %s", err)
	}

	log.Println("channeeeeeeeeeellllllllllllllll valueeeeeeeeeeeeeeeee :", config.CHANNEL_ID)
	configUpdate, err := getConfigUpdate(rsmgmnt, config.CHANNEL_ID, string(configUpdatedJson), config.ORDERER_ENDPOINT)
	if err != nil {
		log.Fatalf("getConfigUpdate returned error: %s", err)
	}

	var bufConfigUpdate bytes.Buffer
	if err := protolator.DeepMarshalJSON(&bufConfigUpdate, configUpdate); err != nil {
		log.Fatalf("DeepMarshalJSON returned error: %s", err)
	}

	log.Println("=========================== Im here")

	// err = ioutil.WriteFile("test-data/ConfigUpdate.json", bufConfigUpdate.Bytes(), 0777)
	// if err != nil {
	// 	panic(err)
	// }

	configEnvelopeBytes, err := getConfigEnvelopeBytes(configUpdate)
	if err != nil {
		response := utils.Response{
			Status:  500,
			Message: "error getting configEnvelopeBytes file",
			Err:     err,
		}
		return response
	}
	log.Println("configEnvelopeBytes: ", configEnvelopeBytes)

	newOrgEnvelopePBFile := "./new-org-envelope.sh"

	var envelopFileJson = filepath.Join("config-envelope/envelopConfig.json")

	// var envelopFilePB = filepath.Join("join-channel-files/config-envelope.pb")
	var envelopFilePB = filepath.Join("config-envelope/config-envelope-final.pb")

	cmd := exec.Command("/bin/sh", newOrgEnvelopePBFile, envelopFileJson, envelopFilePB)

	stdout, err := cmd.Output()
	if err != nil {
		log.Println("inside error", err)
		response := utils.Response{
			Status:  500,
			Message: "error converting json to Pb file",
			Err:     err,
		}
		return response
	}

	result := string(stdout)

	log.Println("string in block file", result)

	response := utils.Response{
		Status:  200,
		Message: "new organization configuration successfully updated",
		Data:    orgInfo,
	}

	return response
}

func AddPeers(body *gin.Context) utils.Response {

	type PeersAdd struct {
		OrgId   int    `json:"org_id"`
		OrgName string `json:"org_name"`
		Peers   []struct {
			Name string `json:"peer_name"`
			Url  string `json:"peer_url"`
			Ip   string `json:"peer_ip"`
			Cert string `json:"peer_cert"`
		} `json:"peers"`
	}

	var getData PeersAdd
	if err := body.BindJSON(&getData); err != nil {
		log.Println("err post route @@@@@@@@@@ ", err)
		response := utils.Response{Status: 422, Message: "payload error", Err: err}
		return response

	} else {
		log.Println("payload from backend  @@@@@@@@@@@@@@@", getData)
	}

	type Peers struct {
		Id        int       `gorm:"primaryKey;autoIncrement"`
		Name      string    `json:"peer_name,omitempty"`
		OrgId     int       `json:"org_id,omitempty"`
		OrgName   string    `json:"org_name,omitempty"`
		Url       string    `json:"peer_url,omitempty"`
		Ip        string    `json:"peer_ip,omitempty"`
		Cert      string    `json:"peer_cert,omitempty"`
		CreatedAt time.Time `json:"created_at,omitempty"`
		UpdatedAt time.Time `json:"updated_at,omitempty" gorm:"autoUpdateTime"`
	}

	tableExists := database.Connector.HasTable(&Peers{})
	//table creationn for peers

	if !tableExists {
		if err := database.Connector.CreateTable(Peers{}).Error; err != nil {
			log.Println("error creating table ", err)
			response := utils.Response{
				Status:  500,
				Message: "peers table does not created successfully",
				// Data:    "",
				Err: err,
			}

			return response
		}
	}

	//to retunr array in response
	var peersData []Peers

	var peers *Peers

	for _, value := range getData.Peers {
		peers = &Peers{Name: value.Name, OrgId: getData.OrgId, OrgName: getData.OrgName, Url: value.Url, Ip: value.Ip, Cert: value.Cert, CreatedAt: time.Now(), UpdatedAt: time.Now()}
		peersData = append(peersData, *peers)
		if err := database.Connector.Create(peers).Error; err != nil {
			response := utils.Response{
				Status:  500,
				Message: "Peers data not succesfully entered",
				Err:     err,
			}
			return response
		}

	}

	response := utils.Response{
		Status:  200,
		Message: "peers data successfully insert",
		Data:    peersData,
	}

	return response

}

//get Organizations list except the login one
func GetOrgs(body *gin.Context) utils.Response {

	type Organizations struct {
		Id             int    `gorm:"primaryKey;autoIncrement"`
		Name           string `json:"name"`
		MspId          string `json:"msp_id"`
		PeersCount     int    `json:"peers_count"`
		Config         string `json:"file" gorm:"type:text"`
		ModifiedConfig string `json:"modified_config" gorm:"type:text"`
		Join_Status    int    `json:"join_status"`
		CreatedAt      string `json:"created_at"`
		UpdatedAt      string `json:"updated_at" gorm:"autoUpdateTime"`
	}

	type OrgId struct {
		OrgId int `json:"org_id"`
	}

	var getData OrgId
	if err := body.BindJSON(&getData); err != nil {
		log.Println("err post route @@@@@@@@@@ ", err)
		response := utils.Response{Status: 422, Message: "payload error", Err: err}
		return response

	} else {
		log.Println("payload from backend  @@@@@@@@@@@@@@@", getData)
	}

	var organizations []Organizations
	if err := database.Connector.Find(&organizations, "id <> ?", getData.OrgId).Error; err != nil {
		log.Println("error in organizations querying ", err)
		response := utils.Response{Status: 404, Message: "organizations not found", Err: err}
		return response

	}

	type OrgSignatures struct {
		Id         int    `json:"id"`
		OrgId      int    `json:"org_id"`
		OrgMsp     string `json:"org_msp"`
		SignbyId   int    `json:"signby_id"`
		SignbyName string `json:"signby_name"`
		Signature  string `json:"signatures"`
		CreatedAt  time.Time
		UpdatedAt  time.Time
	}

	var org_signatures OrgSignatures

	type OrgSignData struct {
		// Id         int    `json:"id"`
		OrgId       int         `json:"org_id"`
		OrgName     string      `json:"org_name"`
		OrgMsp      string      `json:"org_msp"`
		SignedbyOrg interface{} `json:"signedby_org"`
		Join_Status int         `json:"join_status"`
		CreatedAt   string      `json:"created_at"`
	}
	var orgSignData []OrgSignData
	// log.Println("join status in org list :", value.Join_Status)
	for _, value := range organizations {
		log.Println("join status in org list :", value.Join_Status)
		if value.ModifiedConfig == "" {
			orgSignData = append(orgSignData, OrgSignData{OrgId: value.Id, OrgName: value.Name, OrgMsp: value.MspId, SignedbyOrg: nil, Join_Status: value.Join_Status, CreatedAt: value.CreatedAt})
		} else {
			if err := database.Connector.Find(&org_signatures, "signby_id= ? AND  org_id = ?", getData.OrgId, value.Id).Error; err != nil {

				log.Println("error in querying ", err)
				orgSignData = append(orgSignData, OrgSignData{OrgId: value.Id, OrgName: value.Name, OrgMsp: value.MspId, SignedbyOrg: false, Join_Status: value.Join_Status, CreatedAt: value.CreatedAt})

			} else {
				// log.Println("value of true org ", value)
				orgSignData = append(orgSignData, OrgSignData{OrgId: value.Id, OrgName: value.Name, OrgMsp: value.MspId, SignedbyOrg: true, Join_Status: value.Join_Status, CreatedAt: value.CreatedAt})

			}
		}

	}

	response := utils.Response{
		Status:  200,
		Message: "organizations found successfully",
		Data:    orgSignData,
	}
	log.Println("response", response)
	return response
}

//get Organizations list except the login one
func GetSingleOrgs(body *gin.Context) utils.Response {

	type Organizations struct {
		Id             int    `gorm:"primaryKey;autoIncrement"`
		Name           string `json:"name"`
		MspId          string `json:"msp_id"`
		PeersCount     int    `json:"peers_count"`
		Config         string `json:"file" gorm:"type:text"`
		ModifiedConfig string `json:"modified_config" gorm:"type:text"`
		Join_Status    int    `json:"join_status"`
		EnvelopeUrl    string `json:"envelope_url"`
		CreatedAt      string `json:"created_at"`
		UpdatedAt      string `json:"updated_at" gorm:"autoUpdateTime"`
	}

	type OrgId struct {
		OrgId int `json:"org_id"`
	}

	var getData OrgId
	if err := body.BindJSON(&getData); err != nil {
		log.Println("err post route @@@@@@@@@@ ", err)
		response := utils.Response{Status: 422, Message: "payload error", Err: err}
		return response

	} else {
		log.Println("payload from backend  @@@@@@@@@@@@@@@", getData)
	}

	var organizations Organizations
	if err := database.Connector.Table("organizations").Find(&organizations, "id = ?", getData.OrgId).Error; err != nil {
		response := utils.Response{
			Status:  500,
			Message: "Org data not found",
			// Data: ,
			Err: err,
		}
		return response

	}
	response := utils.Response{
		Status:  200,
		Message: "organizations found successfully",
		Data:    organizations,
	}
	log.Println("response", response)
	return response
}

// get org list at admin side
//get Organizations list except the login one
func GetAdminOrgs(body *gin.Context) utils.Response {
	// log.Println("user token is ", userToken)
	log.Println("org id ", orgId)
	log.Println("channel id is :", config.CHANNEL_ID)
	type Organizations struct {
		Id             int    `gorm:"primaryKey;autoIncrement"`
		Name           string `json:"name"`
		MspId          string `json:"msp_id"`
		PeersCount     int    `json:"peers_count"`
		Config         string `json:"file" gorm:"type:text"`
		ModifiedConfig string `json:"modified_config" gorm:"type:text"`
		Join_Status    int    `json:"join_status"`
		CreatedAt      string `json:"created_at"`
		UpdatedAt      string `json:"updated_at" gorm:"autoUpdateTime"`
	}

	var organizations []Organizations
	if err := database.Connector.Find(&organizations, "id <> ?", orgId).Error; err != nil {
		log.Println("error in organizations querying ", err)
		response := utils.Response{Status: 404, Message: "organizations not found", Err: err}
		return response

	}

	response := utils.Response{
		Status:  200,
		Message: "organizations found successfully",
		Data:    organizations,
	}
	return response
}

func GetConfigBlock(orgid int, orgName string, orgMsp string, resmgmtClient *resmgmt.Client, sdk *fabsdk.FabricSDK, clCtx contextApi.ClientProvider, newOrgJsonPath string) (utils.Response, error) {
	log.Println("2ndddddddddddd  channeeeeeeeeeelllllllll valueeeeeeeeeeeeeeeee :", config.CHANNEL_ID)
	channelConfig, _ := getCurrentChannelConfig(resmgmtClient, config.CHANNEL_ID, config.ORDERER_ENDPOINT)

	var buf bytes.Buffer
	if err := protolator.DeepMarshalJSON(&buf, channelConfig); err != nil {
		log.Fatalf("DeepMarshalJSON returned error: %s", err)
	}

	originalChConfigJSON := buf.String()

	// write the whole body at once
	err := ioutil.WriteFile("config_org.json", []byte(originalChConfigJSON), 0777)
	if err != nil {
		response := utils.Response{
			Status:  500,
			Message: "failed writing config json file",
			Err:     err,
		}
		return response, err
	}

	// Prepare the OrgMSP
	log.Printf("Preparing the org MSP %s", config.ORG_MSP)

	type Organizations struct {
		Id             int       `json:"id"`
		Name           string    `json:"name"`
		MspId          string    `json:"msp_id"`
		PeersCount     int       `json:"peers_count"`
		Config         string    `json:"config" gorm:"type:text"`
		ModifiedConfig string    `json:"modified_config" gorm:"type:text"`
		Join_Status    int       `json:"join_status"`
		CreatedAt      time.Time `json:"created_at"`
		UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	}
	//fetching new config from database
	var organizations Organizations
	// database.Connector.Table("organizations").Find(&organizations)
	if err := database.Connector.Find(&organizations, "id = ?", orgid).Error; err != nil {
		log.Println("err founding data", err)
		response := utils.Response{
			Status:  404,
			Message: "Org data not found",
			Err:     err,
		}
		return response, err

	}
	orgConfig := organizations.Config
	// orgId := orgid

	// Prepare the OrgMSP
	log.Printf("org MSP %s", orgMsp)
	//Merge jsons using jq
	modfifiedConfig, err := MergeJson(orgMsp, originalChConfigJSON, orgConfig, newOrgJsonPath)

	//Save
	if err != nil {
		response := utils.Response{
			Status:  500,
			Message: "Failed merging config files",
			Err:     err,
		}
		return response, err
	}

	var newConfigPath = filepath.Join("public/jsonFile/config_update_new.json")
	// // write the whole body at once
	err = ioutil.WriteFile(newConfigPath, []byte(modfifiedConfig), 0777)
	if err != nil {
		response := utils.Response{
			Status:  500,
			Message: "failed writing new config file",
			Err:     err,
		}
		return response, err
	}

	if err := database.Connector.Model(&organizations).Updates(Organizations{ModifiedConfig: modfifiedConfig}).Error; err != nil {
		log.Println("error updating modified config in org table")

	} else {

		log.Println("organization table updated")
	}

	response := utils.Response{
		Status:  200,
		Message: "config data successfully created",
		Data:    true,
	}

	return response, nil

}

func MergeJson(orgMspName string, chConfigJson string, orgJson string, newOrgJsonPath string) (string, error) {

	log.Println("new org json path", newOrgJsonPath)

	app := "./jqshell.sh"

	var chConfigJsonFile = filepath.Join("config_org.json")
	log.Println("chConfigJson file jason path", chConfigJsonFile)

	log.Println("chaconfig path ", chConfigJsonFile)

	log.Println("orgjson path ", newOrgJsonPath)

	err := ioutil.WriteFile(chConfigJsonFile, []byte(chConfigJson), 0777)
	if err != nil {
		panic(err)
	}

	var outputJsonPath = filepath.Join("public/jsonFile/output.json")

	arg0 := orgMspName       //"Org3MSP"
	arg1 := chConfigJsonFile //"config_org3.json"
	arg2 := newOrgJsonPath   //"org3.json"
	arg3 := outputJsonPath

	cmd := exec.Command(app, arg0, arg1, arg2, arg3)

	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println("Error ------------- % ", err)
		return "failed", err
	}

	var updatedConfig = string(stdout)

	fmt.Println("CreateModifiedJson End")

	return updatedConfig, nil

}

// / getCurrentChannelConfig Get the current channel config
func getCurrentChannelConfig(resmgmtClient *resmgmt.Client, channelID, ordererEndPoint string) (*common.Config, error) {
	log.Println("configggggggg varibale value :", config.CHANNEL_ID)
	block, err := resmgmtClient.QueryConfigBlockFromOrderer(config.CHANNEL_ID, resmgmt.WithOrdererEndpoint(config.ORDERER_ENDPOINT))
	if err != nil {
		log.Println(" getCurrentChannelConfig error ", err.Error())
		return nil, err
	}

	return resource.ExtractConfigFromBlock(block)
}

//controller function to sign newly added org
func MConfig(body *gin.Context) utils.Response {

	type Organizations struct {
		Id             int    `gorm:"primaryKey;autoIncrement"`
		Name           string `json:"name"`
		MspId          string `json:"msp_id"`
		PeersCount     int    `json:"peers_count"`
		Config         string `json:"file" gorm:"type:text"`
		ModifiedConfig string `json:"modified_config" gorm:"type:text"`
		Join_Status    int    `json:"join_status"`
		CreatedAt      string `json:"created_at"`
		UpdatedAt      string `json:"updated_at" gorm:"autoUpdateTime"`
	}
	type OrgData struct {
		OrgId int `json:"org_id"`
	}

	var data OrgData
	log.Println("entering Install chaincode controller", body)
	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := body.BindJSON(&data); err != nil {
		log.Println("err post route @@@@@@@@@@ ", err)
		response := utils.Response{
			Status:  422,
			Message: "payload error",
			Err:     err,
		}
		return response

	} else {
		log.Println("payload in signOrg  @@@@@@@@@@@@@@@", data)
	}

	var organizations Organizations
	if err := database.Connector.Find(&organizations, "id = ?", data.OrgId).Error; err != nil {
		log.Println("err founding data", err)
		response := utils.Response{
			Status:  404,
			Message: "Failed to get modified config",
			Err:     err,
		}
		return response

	}

	//fetched modified config after adding new org to configuration
	modfifiedConfigString := organizations.ModifiedConfig

	type ModifiedConfig struct {
		Config string `json:"modified_config" gorm:"type:text"`
	}

	modfifiedConfig := &ModifiedConfig{Config: modfifiedConfigString}
	response := utils.Response{
		Status:  200,
		Message: "Modified config found",
		Data:    modfifiedConfig.Config,
		// Err: err,
	}
	return response

}

//create chaincode by admin
func SaveSign(body *gin.Context) utils.Response {

	//fetching data from payload or calling api
	type SignatureData struct {
		OrgId      int    `json:"org_id"`
		SigningOrg int    `json:"signingorg_id"`
		Signatures string `json:"signatures"`
	}
	var signData SignatureData

	if err := body.BindJSON(&signData); err != nil {
		log.Println("err post route @@@@@@@@@@ ", err)
		response := utils.Response{
			Status:  422,
			Message: "payload error",
			Err:     err,
		}
		return response

	} else {
		log.Println("payload from backend  @@@@@@@@@@@@@@@", signData.OrgId)
	}

	type Organizations struct {
		Id         int       `gorm:"primaryKey;autoIncrement"`
		Name       string    `json:"name"`
		MspId      string    `json:"msp_id"`
		PeersCount int       `json:"peers_count"`
		Config     string    `json:"file" gorm:"type:text"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	}

	var organizations Organizations

	log.Println("signing org org to be signed data is ", signData.OrgId)

	if err := database.Connector.First(&organizations, "id = ?", signData.OrgId).Error; err != nil {
		log.Println("err", err)
		response := utils.Response{
			Status:  500,
			Message: "error finding organization data in organization",
			Err:     err,
		}
		return response

	}
	log.Println("org msp is ", organizations)

	type OrgSignature struct {
		Id int `gorm:"primaryKey;autoIncrement"`
		// ChaincodeId int       `json:"chaincode_id"`
		OrgId     int    `json:"org_id"`
		OrgName   string `json:"org_name"`
		OrgMsp    string `json:"org_msp"`
		SignbyId  int    `json:"signby_id"`
		Signature string `json:"signature" gorm:"type:text"`

		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	}

	var org_signatures OrgSignature

	tableExists := database.Connector.HasTable(&OrgSignature{})

	org_signatures = OrgSignature{OrgId: signData.OrgId, OrgName: organizations.Name, OrgMsp: organizations.MspId, SignbyId: signData.SigningOrg, Signature: signData.Signatures, CreatedAt: time.Now(), UpdatedAt: time.Now()}

	// Create table for chaincode_creation
	if !tableExists {
		if err := database.Connector.CreateTable(OrgSignature{}).Error; err != nil {
			log.Println("error creating table ", err)
			response := utils.Response{
				Status:  500,
				Message: "org_signatures table not created successfully",
				Err:     err,
			}

			return response
		}
	}

	// entries fornewly added signed org for table org_signatures
	if err := database.Connector.Create(&org_signatures).Error; err != nil {
		log.Println("error in saving data ", err)
		response := utils.Response{
			Status:  500,
			Message: "new orga sign data  not saved successfully",
			Err:     err,
		}

		return response
	}

	response := utils.Response{
		Status:  200,
		Message: "new org sign data saved successfully",
		Data:    org_signatures,
	}
	return response

}

//======================================================================

// getConfigUpdate Get the config update from two configs
func getConfigUpdate(resmgmtClient *resmgmt.Client, channelID string, proposedConfigJSON string, ordererEndPoint string) (*common.ConfigUpdate, error) {

	proposedConfig := &common.Config{}

	err := protolator.DeepUnmarshalJSON(bytes.NewReader([]byte(proposedConfigJSON)), proposedConfig)
	if err != nil {
		return nil, err
	}
	channelConfig, err := getCurrentChannelConfig(resmgmtClient, config.CHANNEL_ID, config.ORDERER_ENDPOINT)
	if err != nil {
		return nil, err
	}
	configUpdate, err := resmgmt.CalculateConfigUpdate(config.CHANNEL_ID, channelConfig, proposedConfig)
	if err != nil {
		return nil, err
	}
	configUpdate.ChannelId = config.CHANNEL_ID

	return configUpdate, nil
}

// join channel to org3
func SaveChannelConfig(resmgmtClient *resmgmt.Client, org_id string) utils.Response {
	log.Println("i am in save channel channel config")
	//converting user id into int
	id, err := strconv.Atoi(org_id)
	if err != nil {
		response := utils.Response{Status: 500, Message: "error in enable fx converting org id to int", Err: err}
		return response
	}

	log.Println("org ID in create updtaes fx starting :", id)

	type ConfigSignature struct {
		SignatureHeader      []byte   `protobuf:"bytes,1,opt,name=signature_header,json=signatureHeader,proto3" json:"signature_header,omitempty"`
		Signature            []byte   `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
		XXX_NoUnkeyedLiteral struct{} `json:"-"`
		XXX_unrecognized     []byte   `json:"-"`
		XXX_sizecache        int32    `json:"-"`
	}

	type Signature struct {
		Signature *common.ConfigSignature
	}
	// struct type variable
	var signatures []*common.ConfigSignature
	// if err := database.Connector.Table("org_signatures").Select("signature").Find(&signatures).Error; err != nil {
	if err := database.Connector.Table("org_signatures").Select("signature").Find(&signatures, "org_id = ?", id).Error; err != nil {
		response := utils.Response{
			Status:  500,
			Message: "signature 1 not found",
			// Data: ,
			Err: err,
		}
		return response

	}

	var signaturesData = []*common.ConfigSignature{}

	for _, value := range signatures {
		dataBytes := value.Signature
		log.Println("dataBytes is", dataBytes)

		signOrgData := &common.ConfigSignature{}

		err = protolator.DeepUnmarshalJSON(bytes.NewReader(dataBytes), signOrgData)
		if err != nil {
			log.Println("error in signOrg2Config protolator: ", err)
		}

		signaturesData = append(signaturesData, signOrgData)
	}
	log.Println("signaturesData is :", signaturesData)

	req := resmgmt.SaveChannelRequest{ChannelID: config.CHANNEL_ID, ChannelConfigPath: "/home/riyaz/projects/integra/integra-admin-backend/integra-nock-sdk/config-envelope/config-envelope-final.pb"}
	txID, err := resmgmtClient.SaveChannel(req, resmgmt.WithConfigSignatures(signaturesData...), resmgmt.WithOrdererEndpoint(config.ORDERER_ENDPOINT))
	if err != nil {
		response := utils.Response{
			Status:  500,
			Message: "error saving channel",
			Err:     err,
		}
		return response

	}

	if err := database.Connector.Table("organizations").Where("id = ?", id).Update("join_status", 2).Error; err != nil {
		response := utils.Response{
			Status:  500,
			Message: "orgnizations table does not updated successfully",
			Err:     err,
		}
		return response
	}

	time.Sleep(time.Second * 3)

	log.Println("#################################", txID)

	// return true
	response := utils.Response{
		Status:  200,
		Message: "successfully joined channel",
		Data:    txID,
	}
	return response
}

// getConfigEnvelopeBytes Get Envelope bytes from ConfigUpdate
func getConfigEnvelopeBytes(configUpdate *common.ConfigUpdate) ([]byte, error) {
	log.Println("Inside getConfigEnvelopeBytes fx")
	var buf bytes.Buffer
	if err := protolator.DeepMarshalJSON(&buf, configUpdate); err != nil {
		return nil, err
	}

	channelConfigBytes, err := proto.Marshal(configUpdate)
	if err != nil {
		return nil, err
	}

	var bufConfigUpdate bytes.Buffer
	if err := protolator.DeepMarshalJSON(&bufConfigUpdate, configUpdate); err != nil {
		log.Fatalf("DeepMarshalJSON returned error: %s", err)
	}

	log.Println("=========================== Im here")

	configUpdateEnvelope := &common.ConfigUpdateEnvelope{
		ConfigUpdate: channelConfigBytes,
		Signatures:   nil,
	}
	configUpdateEnvelopeBytes, err := proto.Marshal(configUpdateEnvelope)
	if err != nil {
		return nil, err
	}

	channel_header := &common.ChannelHeader{
		Type:      2,
		ChannelId: "mychannel",
		// Timestamp: timestamppb.Now(),
	}

	channelHeaderBytes, err := proto.Marshal(channel_header)
	if err != nil {
		log.Println("error marshalling channel header")
		return nil, err
	}

	header := &common.Header{
		ChannelHeader: channelHeaderBytes,
		// SignatureHeader: config_sign.SignatureHeader,
		SignatureHeader: nil,
	}

	payload := &common.Payload{
		Header: header,
		Data:   configUpdateEnvelopeBytes,
	}
	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, err
	}
	configEnvelope := &common.Envelope{
		Payload: payloadBytes,
		// Signature: config_sign.Signature,
		Signature: nil,
	}

	var bufEnvConfig bytes.Buffer
	if err := protolator.DeepMarshalJSON(&bufEnvConfig, configEnvelope); err != nil {
		log.Fatalf("DeepMarshalJSON returned error: %s", err)
	}

	log.Println("######### 4")
	time.Sleep(time.Second * 5)
	err = ioutil.WriteFile("config-envelope/envelopConfig.json", bufEnvConfig.Bytes(), 0777)
	if err != nil {
		log.Fatalf("error in wiring modified channell config %s", err)
	}

	return proto.Marshal(configEnvelope)
}

//user home dashboard
func OrgJoinStatus(body *gin.Context) utils.Response {
	log.Println("in admin controller of join Status : ", body)

	type OrgId struct {
		OrgId int `json:"org_id"`
	}

	var getData OrgId
	if err := body.BindJSON(&getData); err != nil {
		log.Println("err post route @@@@@@@@@@ ", err)
		response := utils.Response{Status: 422, Message: "payload error", Err: err}
		return response

	} else {
		log.Println("payload from backend  @@@@@@@@@@@@@@@", getData)
	}
	log.Println("org id in update fx admin controller side :", getData.OrgId)
	if err := database.Connector.Table("organizations").Where("id = ?", getData.OrgId).Update("join_status", 3).Error; err != nil {
		response := utils.Response{
			Status:  500,
			Message: "orgnizations table does not updated successfully",
			Err:     err,
		}
		return response
	}

	response := utils.Response{
		Status:  200,
		Message: "data updated successfully",
	}
	return response

}

//user home dashboard
func OrgSignStatus(body *gin.Context) utils.Response {
	log.Println("in admin controller of join Status : ", body)

	type OrgId struct {
		OrgId int `json:"org_id"`
	}

	var getData OrgId
	if err := body.BindJSON(&getData); err != nil {
		log.Println("err post route @@@@@@@@@@ ", err)
		response := utils.Response{Status: 422, Message: "payload error", Err: err}
		return response

	} else {
		log.Println("payload from backend  @@@@@@@@@@@@@@@", getData)
	}

	type OrgSignature struct {
		Id        int    `json:"id"`
		OrgId     int    `json:"org_id"`
		OrgName   string `json:"org_name"`
		OrgMsp    string `json:"org_msp"`
		SignbyId  int    `json:"signby_id"`
		CreatedAt string `json:"created_at"`
	}

	var org_signatures = []OrgSignature{}
	if err := database.Connector.Table("org_signatures").Find(&org_signatures, "org_id = ?", getData.OrgId).Error; err != nil {
		response := utils.Response{
			Status:  500,
			Message: "chaincode data not found",
			// Data: ,
			Err: err,
		}
		return response

	}

	orgSignLength := len(org_signatures)
	var orgSignLengthData float64 = float64(orgSignLength)
	log.Println("orgSignLengthData length is :", orgSignLengthData)

	type Organizations struct {
		Id             int    `gorm:"primaryKey;autoIncrement"`
		Name           string `json:"name"`
		MspId          string `json:"msp_id"`
		PeersCount     int    `json:"peers_count"`
		Config         string `json:"file" gorm:"type:text"`
		ModifiedConfig string `json:"modified_config" gorm:"type:text"`
		Join_Status    int    `json:"join_status"`
		CreatedAt      string `json:"created_at"`
		UpdatedAt      string `json:"updated_at" gorm:"autoUpdateTime"`
	}

	var organizations []Organizations
	if err := database.Connector.Find(&organizations, "id <> ?", getData.OrgId).Error; err != nil {
		// if err := database.Connector.Table("organizations").Find(&organizations).Error; err != nil {
		response := utils.Response{
			Status:  500,
			Message: "Organization data not found",
			// Data: ,
			Err: err,
		}
		return response

	}

	// log.Println("Organization list at admin side :", organizations)
	orgLength := len(organizations)
	log.Println("Organization length is :", orgLength)

	var orgLengthData float64 = float64(orgLength)

	var midOfOrg float64 = orgLengthData / 2

	log.Println("midOfOrg is :", midOfOrg)

	if orgSignLengthData > midOfOrg {
		log.Println("org id in update fx admin controller side :", getData.OrgId)
		if err := database.Connector.Table("organizations").Where("id = ?", getData.OrgId).Update("join_status", 1).Error; err != nil {
			response := utils.Response{
				Status:  500,
				Message: "orgnizations table does not updated successfully",
				Err:     err,
			}
			return response
		}
		response := utils.Response{
			Status:  200,
			Message: "data updated successfully",
			Data:    organizations,
		}
		return response
	}
	response := utils.Response{
		Status:  200,
		Message: "Requires majority of orgs to be signed",
	}
	return response

}
