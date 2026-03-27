package jadepoolsaas

import (
	"encoding/base64"
	"os"
)

// NewBusinessWithAddr creates a new business instance with server addr, business key, the private key pem file of your service and the public key pem file of xpert.
func NewBusinessWithAddr(addr, businessKey, pemFilePath, pubPemFilePath string) (*Business, error) {
	data, err := os.ReadFile(pemFilePath)
	if err != nil {
		return nil, err
	}

	var pubData []byte
	if len(pubPemFilePath) > 0 {
		pubData, err = os.ReadFile(pubPemFilePath)
		if err != nil {
			return nil, err
		}
	}

	a := &Business{
		Addr:   addr,
		Key:    businessKey,
		Secret: base64.StdEncoding.EncodeToString(data),
	}
	if len(pubData) > 0 {
		a.PubKey = base64.StdEncoding.EncodeToString(pubData)
	}
	a.session = &session{client: a}
	return a, nil
}

// ClientGet fetch client info.
func (b *Business) ClientGet(userID uint) (*BusinessResult, error) {
	return b.session.businessGetWithParams("/api/v1/business/client", map[string]interface{}{
		"userID": userID,
	})
}

// ClientsGet fetch clients.
func (b *Business) ClientsGet(page, amount uint) (*BusinessResult, error) {
	return b.session.businessGetWithParams("/api/v1/business/clients", map[string]interface{}{
		"page":   page,
		"amount": amount,
	})
}

// ClientCardsGet fetch client's cards.
func (b *Business) ClientCardsGet(userID uint) (*BusinessResult, error) {
	return b.session.businessGetWithParams("/api/v1/business/client/cards", map[string]interface{}{
		"userID": userID,
	})
}

// AssetsGet fetch all assets in the wallet.
func (b *Business) AssetsGet() (*BusinessResult, error) {
	return b.session.businessGet("/api/v1/business/assets")
}

// WalletBalancesGet fetch the asset balance in the wallet.
func (b *Business) WalletBalancesGet(userID, assetID uint) (*BusinessResult, error) {
	return b.session.businessGetWithParams("/api/v1/business/wallet/balances", map[string]interface{}{
		"userID":  userID,
		"assetID": assetID,
	})
}

// BalanceSettle settle the balance.
func (b *Business) BalanceSettle(userID, assetID uint, mType, sequence, amount string) (*BusinessResult, error) {
	return b.session.businessPost("/api/v1/business/balance/settle", map[string]interface{}{
		"userID":   userID,
		"assetID":  assetID,
		"type":     mType,
		"sequence": sequence,
		"amount":   amount,
	})
}

// BalanceLock lock the balance.
func (b *Business) BalanceLock(userID, assetID uint, sequence, amount string) (*BusinessResult, error) {
	return b.session.businessPut("/api/v1/business/balance/lock", map[string]interface{}{
		"userID":   userID,
		"assetID":  assetID,
		"sequence": sequence,
		"amount":   amount,
	})
}

// BalanceUnlock unlock the balance.
func (b *Business) BalanceUnlock(userID, assetID uint, sequence, amount string) (*BusinessResult, error) {
	return b.session.businessPut("/api/v1/business/balance/unlock", map[string]interface{}{
		"userID":   userID,
		"assetID":  assetID,
		"sequence": sequence,
		"amount":   amount,
	})
}

// Transfer transfer.
func (b *Business) Transfer(from, to, assetID uint, sequence, amount, note string) (*BusinessResult, error) {
	return b.session.businessPost("/api/v1/business/transfer", map[string]interface{}{
		"from":     from,
		"to":       to,
		"sequence": sequence,
		"assetID":  assetID,
		"amount":   amount,
		"note":     note,
	})
}

// Swap swap.
func (b *Business) Swap(from, fromAssetID, officialAssetID uint, sequence, fromAmount, officialAmount, note string) (*BusinessResult, error) {
	return b.session.businessPost("/api/v1/business/swap", map[string]interface{}{
		"from":            from,
		"fromAssetID":     fromAssetID,
		"officialAssetID": officialAssetID,
		"sequence":        sequence,
		"fromAmount":      fromAmount,
		"officialAmount":  officialAmount,
		"note":            note,
	})
}

// Batch batch.
func (b *Business) Batch(cmd []*BatchCommand) (*BusinessResult, error) {
	return b.session.businessPost("/api/v1/business/batch", map[string]interface{}{
		"cmd": cmd,
	})
}

// BatchCommand ...
type BatchCommand struct {
	Name string      `json:"name"`
	Args interface{} `json:"args"`
}

// OrderGetBySequence fetch the order by the sequence.
func (b *Business) OrderGetBySequence(sequence string) (*BusinessResult, error) {
	return b.session.businessGet("/api/v1/business/order/sequence/" + sequence)
}

// TransactionsGet fetch transactions.
func (b *Business) TransactionsGet(userID uint, mType, status string, page, amount uint) (*BusinessResult, error) {
	return b.session.businessGetWithParams("/api/v1/business/transactions", map[string]interface{}{
		"userID": userID,
		"type":   mType,
		"status": status,
		"page":   page,
		"amount": amount,
	})
}

// Business represents a business instance.
type Business struct {
	Addr   string
	Key    string
	Secret string
	PubKey string

	session *session
}

func (b *Business) getKey() string {
	return b.Key
}

func (b *Business) getKeyHeaderName() string {
	return "X-BUSINESS-KEY"
}

func (b *Business) getSecret() string {
	return b.Secret
}

func (b *Business) getPubKey() string {
	return b.PubKey
}

func (b *Business) getAddr() string {
	return b.Addr
}
