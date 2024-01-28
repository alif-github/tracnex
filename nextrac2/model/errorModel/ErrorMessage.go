package errorModel

import (
	"strconv"
	"strings"
)

var DefaultError map[string]ErrorClass

type ErrorClass struct {
	ErrorCode    string
	ErrorMessage string
}

func GenerateUnauthorizedClientError(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(401, "E-1-TRAC-SRV-001", fileName, funcName)
}
func GenerateInactiveAuditSystem(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-4-TRAC-SRV-001", fileName, funcName)
}
func GenerateOnlyForOfficialPrincipalError(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-4-TRAC-SRV-006", fileName, funcName)
}
func GenerateEmptyFieldError(fileName string, funcName string, fieldName string) ErrorModel {
	errorParam := make([]ErrorParameter, 1)
	errorParam[0].ErrorParameterKey = "FieldName"
	errorParam[0].ErrorParameterValue = fieldName
	return GenerateErrorModelWithErrorParam(400, "E-4-TRAC-DTO-001", fileName, funcName, errorParam)
}
func GenerateDataUsedError(fileName string, funcName string, fieldName string) ErrorModel {
	errorParam := make([]ErrorParameter, 1)
	errorParam[0].ErrorParameterKey = "FieldName"
	errorParam[0].ErrorParameterValue = fieldName
	return GenerateErrorModelWithErrorParam(400, "E-4-TRAC-SRV-002", fileName, funcName, errorParam)
}
func GenerateInternalDBServerError(fileName string, funcName string, causedBy error) ErrorModel {
	return GenerateErrorModel(500, "E-5-TRAC-DBS-001", fileName, funcName, causedBy)
}
func GenerateFormatFieldError(fileName string, funcName string, fieldName string) ErrorModel {
	errorParam := make([]ErrorParameter, 1)
	errorParam[0].ErrorParameterKey = "FieldName"
	errorParam[0].ErrorParameterValue = fieldName
	return GenerateErrorModelWithErrorParam(400, "E-4-TRAC-DTO-002", fileName, funcName, errorParam)
}
func GenerateDataLockedError(fileName string, funcName string, fieldName string) ErrorModel {
	errorParam := make([]ErrorParameter, 1)
	errorParam[0].ErrorParameterKey = "FieldName"
	errorParam[0].ErrorParameterValue = fieldName
	return GenerateErrorModelWithErrorParam(400, "E-4-TRAC-SRV-003", fileName, funcName, errorParam)
}
func GenerateForbiddenAccessClientError(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(403, "E-3-TRAC-SRV-001", fileName, funcName)
}
func GenerateVerifyPasswordInvalidError(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(403, "E-4-TRAC-SRV-008", fileName, funcName)
}

func GenerateUnknownDataError(fileName string, funcName string, fieldName string) ErrorModel {
	errorParam := make([]ErrorParameter, 1)
	errorParam[0].ErrorParameterKey = "FieldName"
	errorParam[0].ErrorParameterValue = fieldName
	return GenerateErrorModelWithErrorParam(400, "E-4-TRAC-SRV-004", fileName, funcName, errorParam)
}

func GenerateParentCustomerMDBDataError(fileName string, funcName string, fieldName string) ErrorModel {
	errorParam := make([]ErrorParameter, 1)
	errorParam[0].ErrorParameterKey = "FieldName"
	errorParam[0].ErrorParameterValue = fieldName
	return GenerateErrorModelWithErrorParam(400, "E-4-FTR-TRAC-SRV-0011", fileName, funcName, errorParam)
}

func GenerateAuthenticationServerError(fileName string, funcName string, code int, errCode string, causedBy error) ErrorModel {
	return GenerateErrorModel(code, errCode, fileName, funcName, causedBy)
}
func GenerateFieldFormatWithRuleError(fileName string, funcName string, ruleName string, fieldName string, additionalInfo string) ErrorModel {
	errorParam := make([]ErrorParameter, 3)
	errorParam[0].ErrorParameterKey = "FieldName"
	errorParam[0].ErrorParameterValue = fieldName
	errorParam[1].ErrorParameterKey = "RuleName"
	errorParam[1].ErrorParameterValue = ruleName
	errorParam[2].ErrorParameterKey = "Other"
	errorParam[2].ErrorParameterValue = additionalInfo
	return GenerateErrorModelWithErrorParam(400, "E-4-TRAC-DTO-008", fileName, funcName, errorParam)
}
func GenerateExpiredTokenError(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(401, "E-1-TRAC-SRV-002", fileName, funcName)
}
func GenerateInvalidSignatureError(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-4-TRAC-SRV-005", fileName, funcName)
}

func GenerateUnsupportedResponseTypeError(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-4-TRAC-DTO-005", fileName, funcName)
}

func GenerateDeleteParentCustomerError(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-4-FTR-TRAC-SRV-0010", fileName, funcName)
}

// haven't used

func GenerateMissingResourceIDError(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-4-TRAC-DTO-005", fileName, funcName)
}
func GenerateExpiredRefreshTokenError(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(401, "E-1-TRAC-SRV-003", fileName, funcName)
}
func GenerateUnknownError(fileName string, funcName string, causedBy error) ErrorModel {
	return GenerateErrorModel(500, "E-5-TRAC-SRV-001", fileName, funcName, causedBy)
}
func GenerateInvalidRequestError(fileName string, funcName string, causedBy error) ErrorModel {
	return GenerateErrorModel(400, "E-4-TRAC-DTO-003", fileName, funcName, causedBy)
}
func GenerateMissingCodeChallengeError(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-4-TRAC-DTO-004", fileName, funcName)
}
func GenerateInvalidJWTCodeError(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(401, "E-4-TRAC-SRV-001", fileName, funcName)
}
func GenerateRecoverError() ErrorModel {
	return GenerateSimpleErrorModel(500, "E-5-TRAC-SRV-001")
}
func GenerateFieldHaveMaxLimitError(fileName string, funcName string, fieldName string, digitMax int) ErrorModel {
	errorParam := make([]ErrorParameter, 2)
	errorParam[0].ErrorParameterKey = "FieldName"
	errorParam[0].ErrorParameterValue = fieldName
	errorParam[1].ErrorParameterKey = "DigitMax"
	errorParam[1].ErrorParameterValue = strconv.Itoa(digitMax)
	return GenerateErrorModelWithErrorParam(400, "E-4-TRAC-DTO-009", fileName, funcName, errorParam)
}
func GenerateAlreadyExistDataError(fileName string, funcName string, primaryField string) ErrorModel {
	errorParam := make([]ErrorParameter, 1)
	errorParam[0].ErrorParameterKey = "PrimaryField"
	errorParam[0].ErrorParameterValue = primaryField
	return GenerateErrorModelWithErrorParam(400, "E-4-TRAC-SRV-002", fileName, funcName, errorParam)
}

func GenerateUsernameAlreadyUsed(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-9FTR-TRAC-SRV-001", fileName, funcName)
}

func GenerateUserStatusActive(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-9FTR-TRAC-SRV-005", fileName, funcName)
}

func GenerateUserIdHasBeenActivated(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-9FTR-TRAC-SRV-0017", fileName, funcName)
}

func GenerateResendOTPStatusActive(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-9FTR-TRAC-SRV-0014", fileName, funcName)
}

func GenerateUserRegDetailNotFound(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-9FTR-TRAC-SRV-0010", fileName, funcName)
}

func GenerateUserNexstarNotFoundInAuthentication(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-9FTR-TRAC-SRV-006", fileName, funcName)
}

func GenerateEmailEmptyAuthNexstarForResendVerification(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-9FTR-TRAC-SRV-009", fileName, funcName)
}

func GenerateUserNexmileNextradeFoundInAuthentication(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-9FTR-TRAC-SRV-007", fileName, funcName)
}

func GenerateUserNexstarHasNoGrochatResource(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-9FTR-TRAC-SRV-008", fileName, funcName)
}

func GenerateUserAuthHasNoTracResource(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-9FTR-TRAC-SRV-0016", fileName, funcName)
}

func GenerateEmailAndPhoneAlreadyRegisteredNextrac(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-9FTR-TRAC-SRV-002", fileName, funcName)
}

func GenerateBothEmailAndPhoneAlreadyRegisteredAuth(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-9FTR-TRAC-SRV-003", fileName, funcName)
}

func GenerateMismatchAuthUserID(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-9FTR-TRAC-SRV-0015", fileName, funcName)
}

func GenerateMandatoryField(fieldName string, fileName string, funcName string) ErrorModel {
	errorParam := make([]ErrorParameter, 1)
	errorParam[0].ErrorParameterKey = "FieldName"
	errorParam[0].ErrorParameterValue = fieldName
	return GenerateErrorModelWithErrorParam(400, "E-9FTR-TRAC-SRV-0011", fileName, funcName, errorParam)
}

func GenerateUserStatusNonactive(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-9FTR-TRAC-SRV-0012", fileName, funcName)
}

func GenerateDifferentAuthUserId(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-9FTR-TRAC-SRV-004", fileName, funcName)
}

func GenerateUnsupportedRequestParam(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-4-TRAC-DTO-001", fileName, funcName)
}
func GenerateChangePasswordNotValidError(fileName string, funcName string, fieldName string) ErrorModel {
	errorParam := make([]ErrorParameter, 1)
	errorParam[0].ErrorParameterKey = "Password"
	errorParam[0].ErrorParameterValue = fieldName
	return GenerateErrorModelWithErrorParam(400, "E-4-TRAC-SRV-006", fileName, funcName, errorParam)
}
func GenerateInvalidActivationCodeError(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-4-TRAC-SRV-004", fileName, funcName)
}
func GenerateActivationCodeExpiredError(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-4-TRAC-SRV-005", fileName, funcName)
}
func GenerateForgetCodeExpiredError(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-4-TRAC-SRV-011", fileName, funcName)
}
func GenerateDifferentAdditionalInformationValueError(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(404, "E-4-TRAC-SRV-009", fileName, funcName)
}
func GenerateParameterNotFoundError(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(401, "E-4-TRAC-SRV-012", fileName, funcName)
}
func GenerateInvalidRegistrationPKCE(fileName string, funcName string, detail []string) ErrorModel {
	return GenerateErrorModelWithAdditionalInformation(403, "E-6-TRAC-SRV-001", fileName, funcName, detail)
}
func GenerateInvalidRegistrationClient(fileName string, funcName string, detail []string) ErrorModel {
	return GenerateErrorModelWithAdditionalInformation(403, "E-6-TRAC-SRV-003", fileName, funcName, detail)
}
func GenerateInvalidRegistrationClientWithDetailData(fileName string, funcName string, detail interface{}) ErrorModel {
	return GenerateErrorModelWithAdditionalInformationWithData(403, "E-6-TRAC-SRV-003", fileName, funcName, detail)
}
func GenerateAuthenticationServerAddResourceError(fileName string, funcName string, detail []string) ErrorModel {
	return GenerateErrorModelWithAdditionalInformation(500, "E-6-TRAC-SRV-002", fileName, funcName, detail)
}
func GenerateDataNotFound(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-4-TRAC-SRV-012", fileName, funcName)
}
func GenerateInvalidAddResourceNexcloud(fileName string, funcName string, detail []string) ErrorModel {
	return GenerateErrorModelWithAdditionalInformation(400, "E-6-TRAC-SRV-004", fileName, funcName, detail)
}
func GenerateInvalidAddResourceNexdrive(fileName string, funcName string, detail []string) ErrorModel {
	return GenerateErrorModelWithAdditionalInformation(400, "E-6-TRAC-SRV-004", fileName, funcName, detail)
}
func GenerateInvalidRegistrationNewBranch(fileName string, funcName string, detail []string) ErrorModel {
	return GenerateErrorModelWithAdditionalInformation(403, "E-6-TRAC-SRV-006", fileName, funcName, detail)
}
func GenerateInvalidUnregistrationPKCE(fileName string, funcName string, detail []string) ErrorModel {
	return GenerateErrorModelWithAdditionalInformation(403, "E-6-TRAC-SRV-007", fileName, funcName, detail)
}
func GenerateFailedChangePassword(fileName string, funcName string, detail []string) ErrorModel {
	return GenerateErrorModelWithAdditionalInformation(403, "E-6-TRAC-SRV-008", fileName, funcName, detail)
}
func GenerateInvalidJSONRequestError(fileName string, funcName string, causedBy error) ErrorModel {
	return GenerateErrorModel(400, "E-6-TRAC-DTO-009", fileName, funcName, causedBy)
}
func GenerateDataDuplicateInDTOError(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-6-TRAC-SRV-010", fileName, funcName)
}
func GenerateRedisError(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(403, "E-RDS-TRAC-SRV-001", fileName, funcName)
}
func GenerateInvalidDifferentCompareData(fileName string, funcName string, fieldNameA string, fieldNameB string) ErrorModel {
	errorParam := make([]ErrorParameter, 2)
	errorParam[0].ErrorParameterKey = "NewPassword"
	errorParam[0].ErrorParameterValue = fieldNameA
	errorParam[1].ErrorParameterKey = "ConfirmationPassword"
	errorParam[1].ErrorParameterValue = fieldNameB
	return GenerateErrorModelWithErrorParam(400, "E-4-TRAC-DTO-004", fileName, funcName, errorParam)
}
func GenerateMultipleErrorAcquired(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-6-TRAC-SRV-011", fileName, funcName)
}
func GenerateEmptyFieldOrZeroValueError(fileName string, funcName string, fieldName string) ErrorModel {
	errorParam := make([]ErrorParameter, 1)
	errorParam[0].ErrorParameterKey = "FieldName"
	errorParam[0].ErrorParameterValue = fieldName
	return GenerateErrorModelWithErrorParam(400, "E-4-TRAC-DTO-006", fileName, funcName, errorParam)
}
func GenerateDataUsedRegisterClientError(fileName string, funcName string, detail []string) ErrorModel {
	return GenerateErrorModelWithAdditionalInformation(400, "E-6-TRAC-SRV-012", fileName, funcName, detail)
}
func GenerateDataUsedRegisterClientDiffClientIDError(fileName string, funcName string, detail []string) ErrorModel {
	return GenerateErrorModelWithAdditionalInformation(400, "E-6-TRAC-SRV-013", fileName, funcName, detail)
}

func GenerateReservedValueString(fileName string, funcName string, fieldName string) ErrorModel {
	errorParam := make([]ErrorParameter, 1)
	errorParam[0].ErrorParameterKey = "FieldName"
	errorParam[0].ErrorParameterValue = fieldName
	return GenerateErrorModelWithErrorParam(400, "E-4-TRAC-SRV-009", fileName, funcName, errorParam)
}

func GenerateErrorFormatJSON(fileName string, funcName string, causedBy error) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-4-TRAC-SRV-010", fileName, funcName)
}

func GenerateDataScopeNotDefinedYet(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(403, "E-3-TRAC-SRV-002", fileName, funcName)
}

func GenerateDuplicateErrorWithParam(fileName string, funcName string, fieldName string) ErrorModel {
	errorParam := make([]ErrorParameter, 1)
	errorParam[0].ErrorParameterKey = "FieldName"
	errorParam[0].ErrorParameterValue = fieldName
	return GenerateErrorModelWithErrorParam(400, "E-6-TRAC-SRV-014", fileName, funcName, errorParam)
}

func GenerateClientIDNotFound(fileName string, funcName string, fieldName string, objectFieldName string, otherFieldName string) ErrorModel {
	errorParam := make([]ErrorParameter, 3)
	errorParam[0].ErrorParameterKey = "FieldName"
	errorParam[0].ErrorParameterValue = fieldName
	errorParam[1].ErrorParameterKey = "ObjectFieldName"
	errorParam[1].ErrorParameterValue = objectFieldName
	errorParam[2].ErrorParameterKey = "OtherFieldName"
	errorParam[2].ErrorParameterValue = otherFieldName
	return GenerateErrorModelWithErrorParam(403, "E-6-TRAC-SRV-015", fileName, funcName, errorParam)
}

func GenerateDateValidateFromThru(fileName string, funcName string, code string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, code, fileName, funcName)
}

func GenerateActivationLicenseError(fileName string, funcName string, message string, causedBy error) ErrorModel {
	errorParam := make([]ErrorParameter, 1)
	errorParam[0].ErrorParameterKey = "FieldName"
	errorParam[0].ErrorParameterValue = message
	return GenerateErrorModelWithErrorParamAndCaused(400, "E-6-TRAC-SRV-017", fileName, funcName, errorParam, causedBy)
}

func GenerateUserLicenseNotFound(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-6-TRAC-SRV-018", fileName, funcName)
}

func GenerateLicenseHasBeenActivated(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-6-TRAC-SRV-019", fileName, funcName)
}

func GenerateLicenseHasNotBeenActivated(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-6-TRAC-SRV-022", fileName, funcName)
}

func GenerateTotalActivatedZeroValue(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-6-TRAC-SRV-024", fileName, funcName)
}

func GenerateLicenseHasBeenDeactivated(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-6-TRAC-SRV-023", fileName, funcName)
}

func GenerateUserLicenseFullFilled(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-6-TRAC-SRV-020", fileName, funcName)
}

func GenerateForbiddenClientCredentialAccess(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(403, "E-6-TRAC-SRV-021", fileName, funcName)
}

func GenerateNotFoundActiveLicense(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-6-TRAC-SRV-025", fileName, funcName)
}

func GenerateStatusHasNotBeenActivated(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-6-TRAC-SRV-026", fileName, funcName)
}

func GenerateClientValidationError(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-6-TRAC-SRV-027", fileName, funcName)
}

func GenerateUnknownAuthUserId(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-6-TRAC-SRV-028", fileName, funcName)
}

func GenerateGetDataNotFound(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-6-TRAC-SRV-018", fileName, funcName)
}

func GenerateDifferentRequestAndDBResult(fileName string, funcName string, compare1 string, compare2 string) ErrorModel {
	errorParam := make([]ErrorParameter, 2)
	errorParam[0].ErrorParameterKey = "NewPassword"
	errorParam[0].ErrorParameterValue = compare1
	errorParam[1].ErrorParameterKey = "ConfirmationPassword"
	errorParam[1].ErrorParameterValue = compare2
	return GenerateErrorModelWithErrorParam(403, "E-4-TRAC-DTO-004", fileName, funcName, errorParam)
}

func GenerateWrongAction(fileName string, funcName string, fieldName string, message string, number int) ErrorModel {
	errorParam := make([]ErrorParameter, 3)
	errorParam[0].ErrorParameterKey = "FieldName"
	errorParam[0].ErrorParameterValue = fieldName
	errorParam[1].ErrorParameterKey = "RuleName"
	errorParam[1].ErrorParameterValue = message
	errorParam[2].ErrorParameterKey = "Number"
	errorParam[2].ErrorParameterValue = strconv.Itoa(number)
	return GenerateErrorModelWithErrorParam(400, "E-6FX-TRAC-SRV-029", fileName, funcName, errorParam)
}

func GenerateWrongActionInstallation(fileName string, funcName string, fieldName string, message string, numberSite int, numberInstallation int) ErrorModel {
	errorParam := make([]ErrorParameter, 4)
	errorParam[0].ErrorParameterKey = "FieldName"
	errorParam[0].ErrorParameterValue = fieldName
	errorParam[1].ErrorParameterKey = "RuleName"
	errorParam[1].ErrorParameterValue = message
	errorParam[2].ErrorParameterKey = "NumberSite"
	errorParam[2].ErrorParameterValue = strconv.Itoa(numberSite)
	errorParam[3].ErrorParameterKey = "NumberInstallation"
	errorParam[3].ErrorParameterValue = strconv.Itoa(numberInstallation)
	return GenerateErrorModelWithErrorParam(400, "E-6FX-TRAC-SRV-030", fileName, funcName, errorParam)
}

func GenerateFieldMustEmptyInstallation(fileName string, funcName string, fieldName string) ErrorModel {
	errorParam := make([]ErrorParameter, 1)
	errorParam[0].ErrorParameterKey = "FieldName"
	errorParam[0].ErrorParameterValue = fieldName
	return GenerateErrorModelWithErrorParam(400, "E-6FX-TRAC-SRV-031", fileName, funcName, errorParam)
}

func GenerateActionMustInsertInstallation(fileName string, funcName string, fieldName string) ErrorModel {
	errorParam := make([]ErrorParameter, 1)
	errorParam[0].ErrorParameterKey = "FieldName"
	errorParam[0].ErrorParameterValue = fieldName
	return GenerateErrorModelWithErrorParam(400, "E-6FX-TRAC-SRV-032", fileName, funcName, errorParam)
}

func GenerateRequestError(fileName string, funcName string, message string) ErrorModel {
	errorParam := make([]ErrorParameter, 1)
	errorParam[0].ErrorParameterKey = "RuleName"
	errorParam[0].ErrorParameterValue = message
	return GenerateErrorModelWithErrorParam(400, "E-4-TRAC-SRV-013", fileName, funcName, errorParam)
}

func GenerateErrorClientTypeNotParent(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-6FX-TRAC-SRV-033", fileName, funcName)
}

func GenerateCannotChangedError(fileName, funcName, fieldName string) ErrorModel {
	errorParam := make([]ErrorParameter, 1)
	errorParam[0].ErrorParameterKey = "FieldName"
	errorParam[0].ErrorParameterValue = fieldName
	return GenerateErrorModelWithErrorParam(400, "E-4-TRAC-SRV-014", fileName, funcName, errorParam)
}

func GenerateErrorParentInstallationNotFound(fileName string, funcName string, fieldName string) ErrorModel {
	errorParam := make([]ErrorParameter, 1)
	errorParam[0].ErrorParameterKey = "FieldName"
	errorParam[0].ErrorParameterValue = fieldName
	return GenerateErrorModelWithErrorParam(400, "E-6FX-TRAC-SRV-034", fileName, funcName, errorParam)
}

func GenerateErrorParentAppUpdatedDeleted(fileName string, funcName string, fieldName string) ErrorModel {
	errorParam := make([]ErrorParameter, 1)
	errorParam[0].ErrorParameterKey = "FieldName"
	errorParam[0].ErrorParameterValue = fieldName
	return GenerateErrorModelWithErrorParam(400, "E-6FX-TRAC-SRV-035", fileName, funcName, errorParam)
}

func GenerateErrorEksternalClientTypeMustHave(fileName string, funcName string, fieldName string) ErrorModel {
	errorParam := make([]ErrorParameter, 1)
	errorParam[0].ErrorParameterKey = "FieldName"
	errorParam[0].ErrorParameterValue = fieldName
	return GenerateErrorModelWithErrorParam(400, "E-4-ETR-TRAC-SRV-001", fileName, funcName, errorParam)
}

func GenerateErrorEksternalClientIDNotValid(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-4-ETR-TRAC-SRV-002", fileName, funcName)
}

func GenerateErrorEksternalSpawnTimeResendOTP(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-4-ETR-TRAC-SRV-003", fileName, funcName)
}

func GenerateErrorCustomActivationCode(fileName, funcName, message string) ErrorModel {
	errorParam := make([]ErrorParameter, 1)
	errorParam[0].ErrorParameterKey = "RuleName"
	errorParam[0].ErrorParameterValue = message
	return GenerateErrorModelWithErrorParam(400, "E-4-ETR-TRAC-SRV-004", fileName, funcName, errorParam)
}

func GenerateInvalidAddBranch(fileName string, funcName string, detail interface{}) ErrorModel {
	return GenerateErrorModelWithAdditionalInformationWithData(403, "E-6-TRAC-SRV-006", fileName, funcName, detail)
}

func GenerateErrorSpawnSynchronize(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-4-ETR-TRAC-SRV-006", fileName, funcName)
}

func GenerateCannotBeUsedError(fileName, funcName, fieldName string) ErrorModel {
	errorParam := make([]ErrorParameter, 1)
	errorParam[0].ErrorParameterKey = "FieldName"
	errorParam[0].ErrorParameterValue = fieldName
	return GenerateErrorModelWithErrorParam(400, "E-4-ETR-TRAC-SRV-007", fileName, funcName, errorParam)
}

func GenerateErrorChildHierarchy(fileName, funcName string, errCount int, detail []string) ErrorModel {
	errorParam := make([]ErrorParameter, 1)
	errorParam[0].ErrorParameterKey = "Count"
	errorParam[0].ErrorParameterValue = strconv.Itoa(errCount)
	return GenerateErrorModelWithAdditionalInformationAndErrorParam(400, "E-4-ETR-TRAC-SRV-008", fileName, funcName, detail, errorParam)
}

func GenerateDataNotFoundWithParam(fileName string, funcName string, fieldName string) ErrorModel {
	errorParam := make([]ErrorParameter, 1)
	errorParam[0].ErrorParameterKey = "FieldName"
	errorParam[0].ErrorParameterValue = fieldName
	return GenerateErrorModelWithErrorParam(400, "E-4-ETR-TRAC-SRV-009", fileName, funcName, errorParam)
}

func GenerateSprintFilterError(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-4-FTR-TRAC-SRV-0012", fileName, funcName)
}

func GenerateValueMustLessThanError(fileName, funcName, fieldName string, max string) ErrorModel {
	errorParam := make([]ErrorParameter, 2)
	errorParam[0].ErrorParameterKey = "FieldName"
	errorParam[0].ErrorParameterValue = fieldName
	errorParam[1].ErrorParameterKey = "Max"
	errorParam[1].ErrorParameterValue = max
	return GenerateErrorModelWithErrorParam(400, "E-4-TRAC-DTO-009", fileName, funcName, errorParam)
}

func GenerateDateContainsWeekendError(fileName string, funcName string, fieldName string) ErrorModel {
	errorParam := make([]ErrorParameter, 1)
	errorParam[0].ErrorParameterKey = "FieldName"
	errorParam[0].ErrorParameterValue = fieldName
	return GenerateErrorModelWithErrorParam(400, "E-4-TRAC-DTO-010", fileName, funcName, errorParam)
}

func GenerateDateMustBeLaterThanError(fileName string, funcName string, fieldName1, fieldName2 string) ErrorModel {
	errorParam := make([]ErrorParameter, 2)
	errorParam[0].ErrorParameterKey = "FieldName1"
	errorParam[0].ErrorParameterValue = fieldName1
	errorParam[1].ErrorParameterKey = "FieldName2"
	errorParam[1].ErrorParameterValue = fieldName2
	return GenerateErrorModelWithErrorParam(400, "E-4-TRAC-DTO-011", fileName, funcName, errorParam)
}

func GenerateDateMustBeLaterThanCurrentDateError(fileName string, funcName string, fieldName string) ErrorModel {
	errorParam := make([]ErrorParameter, 1)
	errorParam[0].ErrorParameterKey = "FieldName"
	errorParam[0].ErrorParameterValue = fieldName
	return GenerateErrorModelWithErrorParam(400, "E-4-TRAC-DTO-014", fileName, funcName, errorParam)
}

func GenerateInvalidFileExtensionError(fileName string, funcName string, fieldName string, validExtensions []string) ErrorModel {
	errorParam := make([]ErrorParameter, 2)
	errorParam[0].ErrorParameterKey = "FieldName"
	errorParam[0].ErrorParameterValue = fieldName
	errorParam[1].ErrorParameterKey = "ValidExtensions"
	errorParam[1].ErrorParameterValue = strings.Join(validExtensions, ", ")
	return GenerateErrorModelWithErrorParam(400, "E-4-TRAC-DTO-012", fileName, funcName, errorParam)
}

func GenerateFileSizeExceedsMaxLimitError(fileName string, funcName string, fieldName string, maxSize int, unit string) ErrorModel {
	errorParam := make([]ErrorParameter, 3)
	errorParam[0].ErrorParameterKey = "FieldName"
	errorParam[0].ErrorParameterValue = fieldName
	errorParam[1].ErrorParameterKey = "MaxSize"
	errorParam[1].ErrorParameterValue = strconv.Itoa(maxSize)
	errorParam[2].ErrorParameterKey = "Unit"
	errorParam[2].ErrorParameterValue = unit
	return GenerateErrorModelWithErrorParam(400, "E-4-TRAC-DTO-013", fileName, funcName, errorParam)
}

func GenerateCannotDeleteData(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-4-FTR-TRAC-SRV-0013", fileName, funcName)
}

func GenerateDateCannotBeLessThanCurrentDate(fileName string, funcName string, fieldName string) ErrorModel {
	errorParam := make([]ErrorParameter, 1)
	errorParam[0].ErrorParameterKey = "FieldName"
	errorParam[0].ErrorParameterValue = fieldName
	return GenerateErrorModelWithErrorParam(400, "E-4-TRAC-DTO-015", fileName, funcName, errorParam)
}

func GenerateDataExpiredError(fileName string, funcName string, fieldName string) ErrorModel {
	errorParam := make([]ErrorParameter, 1)
	errorParam[0].ErrorParameterKey = "FieldName"
	errorParam[0].ErrorParameterValue = fieldName
	return GenerateErrorModelWithErrorParam(400, "E-4-TRAC-DTO-016", fileName, funcName, errorParam)
}

func GenerateDateCannotBeLessThanError(fileName string, funcName string, fieldName1 string, fieldName2 string) ErrorModel {
	errorParam := make([]ErrorParameter, 2)
	errorParam[0].ErrorParameterKey = "FieldName1"
	errorParam[0].ErrorParameterValue = fieldName1
	errorParam[1].ErrorParameterKey = "FieldName2"
	errorParam[1].ErrorParameterValue = fieldName2
	return GenerateErrorModelWithErrorParam(400, "E-4-TRAC-DTO-017", fileName, funcName, errorParam)
}

func GenerateMaxDecimalPlacesError(fileName string, funcName string, fieldName string, maxDigits int) ErrorModel {
	errorParam := make([]ErrorParameter, 2)
	errorParam[0].ErrorParameterKey = "FieldName"
	errorParam[0].ErrorParameterValue = fieldName
	errorParam[1].ErrorParameterKey = "MaxDigits"
	errorParam[1].ErrorParameterValue = strconv.Itoa(maxDigits)
	return GenerateErrorModelWithErrorParam(400, "E-4-TRAC-DTO-018", fileName, funcName, errorParam)
}

func GenerateUnknownUserAuth(fileName string, funcName string) ErrorModel {
	return GenerateErrorModelWithoutCaused(400, "E-6-TRAC-SRV-029", fileName, funcName)
}
func GenerateAmountOfLeaveExceedMaxLeaveError(fileName string, funcName string, maxLeave int) ErrorModel {
	errorParam := make([]ErrorParameter, 1)
	errorParam[0].ErrorParameterKey = "MaxLeave"
	errorParam[0].ErrorParameterValue = strconv.Itoa(maxLeave)

	return GenerateErrorModelWithErrorParam(400, "E-4-TRAC-SRV-015", fileName, funcName, errorParam)
}