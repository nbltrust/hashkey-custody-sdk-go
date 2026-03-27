package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/docopt/docopt-go"
	sdk "github.com/nbltrust/hashkey-custody-sdk-go"
)

func runCommand(arguments docopt.Opts) (*sdk.Result, error) {
	key, _ := arguments.String("<key>")
	secret, _ := arguments.String("<secret>")
	action, _ := arguments.String("<action>")
	params := arguments["<params>"].([]string)
	addr, _ := arguments.String("--address")
	pubKey, _ := arguments.String("--pubkey")

	switch action {
	case "CreateAddress":
		if len(params) != 1 {
			return nil, errors.New("invalid params")
		}
		return getApp(addr, key, secret).CreateAddress(params[0])
	case "CreateAddressWithMode":
		if len(params) != 2 {
			return nil, errors.New("invalid params")
		}
		return getApp(addr, key, secret).CreateAddressWithMode(params[0], params[1])
	case "VerifyAddress":
		if len(params) != 2 {
			return nil, errors.New("invalid params")
		}
		return getApp(addr, key, secret).VerifyAddress(params[0], params[1])
	case "CheckAddress":
		if len(params) != 2 {
			return nil, errors.New("invalid params")
		}
		return getApp(addr, key, secret).CheckAddress(params[0], params[1])
	case "GetAddress":
		if len(params) != 1 {
			return nil, errors.New("invalid params")
		}
		return getApp(addr, key, secret).GetAddress(params[0])
	case "GetAllAssets":
		return getApp(addr, key, secret).GetAllAssets()
	case "GetAssets":
		return getApp(addr, key, secret).GetAssets()
	case "GetAppInfo":
		return getApp(addr, key, secret).GetAppInfo()
	case "AddAsset":
		if len(params) != 1 {
			return nil, errors.New("invalid params")
		}
		return getApp(addr, key, secret).AddAsset(params[0])
	case "GetBalances":
		return getApp(addr, key, secret).GetBalances()
	case "GetBalance":
		if len(params) != 1 {
			return nil, errors.New("invalid params")
		}
		return getApp(addr, key, secret).GetBalance(params[0])

	case "GetOrders":
		if len(params) != 2 {
			return nil, errors.New("invalid params")
		}
		page, err := strconv.Atoi(params[0])
		if err != nil {
			return nil, errors.New("invalid params")
		}
		amount, err := strconv.Atoi(params[1])
		if err != nil {
			return nil, errors.New("invalid params")
		}

		return getApp(addr, key, secret).GetOrders(page, amount)
	case "GetOrder":
		if len(params) != 1 {
			return nil, errors.New("invalid params")
		}
		return getApp(addr, key, secret).GetOrder(params[0])
	case "UpdateOrder":
		if len(params) != 2 {
			return nil, errors.New("invalid params")
		}
		return getApp(addr, key, secret).UpdateOrder(params[0], params[1])
	case "Withdraw":
		if len(params) != 4 {
			return nil, errors.New("invalid params")
		}
		return getApp(addr, key, secret).Withdraw(params[0], params[1], params[2], params[3])
	case "WithdrawWithMemo":
		if len(params) != 5 {
			return nil, errors.New("invalid params")
		}
		return getApp(addr, key, secret).WithdrawWithMemo(params[0], params[1], params[2], params[3], params[4])
	case "Transfer":
		if len(params) != 3 {
			return nil, errors.New("invalid params")
		}
		return getApp(addr, key, secret).Transfer(params[0], params[1], params[2])
	case "Delegate":
		if len(params) != 3 {
			return nil, errors.New("invalid params")
		}
		return getApp(addr, key, secret).Delegate(params[0], params[1], params[2])
	case "UnDelegate":
		if len(params) != 3 {
			return nil, errors.New("invalid params")
		}
		return getApp(addr, key, secret).UnDelegate(params[0], params[1], params[2])
	case "GetValidators":
		if len(params) != 1 {
			return nil, errors.New("invalid params")
		}
		return getApp(addr, key, secret).GetValidators(params[0])
	case "GetStakingInterest":
		if len(params) != 2 {
			return nil, errors.New("invalid params")
		}
		return getApp(addr, key, secret).GetStakingInterest(params[0], params[1])
	case "AddUrgentStakingFunding":
		if len(params) != 4 {
			return nil, errors.New("invalid params")
		}
		expiredAt, err := strconv.ParseInt(params[3], 10, 64)
		if err != nil {
			return nil, errors.New("invalid params")
		}

		return getApp(addr, key, secret).AddUrgentStakingFunding(params[0], params[1], params[2], expiredAt)
	case "GetFundingWallets":
		return getCompany(addr, key, secret).GetFundingWallets()
	case "FundingTransfer":
		if len(params) != 4 {
			return nil, errors.New("invalid params")
		}
		return getCompany(addr, key, secret).FundingTransfer(params[0], params[1], params[2], params[3])
	case "FundingTransferWithMemo":
		if len(params) != 5 {
			return nil, errors.New("invalid params")
		}
		return getCompany(addr, key, secret).FundingTransferWithMemo(params[0], params[1], params[2], params[3], params[4])
	case "GetFundingRecords":
		if len(params) != 2 {
			return nil, errors.New("invalid params")
		}
		page, err := strconv.Atoi(params[0])
		if err != nil {
			return nil, errors.New("invalid params")
		}
		amount, err := strconv.Atoi(params[1])
		if err != nil {
			return nil, errors.New("invalid params")
		}

		return getCompany(addr, key, secret).GetFundingRecords(page, amount)
	case "FilterFundingRecords":
		if len(params) != 8 {
			return nil, errors.New("invalid params")
		}
		page, err := strconv.Atoi(params[0])
		if err != nil {
			return nil, errors.New("invalid params")
		}
		amount, err := strconv.Atoi(params[1])
		if err != nil {
			return nil, errors.New("invalid params")
		}

		return getCompany(addr, key, secret).FilterFundingRecords(page, amount, params[2], params[3], params[4], params[5], params[6], params[7])

	case "CreateWallet":
		if len(params) != 3 {
			return nil, errors.New("invalid params")
		}

		return getCompany(addr, key, secret).CreateWallet(params[0], params[1], params[2])

	case "GetWalletKeys":
		if len(params) != 1 {
			return nil, errors.New("invalid params")
		}

		return getCompany(addr, key, secret).GetWalletKeys(params[0])

	case "GetWalletInfo":
		if len(params) != 1 {
			return nil, errors.New("invalid params")
		}

		return getCompany(addr, key, secret).GetWalletInfo(params[0])

	case "Trade":
		if len(params) != 6 {
			return nil, errors.New("invalid params")
		}

		return getCompany(addr, key, secret).Trade(params[0], params[1], params[2], params[3], params[4], params[5])

	case "GetTradeOrder":
		if len(params) != 3 {
			return nil, errors.New("invalid params")
		}

		return getCompany(addr, key, secret).GetTradeOrder(params[0], params[1], params[2])

	case "UpdateWalletKey":
		if len(params) != 2 {
			return nil, errors.New("invalid params")
		}

		var enable bool
		if params[1] == "true" {
			enable = true
		} else if params[1] == "false" {
			enable = false
		} else {
			return nil, errors.New("invalid params")
		}

		return getCompany(addr, key, secret).UpdateWalletKey(params[0], enable)

	case "OTCSetSymbols":
		if len(params) != 1 {
			return nil, errors.New("invalid params")
		}

		symbols := []map[string]interface{}{}
		err := json.NewDecoder(strings.NewReader(params[0])).Decode(&symbols)
		if err != nil {
			return nil, err
		}
		return getApp(addr, key, secret).OTCSetSymbols(symbols)

	case "OTCGetSymbols":
		return getApp(addr, key, secret).OTCGetSymbols()

	case "OTCDeleteSymbol":
		if len(params) != 2 {
			return nil, errors.New("invalid params")
		}

		base, err := strconv.ParseUint(params[0], 10, 64)
		if err != nil {
			return nil, errors.New("invalid params")
		}
		quote, err := strconv.ParseUint(params[1], 10, 64)
		if err != nil {
			return nil, errors.New("invalid params")
		}

		return getApp(addr, key, secret).OTCDeleteSymbol(uint(base), uint(quote))

	case "OTCGetOrders":
		return getApp(addr, key, secret).OTCGetOrders()

	case "OTCGetPrices":
		return getApp(addr, key, secret).OTCGetPrices()

	case "OTCGetOrder":
		if len(params) != 1 {
			return nil, errors.New("invalid params")
		}

		return getApp(addr, key, secret).OTCGetOrder(params[0])

	case "OTCFeedPrice":
		if len(params) != 4 {
			return nil, errors.New("invalid params")
		}

		invalidAt, err := strconv.ParseInt(params[3], 10, 64)
		if err != nil {
			return nil, errors.New("invalid params")
		}

		return getApp(addr, key, secret).OTCFeedPrice(params[0], params[1], params[2], invalidAt)

	case "OTCGetPrice":
		if len(params) != 1 {
			return nil, errors.New("invalid params")
		}

		return getApp(addr, key, secret).OTCGetPrice(params[0])

	case "OTCClosePrice":
		if len(params) != 1 {
			return nil, errors.New("invalid params")
		}

		return getApp(addr, key, secret).OTCClosePrice(params[0])
	case "OTCTerminatePrice":
		if len(params) != 1 {
			return nil, errors.New("invalid params")
		}

		return getApp(addr, key, secret).OTCTerminatePrice(params[0])
	case "OTCGetPriceByCustomID":
		if len(params) != 1 {
			return nil, errors.New("invalid params")
		}

		return getApp(addr, key, secret).OTCGetPriceByCustomID(params[0])

	case "OTCClosePriceByCustomID":
		if len(params) != 1 {
			return nil, errors.New("invalid params")
		}

		return getApp(addr, key, secret).OTCClosePriceByCustomID(params[0])
	case "OTCTerminatePriceByCustomID":
		if len(params) != 1 {
			return nil, errors.New("invalid params")
		}

		return getApp(addr, key, secret).OTCTerminatePriceByCustomID(params[0])
	case "OTCCustomerGetSymbols":
		return getCompany(addr, key, secret).OTCCustomerGetSymbols()
	case "SystemGetTime":
		return getApp(addr, key, secret).SystemGetTime()
	case "GetMarket":
		if len(params) != 1 {
			return nil, errors.New("invalid params")
		}

		return getApp(addr, key, secret).GetMarket(params[0])
	case "KYCFileUpload":
		if len(params) != 1 {
			return nil, errors.New("invalid params")
		}

		return getKYC(addr, key, secret).FileUpload(params[0])
	case "KYCFileGet":
		if len(params) != 2 {
			return nil, errors.New("invalid params")
		}

		return getKYC(addr, key, secret).FileGet(params[0], params[1])
	case "KYCApplicationCreate":
		if len(params) != 2 {
			return nil, errors.New("invalid params")
		}

		return getKYC(addr, key, secret).ApplicationCreate(params[0], params[1], "demo")
	case "KYCApplicationUpdate":
		if len(params) != 3 {
			return nil, errors.New("invalid params")
		}

		return getKYC(addr, key, secret).ApplicationUpdate(params[0], params[1], params[2])
	case "KYCApplicationGet":
		if len(params) != 1 {
			return nil, errors.New("invalid params")
		}

		return getKYC(addr, key, secret).ApplicationGet(params[0], false)
	case "KYCApplicationSubmit":
		if len(params) != 1 {
			return nil, errors.New("invalid params")
		}

		return getKYC(addr, key, secret).ApplicationSubmit(params[0])
	case "BusinessAssetsGet":
		result, err := getBusiness(addr, key, secret, pubKey).AssetsGet()
		if err != nil {
			return nil, err
		}
		return result.ToResult(), nil
	case "BusinessClientGet":
		if len(params) != 1 {
			return nil, errors.New("invalid params")
		}

		uid, err := strconv.ParseUint(params[0], 10, 64)
		if err != nil {
			return nil, errors.New("invalid params")
		}

		result, err := getBusiness(addr, key, secret, pubKey).ClientGet(uint(uid))
		if err != nil {
			return nil, err
		}
		return result.ToResult(), nil
	case "BusinessClientsGet":
		if len(params) != 2 {
			return nil, errors.New("invalid params")
		}

		page, err := strconv.ParseUint(params[0], 10, 64)
		if err != nil {
			return nil, errors.New("invalid params")
		}

		amount, err := strconv.ParseUint(params[1], 10, 64)
		if err != nil {
			return nil, errors.New("invalid params")
		}

		result, err := getBusiness(addr, key, secret, pubKey).ClientsGet(uint(page), uint(amount))
		if err != nil {
			return nil, err
		}
		return result.ToResult(), nil
	case "BusinessClientCardsGet":
		if len(params) != 1 {
			return nil, errors.New("invalid params")
		}

		uid, err := strconv.ParseUint(params[0], 10, 64)
		if err != nil {
			return nil, errors.New("invalid params")
		}

		result, err := getBusiness(addr, key, secret, pubKey).ClientCardsGet(uint(uid))
		if err != nil {
			return nil, err
		}
		return result.ToResult(), nil
	case "BusinessWalletBalancesGet":
		if len(params) != 2 {
			return nil, errors.New("invalid params")
		}

		uid, err := strconv.ParseUint(params[0], 10, 64)
		if err != nil {
			return nil, errors.New("invalid params")
		}

		aid, err := strconv.ParseUint(params[1], 10, 64)
		if err != nil {
			return nil, errors.New("invalid params")
		}

		result, err := getBusiness(addr, key, secret, pubKey).WalletBalancesGet(uint(uid), uint(aid))
		if err != nil {
			return nil, err
		}
		return result.ToResult(), nil
	case "BusinessBalanceSettle":
		if len(params) != 5 {
			return nil, errors.New("invalid params")
		}

		uid, err := strconv.ParseUint(params[1], 10, 64)
		if err != nil {
			return nil, errors.New("invalid params")
		}

		aid, err := strconv.ParseUint(params[2], 10, 64)
		if err != nil {
			return nil, errors.New("invalid params")
		}

		result, err := getBusiness(addr, key, secret, pubKey).BalanceSettle(uint(uid), uint(aid), params[0], params[3], params[4])
		if err != nil {
			return nil, err
		}
		return result.ToResult(), nil
	case "BusinessBalanceLock":
		if len(params) != 4 {
			return nil, errors.New("invalid params")
		}

		uid, err := strconv.ParseUint(params[1], 10, 64)
		if err != nil {
			return nil, errors.New("invalid params")
		}

		aid, err := strconv.ParseUint(params[2], 10, 64)
		if err != nil {
			return nil, errors.New("invalid params")
		}

		result, err := getBusiness(addr, key, secret, pubKey).BalanceLock(uint(uid), uint(aid), params[0], params[3])
		if err != nil {
			return nil, err
		}
		return result.ToResult(), nil
	case "BusinessBalanceUnlock":
		if len(params) != 4 {
			return nil, errors.New("invalid params")
		}

		uid, err := strconv.ParseUint(params[1], 10, 64)
		if err != nil {
			return nil, errors.New("invalid params")
		}

		aid, err := strconv.ParseUint(params[2], 10, 64)
		if err != nil {
			return nil, errors.New("invalid params")
		}

		result, err := getBusiness(addr, key, secret, pubKey).BalanceUnlock(uint(uid), uint(aid), params[0], params[3])
		if err != nil {
			return nil, err
		}
		return result.ToResult(), nil
	case "BusinessTransfer":
		if len(params) != 5 {
			return nil, errors.New("invalid params")
		}

		aid, err := strconv.ParseUint(params[1], 10, 64)
		if err != nil {
			return nil, errors.New("invalid params")
		}

		fid, err := strconv.ParseUint(params[3], 10, 64)
		if err != nil {
			return nil, errors.New("invalid params")
		}

		tid, err := strconv.ParseUint(params[4], 10, 64)
		if err != nil {
			return nil, errors.New("invalid params")
		}

		result, err := getBusiness(addr, key, secret, pubKey).Transfer(uint(fid), uint(tid), uint(aid), params[0], params[2], "")
		if err != nil {
			return nil, err
		}
		return result.ToResult(), nil
	case "BusinessSwap":
		if len(params) != 6 {
			return nil, errors.New("invalid params")
		}

		uid, err := strconv.ParseUint(params[1], 10, 64)
		if err != nil {
			return nil, errors.New("invalid params")
		}

		aid, err := strconv.ParseUint(params[2], 10, 64)
		if err != nil {
			return nil, errors.New("invalid params")
		}

		oaid, err := strconv.ParseUint(params[4], 10, 64)
		if err != nil {
			return nil, errors.New("invalid params")
		}

		result, err := getBusiness(addr, key, secret, pubKey).Swap(uint(uid), uint(aid), uint(oaid), params[0], params[3], params[5], "")
		if err != nil {
			return nil, err
		}
		return result.ToResult(), nil
	case "BusinessBatch":
		if len(params) != 1 {
			return nil, errors.New("invalid params")
		}

		cmd := []*sdk.BatchCommand{}
		err := json.NewDecoder(strings.NewReader(params[0])).Decode(&cmd)
		if err != nil {
			return nil, err
		}
		result, err := getBusiness(addr, key, secret, pubKey).Batch(cmd)
		if err != nil {
			return nil, err
		}
		return result.ToResult(), nil
	case "BusinessOrderGet":
		if len(params) != 1 {
			return nil, errors.New("invalid params")
		}

		result, err := getBusiness(addr, key, secret, pubKey).OrderGetBySequence(params[0])
		if err != nil {
			return nil, err
		}
		return result.ToResult(), nil
	case "BusinessTransactionsGet":
		if len(params) != 5 {
			return nil, errors.New("invalid params")
		}

		uid, err := strconv.ParseUint(params[0], 10, 64)
		if err != nil {
			return nil, errors.New("invalid params")
		}

		page, err := strconv.ParseUint(params[3], 10, 64)
		if err != nil {
			return nil, errors.New("invalid params")
		}

		amount, err := strconv.ParseUint(params[4], 10, 64)
		if err != nil {
			return nil, errors.New("invalid params")
		}

		result, err := getBusiness(addr, key, secret, pubKey).TransactionsGet(uint(uid), params[1], params[2], uint(page), uint(amount))
		if err != nil {
			return nil, err
		}
		return result.ToResult(), nil
	default:
		return nil, errors.New("unknown action: " + action)
	}
}

func getApp(addr, key, secret string) *sdk.App {
	if len(addr) > 0 {
		return sdk.NewAppWithAddr(addr, key, secret)
	}
	return sdk.NewApp(key, secret)
}

func getCompany(addr, key, secret string) *sdk.Company {
	if len(addr) > 0 {
		return sdk.NewCompanyWithAddr(addr, key, secret)
	}
	return sdk.NewCompany(key, secret)
}

func getKYC(addr, key, secret string) *sdk.KYC {
	return sdk.NewKYCWithAddr(addr, key, secret)
}

func getBusiness(addr, key, secret, pubKey string) *sdk.Business {
	b, _ := sdk.NewBusinessWithAddr(addr, key, secret, pubKey)
	return b
}

func main() {
	usage := `JadePool SAAS control tool.

Usage:
  ctl <key> <secret> <action> [<params>...] [-a <host>] [-p <key>]
  ctl -h | --help

Options:
  -h --help                   Show this screen.
  -a <host>, --address <host> Use custom SAAS server, e.g., http://127.0.0.1:8092
  -p <key>, --pubkey <key> Use the public key pem file for verifying response`

	arguments, _ := docopt.ParseDoc(usage)

	result, err := runCommand(arguments)
	if err != nil {
		fmt.Printf("execute error: %v", err)
		return
	}

	fmt.Println("code:", result.Code)
	fmt.Println("message:", result.Message)
	fmt.Println("sign:", result.Sign)
	fmt.Println("data:")
	printMap(result.Data)
}

func printMap(m map[string]interface{}) {
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Print(string(b))
}
