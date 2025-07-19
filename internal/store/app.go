package store

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/wetee-dao/tee-dsecret/pkg/util"
)

// TEE 封装 appid
func SealAppID(PodId uint64) (string, error) {
	// Add timestamps to prevent id hijacking and misuse
	// 添加时间戳防止id被劫持滥用
	key := fmt.Sprint(PodId) + "-" + fmt.Sprint(time.Now().Unix())
	var val []byte

	val, err := util.SealWithProductKey([]byte(key), nil)
	if err != nil {
		return "", err
	}

	strVal := url.QueryEscape(base64.StdEncoding.EncodeToString(val))
	return strVal, nil
}

// TEE 解封 appid
func UnSealAppID(id string) (uint64, int64, error) {
	var err error
	id, err = url.QueryUnescape(id)
	if err != nil {
		return 0, 0, err
	}

	buf, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return 0, 0, err
	}

	var val []byte
	val, err = util.Unseal(buf, nil)
	if err != nil {
		return 0, 0, err
	}

	str := string(val)
	strs := strings.Split(str, "-")
	if len(strs) != 2 {
		return 0, 0, fmt.Errorf("invalid id")
	}

	t, err := strconv.ParseInt(strs[1], 10, 64)
	if err != nil {
		return 0, 0, err
	}

	podId, err := strconv.ParseUint(strs[0], 10, 64)
	return podId, t, err
}
