package main

import (
	"encoding/json"
	//        "errors"
	//      "strconv"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	//      pb "github.com/hyperledger/fabric/protos/peer"
)

// ============================================================================================================================
// Get Info - get a info asset from ledger  정보를 가져온다...
// ============================================================================================================================
type AuditHistory struct {
	TxId  string      `json:"txId"`
	Value PaymentInfo `json:"value"`
}

func get_payment(stub shim.ChaincodeStubInterface, cusid string) (string, error) {

	var err error
	var paymentinfo PaymentInfo
	InfoAsBytes, err := stub.GetState(cusid) //getState retreives a key/value from the ledger
	if InfoAsBytes == nil {                  // 지금까지 한번도 결제이력이 없을경우
		fmt.Println(" Payment Info is NIL ")
		return "0", err
	}

	// 결제이력이 하나라도 존재하면
	json.Unmarshal(InfoAsBytes, &paymentinfo) //un stringify it aka JSON.parse()

	fmt.Println(cusid + " Customer Exist Payment History .... ")

	return "1", nil
}

// 체크 후, 중복일경우 1, 중복이 아닐경우 2 리턴..?
func check_dup(stub shim.ChaincodeStubInterface, payment_history []AuditHistory, args []string) (string, error) {

	var history []AuditHistory
	var result string
	str_leng := 0 // 현재 저장된 배열수 [ append될때마다 증가 ]
	total := []PaymentInfo{}

	_ = args[0]
	year := args[1]
	month := args[2]
	_ = args[3]
	_ = args[4]
	_ = args[5]
	_ = args[6]

	history = payment_history

	// 히스토리 내역의 value 부분 추출해서 PaymentInfo 타입의 배열에 저장
	for i := 0; i < len(history); i++ { // len 함수로 배열의 길이를 구한 뒤 배열의 길이 만큼 반복
		fmt.Println(history[i])
		value := history[i].Value
		//              tx := history[i].TxId
		//              fmt.Println( " TXID : " + tx)
		fmt.Println(" Year/Month VALUE : " + value.Month + "/" + value.Year)
		fmt.Println(" Payment Amount : " + value.Amount_payment)
		if value.Month == month && value.Year == year {
			result = "2"
			return result, nil
		} else {
			total[str_leng] = history[i].Value
		}
	}
	fmt.Println("Setting RESULT")

	result = "1"

	return result, nil

}
