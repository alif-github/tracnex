package common

const ReadDataAPIMustHave = "read"
const WriteDataAPIMustHave = "write"
const InsertData = "insert"
const ViewData = "view"
const DeleteData = "delete"
const UpdateData = "update"
const UpdateDataOwn = "update-own"
const ChangePassword = "changepassword"
const UpdateDataPermissionMustHave = ":" + UpdateData
const UpdateDataPermissionMustHaveOwn = ":" + UpdateDataOwn
const DeleteDataPermissionMustHave = ":" + DeleteData
const DeleteOwnDataPermissionMustHave = ":delete-own"
const DuplicateDataPermissionMustHave = ":duplicate"
const SyncDataPermissionMustHave = ":synchronize"
const InsertDataPermissionMustHave = ":" + InsertData
const ViewDataPermissionMustHave = ":" + ViewData
const ChangePasswordPermissionMustHave = ":changepassword"
