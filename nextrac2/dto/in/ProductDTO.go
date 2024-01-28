package in

import (
	"errors"
	"fmt"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	util2 "nexsoft.co.id/nextrac2/util"
	"time"
)

type ProductRequest struct {
	AbstractDTO
	ID                 int64              `json:"id"`
	ProductID          string             `json:"product_id"`
	ProductName        string             `json:"product_name"`
	ProductDescription string             `json:"product_description"`
	ProductGroupID     int64              `json:"product_group_id"`
	ClientTypeID       int64              `json:"client_type_id"`
	IsLicense          bool               `json:"is_license"`
	LicenseVariantID   int64              `json:"license_variant_id"`
	LicenseTypeID      int64              `json:"license_type_id"`
	DeploymentMethod   string             `json:"deployment_method"`
	NoOfUser           int64              `json:"no_of_user"`
	IsUserConcurrent   bool               `json:"is_user_concurrent"`
	MaxOfflineDays     int64              `json:"max_offline_days"`
	Module1            int64              `json:"module_1"`
	Module2            int64              `json:"module_2"`
	Module3            int64              `json:"module_3"`
	Module4            int64              `json:"module_4"`
	Module5            int64              `json:"module_5"`
	Module6            int64              `json:"module_6"`
	Module7            int64              `json:"module_7"`
	Module8            int64              `json:"module_8"`
	Module9            int64              `json:"module_9"`
	Module10           int64              `json:"module_10"`
	Component          []ProductComponent `json:"component"`
	UpdatedAtStr       string             `json:"updated_at"`
	UpdatedAt          time.Time
}

type ProductComponent struct {
	ID             int64  `json:"id"`
	ComponentID    int64  `json:"component_id"`
	ComponentValue string `json:"component_value"`
	Deleted        bool   `json:"deleted"`
}

type DeploymentMethod struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

func (input *ProductRequest) ValidateInsertProduct() errorModel.ErrorModel {
	var err errorModel.ErrorModel
	if !util.IsStringEmpty(input.ProductDescription) {
		err = input.ValidateMinMaxString(input.ProductDescription, constanta.ProductDescription, 1, 200)
		if err.Error != nil {
			return err
		}
	}

	return input.validateMandatory()
}

func (input *ProductRequest) ValidateViewProduct() (err errorModel.ErrorModel) {
	fileName := input.fileNameFuncNameProduct()
	funcName := "ValidateViewProduct"

	if input.ID < 1 {
		return errorModel.GenerateUnknownDataError(fileName, funcName, constanta.Product)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *ProductRequest) ValidationDeleteProduct() (err errorModel.ErrorModel) {
	fileName := input.fileNameFuncNameProduct()
	funcName := "ValidationDeleteProduct"

	if input.ID < 1 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.ID)
	}

	if util.IsStringEmpty(input.UpdatedAtStr) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.UpdatedAt)
	}

	input.UpdatedAt, err = TimeStrToTime(input.UpdatedAtStr, constanta.UpdatedAt)
	if err.Error != nil {
		return
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *ProductRequest) ValidationUpdateProduct() (err errorModel.ErrorModel) {
	fileName := input.fileNameFuncNameProduct()
	funcName := "ValidationUpdateProduct"

	if input.ID < 1 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.ID)
	}

	if !util.IsStringEmpty(input.ProductDescription) {
		err = input.ValidateMinMaxString(input.ProductDescription, constanta.ProductDescription, 1, 200)
		if err.Error != nil {
			return err
		}
	}

	if util.IsStringEmpty(input.UpdatedAtStr) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.UpdatedAt)
	}

	input.UpdatedAt, err = TimeStrToTime(input.UpdatedAtStr, constanta.UpdatedAt)
	if err.Error != nil {
		return
	}

	return input.validateMandatory()
}

func (input *ProductRequest) validateMandatory() errorModel.ErrorModel {
	var (
		fileName = input.fileNameFuncNameProduct()
		funcName = "validateMandatory"
		err      errorModel.ErrorModel
	)

	if util.IsStringEmpty(input.ProductID) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ProductID)
	}

	err = util2.ValidateSpecialCharacter(fileName, funcName, constanta.ProductID, input.ProductID)
	if err.Error != nil {
		return err
	}

	err = input.ValidateMinMaxString(input.ProductID, constanta.ProductID, 1, 12)
	if err.Error != nil {
		return err
	}

	if util.IsStringEmpty(input.ProductName) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ProductName)
	}

	err = util2.ValidateSpecialCharacter(fileName, funcName, constanta.ProductName, input.ProductName)
	if err.Error != nil {
		return err
	}

	err = input.ValidateMinMaxString(input.ProductName, constanta.ProductName, 1, 22)
	if err.Error != nil {
		return err
	}

	if input.ProductGroupID < 1 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.ProductGroupID)
	}

	if input.ClientTypeID < 1 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.NewClientType)
	}

	if input.LicenseVariantID < 1 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.LicenseVariantID)
	}

	if input.LicenseTypeID < 1 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.LicenseTypeID)
	}

	if util.IsStringEmpty(input.DeploymentMethod) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.DeploymentMethod)
	}

	if (input.DeploymentMethod != "O") && (input.DeploymentMethod != "C") && (input.DeploymentMethod != "M") {
		return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, constanta.DeploymentMethodRegex, constanta.DeploymentMethod, "")
	}

	//--- Request high authority 15/05/2023
	//if input.NoOfUser < 1 {
	//	return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.NumberOfUser)
	//}

	//--- Request high authority 15/05/2023
	//err = util2.ValidateMinMax(input.NoOfUser, constanta.NumberOfUser, 1, 9999)
	//if err.Error != nil {
	//	return err
	//}

	//--- Request high authority 15/05/2023
	//if input.MaxOfflineDays < 1 {
	//	return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.MaxOfflineDays)
	//}

	//--- Request high authority 15/05/2023
	//err = util2.ValidateMinMax(input.MaxOfflineDays, constanta.MaxOfflineDays, 1, 365)
	//if err.Error != nil {
	//	return err
	//}

	for index, valueComponent := range input.Component {
		prepareErrorCom := fmt.Sprintf(`%s no. %d`, util2.GenerateConstantaI18n(constanta.ComponentID, constanta.IndonesianLanguage, nil), index+1)
		if valueComponent.ID == 0 {
			if valueComponent.Deleted {
				msg := fmt.Sprintf(`rule if id value equal 0, flag deleted must be false`)
				return errorModel.GenerateInvalidRequestError(fileName, funcName, errors.New(msg))
			}
		}

		if valueComponent.ComponentID < 1 {
			return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, prepareErrorCom)
		}

		if !util.IsStringEmpty(valueComponent.ComponentValue) {
			err = input.ValidateMinMaxString(valueComponent.ComponentValue, constanta.ComponentValue, 1, 100)
			if err.Error != nil {
				return err
			}
		}
	}

	return errorModel.GenerateNonErrorModel()
}

func (input ProductRequest) fileNameFuncNameProduct() (fileName string) {
	return "ProductDTO.go"
}
