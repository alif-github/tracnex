package data_scope_test

import (
	"nexsoft.co.id/nextrac2/resource_common_service"
	model2 "nexsoft.co.id/nextrac2/resource_common_service/model"
	token2 "nexsoft.co.id/nextrac2/resource_common_service/token"
	"testing"
)

type userStruct struct {
	UserID     int64
	AuthUserID int64
	ClientID   string
	Name       string
}

func getData() ([]userStruct, []userStruct) {
	assetTaggingUserMapping := []userStruct{{
		UserID:     1,
		AuthUserID: 1,
		ClientID:   "C0",
		Name:       "System",
	}, {
		UserID:     2,
		AuthUserID: 13,
		ClientID:   "C3",
		Name:       "Andi",
	}, {
		UserID:     3,
		AuthUserID: 14,
		ClientID:   "C4",
		Name:       "Siska",
	}, {
		UserID:     4,
		AuthUserID: 0,
		ClientID:   "C6",
		Name:       "ASW",
	}}

	masterDataUserMapping := []userStruct{{
		UserID:     1,
		AuthUserID: 1,
		ClientID:   "C0",
		Name:       "System",
	}, {
		UserID:     2,
		AuthUserID: 11,
		ClientID:   "C1",
		Name:       "Nexchief",
	}, {
		UserID:     3,
		AuthUserID: 12,
		ClientID:   "C2",
		Name:       "AssetTag",
	}, {
		UserID:     4,
		AuthUserID: 13,
		ClientID:   "C3",
		Name:       "Andi",
	}, {
		UserID:     5,
		AuthUserID: 0,
		ClientID:   "C5",
		Name:       "Nestle",
	}}
	return assetTaggingUserMapping, masterDataUserMapping
}

func TestCase1(testing *testing.T) {
	// "cid": "C2","user_client": "C2", "sub": "12"
	token := "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJsb2NhbGUiOiJpZC1JRCIsImNpZCI6IkMyIiwicmVzb3VyY2UiOiJtYXN0ZXIiLCJ2ZXJzaW9uIjoiMS4wLjAiLCJ1c2VyX2NsaWVudCI6IkMyIiwiZXhwIjoxOTE4NTUwMzMzLCJpYXQiOjE2MTg0NjM5MzMsImlzcyI6InRlc3QiLCJzdWIiOiIxMiJ9.fW74OprH-bOwC2gCfPKREfKfV9jI3r8HZDjAEpMv70KJn34KGFQpuQ5cprCHrTt4LdoUddIMy7otBOiqrk8V6Q"
	accessTokenModel := model2.AuthAccessTokenModel{
		ClientID:                   "C2",
		AuthenticationServerUserID: 12,
	}

	createdBy, createdClient := doTask(accessTokenModel, token, testing)
	//expected Result = Created By AssetTag, Created Client AssetTag
	if createdBy != "AssetTag" || createdClient != "AssetTag" {
		testing.Errorf("Not Expected Result At Test Case Internal Token 1")
	}
}

func TestCase2(testing *testing.T) {
	// "cid": "C2","user_client": "C2", "sub": "13"
	token := "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJsb2NhbGUiOiJpZC1JRCIsImNpZCI6IkMyIiwicmVzb3VyY2UiOiJtYXN0ZXIiLCJ2ZXJzaW9uIjoiMS4wLjAiLCJ1c2VyX2NsaWVudCI6IkMyIiwiZXhwIjoxOTE4NTUwMzMzLCJpYXQiOjE2MTg0NjM5MzMsImlzcyI6InRlc3QiLCJzdWIiOiIxMyJ9.av0stD1KP9BtiXAM3S4nOQ2rsW5l5oh7dxaNb0W2fhK50KJOLDFsXw1g2hSxTCY-OJg3jO-3O7wG0DTC49FFvQ"
	accessTokenModel := model2.AuthAccessTokenModel{
		ClientID:                   "C2",
		AuthenticationServerUserID: 13,
	}

	createdBy, createdClient := doTask(accessTokenModel, token, testing)
	//expected Result = Created By Andi, Created Client AssetTag
	if createdBy != "Andi" || createdClient != "AssetTag" {
		testing.Errorf("Not Expected Result At Test Case Internal Token 2")
	}
}

func TestCase3(testing *testing.T) {
	// "cid": "C2","user_client": "C4", "sub": "12"
	token := "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJsb2NhbGUiOiJpZC1JRCIsImNpZCI6IkMyIiwicmVzb3VyY2UiOiJtYXN0ZXIiLCJ2ZXJzaW9uIjoiMS4wLjAiLCJ1c2VyX2NsaWVudCI6IkM0IiwiZXhwIjoxOTE4NTUwMzMzLCJpYXQiOjE2MTg0NjM5MzMsImlzcyI6InRlc3QiLCJzdWIiOiIxMiJ9.NMGDcxG5BT7cPNI-gKqw5-AkANsTLptJlu0SnNzfHoIGh8mdNDVsziBK4kwAaFcRhllajllnc5O_uqZglr-qRA"
	accessTokenModel := model2.AuthAccessTokenModel{
		ClientID:                   "C2",
		AuthenticationServerUserID: 12,
	}

	createdBy, createdClient := doTask(accessTokenModel, token, testing)
	//expected Result = Created By AssetTag, Created Client Siska (Not Defined On Master Data) <- C4
	if createdBy != "AssetTag" || createdClient != "C4" {
		testing.Errorf("Not Expected Result At Test Case Internal Token 3")
	}
}

func TestCase4(testing *testing.T) {
	// "cid": "C2","user_client": "C6", "sub": "12"
	token := "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJsb2NhbGUiOiJpZC1JRCIsImNpZCI6IkMyIiwicmVzb3VyY2UiOiJtYXN0ZXIiLCJ2ZXJzaW9uIjoiMS4wLjAiLCJ1c2VyX2NsaWVudCI6IkM2IiwiZXhwIjoxOTE4NTUwMzMzLCJpYXQiOjE2MTg0NjM5MzMsImlzcyI6InRlc3QiLCJzdWIiOiIxMiJ9.EQbKgdOq4zaw-z0u4Opi3XCycabjC7Etxhv3arCs4ZBNP9Q34NXCuVA7-JwGH0Y5KKPEol4EPY_JmmE7y9QELQ"
	accessTokenModel := model2.AuthAccessTokenModel{
		ClientID:                   "C2",
		AuthenticationServerUserID: 12,
	}

	createdBy, createdClient := doTask(accessTokenModel, token, testing)
	//expected Result = Created By AssetTag, Created Client ASW (Not Defined On Master Data) <- C6
	if createdBy != "AssetTag" || createdClient != "C6" {
		testing.Errorf("Not Expected Result At Test Case Internal Token 4")
	}
}

func doTask(accessTokenModel model2.AuthAccessTokenModel, token string, testing *testing.T) (createdBy string, createdClient string) {
	_, masterData := getData()

	payload, err := token2.ValidateJWTInternal(token, "test3")
	if err.Error != nil {
		testing.Errorf("Not Expected Result At Test Case Internal Token 1")
		return
	}

	result := resource_common_service.ReadAuthTokenAndPayload(accessTokenModel, payload)
	createdClient = result.ClientID

	for i := 0; i < len(masterData); i++ {
		if masterData[i].AuthUserID == result.AuthenticationServerUserID {
			createdBy = masterData[i].Name
		}
		if masterData[i].ClientID == result.ClientID {
			createdClient = masterData[i].Name
		}
	}

	return
}
