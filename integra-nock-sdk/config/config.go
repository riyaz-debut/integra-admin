package config

const (

	// Common for both user and organizations

	CHANNEL_ID       = "mychannel"
	ORDERER_ENDPOINT = "orderer.example.com"
	PEER1            = "peer0.org1.example.com"
	PEER2            = "peer0.org2.example.com"

	// ========== env variables for user1 and organization 1 ================

	ORG_NAME = "Org1"
	ORG_MSP  = "Org1MSP"

	PEER_NAME = PEER1

	CA_INSTANCE = "ca.org1.example.com"
	ORG_ADMIN   = "org1admin"
	SECRET      = "org1adminpw"

	// ========== env variables for user2 and organization 2 ================

	// ORG_MSP  = "Org2MSP"
	// ORG_NAME = "Org2"

	// PEER_NAME = PEER2

	// CA_INSTANCE = "ca.org2.example.com"
	// ORG_ADMIN = "org2admin"
	// SECRET = "org2adminpw"

)
