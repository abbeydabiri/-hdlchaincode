package main

import (
	"errors"
	"fmt"
	"time"
	"strconv"
	"encoding/json"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

//UserCategoryObj ...
type UserCategoryObj struct {
	CatID          int       `json:"catID"`
	CatName        string    `json:"catname"`
	CatDescription string    `json:"catdescription"`
	Created        time.Time `json:"created"`
	Createdby      string    `json:"createdby"`
}

//UserCategoryHistory ...
type UserCategoryHistory struct {
	TxID string `json:"txID"`
	Timestamp time.Time  `json:"timestamp"`
	UserCategory *UserCategoryObj  `json:"usercategory"`
}

//UserCategoryContract for handling writing and reading from the world state
type UserCategoryContract struct {
	contractapi.Contract
}

//Put adds a new key with value to the world state
func (sc *UserCategoryContract) Put(ctx contractapi.TransactionContextInterface, catID int, catname string, catdescription string)	(err error) {

	if catID == 0 {
		err = errors.New("Loan Rating ID can not be empty")
		return
	}
	
	obj := new(UserCategoryObj)
	obj.CatID = catID
	obj.CatName = catname
	obj.CatDescription = catdescription
	

	if obj.Created, err = GetTimestamp(ctx); err != nil {
		return
	}

	if obj.Createdby, err = GetCallerID(ctx); err != nil {
		return
	}

	key := strconv.Itoa(catID)
	objBytes, _ := json.Marshal(obj)	
	err = ctx.GetStub().PutState(key, []byte(objBytes))
    return 
}

//Get retrieves the value linked to a key from the world state
func (sc *UserCategoryContract) Get(ctx contractapi.TransactionContextInterface, key string) (*UserCategoryObj, error) {
	
    existingObj, err := ctx.GetStub().GetState(key)
    if err != nil {
        return nil, err
    }

    if existingObj == nil {
        return nil, fmt.Errorf("Cannot read world state pair with key %s. Does not exist", key)
    }

	userCategoryObj := new(UserCategoryObj)
	if err := json.Unmarshal(existingObj, userCategoryObj); err != nil {
		return nil, fmt.Errorf("Data retrieved from world state for key %s was not of type UserCategoryObj", key)
	}
    return userCategoryObj, nil
}

//History retrieves the history linked to a key from the world state
func (sc *UserCategoryContract) History(ctx contractapi.TransactionContextInterface, key string) ([]UserCategoryHistory, error) {

	iter, err := ctx.GetStub().GetHistoryForKey(key)
	if err != nil {
        return nil, err
	}
	defer func() { _ = iter.Close() }()

	var results []UserCategoryHistory
	for iter.HasNext() {
		state, err := iter.Next()
		if err != nil {
			return nil, err
		}

		entryObj := new(UserCategoryObj)
		if errNew := json.Unmarshal(state.Value, entryObj); errNew != nil {
			return nil, errNew
		}

		entry := UserCategoryHistory{
			TxID:		state.GetTxId(),
			Timestamp:	time.Unix(state.GetTimestamp().GetSeconds(), 0),
			UserCategory:	entryObj,
		}

		results = append(results, entry)
	}
	return results, nil
}
