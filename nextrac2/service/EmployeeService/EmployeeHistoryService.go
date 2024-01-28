package EmployeeService

import (
	"encoding/json"
	"github.com/pkg/errors"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"strconv"
)

func (input employeeService) descriptionEmployeeHistory(recordBefore, recordAfter []repository.AuditSystemModel) (history string) {
	var (
		fileName       = "EmployeeHistoryService.go"
		funcName       = "descriptionEmployeeHistory"
		err            errorModel.ErrorModel
		errorS         error
		employeeBefore in.EmployeeJSONDB
		employeeAfter  in.EmployeeJSONDB
		serverVersion  = config.ApplicationConfiguration.GetServerVersion()
		resourceID     = config.ApplicationConfiguration.GetServerResourceID()
	)

	defer func() {
		if err.Error != nil {
			logModel := applicationModel.GenerateLogModel(serverVersion, resourceID)
			logModel.Message = err.CausedBy.Error()
			logModel.Status = err.Code
			util.LogError(logModel.ToLoggerObject())
		}
	}()

	//--- Check if record empty then error
	if len(recordBefore) < 1 || len(recordAfter) < 1 {
		err = errorModel.GenerateUnknownError(fileName, funcName, errors.New("record empty"))
		return
	}

	//--- Record Before Unmarshal
	errorS = json.Unmarshal([]byte(recordBefore[0].Data.String), &employeeBefore)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errorS)
		return
	}

	//--- Record After Unmarshal
	errorS = json.Unmarshal([]byte(recordAfter[0].Data.String), &employeeAfter)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errorS)
		return
	}

	historyTemp := input.mappingAndAppendToHistory(employeeBefore, employeeAfter)
	if len(historyTemp) > 0 {
		byteTemp, errs := json.Marshal(historyTemp)
		if errs != nil {
			err = errorModel.GenerateUnknownError(fileName, funcName, errs)
			return
		}
		history = string(byteTemp)
	}

	return
}

func (input employeeService) mappingAndAppendToHistory(before, after in.EmployeeJSONDB) (history []in.EmployeeHistory) {
	input.addToHistoryStruct("Kartu ID", "ID Card", before.IDCard, after.IDCard, &history)                                               //--- Check ID Card
	input.addToHistoryStruct("NPWP", "NPWP", before.NPWP, after.NPWP, &history)                                                          //--- Check NPWP
	input.addToHistoryStruct("Nama Depan", "Firstname", before.FirstName, after.FirstName, &history)                                     //--- Check First Name
	input.addToHistoryStruct("Nama Belakang", "Lastname", before.LastName, after.LastName, &history)                                     //--- Check Last Name
	input.addToHistoryStruct("Departemen", "Department", before.DepartmentID, after.DepartmentID, &history)                              //--- Check Department
	input.addToHistoryStruct("Email", "Email", before.Email, after.Email, &history)                                                      //--- Check Email
	input.addToHistoryStruct("Telepon", "Phone", before.Phone, after.Phone, &history)                                                    //--- Check Phone
	input.addToHistoryStruct("Jenis Kelamin", "Gender", before.Gender, after.Gender, &history)                                           //--- Check Gender
	input.addToHistoryStruct("Tempat Lahir", "Place of Birth", before.PlaceOfBirth, after.PlaceOfBirth, &history)                        //--- Check Place Of Birth
	input.addToHistoryStruct("Tanggal Lahir", "Date of Birth", before.DateOfBirth, after.DateOfBirth, &history)                          //--- Check Date Of Birth
	input.addToHistoryStruct("Alamat Tinggal", "Address Residence", before.AddressResidence, after.AddressResidence, &history)           //--- Check Address Residence
	input.addToHistoryStruct("Alamat Pajak", "Address Tax", before.AddressTax, after.AddressTax, &history)                               //--- Check Address Tax
	input.addToHistoryStruct("Tanggal Bergabung", "Date Join", before.DateJoin, after.DateJoin, &history)                                //--- Check Date Join
	input.addToHistoryStruct("Tanggal Keluar", "Date Out", before.DateOut, after.DateOut, &history)                                      //--- Check Date Out
	input.addToHistoryStruct("Agama", "Religion", before.Religion, after.Religion, &history)                                             //--- Check Religion
	input.addToHistoryStruct("Tipe", "Type", before.Type, after.Type, &history)                                                          //--- Check Type
	input.addToHistoryStruct("Status", "Status", before.Status, after.Status, &history)                                                  //--- Check Status
	input.addToHistoryStruct("Posisi", "Position", before.EmployeePositionID, after.EmployeePositionID, &history)                        //--- Check Position
	input.addToHistoryStruct("Status Pernikahan", "Marital Status", before.MaritalStatus, after.MaritalStatus, &history)                 //--- Check Marital Status
	input.addToHistoryStruct("Pendidikan", "Education", before.Education, after.Education, &history)                                     //--- Check Education
	input.addToHistoryStruct("Mothers Maiden", "Mothers Maiden", before.MothersMaiden, after.MothersMaiden, &history)                    //--- Check Mothers Maiden
	input.addToHistoryStruct("Jumlah Tanggungan", "Number of Dependents", before.NumberOfDependents, after.NumberOfDependents, &history) //--- Check Numbers Of Dependents
	input.addToHistoryStruct("Kebangsaan", "Nationality", before.Nationality, after.Nationality, &history)                               //--- Check Nationality
	input.addToHistoryStruct("Metode Pajak", "Tax Method", before.TaxMethod, after.TaxMethod, &history)                                  //--- Check Tax Method
	input.addToHistoryStruct("Alasan Resign", "Reason Resignation", before.ReasonResignation, after.ReasonResignation, &history)         //--- Check Reason Resignation
	input.addToHistoryStruct("Aktif", "Active", before.Active, after.Active, &history)                                                   //--- Check Active
	input.addToHistoryStruct("Foto", "Photo", before.Photo, after.Photo, &history)                                                       //--- Check Photo
	input.addToHistoryStruct("Sebagai Atasan?", "As a Lead?", before.IsHaveMember, after.IsHaveMember, &history)                         //--- Check Is Have Member
	input.addToHistoryStruct("Anggota", "Member", before.Member, after.Member, &history)                                                 //--- Check Member (*)
	input.addToHistoryStruct("No BPJS", "No BPJS", before.NoBpjs, after.NoBpjs, &history)                                                //--- Check No BPJS
	input.addToHistoryStruct("No BPJS Tenaga Kerja", "No BPJS-TK", before.NoBpjsTk, after.NoBpjsTk, &history)                            //--- Check No BPJS TK
	input.addToHistoryStruct("Level", "Level", before.Level, after.Level, &history)                                                      //--- Check Level ID
	input.addToHistoryStruct("Grade", "Grade", before.Grade, after.Grade, &history)                                                      //--- Check Grade ID
	input.addToHistoryStruct("Foto", "Photo", before.FileUploadID, after.FileUploadID, &history)
	return
}

func (input employeeService) addToHistoryStruct(keyID, keyEn string, before, after interface{}, history *[]in.EmployeeHistory) {
	var emptyHistory in.EmployeeHistory
	temp := input.checkAndCompareHistory(keyID, keyEn, before, after)
	if temp != emptyHistory {
		*history = append(*history, temp)
	}
}

func (input employeeService) checkAndCompareHistory(keyId, KeyEn string, before, after interface{}) in.EmployeeHistory {
	var (
		recordBefore string
		recordAfter  string
	)

	recordBefore = input.parseInterfaceToString(before)
	recordAfter = input.parseInterfaceToString(after)
	if recordBefore == recordAfter {
		return in.EmployeeHistory{}
	}

	return in.EmployeeHistory{
		KeyId:  keyId,
		KeyEn:  KeyEn,
		Before: recordBefore,
		After:  recordAfter,
	}
}

func (input employeeService) parseInterfaceToString(value interface{}) (record string) {
	switch value.(type) {
	case int64:
		num := value.(int64)
		record = strconv.Itoa(int(num))
	case string:
		record = value.(string)
	case bool:
		bol := value.(bool)
		record = strconv.FormatBool(bol)
	default:
	}

	return
}
