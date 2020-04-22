package main

import (
	"errors"
	"fmt"
	"time"
	"strconv"
	"encoding/json"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

//UserObj ...
type UserObj struct {
	UserID        int       `json:"userID"`
	UserType      string    `json:"usertype"`
	AccountStatus string    `json:"accountstatus"`
	UserCategory  string    `json:"usercategory"`
	FirstName     string    `json:"firstname"`
	LastName      string    `json:"lastname"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	Password      string    `json:"password"`
	RegDate       string    `json:"regdate"`
	Created       time.Time `json:"created"`
	Createdby     string    `json:"createdby"`
}

//UserHistory ...
type UserHistory struct {
	TxID string `json:"txID"`
	Timestamp time.Time  `json:"timestamp"`
	User *UserObj  `json:"user"`
}

//UserContract for handling writing and reading from the world state
type UserContract struct {
	contractapi.Contract
}



//Put adds a new key with value to the world state
func (sc *UserContract) Put(ctx contractapi.TransactionContextInterface, userID int, usertype string, accountstatus string, usercategory string, firstname string, lastname string, email string, phone string, password string, regdate string)	(err error) {

	if userID == 0 {
		err = errors.New("User ID can not be empty")
		return
	}
	
	obj := new(UserObj)
	obj.UserID = userID
	obj.UserType = usertype
	obj.AccountStatus = accountstatus
	obj.UserCategory = usercategory
	obj.FirstName = firstname
	obj.LastName = lastname
	obj.Email = email
	obj.Phone = phone
	obj.Password = password
	obj.RegDate = regdate	

	if obj.Created, err = GetTimestamp(ctx); err != nil {
		return
	}

	if obj.Createdby, err = GetCallerID(ctx); err != nil {
		return
	}

	key := strconv.Itoa(userID)
	objBytes, _ := json.Marshal(obj)	
	err = ctx.GetStub().PutState(key, []byte(objBytes))
    return 
}

//Get retrieves the value linked to a key from the world state
func (sc *UserContract) Get(ctx contractapi.TransactionContextInterface, key string) (*UserObj, error) {
	
    existingObj, err := ctx.GetStub().GetState(key)
    if err != nil {
        return nil, err
    }

    if existingObj == nil {
        return nil, fmt.Errorf("Cannot read world state pair with key %s. Does not exist", key)
    }

	userObj := new(UserObj)
	if err := json.Unmarshal(existingObj, userObj); err != nil {
		return nil, fmt.Errorf("Data retrieved from world state for key %s was not of type UserObj", key)
	}
    return userObj, nil
}

//History retrieves the history linked to a key from the world state
func (sc *UserContract) History(ctx contractapi.TransactionContextInterface, key string) ([]UserHistory, error) {

	iter, err := ctx.GetStub().GetHistoryForKey(key)
	if err != nil {
        return nil, err
	}
	defer func() { _ = iter.Close() }()

	var results []UserHistory
	for iter.HasNext() {
		state, err := iter.Next()
		if err != nil {
			return nil, err
		}

		entryObj := new(UserObj)
		if errNew := json.Unmarshal(state.Value, entryObj); errNew != nil {
			return nil, errNew
		}

		entry := UserHistory{
			TxID:		state.GetTxId(),
			Timestamp:	time.Unix(state.GetTimestamp().GetSeconds(), 0),
			User:	entryObj,
		}

		results = append(results, entry)
	}
	return results, nil
}
