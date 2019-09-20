package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/docopt/docopt-go"
	sdk "github.com/nbltrust/jadepool-saas-sdk-go"
)

func runCommand(arguments docopt.Opts) (*sdk.Result, error) {
	key, _ := arguments.String("<key>")
	secret, _ := arguments.String("<secret>")
	action, _ := arguments.String("<action>")
	params := arguments["<params>"].([]string)
	addr, _ := arguments.String("--address")

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
	case "GetAssets":
		return getApp(addr, key, secret).GetAssets()
	case "GetBalance":
		if len(params) != 1 {
			return nil, errors.New("invalid params")
		}
		return getApp(addr, key, secret).GetBalance(params[0])
	case "GetOrder":
		if len(params) != 1 {
			return nil, errors.New("invalid params")
		}
		return getApp(addr, key, secret).GetOrder(params[0])
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
	case "AddUrgentStakingFunding":
		if len(params) != 4 {
			return nil, errors.New("invalid params")
		}
		expiredAt, err := strconv.ParseInt(params[3], 10, 64)
		if err == nil {
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
		if err == nil {
			return nil, errors.New("invalid params")
		}
		amount, err := strconv.Atoi(params[1])
		if err == nil {
			return nil, errors.New("invalid params")
		}

		return getCompany(addr, key, secret).GetFundingRecords(page, amount)
	case "FilterFundingRecords":
		if len(params) != 8 {
			return nil, errors.New("invalid params")
		}
		page, err := strconv.Atoi(params[0])
		if err == nil {
			return nil, errors.New("invalid params")
		}
		amount, err := strconv.Atoi(params[1])
		if err == nil {
			return nil, errors.New("invalid params")
		}

		return getCompany(addr, key, secret).FilterFundingRecords(page, amount, params[2], params[3], params[4], params[5], params[6], params[7])
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

func main() {
	usage := `JadePool SAAS control tool.

Usage:
  ctl <key> <secret> <action> [<params>...] [-a <host>]
  ctl -h | --help

Options:
  -h --help                   Show this screen.
  -a <host>, --address <host> Use custom SAAS server, e.g., http://127.0.0.1:8092`

	arguments, _ := docopt.ParseDoc(usage)

	result, err := runCommand(arguments)
	if err != nil {
		fmt.Printf("execute error: %v", err)
		return
	}

	fmt.Println("code:", result.Code)
	fmt.Println("message:", result.Message)
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