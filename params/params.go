package params

import (
	"crypto/aes"
	"encoding/base64"
	"encoding/json"
	"errors"
	stdUrl "net/url"
)

type Params struct {
	Tag    int    `json:"tag"`    // 标签
	Url    string `json:"url"`    // 链接
	Expire int64  `json:"expire"` // 过期
}

func EncParams(params Params, key string) (string, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return "", err
	}
	enc, err := enc(string(data), key)
	if err != nil {
		return "", err
	}
	return enc, nil
}

func DecParams(data string, key string) (Params, error) {
	str, err := dec(data, key)
	if err != nil {
		return Params{}, err
	}
	content := Params{}
	if err := json.Unmarshal([]byte(str), &content); err != nil {
		return Params{}, err
	}
	if len(content.Url) == 0 {
		return Params{}, errors.New("链接不存在")
	}
	targetUrl, err := stdUrl.Parse(content.Url)
	if err != nil {
		return Params{}, err
	}
	content.Url = targetUrl.String()
	return content, nil
}

func enc(data string, key string) (string, error) {
	bs, err := aesEncECB([]byte(data), []byte(key))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bs), nil
}

func dec(data string, key string) (string, error) {
	bs, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	if len(bs) == 0 {
		return "", errors.New("data is nil")
	}
	bs, err = aesDecECB(bs, []byte(key))
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func aesEncECB(origData []byte, key []byte) (encrypted []byte, err error) {
	cipher, err := aes.NewCipher(generateKey(key))
	if err != nil {
		return nil, err
	}
	length := (len(origData) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, origData)
	pad := byte(len(plain) - len(origData))
	for i := len(origData); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted = make([]byte, len(plain))
	for bs, be := 0, cipher.BlockSize(); bs <= len(origData); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Encrypt(encrypted[bs:be], plain[bs:be])
	}
	return encrypted, nil
}

func aesDecECB(encrypted []byte, key []byte) (decrypted []byte, err error) {
	cipher, err := aes.NewCipher(generateKey(key))
	if err != nil {
		return nil, err
	}
	decrypted = make([]byte, len(encrypted))
	for bs, be := 0, cipher.BlockSize(); bs < len(encrypted); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}
	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}
	return decrypted[:trim], nil
}

func generateKey(key []byte) (genKey []byte) {
	genKey = make([]byte, 16)
	copy(genKey, key)
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}
