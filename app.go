package jadepoolsaas

import "errors"

// NewApp creates a new wallet with key and secret.
func NewApp(appKey, appSecret string) *App {
	return NewAppWithAddr(defaultAddr, appKey, appSecret)
}

// NewAppWithAddr creates a new wallet with server addr, key and secret.
func NewAppWithAddr(addr, appKey, appSecret string) *App {
	a := &App{
		Addr:   addr,
		Key:    appKey,
		Secret: appSecret,
	}
	a.session = &session{client: a}
	return a
}

// CreateAddress request new address.
func (a *App) CreateAddress(coinType string) (*Result, error) {
	return a.CreateAddressWithMode(coinType, "")
}

// CreateAddressWithMode request new address for the coin with specified memo.
func (a *App) CreateAddressWithMode(coinType, mode string) (*Result, error) {
	if len(coinType) == 0 {
		return nil, errors.New("coinType is empty")
	}

	return a.session.post("/api/v1/address/"+coinType+"/new", map[string]interface{}{
		"mode": mode,
	})
}

// VerifyAddress verify if an address is valid for specified coin.
func (a *App) VerifyAddress(coinType, address string) (*Result, error) {
	if len(coinType) == 0 || len(address) == 0 {
		return nil, errors.New("coinType or address is empty")
	}

	return a.session.post("/api/v1/address/"+coinType+"/verify", map[string]interface{}{
		"address": address,
	})
}

// CheckAddress anti-money laundering check.
func (a *App) CheckAddress(coinType, address string) (*Result, error) {
	if len(coinType) == 0 || len(address) == 0 {
		return nil, errors.New("coinType or address is empty")
	}

	return a.session.post("/api/v1/address/"+coinType+"/check", map[string]interface{}{
		"address": address,
	})
}

// GetAddress request an address, create if not exist.
func (a *App) GetAddress(coinType string) (*Result, error) {
	if len(coinType) == 0 {
		return nil, errors.New("coinType or address is empty")
	}

	return a.session.get("/api/v1/address/" + coinType)
}

// GetAllAssets fetch all available assets in the wallet.
func (a *App) GetAllAssets() (*Result, error) {
	return a.session.get("/api/v1/app/allAssets")
}

// GetAssets fetch all assets in the wallet.
func (a *App) GetAssets() (*Result, error) {
	return a.session.get("/api/v1/app/assetsWithID")
}

// GetAppInfo get the wallet's attributes.
func (a *App) GetAppInfo() (*Result, error) {
	return a.session.get("/api/v1/app/info")
}

// AddAsset add asset into the wallet.
func (a *App) AddAsset(coinName string) (*Result, error) {
	return a.session.post("/api/v1/app/assets", map[string]interface{}{
		"coinName": coinName,
	})
}

// GetBalances fetch all asset balances in the wallet.
func (a *App) GetBalances() (*Result, error) {
	return a.session.get("/api/v1/app/balances")
}

// GetBalance get the balance for specified coin.
func (a *App) GetBalance(coinType string) (*Result, error) {
	if len(coinType) == 0 {
		return nil, errors.New("coinType is empty")
	}

	return a.session.get("/api/v1/app/balance/" + coinType)
}

// GetOrders get orders in the wallet.
func (a *App) GetOrders(page, amount int) (*Result, error) {
	return a.session.getWithParams("/api/v1/app/orders", map[string]interface{}{
		"page":   page,
		"amount": amount,
	})
}

// GetOrder get order by id.
func (a *App) GetOrder(id string) (*Result, error) {
	if len(id) == 0 {
		return nil, errors.New("id is empty")
	}

	return a.session.get("/api/v1/app/order/" + id)
}

// UpdateOrder update order's note.
func (a *App) UpdateOrder(id string, note string) (*Result, error) {
	return a.session.put("/api/v1/app/order/"+id, map[string]interface{}{
		"note": note,
	})
}

// Withdraw request withdrawal.
func (a *App) Withdraw(id, coinType, to, value string) (*Result, error) {
	return a.WithdrawWithMemo(id, coinType, to, value, "")
}

// WithdrawWithMemo request withdrawal for the coin with specified memo.
func (a *App) WithdrawWithMemo(id, coinType, to, value, memo string) (*Result, error) {
	if len(coinType) == 0 || len(id) == 0 || len(to) == 0 || len(value) == 0 {
		return nil, errors.New("id or coinType or to or value is empty")
	}

	return a.session.post("/api/v1/app/"+coinType+"/withdraw", map[string]interface{}{
		"to":    to,
		"value": value,
		"memo":  memo,
		"id":    id,
	})
}

// Transfer transfer the funding to the specified wallet.
func (a *App) Transfer(to, coinType, value string) (*Result, error) {
	if len(coinType) == 0 || len(to) == 0 || len(value) == 0 {
		return nil, errors.New("coinType or to or value is empty")
	}

	return a.session.post("/api/v1/app/"+coinType+"/transfer", map[string]interface{}{
		"to":      to,
		"value":   value,
		"note":    "",
		"message": "",
	})
}

// Delegate request delegation.
func (a *App) Delegate(id, coinType, value string) (*Result, error) {
	if len(coinType) == 0 || len(id) == 0 || len(value) == 0 {
		return nil, errors.New("id or coinType or value is empty")
	}

	return a.session.post("/api/v1/staking/"+coinType+"/delegate", map[string]interface{}{
		"value": value,
		"id":    id,
	})
}

// UnDelegate request undelegation.
func (a *App) UnDelegate(id, coinType, value string) (*Result, error) {
	if len(coinType) == 0 || len(id) == 0 || len(value) == 0 {
		return nil, errors.New("id or coinType or value is empty")
	}

	return a.session.post("/api/v1/staking/"+coinType+"/undelegate", map[string]interface{}{
		"value": value,
		"id":    id,
	})
}

// GetValidators fetch all validators of specified coin.
func (a *App) GetValidators(coinType string) (*Result, error) {
	if len(coinType) == 0 {
		return nil, errors.New("coinType is empty")
	}

	return a.session.get("/api/v1/staking/" + coinType + "/validators")
}

// GetStakingInterest fetch one day interest for one cointype
func (a *App) GetStakingInterest(coinType, date string) (*Result, error) {
	if len(coinType) == 0 {
		return nil, errors.New("coinType is empty")
	}

	return a.session.getWithParams("/api/v1/staking/"+coinType+"/interest", map[string]interface{}{
		"date": date,
	})
}

// AddUrgentStakingFunding add urgent staking funding.
func (a *App) AddUrgentStakingFunding(id, coinType, value string, expiredAt int64) (*Result, error) {
	if len(coinType) == 0 || len(id) == 0 || len(value) == 0 {
		return nil, errors.New("id or coinType or value is empty")
	}

	return a.session.post("/api/v1/staking/"+coinType+"/funding", map[string]interface{}{
		"value":     value,
		"id":        id,
		"expiredAt": expiredAt,
	})
}

// OTCSetSymbols set otc symbols.
func (a *App) OTCSetSymbols(symbols []map[string]interface{}) (*Result, error) {
	return a.session.post("/api/v1/otc/symbols", map[string]interface{}{
		"symbols": symbols,
	})
}

// OTCGetSymbols get otc symbols.
func (a *App) OTCGetSymbols() (*Result, error) {
	return a.session.get("/api/v1/otc/symbols")
}

// OTCDeleteSymbol delete otc symbol.
func (a *App) OTCDeleteSymbol(baseCoinID, quoteCoinID uint) (*Result, error) {
	return a.session.deleteWithParams("/api/v1/otc/symbol", map[string]interface{}{
		"baseCoinID":  baseCoinID,
		"quoteCoinID": quoteCoinID,
	})
}

// OTCGetOrders get opening quote orders without feeding price.
func (a *App) OTCGetOrders() (*Result, error) {
	return a.session.get("/api/v1/otc/orders")
}

// OTCGetPrices get opening prices the app feed.
func (a *App) OTCGetPrices() (*Result, error) {
	return a.session.get("/api/v1/otc/prices")
}

// OTCGetOrder get order by id.
func (a *App) OTCGetOrder(orderID string) (*Result, error) {
	return a.session.get("/api/v1/otc/order/" + orderID)
}

// OTCFeedPrice feed otc price.
func (a *App) OTCFeedPrice(orderID, price, customID string, invalidAt int64) (*Result, error) {
	return a.session.post("/api/v1/otc/orders/"+orderID+"/price", map[string]interface{}{
		"price":     price,
		"customID":  customID,
		"invalidAt": invalidAt,
	})
}

// OTCGetPrice get the latest status of price.
func (a *App) OTCGetPrice(priceID string) (*Result, error) {
	if len(priceID) == 0 {
		return nil, errors.New("priceID is empty")
	}

	return a.session.get("/api/v1/otc/price/" + priceID)
}

// OTCClosePrice make a deal with the price.
func (a *App) OTCClosePrice(priceID string) (*Result, error) {
	if len(priceID) == 0 {
		return nil, errors.New("priceID is empty")
	}

	return a.session.get("/api/v1/otc/price/" + priceID + "/close")
}

// OTCTerminatePrice reject the price.
func (a *App) OTCTerminatePrice(priceID string) (*Result, error) {
	if len(priceID) == 0 {
		return nil, errors.New("priceID is empty")
	}

	return a.session.get("/api/v1/otc/price/" + priceID + "/terminate")
}

// OTCGetPriceByCustomID get the latest status of price.
func (a *App) OTCGetPriceByCustomID(customID string) (*Result, error) {
	if len(customID) == 0 {
		return nil, errors.New("customID is empty")
	}

	return a.session.get("/api/v1/otc/price/custom/" + customID)
}

// OTCClosePriceByCustomID make a deal with the price.
func (a *App) OTCClosePriceByCustomID(customID string) (*Result, error) {
	if len(customID) == 0 {
		return nil, errors.New("customID is empty")
	}

	return a.session.get("/api/v1/otc/price/custom/" + customID + "/close")
}

// OTCTerminatePriceByCustomID reject the price.
func (a *App) OTCTerminatePriceByCustomID(customID string) (*Result, error) {
	if len(customID) == 0 {
		return nil, errors.New("customID is empty")
	}

	return a.session.get("/api/v1/otc/price/custom/" + customID + "/terminate")
}

// SystemGetTime get the system timestamp.
func (a *App) SystemGetTime() (*Result, error) {
	return a.session.get("/api/v1/system/time")
}

// GetMarket get the system timestamp.
func (a *App) GetMarket(coinType string) (*Result, error) {
	return a.session.get("/api/v1/market/" + coinType)
}

// App represent a wallet.
type App struct {
	Addr   string
	Key    string
	Secret string

	session *session
}

func (a *App) getKey() string {
	return a.Key
}

func (a *App) getKeyHeaderName() string {
	return "X-App-Key"
}

func (a *App) getSecret() string {
	return a.Secret
}

func (a *App) getPubKey() string {
	return ""
}

func (a *App) getAddr() string {
	return a.Addr
}
