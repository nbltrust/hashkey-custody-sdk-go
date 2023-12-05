package jadepoolsaas

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"math/rand"
)

// NewCompany creates a new company with key and secret.
func NewCompany(key, secret string) *Company {
	return NewCompanyWithAddr(defaultAddr, key, secret)
}

// NewCompanyWithAddr creates a new company with server addr, key and secret.
func NewCompanyWithAddr(addr, key, secret string) *Company {
	a := &Company{
		Addr:   addr,
		Key:    key,
		Secret: secret,
	}
	a.session = &session{client: a}
	return a
}

// GetFundingWallets get all wallets in the company.
func (c *Company) GetFundingWallets() (*Result, error) {
	return c.session.get("/api/v1/funding/balances")
}

// FundingTransfer transfer the funding between wallets.
func (c *Company) FundingTransfer(from, to, coinType, value string) (*Result, error) {
	return c.FundingTransferWithMemo(from, to, coinType, value, "")
}

// FundingTransferWithMemo transfer the funding between wallets for the coin with specified memo.
func (c *Company) FundingTransferWithMemo(from, to, coinType, value, memo string) (*Result, error) {
	if len(coinType) == 0 || len(from) == 0 || len(to) == 0 || len(value) == 0 {
		return nil, errors.New("from or coinType or to or value is empty")
	}

	return c.session.post("/api/v1/funding/transfer", map[string]interface{}{
		"from":      from,
		"to":        to,
		"value":     value,
		"assetName": coinType,
		"memo":      memo,
		"message":   "",
	})
}

// GetFundingRecords get funding records.
func (c *Company) GetFundingRecords(page, amount int) (*Result, error) {
	return c.FilterFundingRecords(page, amount, "DESC", "", "", "", "", "created_at")
}

// FilterFundingRecords get funding records with filters.
func (c *Company) FilterFundingRecords(page, amount int, sort, coins, froms, toes, coinType, orderBy string) (*Result, error) {
	if len(sort) == 0 || len(orderBy) == 0 {
		return nil, errors.New("sort or orderBy is empty")
	}

	if page <= 0 {
		page = 1
	}
	if amount <= 0 {
		amount = 10
	}

	return c.session.getWithParams("/api/v1/funding/records", map[string]interface{}{
		"page":    page,
		"amount":  amount,
		"sort":    sort,
		"coins":   coins,
		"froms":   froms,
		"toes":    toes,
		"type":    coinType,
		"orderBy": orderBy,
	})
}

// CreateWallet create wallet.
func (c *Company) CreateWallet(name, password, webHook string) (*Result, error) {
	if len(name) == 0 || len(password) == 0 {
		return nil, errors.New("name or password is empty")
	}

	aesIV := make([]byte, 16)
	rand.Read(aesIV)
	mkey := sha256.Sum256([]byte(c.Secret))
	encryptPassword, err := aesEncryptStr(password, mkey[:], aesIV)
	if err != nil {
		return nil, err
	}

	ret, err := c.session.post("/api/v1/app", map[string]interface{}{
		"name":     name,
		"password": encryptPassword,
		"webHook":  webHook,
		"aesIV":    base64.StdEncoding.EncodeToString(aesIV),
	})
	if err != nil {
		return nil, err
	}

	appSecret, err := aesDecryptStr(ret.Data["encryptedAppSecret"].(string), mkey[:], aesIV)
	if err != nil {
		return ret, err
	}
	ret.Data["appSecret"] = appSecret
	return ret, nil
}

// GetWalletKeys get all keys for the specified wallet.
func (c *Company) GetWalletKeys(walletID string) (*Result, error) {
	if len(walletID) == 0 {
		return nil, errors.New("walletID is empty")
	}

	aesIV := make([]byte, 16)
	rand.Read(aesIV)
	mkey := sha256.Sum256([]byte(c.Secret))

	ret, err := c.session.getWithParams("/api/v1/app/"+walletID+"/keys", map[string]interface{}{
		"aesIV": base64.StdEncoding.EncodeToString(aesIV),
	})
	if err != nil {
		return nil, err
	}

	for _, key := range ret.Data["keys"].([]interface{}) {
		keymap := key.(map[string]interface{})
		appSecret, err := aesDecryptStr(keymap["encryptedAppSecret"].(string), mkey[:], aesIV)
		if err != nil {
			return ret, err
		}
		keymap["appSecret"] = appSecret
	}
	return ret, nil
}

// GetWalletInfo get attributes for the specified wallet.
func (c *Company) GetWalletInfo(walletID string) (*Result, error) {
	if len(walletID) == 0 {
		return nil, errors.New("walletID is empty")
	}
	return c.session.get("/api/v1/app/" + walletID + "/info")
}

// Trade create a new trade order in the API wallet.
func (c *Company) Trade(walletID, symbol, mType, side, amount, amountCoin string) (*Result, error) {
	if len(walletID) == 0 {
		return nil, errors.New("walletID is empty")
	}
	return c.session.post("/api/v1/app/"+walletID+"/trade", map[string]interface{}{
		"symbol":     symbol,
		"type":       mType,
		"side":       side,
		"amount":     amount,
		"amountCoin": amountCoin,
	})
}

// GetTradeOrder get trade order.
func (c *Company) GetTradeOrder(walletID, symbol, tradeID string) (*Result, error) {
	if len(walletID) == 0 {
		return nil, errors.New("walletID is empty")
	}
	return c.session.getWithParams("/api/v1/app/"+walletID+"/trade/"+tradeID, map[string]interface{}{
		"symbol": symbol,
	})
}

// UpdateWalletKey update app key attributes.
func (c *Company) UpdateWalletKey(appKey string, enable bool) (*Result, error) {
	if len(appKey) == 0 {
		return nil, errors.New("walletID is empty")
	}

	return c.session.put("/api/v1/appKey/"+appKey, map[string]interface{}{
		"enable": enable,
	})
}

// OTCCustomerGetSymbols get all otc symbols
func (c *Company) OTCCustomerGetSymbols() (*Result, error) {
	return c.session.get("/api/v1/otc/customer/symbols")
}

// Company represent a company.
type Company struct {
	Addr   string
	Key    string
	Secret string

	session *session
}

func (c *Company) getKey() string {
	return c.Key
}

func (c *Company) getKeyHeaderName() string {
	return "X-Company-Key"
}

func (c *Company) getSecret() string {
	return c.Secret
}

func (c *Company) getPubKey() string {
	return ""
}

func (c *Company) getAddr() string {
	return c.Addr
}
