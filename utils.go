package main

import (
	"fmt"
	"encoding/base64"
	"errors"
	
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

//GetTimestamp ...
func GetTimestamp(ctx contractapi.TransactionContextInterface) (time.Time, error) {
	epochTime, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return nil, err
	}

	created = time.Unix(epochTime.GetSeconds(), 0)
	return created, nil
}

//GetCallerID ...
func GetCallerID(ctx contractapi.TransactionContextInterface) (string, error ){
	callerIDBase64, err := ctx.GetClientIdentity().GetID
	if err != nil {
		return "", err
	}

	callerID, err := base64.StdEncoding.DecodeString(callerIDBase64)
	if err != nil {
		return "", err
	}
	
	callerIDList := strings.Split(string(callerID), "::")
	return callerIDList[1], nil
}