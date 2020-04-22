package main

import (
	"os"
	"testing"
	"encoding/json"
	"strconv"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func  TestBankContract(t *testing.T) {
	os.Setenv("MODE","TEST")
	
	assert := assert.New(t)
	uid := uuid.New().String()

	cc, err := contractapi.NewChaincode(new(BankContract))
	assert.Nil(err, "error should be nil")

	stub := shimtest.NewMockStub("TestStub", cc)
	assert.NotNil(stub, "Stub is nil, TestStub creation failed")

	// - - - test BankContract:Put function - - - 
	putResp := stub.MockInvoke(uid,[][]byte{
		[]byte("BankContract:Put"),
		[]byte("1"),
		[]byte("Test Bank 1"),
		[]byte("Test Hq 1"),
		[]byte("Test BankCategory 1"),
		[]byte("1"),
		[]byte("Test Location 1"),
		[]byte("Test LocationLat 1"),
		[]byte("Test LocationLong 1"),
	})
	assert.EqualValues(OK, putResp.GetStatus(), putResp.GetMessage())
	

	// - - - test BankContract:Get function - - - 
	testID := "1"
	getResp := stub.MockInvoke(uid, [][]byte{
		[]byte("BankContract:Get"),
		[]byte(testID),
	})
	assert.EqualValues(OK, getResp.GetStatus(), getResp.GetMessage())
	assert.NotNil(getResp.Payload, "getResp.Payload should not be nil")
	
	bankObj := new(BankObj)
	err = json.Unmarshal(getResp.Payload, bankObj)
	assert.Nil(err, "json.Unmarshal error should be nil")
	assert.NotNil(bankObj, "bankObj should not be nil")

	retrievedID := strconv.Itoa(bankObj.BankID)
	assert.EqualValues(testID, retrievedID, "testID and retrievedID mismatch")
}