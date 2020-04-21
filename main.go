package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	
	bankContract := new(Bank)

    cc, err := contractapi.NewChaincode(bankContract)

    if err != nil {
        panic(err.Error())
    }

}
