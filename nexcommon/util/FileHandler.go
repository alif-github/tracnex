package util

import (
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
)

var suffixes = []string{"B", "KB", "MB", "GB", "TB"}

func MoveFile(oldPath string, newPath string) (err error) {
	newPath = strings.Replace(newPath, "\\", "/", -1)
	newPathSplit := strings.Split(newPath, "/")
	newPathMkdir := strings.Join(newPathSplit[0:len(newPathSplit)-1], "/")
	_ = os.MkdirAll(newPathMkdir, 0770)
	err = os.Rename(oldPath, newPath)
	if err != nil {
		var byteArr []byte
		byteArr, err = ioutil.ReadFile(oldPath)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(newPath, byteArr, 0660)
		if err == nil {
			_ = os.Remove(oldPath)
		} else {
			_ = os.Remove(newPath)
		}

		return err
	}

	return err
}

func roundByte(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

func ConvertByteToSize(byteSize int64) string {
	if byteSize != 0 {
		base := math.Log(float64(byteSize)) / math.Log(1024)
		getSize := roundByte(math.Pow(1024, base-math.Floor(base)), .5, 2)
		suffixIdx := int(math.Floor(base))
		getSuffix := suffixes[suffixIdx]
		return strconv.FormatFloat(getSize, 'f', -1, 64) + " " + getSuffix
	}
	return "0 B"
}
