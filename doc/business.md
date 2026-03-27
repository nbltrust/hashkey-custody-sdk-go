git clone git@github.com:nbltrust/hashkey-custody-sdk-go.git  

git checkout business  

go mod tidy  

prepare pri_hashkey-hub.pem and pub_xpert_238.pem  

go run cmd/ctl/main.go hashkey-hub pri_hashkey-hub.pem BusinessAssetsGet -a "https://develop-saas.nbltrust.com/saas-business" -p pub_xpert_238.pem  

```json
{
  "assets": [
    {
      "decimal": 18,
      "id": 1,
      "name": "ETH",
      "switch": true
    },
    {
      "decimal": 8,
      "id": 2,
      "name": "BTC",
      "switch": true
    }
  ]
}
```

go run cmd/ctl/main.go hashkey-hub pri_hashkey-hub.pem BusinessClientGet 235 -a "https://develop-saas.nbltrust.com/saas-business" -p pub_xpert_238.pem

```json
{
  "email": "cob63176@tuofs.com",
  "id": 235,
  "kycLevel": 1,
  "name": "Christina",
  "phone": "+86-13817572905"
}
```

go run cmd/ctl/main.go hashkey-hub pri_hashkey-hub.pem BusinessWalletBalancesGet 435 1 -a "https://develop-saas.nbltrust.com/saas-business" -p pub_xpert_238.pem

```json
{
  "balances": [
    {
      "assetID": 1,
      "assetName": "ETH",
      "available": "1000.000000000000000000",
      "locked": "0.000000000000000000",
      "total": "1000.000000000000000000"
    }
  ]
}
```

go run cmd/ctl/main.go hashkey-hub pri_hashkey-hub.pem BusinessBalanceLock 165045845511 435 1 11 -a "https://develop-saas.nbltrust.com/saas-business" -p pub_xpert_238.pem

```json
{
  "available": "989.000000000000000000",
  "total": "1000.000000000000000000"
}
```

go run cmd/ctl/main.go hashkey-hub pri_hashkey-hub.pem BusinessBalanceUnlock 165045845512 435 1 6 -a "https://develop-saas.nbltrust.com/saas-business" -p pub_xpert_238.pem

```json
{
  "available": "995.000000000000000000",
  "total": "1000.000000000000000000"
}
```

go run cmd/ctl/main.go hashkey-hub pri_hashkey-hub.pem BusinessTransfer 165045845513 1 4.1 435 235 -a "https://develop-saas.nbltrust.com/saas-business" -p pub_xpert_238.pem

```json
{
  "id": 4
}
```

go run cmd/ctl/main.go hashkey-hub pri_hashkey-hub.pem BusinessSwap 165045845515 435 2 0.11 1 1.22 -a "https://develop-saas.nbltrust.com/saas-business" -p pub_xpert_238.pem

```json
{
  "id": 7
}
```

go run cmd/ctl/main.go hashkey-hub pri_hashkey-hub.pem BusinessOrderGet 165045845515 -a "https://develop-saas.nbltrust.com/saas-business" -p pub_xpert_238.pem

```json
{
  "detail": {
    "from": 435,
    "fromAmount": "0.11",
    "fromAssetID": 2,
    "fromWalletID": 1339,
    "note": "",
    "officialAmount": "1.22",
    "officialAssetID": 1,
    "officialWalletID": 1307,
    "sequence": "165045845515"
  },
  "id": 7,
  "status": "DONE",
  "type": "SWAP"
}
```
