package params

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	stdUrl "net/url"
)

func Gen(url string, fileName string, params Params, key string) (string, error) {
	p, err := EncParams(params, key)
	if err != nil {
		return "", err
	}
	v := stdUrl.Values{}
	v.Set("p", p)
	if len(fileName) > 0 {
		h := md5.New()
		h.Write([]byte(fileName))
		v.Set("f", hex.EncodeToString(h.Sum(nil)))
	}
	return fmt.Sprintf("%s?%s", url, v.Encode()), nil
}
