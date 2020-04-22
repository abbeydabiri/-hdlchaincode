package main

import (
	"errors"
	"fmt"
	"time"
	"strconv"
	"encoding/json"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

//PermissionObj ...
type PermissionObj struct {
	PermID     int       `json:"permID"`
	PermRoleID int       `json:"permroleID"`
	PermName   string    `json:"permname"`
	PermModule string    `json:"permmodule"`
	Created    time.Time `json:"created"`
	Createdby  string    `json:"createdby"`
}

//PermissionHistory ...
type PermissionHistory struct {
	TxID string `json:"txID"`
	Timestamp time.Time  `json:"timestamp"`
	Permission *PermissionObj  `json:"permission"`
}

//PermissionContract for handling writing and reading from the world state
type PermissionContract struct {
	contractapi.Contract
}

//Put adds a new key with value to the world state
func (sc *PermissionContract) Put(ctx contractapi.TransactionContextInterface, permID int, permroleID int, permName string, permModule string)	(err error) {

	if permID == 0 {
		err = errors.New("Permission ID can not be empty")
		return
	}
	
	obj := new(PermissionObj)
	obj.PermID = permID
	obj.PermRoleID = permroleID
	obj.PermName = permName
	obj.PermModule = permModule

	if obj.Created, err = GetTimestamp(ctx); err != nil {
		return
	}

	if obj.Createdby, err = GetCallerID(ctx); err != nil {
		return
	}

	key := strconv.Itoa(permID)
	objBytes, _ := json.Marshal(obj)	
	err = ctx.GetStub().PutState(key, []byte(objBytes))
    return 
}

//Get retrieves the value linked to a key from the world state
func (sc *PermissionContract) Get(ctx contractapi.TransactionContextInterface, key string) (*PermissionObj, error) {
	
    existingObj, err := ctx.GetStub().GetState(key)
    if err != nil {
        return nil, err
    }

    if existingObj == nil {
        return nil, fmt.Errorf("Cannot read world state pair with key %s. Does not exist", key)
    }

	permissionObj := new(PermissionObj)
	if err := json.Unmarshal(existingObj, permissionObj); err != nil {
		return nil, fmt.Errorf("Data retrieved from world state for key %s was not of type PermissionObj", key)
	}
    return permissionObj, nil
}

//History retrieves the history linked to a key from the world state
func (sc *PermissionContract) History(ctx contractapi.TransactionContextInterface, key string) ([]PermissionHistory, error) {

	iter, err := ctx.GetStub().GetHistoryForKey(key)
	if err != nil {
        return nil, err
	}
	defer func() { _ = iter.Close() }()

	var results []PermissionHistory
	for iter.HasNext() {
		state, err := iter.Next()
		if err != nil {
			return nil, err
		}

		entryObj := new(PermissionObj)
		if errNew := json.Unmarshal(state.Value, entryObj); errNew != nil {
			return nil, errNew
		}

		entry := PermissionHistory{
			TxID:		state.GetTxId(),
			Timestamp:	time.Unix(state.GetTimestamp().GetSeconds(), 0),
			Permission:	entryObj,
		}

		results = append(results, entry)
	}
	return results, nil
}
