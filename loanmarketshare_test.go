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

func  TestLoanMarketShareContract(t *testing.T) {
	os.Setenv("MODE","TEST")
	
	assert := assert.New(t)
	uid := uuid.New().String()

	cc, err := contractapi.NewChaincode(new(LoanMarketShareContract))
	assert.Nil(err, "error should be nil")

	stub := shimtest.NewMockStub("TestStub", cc)
	assert.NotNil(stub, "Stub is nil, TestStub creation failed")

	// - - - test LoanMarketShareContract:Put function - - - 
	putResp := stub.MockInvoke(uid,[][]byte{
		[]byte("LoanMarketShareContract:Put"),
		[]byte("1"),
		[]byte("TitleHolder"),
		[]byte("1.1"),
		[]byte("2.2"),
		[]byte("Statutes"),
		[]byte("3.3"),
		[]byte("Status"),
	})
	assert.EqualValues(OK, putResp.GetStatus(), putResp.GetMessage())
	

	// - - - test LoanMarketShareContract:Get function - - - 
	testID := "1"
	getResp := stub.MockInvoke(uid, [][]byte{
		[]byte("LoanMarketShareContract:Get"),
		[]byte(testID),
	})
	assert.EqualValues(OK, getResp.GetStatus(), getResp.GetMessage())
	assert.NotNil(getResp.Payload, "getResp.Payload should not be nil")
	
	loanMarketShareObj := new(LoanMarketShareObj)
	err = json.Unmarshal(getResp.Payload, loanMarketShareObj)
	assert.Nil(err, "json.Unmarshal error should be nil")
	assert.NotNil(loanMarketShareObj, "loanMarketShareObj should not be nil")

	retrievedID := strconv.Itoa(loanMarketShareObj.ShareID)
	assert.EqualValues(testID, retrievedID, "testID and retrievedID mismatch")
}