package store

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/wetee-dao/tee-dsecret/chains/pallets/generated/types"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"wetee.app/worker/util"
)

func SealAppID(WorkID types.WorkId) (string, error) {
	// Add timestamps to prevent id hijacking and misuse
	// 添加时间戳防止id被劫持滥用
	key := util.GetWorkTypeStr(WorkID) + "-" + fmt.Sprint(WorkID.Id) + "-" + fmt.Sprint(time.Now().Unix())
	var val []byte

	val, err := model.SealWithProductKey([]byte(key), nil)
	if err != nil {
		return "", err
	}

	strVal := url.QueryEscape(base64.StdEncoding.EncodeToString(val))
	return strVal, nil
}

func UnSealAppID(id string) (types.WorkId, error) {
	var err error
	id, err = url.QueryUnescape(id)
	if err != nil {
		return types.WorkId{}, err
	}
	buf, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return types.WorkId{}, err
	}

	var val []byte
	val, err = model.Unseal(buf, nil)
	if err != nil {
		return types.WorkId{}, err
	}

	str := string(val)
	strs := strings.Split(str, "-")
	if len(strs) != 3 {
		return types.WorkId{}, fmt.Errorf("invalid id")
	}

	wid, err := strconv.ParseUint(strs[1], 10, 64)
	if err != nil {
		return types.WorkId{}, err
	}
	return types.WorkId{
		Wtype: util.GetWorkType(strs[0]),
		Id:    wid,
	}, nil
}
