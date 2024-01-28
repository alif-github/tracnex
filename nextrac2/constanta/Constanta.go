package constanta

// --------------------------- Header Request Constanta ------------------------------------
const RequestIDConstanta = "X-Request-ID"
const RequestCacheControl = "Cache-Control"
const IPAddressConstanta = "X-Forwarded-For"
const SourceConstanta = "X-Source"
const TokenHeaderNameConstanta = "Authorization"
const RedirectResponseValidateConstanta = "X-NEXREDIRECT"
const TokenResponseValidateConstanta = "X-NEXTOKEN"
const RefreshResponseValidateConstanta = "X-NEXREFRESH"
const StateResponseValidateConstanta = "X-NEXSTATE"
const CodeResponseValidateConstanta = "X-Nexcode"
const DefaultTokenKeyConstanta = "X-NEXTOKEN"
const DefaultTokenValueConstanta = "MasterData2021"
const TimestampSignatureHeaderNameConstanta = "X-Timestamp"
const SignatureHeaderNameConstanta = "X-Signature"
const X_NEXCODE = "X-Nexcode"

// --------------------------------- Expired Time Constanta ---------------------------------------------------------
const ExpiredTokenOnRedisConstanta = 0

// ---------------------------------- Context Name Constanta --------------------------------------------------------
const ApplicationContextConstanta = "application_context"

// ---------------------------------- Default Variable --------------------------------------------------------------
const DefaultApplicationsLanguage = "id-ID"
const DefaultTimeFormat = "2006-01-02T15:04:05Z"
const DefaultTimeWithNanoFormat = "2006-01-02T15:04:05.000000Z"
const DefaultDBTimeFormat = "2006-01-02T15:04:05"
const DefaultDBSQLTimeFormat = "2006-01-02 15:04:05"
const DefaultTimeSnoozeNotification = "02-01-2006 03:04 PM"
const DefaultTimeSprintFormat = "20060102"
const DefaultRoleUser = "Pending"
const DefaultGroupUser = "Pending"

const ActionAuditDeleteConstanta = 0
const ActionAuditInsertConstanta = 1
const ActionAuditUpdateConstanta = 2

const SystemID = 1
const SystemClient = "SYSTEM"
const RandTokenForDeleteLength = 8

// ---- Job Process
const UpdateLastUpdateTimeInMinute = 3
const JobProcessErrorStatus = "ERROR"
const JobProcessOnProgressStatus = "ONPROGRESS"
const JobProcessOnProgressErrorStatus = "ONPROGRESS-ERROR"
const JobProcessDoneStatus = "OK"
const JobProcessAddResourceType = "Add Resource"
const JobProcessUpdateType = "Update"
const JobProcessElasticType = "Elastic"
const JobProcessExpirationType = "Expiration"
const JobProcessSendEmailType = "Send Email"
const JobProcessMasterDataType = "Master Data"

const JobProcessUserGroup = "User"
const JobProcessSynchronizeGroup = "Synchronize"
const JobProcessCheckGroup = "Check"

const JobProcessSyncElasticBank = JobProcessSynchronizeGroup + " " + JobProcessElasticType + " Bank"
const JobProcessSyncElasticPersonProfile = JobProcessSynchronizeGroup + " " + JobProcessElasticType + " Person Profile"
const JobProcessSyncElasticCompanyTitle = JobProcessSynchronizeGroup + " " + JobProcessElasticType + " Company Title"
const JobProcessSyncElasticDistrict = JobProcessSynchronizeGroup + " " + JobProcessElasticType + " District"
const JobProcessSyncElasticCompanyProfile = JobProcessSynchronizeGroup + " " + JobProcessElasticType + " Company Profile"
const JobProcessSyncElasticProvince = JobProcessSynchronizeGroup + " " + JobProcessElasticType + " Province"
const JobProcessSyncElasticPersonTitle = JobProcessSynchronizeGroup + " " + JobProcessElasticType + " Person Title"
const JobProcessSyncElasticIsland = JobProcessSynchronizeGroup + " " + JobProcessElasticType + " Island"
const JobProcessSyncElasticPostalCode = JobProcessSynchronizeGroup + " " + JobProcessElasticType + " Postal Code"
const JobProcessSyncElasticSubDistrict = JobProcessSynchronizeGroup + " " + JobProcessElasticType + " Sub District"
const JobProcessSyncElasticUrbanVillage = JobProcessSynchronizeGroup + " " + JobProcessElasticType + " Urban Village"
const JobProcessSyncElasticCountry = JobProcessSynchronizeGroup + " " + JobProcessElasticType + " Country"
const JobProcessAddResourceNexcloudName = JobProcessUserGroup + " " + JobProcessAddResourceType + " Nexcloud"
const JobProcessUpdateLogAfterAddResourceNexcloudName = JobProcessUserGroup + " " + JobProcessUpdateType + " Log After Add Resource"
const JobProcessCheckExpirationProductLicense = JobProcessCheckGroup + " " + JobProcessExpirationType + " Product License"
const JobProcessSendEmailReminder = JobProcessUserGroup + " " + JobProcessSendEmailType + " Reminder"
const JobProcessSynchronizeRegional = JobProcessSynchronizeGroup + " " + JobProcessMasterDataType + " Regional"
const JobProcessSynchronizeProvince = JobProcessSynchronizeGroup + " " + JobProcessMasterDataType + " Province"
const JobProcessSynchronizeDistrict = JobProcessSynchronizeGroup + " " + JobProcessMasterDataType + " District"
const JobProcessSynchronizeSubDistrict = JobProcessSynchronizeGroup + " " + JobProcessMasterDataType + " Sub District"
const JobProcessSynchronizePostalCode = JobProcessSynchronizeGroup + " " + JobProcessMasterDataType + " Postal Code"
const JobProcessSynchronizeUrbanVillage = JobProcessSynchronizeGroup + " " + JobProcessMasterDataType + " Urban Village"

// ---- Photo
const PersonProfilePhotoPrefix = "photo/personprofile/"
const PersonProfileMaximumPhoto = 1
const PersonProfileMaximumPhotoSize = 500000

const CompanyProfileLogoPrefix = "logo/companyprofile/"
const CompanyProfileMaximumLogo = 1
const CompanyProfileMaximumPhotoSize = 500000

const PrincipalPhotoPrefix = "photo/principal/"
const PrincipalMaximumPhoto = 1
const PrincipalMaximumPhotoSize = 500000

const ProductPhotoPrefix = "photo/product/"
const ProductMaximumPhoto = 5
const ProductMaximumPhotoSize = 500000

// ---- Channel
const TotalDataProductLicensePerChannel = 100

// ---- Client
const ScopeClient = "read write"
const GrantTypes = "client_credentials"
const AccessTokenValidity = 86400000000000
const RefreshTokenValidity = 2592000000000000
const MaxAuthFail = 10
const AuthUserNonPKCE = 0
const RoleUserND6 = 2
const StatusActive = "A"
const StatusPending = "P"

// ---- Activation GroChat
const StatusMessage = "success"
const MessageSent = "Message sent"

// ---- PKCE
const IndonesianCodeNumber = "+62"
const IndonesianCodeNumberWithDash = "+62-"
const IndonesianLanguage = "id-ID"
const RoleUserNexMile = 3
const PendingOnApproval = "P"
const UserTypeNexmile = 1
const UserTypeNexstar = 2

// ---- Source system
const ND6 = "ND6"
const Nexdistribution = "Nexdistribution"
const Nexmile = "NexMile"
const Nexstar = "NexStar"
const Nextrade = "NexTrade"
const Resend = "Resend"
const AuthDestination = "auth"

// ---- External Resource ID name
const NexCloudResourceID = "api"
const NexdriveResourceID = "drive"
const NexmileResourceID = "mile"
const NexstarResourceID = "star"
const GroChatResourceID = "chat"
const Nd6ResourceID = "nd6"

// ---- Error cause by external nexcloud
const ErrorClientIDExistNexcloud = "E-4-NC2-SRV-002"

// ---- Important constanta for import file customer list
const SizeMaximumRequireFileImport = 10 * 1024 * 1024
const GetTypeImportFile = "Content-Type"
const DefineTypeImportFile = "text/csv"
const NameFileLevel2 = "list"
const NameFileLevel3 = "customer"
const TanggerangDataGoup = "Tangerang"

// ---------------------------------- Log Level --------------------------------------------------------
const LogLevelNotSet = "0"
const LogLevelDebug = "10"
const LogLevelInfo = "20"
const LogLevelWarning = "30"
const LogLevelError = "40"
const LogLevelCritical = "50"

// --------------------------------- Layout Date for file -----------------------------------------------
const DefaultTimeFormatForFile = "2006/01/02"

// ---------------------------------- Resource ID in DB -------------------------------------------------
const ResourceNexmileID = 2
const ResourceTestingNexmileID = 5
const ResourceTestingNexchiefMobileID = 6
const ResourceND6ID = 1
const ResourceNexChiefID = 4
const ResourceNexstarID = 3
const ResourceNextradeID = 4

// ---------------------------------- Role ID in DB -----------------------------------------------------
const RoleSuperAdminName = "super_admin"

// ------------------------------------ Table Name -----------------------------------------------------
const TableNameMenuSysAdmin = "menu_sysadmin"
const TableNameMenuParent = "parent_menu"
const TableMenuService = "service_menu"
const TableMenuItem = "menu_item"
const TableNameMasterDataCompanyProfile = "company_profile"

// ----------------------------------- Master Data -----------------------------------------------------
const ResourceMasterData = "master"
const Issue = "Access from trac"

// ---------------------------- Format Time Installation -----------------------------------------------
const DefaultInstallationTimeFormat = "2006-01-02"

// ----------------------- Constanta For Update With Action --------------------------------------------
const ActionInsertCode = 1
const ActionUpdateCode = 2
const ActionDeleteCode = 3
const ActionNoActionCode = 4

// ----------------------- Constanta For Product License Status --------------------------------------------
const ProductLicenseStatusActive = 1
const ProductLicenseStatusExpired = 2
const ProductLicenseStatusTerminated = 3
const ProductEncryptAction = "ProductEncrypt"
const ProductDecryptAction = "ProductDecrypt"

// ------------------------------ No Rows Result in DB -----------------------------------------------------
const NoRowsInDB = "sql: no rows in result set"
const ErrorMDBDataNotFound = "E-4-MAD-SRV-004"

// ----------------------- Constanta For User Registration Detail Status -----------------------------------------------
const NonactiveUser = "N"

const StatusRegistered = "R"
const StatusNonActive = "N"
const FlagStatusTrue = "Y"
const TagNameJSON = "json"
const Mandatory = "mandatory"
const Optional = "optional"

// ---------------------------------- Constanta For TodoList -----------------------------------------------------------
const TodoListFormatTime = "2006-01-02 15:04:05"
const RepeatTypeDaily = "D"
const RepeatTypeWeekly = "W"
const RepeatTypeWeekday = "WD"
const RepeatTypeMonthly = "M"
const RepeatTypeYearly = "Y"
const TaskMyDay = "task_my_day"
const TaskPlanned = "task_planned"
const Task = "task"
const TaskImportant = "task_important"
const TaskListCustom = "task_list_custom"
const RemarkTodoListReminderType = "Reminder Type"
const RemarkTodoListRepeatType = "Repeat Type"
const Reminder0MinutesBeforeID = 5
const Reminder5MinutesBeforeID = 6
const Reminder15MinutesBeforeID = 7
const Reminder30MinutesBeforeID = 8
const Reminder1HourBeforeID = 9
const Reminder2HourBeforeID = 10
const Reminder12HourBeforeID = 11
const Reminder1DayBeforeID = 12
const Reminder2DayBeforeID = 13
const Reminder1WeekBeforeID = 14

// ----------------------------------------- Constanta For Parameter ---------------------------------------------------
const NexTracParameterPermission = "nextrac"

const ParameterLabelName = "label_name"
const ParameterCategoryName = "Category Alert Engine"
const ParameterPenerima = "Send Notification"
const ParameterChannel = "Channel Notification"
const ParameterWaktuNotifikasi = "Reminder Alert Engine"
const ParameterValidation = "Validation"
const ParameterExpiredMedicalClaim = "expiredMedicalClaim"

// ------------------------------ Field Name Constanta Account Registration --------------------------------------------
const CategoryNameAccountRegistCompletion = "Account Registration Completion"
const CategoryNameAccountRegistDocChecklist = "Account Registration Document Checklist"
const CompanyTitleIDStruct = "CompanyTitleId"
const ProvinceIDStruct = "ProvinceID"
const DistrictIDStruct = "DistrictID"
const SubDistrictIDStruct = "SubDistrictID"
const UrbanVillageIDStruct = "UrbanVillageID"
const PostalCodeIDStruct = "PostalCodeID"
const PositionIDStruct = "PositionID"
const AddressStruct = "Address"
const FaxStruct = "Fax"
const CompanyEmailStruct = "CompanyEmail"
const ContactPhoneStruct = "ContactPhone"
const TaxNameStruct = "TaxName"
const NpwpStruct = "Npwp"
const TaxAddressStruct = "TaxAddress"
const FileNpwpStruct = "FileNPWP"
const FileSppkpStruct = "FileSPPKP"
const FileNPWPContentType = "npwp"
const FileSPPKPContentType = "sppkp"
const FileFieldName = "file"

// ----------------------------------------- Drop Down List Audit System -----------------------------------------------
const TableNameDDLName = "table_name"
const PrimaryKeyDDLName = "primary_key"
const SchemaDDLName = "schema_name"
const CreatedByDDLName = "created_by"
const CreatedClientDDLName = "created_client"
const FieldParamDDLName = "field_name"
const MenuCodeDDLName = "menu_code"

// --------------------------------------------- Alert Engine ----------------------------------------------------------
const FixAlert = "fix"
const CustomAlert = "custom"
const StatusComplete = "C"
const StatusInactive = "I"
const KeyProvinceIDValidation = "province_id"
const KeyDistrictIDValidation = "district_id"
const KeySubDistrictIDValidation = "sub_district_id"
const KeyUrbanVillageIDValidation = "urban_village_id"
const KeyPostalCodeIDValidation = "postal_code_id"

//--------------------------------------------- Unit Test ----------------------------------------------------------
const ContentType = "Content-Type"
const ApplicationJSON = "application/json"

// -------------------------------------------- Common Messages --------------------------------------------------------
const CommonSuccessInsertMessages = "SUCCESS_INSERT_MESSAGE"
const CommonSuccessUpdateMessages = "SUCCESS_UPDATE_MESSAGE"
const CommonSuccessDeleteMessages = "SUCCESS_DELETE_MESSAGE"
const CommonSuccessGetListMessages = "SUCCESS_GET_LIST_MESSAGE"
const CommonSuccessGetViewMessages = "SUCCESS_VIEW_MESSAGE"
const CommonSuccessGetInitiateMessages = "SUCCESS_INITIATE_MESSAGE"

// -------------------------------------- Status History Alert Front End -----------------------------------------------
const StatusNotSent = "Not Sent"
const StatusHold = "Hold"
const StatusSent = "Sent"
const StatusCancelled = "Cancelled"

// ---------------------------------------- Forget Password Messages ---------------------------------------------------
//const EmailResetPassword = "Hi {{.USER_ID}}," +
//	"Your Requested to reset the password for your Nextrac Account with the E-Mail Address ({{.EMAIL_LINK}})" +
//	"\nPlease click this link to reset the password." +
//	"\n{{.EMAIL_LINK}}" +
//	"\n\nRegards," +
//	"\nNextrac"

const EmailResetPassword = "Ini adalah pesan lupa password dari Nextrac\n." +
	"\nAlamat URL : {{.EMAIL_LINK}}, \n User ID : {{.USER_ID}} \n Kode Lupa password : {{.FORGET_CODE}}"

// ----------------------------------------------- Administrator -------------------------------------------------------
const AdminName = "admin"
const AdminClientID = "3e3cb40e14d645eb8783f53a30c822d4"

// ----------------------------------------------- URL Link Parameter --------------------------------------------------
const ClientIDQueryParam = "client_id"
const UniqueID1QueryParam = "unique_id_1"
const UniqueID2QueryParam = "unique_id_2"
const SalesmanIDQueryParam = "salesman_id"
const UserQueryParam = "user"
const UserIDQueryParam = "user_id"
const ActivationCodeQueryParam = "activation_code"
const EmailQueryParam = "email"
const OTPQueryParam = "otp"
const PhoneQueryParam = "no_telp"
const SubjectActivationEmail = "Activation"

// --------------------------------------------- Table Parameter Message -----------------------------------------------
const PurposeTableParam = "Purpose"
const NameTableParam = "Name"
const ClientTypeTableParam = "ClientType"
const CompanyIDTableParam = "UniqueID1"
const CompanyNameTableParam = "CompanyName"
const BranchIDTableParam = "UniqueID2"
const BranchNameTableParam = "BranchName"
const SalesmanIDTableParam = "SalesmanID"
const UserTableParam = "UserID"
const PasswordTableParam = "Password"
const OTPTableParam = "ACTIVATION_CODE"
const EmailTableParam = "Email"
const ClientIDTableParam = "ClientID"
const RegistrationIDTableParam = "USER_ID"
const LinkTableParam = "ACTIVATION_LINK"
const AuthenticationDataNotFound = "Data User tidak diketahui"
const PhoneMessageEmptyDefault = "not_send_otp"

// ------------------------------------------- Auth Message ------------------------------------------------------------
const ActivationCodeError = "E-4-AUT-SRV-004"
const EmailPhoneUnknownDataError = "E-4-AUT-SRV-003"
const ExpiredActivationCodeError = "E-4-AUT-SRV-005"
const FormatPhoneActivationCodeError = "E-4-AUT-DTO-007"
const OTPCodeActivationError = "E-4-AUT-SRV-007"
const ActivationCodeMessage = "Activation Code Salah"
const EmailUnknownMessage = "Data User & Email tidak diketahui"
const PhoneUnknownMessage = "Data User & Phone tidak diketahui"
const ExpiredActivationCodeMessage = "Activation Code sudah expired"
const FormatPhoneActivationCodeMessage = "Format Phone salah"
const OTPCodeActivationMessage = "Format OTPCode salah"

// ------------------------------------------- MDB Message ------------------------------------------------------------
const MasterDataUnknownDataErrorCode = "E-4-MAD-SRV-004"

// ---------------------------------------------- Role -----------------------------------------------------------------
const RoleIDUserND6 = "user_nd6"
const RoleIDUserNexmile = "user_nexmile"
const VersionRedesign = 3

// ---------------------------------------- Filter Key Report ----------------------------------------------------------
const FilterDepartment = "department"
const FilterEmployee = "employee"
const FilterSprint = "sprint"

// ---------------------------------------------- Backlog --------------------------------------------------------------
const StatusDone = "Done"
const StatusNotYet = "Not Yet"
const StatusNew = "New"
const StatusReadyToDev = "ReadyToDev"
const StatusOnProgress = "On Progress"
const StatusInProgress = "In Progress"
const StatusCompleteDev = "Completed Dev"
const StatusNeedMoreReq = "Need More Requirement"
const StatusReadyToTest = "Ready To Test"
const StatusClosed = "Closed"
const StatusReOpen = "ReOpen"
const StatusInTest = "In Test"
const StatusDoneTesting = "Done Testing"
const StatusReleased = "Released"
const Development = "Development"
const QAQC = "QA/QC"
const DepartmentDeveloper = "developer"
const DepartmentQAQC = "qaqc"
const DepartmentDevOps = "devops"
const DepartmentUIUX = "uiux"
const DepartmentInfra = "infra"
const TrackerTask = "Task"
const TrackerAuto = "Automation"
const TrackerManual = "Manual"
const DaysDefaultMandays = 8
const IDSprintOnRedmineCustomFields = 9
const IDPaymentOnRedmineCustomFields = 7
const ContainerBacklog = "backlog/attachment/"
const DeveloperDepartmentID = 1
const QAQCDepartmentID = 2
const InfraDepartmentID = 3
const DevOpsDepartmentID = 4
const UIUXDepartmentID = 5
const StartTicket = "start"
const PauseTicket = "pause"
const EndTicket = "end"

//Attacment Logo
const InternalCompanyAttachmentPrefix = "internal-company/attachment/"
const EmployeeAttachmentPrefix = "profile-employee/photo/"
const InternalCompanyAttachment = "COMPANY_ATTACHMENT"
const InternalCompanyAttachmentMaximumLogo = 1
const InternalCompanyAttachmentMaximumPhotoSize = 10485760

// ---------------------------------------------- Discord --------------------------------------------------------------
const TokenDiscord = "Bot MTE1MjEzNTM2OTY0MzU0MDUxMQ.Gh3jIo.JWXad8szTcWik2xZb5f14KZwMy8eDXQTEO4Rg4"

// ---------------------------------------------- Employee Request --------------------------------------------------------------
const PendingRequestStatus = "Pending"
const ApprovedRequestStatus = "Approved"
const RejectedRequestStatus = "Rejected"
const CancelledRequestStatus = "Cancelled"
const PendingCancellationRequestStatus = "Pending Cancellation"

// ---------------------------------------------- Employee Leave --------------------------------------------------------------
const PendingLeaveRequestType = "Pending"
const ApprovedLeaveRequestType = "Approved"
const RejectedLeaveRequestType = "Rejected"
const CancelledLeaveRequestType = "Cancelled"

const ContainerEmployeeLeave = "employee-leave/attachment/"
const ContainerReportEmployeeLeave = "employee-leave/report/"

const LeaveType = "leave"
const PermitType = "permit"
const SickLeaveType = "sick-leave"
const ReimbursementType = "reimbursement"

const LeaveAllowanceType = "leave"
const PermitAllowanceType = "permit"
const SickLeaveAllowanceType = "sick"

const AnnualLeaveAllowanceTypeKeyword = "annual leave"
const CutiTahunanAllowanceTypeKeyword = "cuti tahunan"

const EmployeeLeaveReportFileName = "Pengajuan_Izin_&_Cuti.xlsx"
const EmployeeLeaveReportSheetName = "Pengajuan Izin & Cuti"

const LeaveTypeAlias = "Cuti"
const PermitTypeAlias = "Izin"
const SickLeaveTypeAlias = "Sakit"

// ---------------------------------------------- Employee Reimbursement --------------------------------------------------------------
const PendingReimbursementRequestType = "Pending"
const ApprovedReimbursementRequestType = "Approved"
const RejectedReimbursementRequestType = "Rejected"

const VerifiedReimbursementVerification = "Verified"
const UnverifiedReimbursementVerification = "Unverified"
const PendingReimbursementVerification = "Pending"
const CancelReimbursementVerification = "Cancelled"

const ContainerEmployeeReimbursement = "employee-reimbursement/attachment/"

const MedicalBenefitType = "medical"

// ---------------------------------------------- Leave & Reimbursement Keywords --------------------------------------------------------------
const LeaveKeyword = "%leave%"
const CutiKeyword = "%cuti%"
const MedicalKeyword = "%medical%"

// ---------------------------------------------- Employee Notification --------------------------------------------------------------
const EmployeeRequestApprovedMessageTitle = "Your %s has been approved"
const EmployeeRequestApprovedMessageBody = "Your %s on %s has been approved"

const EmployeeRequestRejectedMessageTitle = "Your %s was rejected"
const EmployeeRequestRejectedMessageBody = "Your %s on %s was rejected"

const EmployeeCancelRequestRejectedMessageTitle = "Your request to cancel %s was rejected"
const EmployeeCancelRequestRejectedMessageBody = "Your request to cancel %s on %s was rejected"

const EmployeeCancelRequestApprovedMessageTitle = "Your request to cancel %s has been approved"
const EmployeeCancelRequestApprovedMessageBody = "Your request to cancel %s on %s has been approved"

const EmployeeRequestPendingMessageTitle = "Request for approval"
const EmployeeRequestPendingMessageBody = "%s sent a %s on %s"

const EmployeeCancelRequestPendingMessageTitle = "Request for cancellation"
const EmployeeCancelRequestPendingMessageBody = "%s sent a %s cancellation on %s"

const EmployeeRequestVerifiedMessageTitle = "Your %s has been verified"
const EmployeeRequestVerifiedMessageBody = "Your %s on %s has been verified"

const LeaveTypeNotification = "leave"
const PermitTypeNotification = "permission to leave early"
const SickLeaveTypeNotification = "sick leave"
const MedicalTypeNotification = "medical claim"

// ---------------------------------------------- GroChat --------------------------------------------------------------
const GroChatUserTypeRegular = "R"

const InvitationEmailBody = `<p>Dear</p><br><p>Anda telah diundang untuk bergabung dengan NexTrac. Silahkan klik tombol dibawah ini untuk melanjutkan : </p><a href="{{.INVITATION_LINK}}"><button type="button">Lanjut</button></a>`

// ---------------------------------------------- GroChat WS --------------------------------------------------------------
const MessageTypeChat = "chat"

const ChatTypeSend = "send"
const ChatTypeAck = "ack"

// ---------------------------------------------- HRIS Email --------------------------------------------------------------
const HRISSubject = "Notifikasi HRIS"
const RequestApprovalHRISSubject = "Request Approval HRIS"
const CancelRequestApprovalHRISSubject = "Cancel Request Approval HRIS"

const RequestApprovedEmailBody = `Dear %s

Pengajuan %s pada tanggal %s kamu telah disetujui. Untuk mengetahui detail dari pengajuan silahkan lihat pada NexSoft Apps.

Terima Kasih,
Hormat Kami.

NexSoft System
`

const RequestRejectedEmailBody = `Dear %s

Pengajuan %s pada tanggal %s kamu ditolak. Untuk mengetahui detail dari pengajuan silahkan lihat pada NexSoft Apps.

Terima Kasih,
Hormat Kami.

NexSoft System
`

const RequestVerifiedEmailBody = `Dear %s

Pengajuan %s pada tanggal %s kamu telah diverifikasi. Untuk mengetahui detail dari pengajuan silahkan lihat pada NexSoft Apps.

Terima Kasih,
Hormat Kami.

NexSoft System
`

const LeaveApprovalRequestEmailBody = `Dear %s

%s mengajukan %s untuk tanggal %s. Untuk melakukan approval silahkan masuk ke NexSoft App > HRIS > Approval.

Terima Kasih,
Hormat Kami.

NexSoft System
`

const ReimbursementApprovalRequestEmailBody = `Dear %s

%s mengajukan %s pada tanggal %s. Untuk melakukan approval silahkan masuk ke NexSoft App > HRIS > Approval.

Terima Kasih,
Hormat Kami.

NexSoft System
`

const CancellationRequestEmailBody = `Dear %s

%s mengajukan pembatalan %s pada tanggal %s. Untuk menyetujui pembatalan silahkan masuk ke NexSoft App > HRIS > Approval.

Terima Kasih,
Regards

NexSoft System
`

const CancellationRequestApprovedEmailBody = `Dear %s

Pengajuan pembatalan %s pada tanggal %s telah disetujui. Untuk mengetahui detail dari pengajuan silahkan lihat pada NexSoft Apps.

Terima Kasih,
Regards

NexSoft System
`

const CancellationRequestRejectedEmailBody = `Dear %s

Pengajuan pembatalan %s pada tanggal %s ditolak. Untuk mengetahui detail dari pengajuan silahkan lihat pada NexSoft Apps.

Terima Kasih,
Regards

NexSoft System
`

const TypeLeaveEmail = "Cuti"
const TypePermitEmail = "Izin Keluar"
const TypeSickLeaveEmail = "Izin Sakit"
const TypeMedicalEmail = "Medical Claim"

// -------------------------------------------------------------- Platform Devices --------------------------------------------------------------
const PlatformWebsite = "Website"

// -------------------------------------------------------------- Currencies --------------------------------------------------------------
const CurrencyIDR = "IDR"