package main

import (
	"errors"
	"fmt"
	"time"
	"strconv"
	"encoding/json"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

//LoanDocObj ...
type LoanDocObj struct {
	DocID     int       `json:"docID"`
	LoanID    int       `json:"loanID"`
	DocName   string    `json:"docname"`
	DocDesc   string    `json:"docdesc"`
	DocLink   string    `json:"doclink"`
	Created   time.Time `json:"created"`
	Createdby string    `json:"createdby"`
}

//LoanDocHistory ...
type LoanDocHistory struct {
	TxID string `json:"txID"`
	Timestamp time.Time  `json:"timestamp"`
	LoanDoc *LoanDocObj  `json:"loandoc"`
}

//LoanDocContract for handling writing and reading from the world state
type LoanDocContract struct {
	contractapi.Contract
}

//Put adds a new key with value to the world state
func (sc *LoanDocContract) Put(ctx contractapi.TransactionContextInterface, docID int, loanID int, docname string, docdesc string, doclink string)	(err error) {

	if docID == 0 {
		err = errors.New("Doc ID can not be empty")
		return
	}

	if loanID == 0 {
		err = errors.New("Loan ID can not be empty")
		return
	}

	
	obj := new(LoanDocObj)
	obj.DocID = docID
	obj.LoanID = loanID
	obj.DocName = docname
	obj.DocDesc = docdesc
	obj.DocLink = doclink	

	if obj.Created, err = GetTimestamp(ctx); err != nil {
		return
	}

	if obj.Createdby, err = GetCallerID(ctx); err != nil {
		return
	}

	key := strconv.Itoa(docID)
	objBytes, _ := json.Marshal(obj)	
	err = ctx.GetStub().PutState(key, []byte(objBytes))
    return 
}

//Get retrieves the value linked to a key from the world state
func (sc *LoanDocContract) Get(ctx contractapi.TransactionContextInterface, key string) (*LoanDocObj, error) {
	
    existingObj, err := ctx.GetStub().GetState(key)
    if err != nil {
        return nil, err
    }

    if existingObj == nil {
        return nil, fmt.Errorf("Cannot read world state pair with key %s. Does not exist", key)
    }

	loanDocObj := new(LoanDocObj)
	if err := json.Unmarshal(existingObj, loanDocObj); err != nil {
		return nil, fmt.Errorf("Data retrieved from world state for key %s was not of type LoanDocObj", key)
	}
    return loanDocObj, nil
}

//History retrieves the history linked to a key from the world state
func (sc *LoanDocContract) History(ctx contractapi.TransactionContextInterface, key string) ([]LoanDocHistory, error) {

	iter, err := ctx.GetStub().GetHistoryForKey(key)
	if err != nil {
        return nil, err
	}
	defer func() { _ = iter.Close() }()

	var results []LoanDocHistory
	for iter.HasNext() {
		state, err := iter.Next()
		if err != nil {
			return nil, err
		}

		entryObj := new(LoanDocObj)
		if errNew := json.Unmarshal(state.Value, entryObj); errNew != nil {
			return nil, errNew
		}

		entry := LoanDocHistory{
			TxID:		state.GetTxId(),
			Timestamp:	time.Unix(state.GetTimestamp().GetSeconds(), 0),
			LoanDoc:	entryObj,
		}

		results = append(results, entry)
	}
	return results, nil
}
