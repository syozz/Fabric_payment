package main

import (
	"encoding/json"
	"fmt"
	_ "strconv"
	_ "strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// ============================================================================================================================
// write() - genric write variable into ledger
//
// Shows Off PutState() - writting a key/value into the ledger
//
// Inputs - Array of strings
//    0   ,    1
//   key  ,  value
//  "abc" , "test"
// ============================================================================================================================
func write(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var key string
	var err error
	fmt.Println("starting write")

	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 2. key of the variable and value to set")
	}

	key = args[0] //rename for funsies
	res := PaymentInfo{}

	res.CusId = args[0]
	res.Year = args[1]
	res.Month = args[2]
	res.Payment_plan = args[3]
	res.Extra_plan = args[4]
	res.Amount_payment = args[5]
	res.Method_payment = args[6]

	jsonAsBytes, _ := json.Marshal(res)

	err = stub.PutState(key, jsonAsBytes) //write the variable into the ledger
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end write")
	return shim.Success(nil)
}

// ============================================================================================================================
// delete_info() - remove a marble from state and from marble index
// ============================================================================================================================
//func delete_info(stub shim.ChaincodeStubInterface, args []string) (pb.Response) {
//        fmt.Println("starting delete_info")
//
//       id := args[0]
//
//        // get the info
//       err := get_payment(stub, id) // _ 를 변수로 한 이유는 에러확인 유무만 하고 해당 변수를 사용하지 않기 때문이다.
//        if err != nil{
//                fmt.Println("Failed to find info by id " + id)
//                return shim.Error(err.Error())
//        }
//
//        // remove the marble
//        err = stub.DelState(id)                                                 //remove the key from chaincode state
//        if err != nil {
//                return shim.Error("Failed to delete state")
//        }
//
//        fmt.Println("- end delete_info")
//        return shim.Success(nil)
//}

// ============================================================================================================================
// Init Marble - create a new marble, store into chaincode state
//
// Shows off building a key's JSON value manually
//
// Inputs - Array of strings
// type PaymentInfo struct {
//         CusId                string   `json:"cusid"`  // 고객아이디
//         Year                 string   `json:"year"`  // 년도
//         Month                string   `json:"month"`  // 결제 월
//         Payment_plan         string   `json:"month"`  // 가입요금제
//         Extra_plan           string   `json:“extra_plan＂`   //  부가서비스
//         Amount_payment       string   `json:"amount_payment"`  // 결제금액
//         Method_payment       string   `json:“method_payment"`  // 결제방식 [ 카드, 계좌이체, 무통장 ]
//         Result               string   `json:“result"`  // 결제 결과 [ 성공, 실패 ]
// }
// ============================================================================================================================
func payment(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	//      var result PaymentInfo
	fmt.Println("starting Payment")
	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}

	//build the info  json string manually
	res := PaymentInfo{}

	cusid := args[0]
	res.CusId = args[0]
	res.Year = args[1]
	res.Month = args[2]
	res.Payment_plan = args[3]
	res.Extra_plan = args[4]
	res.Amount_payment = args[5]
	res.Method_payment = args[6]

	jsonAsBytes, _ := json.Marshal(res)

	// 고객 아이디를 기준으로 원장을 검색해서 이전에 동일한 결제내역이 있는지 파악....
	// 같은 월에 같은 요금제, 부가서비스의 결제금액이 존재 하는지 체크,,, [중복결제 확인]
	//      res , err := get_payment(stub, cusid)
	//        if res == "0" { // 한번도 결제한 이력이 없을경우 무조건 결제 승인
	//              fmt.Println(" ##### Data String : " + str )
	//              err = stub.PutState(cusid, jsonAsBytes)
	//              if err != nil {
	//                      return shim.Error(err.Error())
	//              }
	//              fmt.Println(" ##### Insert First payment data ")
	//              return shim.Success(nil)
	//        }

	//        fmt.Println( cusid + " Customer Exist Payment History .... ")

	// 고객 아이디를 기준으로 지금까지의 결제내역을 추출,,,
	// txid, value(결제이력) 으로 구성된 배열을 리턴,,, []AuditHistory
	fmt.Println(" #### getHistory START ")
	payment_history, err := getHistory(stub, cusid)

	fmt.Println(" #### check_dup START ")
	check_history, err := check_dup(stub, payment_history, args)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println(" check_dup RESULT : " + check_history)
	if check_history == "1" { // 중복이 아님, put 해줌.
		fmt.Println(" Start PutState. CusID : " + res.CusId)
		fmt.Println("==========================")
		fmt.Println("CusId : " + res.CusId)
		fmt.Println("Year : " + res.Year)
		fmt.Println("MONTH : " + res.Month)
		fmt.Println("Payment_Plan : " + res.Payment_plan)
		fmt.Println("Extra_Plan : " + res.Extra_plan)
		fmt.Println("Amount_Payment : " + res.Amount_payment)
		fmt.Println("Method_Payment : " + res.Method_payment)
		fmt.Println("==========================")

		err = stub.PutState(cusid, jsonAsBytes)
		if err != nil {
			return shim.Error(err.Error())
		}
	} else if check_history == "2" {
		fmt.Println(" Already Exist Payment History ")
	}

	fmt.Println("- end Payment")
	return shim.Success(nil)
}
