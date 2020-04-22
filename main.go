package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	
	bankContract := new(BankContract)
	bankBranchContract := new(BankBranchContract)
	branchConfigContract := new(BranchConfigContract)
	buyerContract := new(BuyerContract)
	loanContract := new(LoanContract)
	loanBuyerContract := new(LoanBuyerContract)
	loanDocContract := new(LoanDocContract)
	loanMarketShareContract := new(LoanMarketShareContract)
	loanRatingContract := new(LoanRatingContract)
	messageContract := new(MessageContract)
	permissionContract := new(PermissionContract)
	propertyContract := new(PropertyContract)
	roleContract := new(RoleContract)
	sellerContract := new(SellerContract)
	transactionContract := new(TransactionContract)
	userContract := new(UserContract)
	userCategoryContract := new(UserCategoryContract)
	


	cc, err := contractapi.NewChaincode(bankContract, bankBranchContract, branchConfigContract, buyerContract, loanContract, loanBuyerContract, loanDocContract, 
		loanMarketShareContract, loanRatingContract, messageContract, permissionContract, propertyContract, roleContract, sellerContract, transactionContract, 
		userContract, userCategoryContract )

    if err != nil {
        panic(err.Error())
    }

	if err := cc.Start(); err != nil {
        panic(err.Error())
    }
}
