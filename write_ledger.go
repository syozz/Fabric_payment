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
	var key, value string
	var err error
	fmt.Println("starting write")

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2. key of the variable and value to set")
	}

	key = args[0] //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the ledger
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end write")
	return shim.Success(nil)
}

// ============================================================================================================================
// delete_info() - remove a marble from state and from marble index
// ============================================================================================================================
func delete_info(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("starting delete_info")

	id := args[0]

	// get the info
	_, err := get_payment(stub, id) // _ 를 변수로 한 이유는 에러확인 유무만 하고 해당 변수를 사용하지 않기 때문이다.
	if err != nil {
		fmt.Println("Failed to find info by id " + id)
		return shim.Error(err.Error())
	}

	// remove the marble
	err = stub.DelState(id) //remove the key from chaincode state
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	fmt.Println("- end delete_info")
	return shim.Success(nil)
}

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
	fmt.Println("starting Payment")

	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}

	cusid := args[0]
	year := args[1]
	month := args[2]
	payment_plan := args[3]
	extra_plan := args[4]
	amount_payment := args[5]
	method_payment := args[6]
	result := "1"

	//build the info  json string manually
	str := `{
                "cusid": "` + cusid + `",
                "year": "` + year + `",
                "month": "` + month + `",
                "payment_plan": "` + payment_plan + `",
                "extra_plan": "` + extra_plan + `"
                "amount_payment": "` + amount_payment + `"
                "method_payment": "` + method_payment + `"
                "result": "` + result + `"
        }`

	// 고객 아이디를 기준으로 원장을 검색해서 이전에 동일한 결제내역이 있는지 파악....
	// 같은 월에 같은 요금제, 부가서비스의 결제금액이 존재 하는지 체크,,, [중복결제 확인]
	result, err := get_payment(stub, cusid)
	if err != nil { // 한번도 결제한 이력이 없을경우 무조건 결제 승인
		err = stub.PutState(cusid, []byte(str))
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(cusid + "Payment : " + amount_payment) // 정보가 존재하면 중지
	}

	fmt.Println(cusid + " Customer Exist Payment History .... ")

	payment_history, err := gethistory(cusid)

	check_history, err := check_dup(payment_history)
	if err != nil {
		return shim.Error(err.Error())
	}

	if check_history == "2" { // 중복이 아님, put 해줌.
		err = stub.PutState(cusid, []byte(str))
		if err != nil {
			return shim.Error(err.Error())
		}
	} else if check_history == "1" {
		fmt.Println(" Already Exist Payment History ")
	}

	fmt.Println("- end Payment")
	return shim.Success(nil)
}

// ============================================================================================================================
// modify   //   id, 수정할 내용, 데이터 를 입력 받는다..   변수 3개를 받는다.
// ============================================================================================================================
func modify(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting modify")

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	var id = args[0]
	var type_d = args[1]
	var data = args[2]

	fmt.Println("Modify " + type_d + " Data -> " + data + " ...")

	// 현재 ID의 정보를 가져옴
	infoAsBytes, err := stub.GetState(id)
	if err != nil {
		return shim.Error("Failed to get Info")
	}
	res := Info{}
	json.Unmarshal(infoAsBytes, &res) //un stringify it aka JSON.parse()

	// transfer
	if type_d == "name" {
		fmt.Println("Current " + type_d + " : " + res.Name + " => " + data + "...")
		res.Name = data
	} else if type_d == "phone" {
		fmt.Println("Current " + type_d + " : " + res.Phone + " => " + data + "...")
		res.Phone = data
	} else if type_d == "address" {
		fmt.Println("Current " + type_d + " : " + res.Address + " => " + data + "...")
		res.Address = data
	}

	jsonAsBytes, _ := json.Marshal(res) //convert to array of bytes
	err = stub.PutState(id, jsonAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end set " + type_d + " data.")
	return shim.Success(nil)
}
