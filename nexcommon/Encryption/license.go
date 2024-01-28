package Encryption

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)


type ProductConfiguration struct {
	ProductKey string
}


func AESEncrypt(plainText string, key string, productConfiguration ProductConfiguration) []byte {
	var stdout []byte
	var err error
	path := os.Getenv("GOPATH")
	paths := strings.Split(path, ";")
	final := strings.Replace(paths[0], "\\", "/", -1)
	if runtime.GOOS == "windows" {
		stdout, err = exec.Command( final+"/src/nexsoft.co.id/nexcommon/Encryption/nextrack-library-go.exe", "AESEncrypt", plainText, key, key).Output()
	}
	//else {
	//	stdout, err = exec.Command("bash", "-c", "git config user.name").Output()
	//}

	//cmd := exec.Command(final+"/src/nexsoft.co.id/nexcommon/Encryption/nextrack-library-go.exe", "AESEncrypt", plainText, key)
	//stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return stdout
}

func AESDecrypt(plainText string, key string, productConfiguration string) []byte {
	var stdout []byte
	path := os.Getenv("GOPATH")
	paths := strings.Split(path,";")
	final := strings.Replace(paths[0],"\\","/",-1)
	if runtime.GOOS == "windows" {
		stdout, _ = exec.Command(final+"/src/nexsoft.co.id/nexcommon/Encryption/nextrack-library-go.exe", "AESDecrypt", plainText, key).Output()
	}
	return stdout
}


