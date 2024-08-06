package bridge

import (
	"github.com/go-resty/resty/v2"
	gtypes "github.com/wetee-dao/go-sdk/pallet/types"
)

func CallTeeApp(call *gtypes.TEECall, meta *gtypes.ApiMeta, body []byte) ([]byte, error) {
	client := resty.New()
	req := client.R().SetBody(call.Args)

	api := meta.Apis[call.Method]
	url := string(api.Url)

	// 0: get, 1: post, 2: put, 3: delete
	switch api.Method {
	case 0:
		resp, err := req.Get(url)
		if err != nil {
			return nil, err
		}
		return resp.Body(), nil
	case 1:
		resp, err := req.Post(url)
		if err != nil {
			return nil, err
		}
		return resp.Body(), nil
	case 2:
		resp, err := req.Put(url)
		if err != nil {
			return nil, err
		}
		return resp.Body(), nil
	case 3:
		resp, err := req.Delete(url)
		if err != nil {
			return nil, err
		}
		return resp.Body(), nil
	default:
		return nil, nil
	}
}
