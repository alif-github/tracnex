package cryptoModel

import (
	"github.com/Azure/go-autorest/autorest/date"
)

type JSONFileActivationLicenseModel struct {
	InstallationID      int64               `json:"installationId"`
	ClientID            string              `json:"clientId"`
	ProductID           string              `json:"productId"`
	LicenseVariantName  string              `json:"licenseVariantName"`
	LicenseTypeName     string              `json:"licenseTypeName"`
	DeploymentMethod    string              `json:"deploymentMethod"`
	NumberOfUser        int64               `json:"noOfUser"`
	UniqueID1           string              `json:"uniqueId1"`
	UniqueID2           string              `json:"uniqueId2"`
	ProductValidFrom    date.Date           `json:"-"`
	ProductValidFromStr string              `json:"productValidFrom"`
	ProductValidThru    date.Date           `json:"-"`
	ProductValidThruStr string              `json:"productValidThru"`
	LicenseStatus       int64               `json:"licenseStatus"`
	ModuleName1         string              `json:"moduleName1"`
	ModuleName2         string              `json:"moduleName2"`
	ModuleName3         string              `json:"moduleName3"`
	ModuleName4         string              `json:"moduleName4"`
	ModuleName5         string              `json:"moduleName5"`
	ModuleName6         string              `json:"moduleName6"`
	ModuleName7         string              `json:"moduleName7"`
	ModuleName8         string              `json:"moduleName8"`
	ModuleName9         string              `json:"moduleName9"`
	ModuleName10        string              `json:"moduleName10"`
	MaxOfflineDays      int64               `json:"maxOfflineDays"`
	IsConcurrentUser    string              `json:"concurrentUser"`
	ProductComponent    []ProductComponents `json:"component"`
}

type ProductComponents struct {
	ComponentName  string `json:"name"`
	ComponentValue string `json:"value"`
}

type SalesmanList struct {
	ID         int64  `json:"id"`
	AuthUserID int64  `json:"auth_user_id"`
	UserID     string `json:"user_id"`
	Status     string `json:"status"`
}
