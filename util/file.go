package util

import (
	"encoding/base64"
	"os"
)

func IsFileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func readFileBase4(filename string) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		LogWithRed("read file error", err)
		return ""
	}
	return base64.StdEncoding.EncodeToString(data)
}

func GetSslRoot() []string {
	if !IsFileExists(WORK_DIR+"/ser.pem") || !IsFileExists(WORK_DIR+"/ser.key") {
		return []string{"", ""}
	}

	return []string{readFileBase4(WORK_DIR + "/ser.key"), readFileBase4(WORK_DIR + "/ser.pem")}
}
