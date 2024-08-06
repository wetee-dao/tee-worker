package store

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/wetee-dao/go-sdk/pallet/types"
	"wetee.app/worker/util"
)

const SecretBucket = "secret"

type Secrets struct {
	Files map[string]string
	Env   map[string]string
}

func SealAppID(WorkID types.WorkId) (string, error) {
	// Add timestamps to prevent id hijacking and misuse
	// 添加时间戳防止id被劫持滥用
	key := util.GetWorkTypeStr(WorkID) + "-" + fmt.Sprint(WorkID.Id) + "-" + fmt.Sprint(time.Now().Unix())
	var val []byte

	val, err := SealWithProductKey([]byte(key), nil)
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
	val, err = Unseal(buf, nil)
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

func SetSecrets(id types.WorkId, secrets *Secrets) error {
	key := []byte(util.GetWorkTypeStr(id) + "-" + fmt.Sprint(id.Id))
	val, err := json.Marshal(secrets)
	if err != nil {
		return err
	}

	return SealSave(SecretBucket, key, val)
}

func GetSecrets(id types.WorkId) (*Secrets, error) {
	key := []byte(util.GetWorkTypeStr(id) + "-" + fmt.Sprint(id.Id))
	val, err := SealGet(SecretBucket, key)
	if err != nil {
		return nil, err
	}
	s := &Secrets{}
	err = json.Unmarshal(val, s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// func GetSetAppSignerAddress(id types.WorkId, address string) (string, error) {
// 	key := []byte(util.GetWorkTypeStr(id) + "-" + fmt.Sprint(id.Id) + "-" + "signer")
// 	val := []byte(address)

// 	// 加密数据
// 	val, errr := SealWithProductKey(val, nil)
// 	if errr != nil {
// 		return "", errr
// 	}

// 	var data []byte = []byte{}
// 	err := DB.Update(
// 		func(tx *nutsdb.Tx) error {
// 			// 检查是否存在bucket
// 			if !tx.ExistBucket(nutsdb.DataStructureBTree, SecretBucket) {
// 				if err := tx.NewBucket(nutsdb.DataStructureBTree, SecretBucket); err != nil {
// 					return err
// 				}
// 				tx.SubmitBucket()
// 			}

// 			// 如果不存在则写入数据
// 			err := tx.PutIfNotExists(SecretBucket, key, val, 0)
// 			if err != nil {
// 				return err
// 			}

// 			// 获取数据
// 			val, err := tx.Get(SecretBucket, key)
// 			if err != nil {
// 				return err
// 			}

// 			// 解析数据
// 			val, err = Unseal(val, nil)
// 			if err != nil {
// 				return err
// 			}
// 			data = val
// 			return nil
// 		},
// 	)
// 	if err != nil {
// 		return "", err
// 	}
// 	return string(data), nil
// }
