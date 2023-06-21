package jadepoolsaas

import (
	"bytes"
	"github.com/imroc/req"
)

// NewKYCWithAddr creates a new kyc instance with server addr, key and secret.
func NewKYCWithAddr(addr, appKey, appSecret string) *KYC {
	a := &KYC{
		Addr:   addr,
		Key:    appKey,
		Secret: appSecret,
	}
	a.session = &session{client: a}
	return a
}

// GeneralSettingsGet get the general settings.
func (k *KYC) GeneralSettingsGet() (*Result, error) {
	return k.session.get("/api/v1/generalSettings")
}

// FileUpload upload file.
func (k *KYC) FileUpload(filePath string) (*Result, error) {
	return k.session.postFile("/api/v1/file", filePath)
}

// FileUpload2 upload file.
func (k *KYC) FileUpload2(applicationID, fileName string, file *bytes.Reader) (*Result, error) {
	return k.session.postFile2("/api/v1/file", fileName, file, map[string]interface{}{
		"applicationID": applicationID,
	})
}

// FileGet get file.
func (k *KYC) FileGet(fileID, filePath string) (*Result, error) {
	return k.session.getFile("/api/v1/file/"+fileID, filePath)
}

// FileGet2 get file.
func (k *KYC) FileGet2(fileID, applicationID string) (*req.Resp, error) {
	return k.session.getFile2("/api/v1/file/"+fileID, map[string]interface{}{
		"applicationID": applicationID,
	})
}

// ApplicationCreate create an application.
func (k *KYC) ApplicationCreate(mType, identifier, operator string, identifierID uint) (*Result, error) {
	return k.session.post("/api/v1/application", map[string]interface{}{
		"type":         mType,
		"identifier":   identifier,
		"identifierID": identifierID,
		"operator":     operator,
	})
}

// ApplicationUpdate update the application.
func (k *KYC) ApplicationUpdate(applicationID, key, value string) (*Result, error) {
	return k.session.patch("/api/v1/application/"+applicationID, map[string]interface{}{
		key: value,
	})
}

// ApplicationUpdate2 update the application.
func (k *KYC) ApplicationUpdate2(applicationID string, content map[string]interface{}) (*Result, error) {
	return k.session.patch("/api/v1/application/"+applicationID, content)
}

// ApplicationGet get the application.
func (k *KYC) ApplicationGet(applicationID string, expand bool) (*Result, error) {
	return k.session.getWithParams("/api/v1/application/"+applicationID, map[string]interface{}{
		"expand": expand,
	})
}

// ApplicationJumioGet get the application's jumio info.
func (k *KYC) ApplicationJumioGet(applicationID, locale, id string) (*Result, error) {
	return k.session.getWithParams("/api/v1/application/"+applicationID+"/jumio", map[string]interface{}{
		"locale": locale,
		"id":     id,
	})
}

// ApplicationGetByIdentifier get the application.
func (k *KYC) ApplicationGetByIdentifier(mType, identifier string, expand bool) (*Result, error) {
	return k.session.getWithParams("/api/v1/application/identifier/"+mType+"/"+identifier, map[string]interface{}{
		"expand": expand,
	})
}

// ApplicationSubmit submit the application.
func (k *KYC) ApplicationSubmit(applicationID string) (*Result, error) {
	return k.session.put("/api/v1/application/"+applicationID, map[string]interface{}{})
}

// ApplicationSettingsUpdate update the settings of application.
func (k *KYC) ApplicationSettingsUpdate(applicationID string, settings map[string]interface{}) (*Result, error) {
	return k.session.put("/api/v1/application/"+applicationID+"/settings", settings)
}

// JumioPost post jumio result with the application.
func (k *KYC) JumioPost(applicationID string, content map[string]interface{}) (*Result, error) {
	return k.session.post("/api/v1/application/"+applicationID+"/jumio", content)
}

// FiatCreate create a fiat with the application.
func (k *KYC) FiatCreate(applicationID string, content map[string]interface{}) (*Result, error) {
	return k.session.post("/api/v1/application/"+applicationID+"/fiat", content)
}

// FiatsGet get fiats with the application.
func (k *KYC) FiatsGet(applicationID string) (*Result, error) {
	return k.session.get("/api/v1/application/" + applicationID + "/fiats")
}

// FiatUpdate update the fiat.
func (k *KYC) FiatUpdate(fiatID string, content map[string]interface{}) (*Result, error) {
	return k.session.put("/api/v1/fiat/"+fiatID, content)
}

// FiatDelete delete the fiat.
func (k *KYC) FiatDelete(fiatID string) (*Result, error) {
	return k.session.delete("/api/v1/fiat/" + fiatID)
}

// ApplicationHistoriesGet get the change histories with the application.
func (k *KYC) ApplicationHistoriesGet(applicationID string) (*Result, error) {
	return k.session.get("/api/v1/application/" + applicationID + "/histories")
}

// KYC represents a kyc instance.
type KYC struct {
	Addr   string
	Key    string
	Secret string

	session *session
}

func (k *KYC) getKey() string {
	return k.Key
}

func (k *KYC) getKeyHeaderName() string {
	return "X-API-Key"
}

func (k *KYC) getSecret() string {
	return k.Secret
}

func (k *KYC) getAddr() string {
	return k.Addr
}
