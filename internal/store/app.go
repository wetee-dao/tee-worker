package store

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

func SealAppID(PodId uint64) (string, error) {
	// Add timestamps to prevent id hijacking and misuse
	// 添加时间戳防止id被劫持滥用
	key := fmt.Sprint(PodId) + "-" + fmt.Sprint(time.Now().Unix())
	var val []byte

	val, err := model.SealWithProductKey([]byte(key), nil)
	if err != nil {
		return "", err
	}

	strVal := url.QueryEscape(base64.StdEncoding.EncodeToString(val))
	return strVal, nil
}

func UnSealAppID(id string) (uint64, error) {
	var err error
	id, err = url.QueryUnescape(id)
	if err != nil {
		return 0, err
	}
	buf, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return 0, err
	}

	var val []byte
	val, err = model.Unseal(buf, nil)
	if err != nil {
		return 0, err
	}

	str := string(val)
	strs := strings.Split(str, "-")
	if len(strs) != 3 {
		return 0, fmt.Errorf("invalid id")
	}

	return strconv.ParseUint(strs[1], 10, 64)
}
