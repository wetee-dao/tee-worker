package store

import (
	"fmt"
	"strconv"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

const MinTimeKey = "last_mint_"

// set pod last mint time
func SetLastMintTime(id string, value int64) error {
	val := fmt.Append(nil, value)
	return model.SetKey(PodBucket, MinTimeKey+id, val)
}

// get pod last mint time
func GetLastMintTime(id string) (int64, error) {
	val, err := model.GetKey(PodBucket, MinTimeKey+id)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(string(val), 10, 64)
}
