package main

import (
	"encoding/json"
	"errors"
	//      "strconv"
	//
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// ============================================================================================================================
// Get Info - get a info asset from ledger  정보를 가져온다...
// ============================================================================================================================
type AuditHistory struct {
	TxId  string      `json:"txId"`
	Value PaymentInfo `json:"value"`
}

func get_payment(stub shim.ChaincodeStubInterface, cusid string) ([]byte, error) {

	var paymentinfo PaymentInfo
	InfoAsBytes, err := stub.GetState(cusid) //getState retreives a key/value from the ledger
	if err == nil {                          // 지금까지 한번도 결제이력이 없을경우
		return InfoAsBytes, errors.New(cusid + " is No payment history. ")
	}

	json.Unmarshal(InfoAsBytes, &paymentinfo) //un stringify it aka JSON.parse()

	fmt.Println(cusid + " Customer Exist Payment History .... ")

	return InfoAsBytes, nil
}

// 체크 후, 중복일경우 1, 중복이 아닐경우 2 리턴..?
func check_dup(stub shim.ChaincodeStubInterface, payment_history []AuditHistory) (string, error) {

	var history []AuditHistory
	var info []PaymentInfo
	var result string

	history = payment_history

	//      json.Unmarshal(history.Value, &[]info)

	//        json.Unmarshal(payment_history, &history)     //convert to array of bytes

	//        json.Unmarshal(history, &info)     //convert to array of bytes

	for i := 0; i < len(history); i++ { // len 함수로 배열의 길이를 구한 뒤 배열의 길이 만큼 반복
		fmt.Println(history[i])
		var data AuditHistory
		var paymentinfo PaymentInfo
		data = history[i]
		value := data.Value
		tx := data.TxId
		fmt.Println(" TXID : " + tx)
		fmt.Println(" Month VALUE : " + value.Month)
	}

	//
	result = "1"

	return result, nil

}
