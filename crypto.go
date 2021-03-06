package jadepoolsaas

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

func signHMACSHA256(data interface{}, secret string) (string, error) {
	buf, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	decoder := json.NewDecoder(bytes.NewReader(buf))
	decoder.UseNumber()
	obj := make(map[string]interface{})
	err = decoder.Decode(&obj)
	if err != nil {
		return "", err
	}

	msgStr := buildMsg(obj, "=", "&")
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(msgStr))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha, nil
}

func buildMsg(val interface{}, keyValSeparator, groupSeparator string) string {
	if val == nil {
		return ""
	}

	msg := ""
	switch reflect.TypeOf(val).Kind() {
	case reflect.Struct:
		buf, err := json.Marshal(val)
		if err != nil {
			return ""
		}
		decoder := json.NewDecoder(bytes.NewReader(buf))
		decoder.UseNumber()
		m := make(map[string]interface{})
		err = decoder.Decode(&m)
		if err != nil {
			return ""
		}
		msg = buildMsg(m, keyValSeparator, groupSeparator)
	case reflect.Map:
		buf, err := json.Marshal(val)
		if err != nil {
			return ""
		}
		decoder := json.NewDecoder(bytes.NewReader(buf))
		decoder.UseNumber()
		obj := make(map[string]interface{})
		err = decoder.Decode(&obj)

		keyVals := make(map[string]string)
		keys := make([]string, 0, len(obj))

		for k, v := range obj {
			_msg := buildMsg(v, keyValSeparator, groupSeparator)
			keyVals[k] = _msg
			keys = append(keys, k)
		}
		sort.Strings(keys)
		groupStrs := make([]string, 0, len(keys))
		for _, key := range keys {
			groupStrs = append(groupStrs, key+keyValSeparator+keyVals[key])
		}
		msg += strings.Join(groupStrs, groupSeparator)
	case reflect.Slice:
		arr := val.([]interface{})
		keyVals := make(map[string]string)
		keys := make([]string, 0, len(arr))

		for i, v := range arr {
			key := strconv.Itoa(i)
			keys = append(keys, key)
			keyVals[key] = buildMsg(v, keyValSeparator, groupSeparator)
		}
		sort.Strings(keys)

		groupStrs := make([]string, 0, len(keys))
		for _, key := range keys {
			groupStrs = append(groupStrs, key+keyValSeparator+keyVals[key])
		}
		msg += strings.Join(groupStrs, groupSeparator)
	default:
		msg = fmt.Sprintf("%v", val)
	}

	return msg
}

func aesEncryptStr(src string, key, iv []byte) (encmess string, err error) {
	ciphertext, err := aesEncrypt([]byte(src), key, iv)
	if err != nil {
		return
	}

	encmess = base64.StdEncoding.EncodeToString(ciphertext)
	return
}

func aesEncrypt(src []byte, key []byte, iv []byte) ([]byte, error) {
	if len(iv) == 0 {
		iv = key[:16]
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	src = padding(src, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, iv)
	blockMode.CryptBlocks(src, src)
	return src, nil
}

func padding(src []byte, blockSize int) []byte {
	padNum := blockSize - len(src)%blockSize
	pad := bytes.Repeat([]byte{byte(padNum)}, padNum)
	return append(src, pad...)
}

func aesDecryptStr(src string, key, iv []byte) (string, error) {
	bsrc, err := base64.StdEncoding.DecodeString(src)
	bret, err := aesDecrypt(bsrc, key, iv)
	if err != nil {
		return "", err
	}
	return string(bret), nil
}

func aesDecrypt(src []byte, key []byte, iv []byte) ([]byte, error) {
	if len(iv) == 0 {
		iv = key[:16]
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	blockMode.CryptBlocks(src, src)
	src = unpadding(src)
	return src, nil
}

func unpadding(src []byte) []byte {
	n := len(src)
	unPadNum := int(src[n-1])
	return src[:n-unPadNum]
}
