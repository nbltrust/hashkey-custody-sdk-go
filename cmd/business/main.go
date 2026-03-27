package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/btcsuite/btcd/btcec"
	"github.com/labstack/echo/v4"
)

// ParseEcdsaPubKeyFromPem ...
func ParseEcdsaPubKeyFromPem(pemContent []byte) (*btcec.PublicKey, error) {
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

// ParseEcdsaPrivateKeyFromPem ...
func ParseEcdsaPrivateKeyFromPem(pemContent []byte) (*btcec.PrivateKey, error) {
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

func genShareKey() string {
	apiGatewayPubKey := string(`LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUZZd0VBWUhLb1pJemowQ0FRWUZLNEVFQUFvRFFnQUVyb0QxMG53SzJkcElqYSszb1pncGNtbk40MWFCc0FFSQpsSE9MMEpublJad3pRa3k5cmFDSW5iSk9YcGtzcDBFUVZqdDVkdkJUMEw3b2pXQXFVSlk3b1E9PQotLS0tLUVORCBQVUJMSUMgS0VZLS0tLS0K`)
	myPriKey := string(`LS0tLS1CRUdJTiBFQyBQUklWQVRFIEtFWS0tLS0tCk1IUUNBUUVFSUNrOGw2UFJ6QTRhWjVVZkl1aXViaGplUDJBSTEvVENJajE2OUwzWU9xekFvQWNHQlN1QkJBQUsKb1VRRFFnQUVvVVpwRXNCRFpWTC9TQ29DanlreXAwdXRNMDc2b29YNUU2eEtzQW5TZlpOMFEwM3VlbGVVL09aeAovMmxsUXFBdU9aVlFLNE9aSGdFODh1c3RWVkY3YWc9PQotLS0tLUVORCBFQyBQUklWQVRFIEtFWS0tLS0tCg==`)

	apiGatewayPubKeyPem, err := base64.StdEncoding.DecodeString(apiGatewayPubKey)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	pubKey, err := ParseEcdsaPubKeyFromPem(apiGatewayPubKeyPem)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	priKeyPem, err := base64.StdEncoding.DecodeString(myPriKey)
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
	priKey, err := ParseEcdsaPrivateKeyFromPem(priKeyPem)
	if err != nil {
		fmt.Println(err)
		os.Exit(4)
	}

	aesKey := btcec.GenerateSharedSecret(priKey, pubKey)
	fmt.Println(string(aesKey))

	return base64.StdEncoding.EncodeToString(aesKey)
}

type encryptedReqBody struct {
	Encrypted string `json:"encrypted"`
}

type encryptedResBody struct {
	Encrypted string `json:"encrypted"`
	Iv        string `json:"iv"`
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)

	if err != nil {
		return nil, err
	}

	return b, nil
}

func encryptResBody(msg, aesKey string) ([]byte, []byte, error) {
	iv, err := GenerateRandomBytes(16)
	if err != nil {
		return nil, nil, errors.New("iv generaton error")
	}

	mKey, err := base64.StdEncoding.DecodeString(aesKey)
	if err != nil {
		return nil, nil, errors.New("invalid aes key")
	}

	encrypted, err := AESEncryptStr(msg, mKey, iv)
	if err != nil {
		return nil, nil, err
	}
	return []byte(encrypted), iv, nil
}

func decryptReqBody(encrypted *encryptedReqBody, aesKey, iv string) ([]byte, error) {
	aesIV, err := base64.StdEncoding.DecodeString(iv)
	if err != nil || len(aesIV) != 16 {
		return nil, errors.New("invalid aes iv")
	}

	mKey, err := base64.StdEncoding.DecodeString(aesKey)
	if err != nil {
		return nil, errors.New("invalid aes key")
	}

	decrypted, err := AESDecryptStr(encrypted.Encrypted, mKey, aesIV)
	if err != nil {
		return nil, err
	}
	return []byte(decrypted), nil
}

// AESEncryptStr base64加密字符串
func AESEncryptStr(src string, key, iv []byte) (encmess string, err error) {
	ciphertext, err := AESEncrypt([]byte(src), key, iv)
	if err != nil {
		return
	}

	encmess = base64.StdEncoding.EncodeToString(ciphertext)
	return
}

// AESEncrypt 加密
func AESEncrypt(src []byte, key []byte, iv []byte) ([]byte, error) {
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

// 填充数据
func padding(src []byte, blockSize int) []byte {
	padNum := blockSize - len(src)%blockSize
	pad := bytes.Repeat([]byte{byte(padNum)}, padNum)
	return append(src, pad...)
}

// 去掉填充数据
func unpadding(src []byte) []byte {
	n := len(src)
	unPadNum := int(src[n-1])
	return src[:n-unPadNum]
}

func AESDecryptStr(src string, key, iv []byte) (string, error) {
	bsrc, err := base64.StdEncoding.DecodeString(src)
	bret, err := AESDecrypt(bsrc, key, iv)
	if err != nil {
		return "", err
	}
	return string(bret), nil
}

// AESDecrypt 解密
func AESDecrypt(src []byte, key []byte, iv []byte) ([]byte, error) {
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

func main() {
	aesKey := genShareKey()

	e := echo.New()
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("/pong", func(c echo.Context) error {
		requestBody, _ := ioutil.ReadAll(c.Request().Body)
		fmt.Println(string(requestBody))

		for key, value := range c.Request().Header {
			fmt.Print(key)

			for _, v := range value {
				fmt.Printf("\t%v", v)
			}

			fmt.Println()
		}

		iv := c.Request().Header.Get("X-Encrypt-Iv")
		userId := c.Request().Header.Get("X-User-Id")
		var encrypted encryptedReqBody
		err := json.Unmarshal(requestBody, &encrypted)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}

		if err == nil && len(encrypted.Encrypted) > 0 {
			fmt.Printf("request msg: %v\n", encrypted.Encrypted)
			// fmt.Printf("%v\n", aesKey)
			fmt.Printf("request iv: %v\n", iv)

			request_plain_msg, err := decryptReqBody(&encrypted, aesKey, iv)
			if err != nil {
				fmt.Printf("err: %v\n", err)
			}

			fmt.Printf("request decrypted msg: %v\n", string(request_plain_msg))
		}

		timeStr := time.Now().Format("2006-01-02 15:04:05")
		response_plain_msg := fmt.Sprintf(`{"kkk":"hahaha","userID":"%v","ts":%v}`, userId, timeStr)
		msg, ivNew, err2 := encryptResBody(response_plain_msg, aesKey)
		if err2 != nil {
			fmt.Printf("%v\n", err2)
		}

		ivBase64 := base64.StdEncoding.EncodeToString(ivNew)

		fmt.Printf("response msg: %v\n", string(msg))
		fmt.Printf("response iv: %v\n", ivBase64)

		u := &encryptedResBody{
			Encrypted: string(msg),
			Iv:        ivBase64,
		}

		c.Response().Header().Set("X-Encrypted", "true")

		return c.JSON(http.StatusOK, u)
	})
	e.Logger.Fatal(e.Start(":8080"))
}
