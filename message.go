package main

import (
	"errors"
	"fmt"
	"time"
	"strconv"
	"encoding/json"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

//MessageObj ...
type MessageObj struct {
	MessageID   int       `json:"messageID"`
	From        string    `json:"from"`
	To          string    `json:"to"`
	Subject     string    `json:"subject"`
	Message     string    `json:"message"`
	MessageDate string    `json:"messagedate"`
	Created     time.Time `json:"created"`
	Createdby   string    `json:"createdby"`
}

//MessageHistory ...
type MessageHistory struct {
	TxID string `json:"txID"`
	Timestamp time.Time  `json:"timestamp"`
	Message *MessageObj  `json:"message"`
}

//MessageContract for handling writing and reading from the world state
type MessageContract struct {
	contractapi.Contract
}

//Put adds a new key with value to the world state
func (sc *MessageContract) Put(ctx contractapi.TransactionContextInterface, messageID int, from string, to string, subject string, message string, messagedate string)	(err error) {

	if messageID == 0 {
		err = errors.New("Message ID can not be empty")
		return
	}
	
	obj := new(MessageObj)
	obj.MessageID = messageID
	obj.From = from
	obj.To = to
	obj.Subject = subject
	obj.Message = message
	obj.MessageDate = messagedate

	if obj.Created, err = GetTimestamp(ctx); err != nil {
		return
	}

	if obj.Createdby, err = GetCallerID(ctx); err != nil {
		return
	}

	key := strconv.Itoa(messageID)
	objBytes, _ := json.Marshal(obj)	
	err = ctx.GetStub().PutState(key, []byte(objBytes))
    return 
}

//Get retrieves the value linked to a key from the world state
func (sc *MessageContract) Get(ctx contractapi.TransactionContextInterface, key string) (*MessageObj, error) {
	
    existingObj, err := ctx.GetStub().GetState(key)
    if err != nil {
        return nil, err
    }

    if existingObj == nil {
        return nil, fmt.Errorf("Cannot read world state pair with key %s. Does not exist", key)
    }

	messageObj := new(MessageObj)
	if err := json.Unmarshal(existingObj, messageObj); err != nil {
		return nil, fmt.Errorf("Data retrieved from world state for key %s was not of type MessageObj", key)
	}
    return messageObj, nil
}

//History retrieves the history linked to a key from the world state
func (sc *MessageContract) History(ctx contractapi.TransactionContextInterface, key string) ([]MessageHistory, error) {

	iter, err := ctx.GetStub().GetHistoryForKey(key)
	if err != nil {
        return nil, err
	}
	defer func() { _ = iter.Close() }()

	var results []MessageHistory
	for iter.HasNext() {
		state, err := iter.Next()
		if err != nil {
			return nil, err
		}

		entryObj := new(MessageObj)
		if errNew := json.Unmarshal(state.Value, entryObj); errNew != nil {
			return nil, errNew
		}

		entry := MessageHistory{
			TxID:		state.GetTxId(),
			Timestamp:	time.Unix(state.GetTimestamp().GetSeconds(), 0),
			Message:	entryObj,
		}

		results = append(results, entry)
	}
	return results, nil
}
