package util

import "nexsoft.co.id/nexcommon/util"

func HashingPassword(password string, salt string) string {
	return util.CheckSumWithSha512([]byte(password + salt))
}

func CheckIsPasswordMatch(passwordInput string, passwordDB string, salt string) bool {
	return HashingPassword(passwordInput, salt) == passwordDB
}
