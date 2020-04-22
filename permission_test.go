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

func  TestPermissionContract(t *testing.T) {
	os.Setenv("MODE","TEST")
	
	assert := assert.New(t)
	uid := uuid.New().String()

	cc, err := contractapi.NewChaincode(new(PermissionContract))
	assert.Nil(err, "error should be nil")

	stub := shimtest.NewMockStub("TestStub", cc)
	assert.NotNil(stub, "Stub is nil, TestStub creation failed")

	// - - - test PermissionContract:Put function - - - 
	putResp := stub.MockInvoke(uid,[][]byte{
		[]byte("PermissionContract:Put"),
		[]byte("1"),
		[]byte("1"),
		[]byte("Perm Name"),
		[]byte("Perm Module"),
	})
	assert.EqualValues(OK, putResp.GetStatus(), putResp.GetMessage())
	

	// - - - test PermissionContract:Get function - - - 
	testID := "1"
	getResp := stub.MockInvoke(uid, [][]byte{
		[]byte("PermissionContract:Get"),
		[]byte(testID),
	})
	assert.EqualValues(OK, getResp.GetStatus(), getResp.GetMessage())
	assert.NotNil(getResp.Payload, "getResp.Payload should not be nil")
	
	permissionObj := new(PermissionObj)
	err = json.Unmarshal(getResp.Payload, permissionObj)
	assert.Nil(err, "json.Unmarshal error should be nil")
	assert.NotNil(permissionObj, "permissionObj should not be nil")

	retrievedID := strconv.Itoa(permissionObj.PermID)
	assert.EqualValues(testID, retrievedID, "testID and retrievedID mismatch")
}