package main

import (
	"errors"
	"fmt"
	"time"
	"strconv"
	"encoding/json"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

//BankBranchObj ...
type BankBranchObj struct {
	BranchID            int       `json:"branchID"`
	BankID              int       `json:"bankID"`
	BranchName          string    `json:"branchname"`
	BranchManagerUserID int       `json:"branchmanageruserID"`
	BranchManagerRoleID int       `json:"branchmanagerroleID"`
	Location            string    `json:"location"`
	LocationLat         string    `json:"locationlat"`
	LocationLong        string    `json:"locationlong"`
	Created             time.Time `json:"created"`
	Createdby           string    `json:"createdby"`
}

//BankBranchHistory ...
type BankBranchHistory struct {
	TxID string `json:"txID"`
	Timestamp time.Time  `json:"timestamp"`
	BankBranch *BankBranchObj  `json:"bankbranch"`
}

//BankBranchContract for handling writing and reading from the world state
type BankBranchContract struct {
	contractapi.Contract
}

//Put adds a new key with value to the world state
func (sc *BankBranchContract) Put(ctx contractapi.TransactionContextInterface, branchID int, bankID int, branchname string, branchmanageruserID int, branchmanagerroleID int, location string, locationlat string, locationlong string)	(err error) {

	if branchID == 0 {
		err = errors.New("Branch ID can not be empty")
		return
	}

	if bankID == 0 {
		err = errors.New("Bank ID can not be empty")
		return
	}

	
	obj := new(BankBranchObj)
	obj.BranchID = branchID
	obj.BranchName = branchname

	obj.BranchManagerUserID = branchmanageruserID
	obj.BranchManagerRoleID = branchmanagerroleID

	obj.Location = location
	obj.LocationLat = locationlat
	obj.LocationLong = locationlong

	if obj.Created, err = GetTimestamp(ctx); err != nil {
		return
	}

	if obj.Createdby, err = GetCallerID(ctx); err != nil {
		return
	}

	key := strconv.Itoa(branchID)
	objBytes, _ := json.Marshal(obj)	
	err = ctx.GetStub().PutState(key, []byte(objBytes))
    return 
}

//Get retrieves the value linked to a key from the world state
func (sc *BankBranchContract) Get(ctx contractapi.TransactionContextInterface, key string) (*BankBranchObj, error) {
	
    existingObj, err := ctx.GetStub().GetState(key)
    if err != nil {
        return nil, err
    }

    if existingObj == nil {
        return nil, fmt.Errorf("Cannot read world state pair with key %s. Does not exist", key)
    }

	bankBranchObj := new(BankBranchObj)
	if err := json.Unmarshal(existingObj, bankBranchObj); err != nil {
		return nil, fmt.Errorf("Data retrieved from world state for key %s was not of type BankBranchObj", key)
	}
    return bankBranchObj, nil
}

//History retrieves the history linked to a key from the world state
func (sc *BankBranchContract) History(ctx contractapi.TransactionContextInterface, key string) ([]BankBranchHistory, error) {

	iter, err := ctx.GetStub().GetHistoryForKey(key)
	if err != nil {
        return nil, err
	}
	defer func() { _ = iter.Close() }()

	var results []BankBranchHistory
	for iter.HasNext() {
		state, err := iter.Next()
		if err != nil {
			return nil, err
		}

		entryObj := new(BankBranchObj)
		if errNew := json.Unmarshal(state.Value, entryObj); errNew != nil {
			return nil, errNew
		}

		entry := BankBranchHistory{
			TxID:		state.GetTxId(),
			Timestamp:	time.Unix(state.GetTimestamp().GetSeconds(), 0),
			BankBranch:	entryObj,
		}

		results = append(results, entry)
	}
	return results, nil
}
