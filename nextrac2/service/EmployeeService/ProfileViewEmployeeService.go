package EmployeeService

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
)

func (input employeeService) ViewEmployee(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.EmployeeRequest
	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateViewEmployee)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doViewEmployee(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_VIEW_MESSAGE", contextModel)
	return
}

func (input employeeService) doViewEmployee(inputStruct in.EmployeeRequest, contextModel *applicationModel.ContextModel) (result interface{}, err errorModel.ErrorModel) {
	var (
		fileName     = "ProfileViewEmployeeService.go"
		funcName     = "doViewEmployee"
		db           = serverconfig.ServerAttribute.DBConnection
		employeeOnDB repository.EmployeeModel
		scope        map[string]interface{}
	)

	mappingScopeDB := make(map[string]applicationModel.MappingScopeDB)
	mappingScopeDB[constanta.EmployeeDataScope] = applicationModel.MappingScopeDB{
		View:  "e.id",
		Count: "e.id",
	}

	createdBy := contextModel.LimitedByCreatedBy       //--- Add userID when have own permission
	scope, err = input.validateDataScope(contextModel) //--- Get scope
	if err.Error != nil {
		return
	}

	employeeOnDB, err = input.EmployeeDAO.ViewEmployee(db, repository.EmployeeModel{ID: sql.NullInt64{Int64: inputStruct.ID}}, createdBy, scope, mappingScopeDB)
	if err.Error != nil {
		return
	}

	if employeeOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.EmployeeConstanta)
		return
	}

	err = input.CheckUserLimitedByOwnAccess(contextModel, employeeOnDB.CreatedBy.Int64)
	if err.Error != nil {
		return
	}

	var (
		memberTemp   in.MemberList
		isAllMember  bool
		memberID     []int64
		resultMember []interface{}
	)

	if employeeOnDB.Member.String != "" {
		errorS := json.Unmarshal([]byte(employeeOnDB.Member.String), &memberTemp)
		if errorS != nil {
			err = errorModel.GenerateUnknownError(fileName, funcName, errorS)
			return
		}
	}

	if len(memberTemp.MemberID) > 0 {
		for _, itemMemberTemp := range memberTemp.MemberID {
			if itemMemberTemp == "all" {
				isAllMember = true
				break
			}

			num, errorS := strconv.Atoi(itemMemberTemp)
			if errorS != nil {
				err = errorModel.GenerateUnknownError(fileName, funcName, errorS)
				return
			}

			memberID = append(memberID, int64(num))
		}
	}

	if isAllMember {
		//--- Get All Employee
		resultMember, err = dao.EmployeeDAO.GetListMember(db, inputStruct.ID, memberID, isAllMember)
		if err.Error != nil {
			return
		}
	} else {
		if len(memberID) > 0 {
			//--- Todo Get Include Employee
			resultMember, err = dao.EmployeeDAO.GetListMember(db, inputStruct.ID, memberID, isAllMember)
			if err.Error != nil {
				return
			}
		}
	}

	result = input.convertModelToResponseDetail(employeeOnDB, resultMember)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeService) convertModelToResponseDetail(inputModel repository.EmployeeModel, resultMember []interface{}) out.ViewEmployeeResponse {
	var memberList []out.MemberList
	for _, itemResultMember := range resultMember {
		r := itemResultMember.(repository.MemberList)
		memberList = append(memberList, out.MemberList{
			ID:        r.ID.Int64,
			FirstName: r.FirstName.String,
			LastName:  r.LastName.String,
		})
	}

	return out.ViewEmployeeResponse{
		ID:                 inputModel.ID.Int64,
		IDCard:             inputModel.IDCard.String,
		FirstName:          inputModel.FirstName.String,
		LastName:           inputModel.LastName.String,
		Gender:             inputModel.Gender.String,
		PlaceOfBirth:       inputModel.PlaceOfBirth.String,
		DateOfBirth:        inputModel.DateOfBirth.Time,
		AddressResidence:   inputModel.AddressResidence.String,
		NPWP:               inputModel.NPWP.String,
		AddressTax:         inputModel.AddressTax.String,
		Email:              inputModel.Email.String,
		Phone:              inputModel.Phone.String,
		Religion:           inputModel.Religion.String,
		DateJoin:           inputModel.DateJoin.Time,
		DateOut:            inputModel.DateOut.Time,
		ReasonResignation:  inputModel.ReasonResignation.String,
		Type:               inputModel.Type.String,
		Status:             inputModel.Status.String,
		DepartmentID:       inputModel.DepartmentId.Int64,
		DepartmentName:     inputModel.DepartmentName.String,
		PositionID:         inputModel.PositionID.Int64,
		Position:           inputModel.Position.String,
		BPJSNo:             inputModel.BPJS.String,
		BPJSTkNo:           inputModel.BPJSTk.String,
		MaritalStatus:      inputModel.MaritalStatus.String,
		NumberOfDependents: inputModel.NumberOfDependents.Int64,
		Education:          inputModel.Education.String,
		Nationality:        inputModel.Nationality.String,
		MothersMaiden:      inputModel.MothersMaiden.String,
		TaxMethod:          inputModel.TaxMethod.String,
		Photo:              inputModel.Photo.String,
		IsHaveMember:       inputModel.IsHaveMember.Bool,
		MemberList:         memberList,
		LevelID:            inputModel.LevelID.Int64,
		Level:              inputModel.Level.String,
		GradeID:            inputModel.GradeID.Int64,
		Grade:              inputModel.Grade.String,
		Active:             inputModel.Active.Bool,
		CreatedAt:          inputModel.CreatedAt.Time,
		UpdatedAt:          inputModel.UpdatedAt.Time,
		CreatedName:        inputModel.CreatedName.String,
		UpdatedName:        inputModel.UpdatedName.String,
	}
}

func (input employeeService) validateViewEmployee(inputStruct *in.EmployeeRequest) errorModel.ErrorModel {
	return inputStruct.ValidateViewEmployee()
}
