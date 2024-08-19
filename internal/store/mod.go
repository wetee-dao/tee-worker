package store

import (
	"encoding/base64"
	"flag"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/edgelesssys/ego/ecrypto"
	"github.com/nutsdb/nutsdb"
	"github.com/wetee-dao/go-sdk/pallet/types"
	"wetee.app/worker/util"
)

var DB *nutsdb.DB

func DBInit(path string) error {
	var err error
	DB, err = nutsdb.Open(
		nutsdb.DefaultOptions,
		nutsdb.WithDir(path),
	)

	return err
}

func DBClose() {
	DB.Close()
}

func SealSave(bucket string, key []byte, val []byte) error {
	val, err := SealWithProductKey(val, nil)
	if err != nil {
		return err
	}

	err = checkBucket(bucket, nutsdb.DataStructureBTree)
	if err != nil {
		return err
	}

	return DB.Update(
		func(tx *nutsdb.Tx) error {
			err := tx.Put(bucket, key, val, 0)
			return err
		},
	)
}

func SealGet(bucket string, key []byte) ([]byte, error) {
	var data []byte = []byte{}
	err := checkBucket(bucket, nutsdb.DataStructureBTree)
	if err != nil {
		return nil, err
	}
	err = DB.View(
		func(tx *nutsdb.Tx) error {
			val, err := tx.Get(bucket, key)
			if err != nil {
				return err
			}

			if flag.Lookup("test.v") != nil {
				data = val
			} else {
				val, err = ecrypto.Unseal(val, nil)
				if err != nil {
					return err
				}
			}

			data = val
			return nil
		},
	)
	return data, err
}

func checkBucket(bucket string, ds uint16) error {
	return DB.Update(
		func(tx *nutsdb.Tx) error {
			if !tx.ExistBucket(ds, bucket) {
				err := tx.NewBucket(ds, bucket)
				return err
			}
			return nil
		},
	)
}

func SealWithProductKey(val []byte, additionalData []byte) ([]byte, error) {
	if flag.Lookup("test.v") == nil {
		return ecrypto.SealWithProductKey(val, additionalData)
	}
	return val, nil
}

func Unseal(ciphertext []byte, additionalData []byte) ([]byte, error) {
	if flag.Lookup("test.v") == nil {
		return ecrypto.Unseal(ciphertext, additionalData)
	}
	return ciphertext, nil
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
