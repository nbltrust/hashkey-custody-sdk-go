package jadepoolsaas

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"golang.org/x/crypto/sha3"
	"math/big"
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

func parseEcdsaPrivateKeyFromPem(pemContent []byte) (*btcec.PrivateKey, error) {
	block, _ := pem.Decode(pemContent)
	if block == nil {
		return nil, errors.New("invalid pem")
	}

	var ecp ecPrivateKey
	_, err := asn1.Unmarshal(block.Bytes, &ecp)
	if err != nil {
		return nil, err
	}

	priKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), ecp.PrivateKey)
	return priKey, nil
}

func signECCDataStr(priKeyPemBase64 string, msgStr string, sigEncode string) (*eccSig, error) {
	priKeyPem, err := base64.StdEncoding.DecodeString(priKeyPemBase64)
	if err != nil {
		return nil, err
	}
	priKey, err := parseEcdsaPrivateKeyFromPem(priKeyPem)
	if err != nil {
		return nil, err
	}

	sha3Hash := sha3.NewLegacyKeccak256()
	_, err = sha3Hash.Write([]byte(msgStr))
	if err != nil {
		return nil, err
	}
	msgBuf := sha3Hash.Sum(nil)
	sig, err := priKey.Sign(msgBuf)
	if err != nil {
		return nil, err
	}

	r := sig.R.Bytes()
	s := sig.S.Bytes()
	if len(r) < 32 {
		preArr := []byte{}
		for i := len(r) + 1; i <= 32; i++ {
			preArr = append(preArr, 0)
		}
		r = append(preArr, r...)
	}
	if len(s) < 32 {
		preArr := []byte{}
		for i := len(s) + 1; i <= 32; i++ {
			preArr = append(preArr, 0)
		}
		s = append(preArr, s...)
	}
	_sig := &eccSig{}
	if sigEncode == "hex" {
		_sig.R = hex.EncodeToString(r)
		_sig.S = hex.EncodeToString(s)
	} else if sigEncode == "base64" {
		_sig.R = base64.StdEncoding.EncodeToString(r)
		_sig.S = base64.StdEncoding.EncodeToString(s)
	}
	return _sig, nil
}

func parseEcdsaPubKeyFromPem(pemContent []byte) (*btcec.PublicKey, error) {
	block, _ := pem.Decode(pemContent)
	if block == nil {
		return nil, errors.New("invalid pem")
	}

	var ecp ecPublicKey
	_, err := asn1.Unmarshal(block.Bytes, &ecp)
	if err != nil {
		return nil, err
	}

	return btcec.ParsePubKey(ecp.PublicKey.RightAlign(), btcec.S256())
}

func verifyECCSign(pubKeyPemBase64 string, obj map[string]interface{}, sign *eccSig, sigEncode string) (bool, error) {
	pubKeyPem, err := base64.StdEncoding.DecodeString(pubKeyPemBase64)
	if err != nil {
		return false, err
	}

	pubKey, err := parseEcdsaPubKeyFromPem(pubKeyPem)
	if err != nil {
		return false, err
	}

	msgStr := buildMsg(obj, "", "")
	return verifyECCSignStr(msgStr, sign, pubKey, sigEncode)
}

func verifyECCSignStr(msgStr string, sign *eccSig, pubKey *btcec.PublicKey, sigEncode string) (bool, error) {
	sha3Hash := sha3.NewLegacyKeccak256()
	_, err := sha3Hash.Write([]byte(msgStr))
	if err != nil {
		return false, err
	}
	msgBuf := sha3Hash.Sum(nil)

	var decodedR, decodedS []byte
	if sigEncode == "hex" {
		decodedR, err = hex.DecodeString(sign.R)
		if err != nil {
			return false, err
		}
		decodedS, err = hex.DecodeString(sign.S)
		if err != nil {
			return false, err
		}
	} else if sigEncode == "base64" {
		decodedR, err = base64.StdEncoding.DecodeString(sign.R)
		if err != nil {
			return false, err
		}
		decodedS, err = base64.StdEncoding.DecodeString(sign.S)
		if err != nil {
			return false, err
		}
	}

	signature := btcec.Signature{
		R: new(big.Int).SetBytes(decodedR),
		S: new(big.Int).SetBytes(decodedS),
	}
	return signature.Verify(msgBuf, pubKey), nil
}

func unpadding(src []byte) []byte {
	n := len(src)
	unPadNum := int(src[n-1])
	return src[:n-unPadNum]
}

//This type provides compatibility with the btcec package
type ecPrivateKey struct {
	Version       int
	PrivateKey    []byte
	NamedCurveOID asn1.ObjectIdentifier `asn1:"optional,explicit,tag:0"`
	PublicKey     asn1.BitString        `asn1:"optional,explicit,tag:1"`
}

//This type provides compatibility with the btcec package
type ecPublicKey struct {
	Raw       asn1.RawContent
	Algorithm pkix.AlgorithmIdentifier
	PublicKey asn1.BitString
}

type eccSig struct {
	R string `json:"r"`
	S string `json:"s"`
}
