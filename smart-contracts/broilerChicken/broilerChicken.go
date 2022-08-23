package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	sc "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/common/flogging"
)

// SmartContract Define the Smart Contract structure
type SmartContract struct {
}

type GrowthInformation struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Instruction string `json:"instruction"`
}

// Chicken :  Define the Chicken structure, with 4 properties.  Structure tags are used by encoding/json library
type Chicken struct {
	Birthday  string            `json:"birthday"`
	Breed     string            `json:"breed"`
	Price     string            `json:"price"`
	Owner     string            `json:"owner"`
	GrowthInf GrowthInformation `json:"growthInformation"`
}

type HistoryModel struct {
	TxId      string `json:"TxId"`
	Value     *Chicken
	Timestamp string `json:"Timestamp"`
	IsDelete  bool   `json:"IsDelete"`
}

// Init ;  Method for initializing smart contract
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

var logger = flogging.MustGetLogger("chicken_cc")

// Invoke :  Method for INVOKING smart contract
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	function, args := APIstub.GetFunctionAndParameters()
	logger.Infof("Function name is:  %d", function)
	logger.Infof("Args length is : %d", len(args))

	switch function {
	case "queryChicken":
		return s.queryChicken(APIstub, args)
	case "initLedger":
		return s.initLedger(APIstub)
	case "createChicken":
		return s.createChicken(APIstub, args)
	case "queryAllChickens":
		return s.queryAllChickens(APIstub)
	case "changeChickenOwner":
		return s.changeChickenOwner(APIstub, args)
	case "getHistoryForAsset":
		return s.getHistoryForAsset(APIstub, args)
	case "queryChickensByOwner":
		return s.queryChickensByOwner(APIstub, args)
	default:
		return shim.Error("Invalid Smart Contract function name.")
	}

	// return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryChicken(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	chickenAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(chickenAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {

	gi := GrowthInformation{
		Key:         "Vaccine",
		Value:       "VC123",
		Instruction: "This is the first vaccination.",
	}

	chickens := []Chicken{
		Chicken{Birthday: "1401/5/18", Breed: "Ross", Price: "12.0", Owner: "Sadegh", GrowthInf: gi},
		Chicken{Birthday: "1401/5/18", Breed: "Cobb", Price: "10.0", Owner: "Sadegh", GrowthInf: gi},
		Chicken{Birthday: "1401/5/18", Breed: "Cobb", Price: "10.0", Owner: "Sadegh", GrowthInf: gi},
		Chicken{Birthday: "1401/5/18", Breed: "Ross", Price: "12.0", Owner: "Sadegh", GrowthInf: gi},
	}

	i := 0
	for i < len(chickens) {
		chickenAsBytes, _ := json.Marshal(chickens[i])
		APIstub.PutState("CHICKEN"+strconv.Itoa(i), chickenAsBytes)
		// APIstub.PutState(uuid.New().String(), chickenAsBytes)
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) createChicken(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var chicken = Chicken{Birthday: args[1], Breed: args[2], Price: args[3], Owner: args[4]}

	chickenAsBytes, _ := json.Marshal(chicken)
	APIstub.PutState(args[0], chickenAsBytes)

	// indexName := "owner~key"
	// indexKey, err := APIstub.CreateCompositeKey(indexName, []string{chicken.Owner, args[0]})
	// //indexKey, err := APIstub.CreateCompositeKey(chicken.Owner, []string{args[0]})
	// if err != nil {
	// 	return shim.Error(err.Error())
	// }
	// value := []byte{0x00}
	// APIstub.PutState(indexKey, value)

	return shim.Success(chickenAsBytes)
}

func (S *SmartContract) queryChickensByOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments")
	}
	owner := args[0]

	startKey := ""
	endKey := ""

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		chicken := Chicken{}
		json.Unmarshal(queryResponse.Value, &chicken)

		if chicken.Owner == owner {
			// Add a comma before array members, suppress it for the first array member
			if bArrayMemberAlreadyWritten == true {
				buffer.WriteString(",")
			}
			buffer.WriteString("{\"Key\":")
			buffer.WriteString("\"")
			buffer.WriteString(queryResponse.Key)
			buffer.WriteString("\"")

			buffer.WriteString(", \"Record\":")
			// Record is a JSON object, so we write as-is
			buffer.WriteString(string(queryResponse.Value))
			buffer.WriteString("}")
			bArrayMemberAlreadyWritten = true
		}

	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllChickens:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) queryAllChickens(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := ""
	endKey := ""

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllChickens:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) changeChickenOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	chickenAsBytes, _ := APIstub.GetState(args[0])
	chicken := Chicken{}

	json.Unmarshal(chickenAsBytes, &chicken)
	chicken.Owner = args[1]

	chickenAsBytes, _ = json.Marshal(chicken)
	APIstub.PutState(args[0], chickenAsBytes)

	return shim.Success(chickenAsBytes)
}

func (t *SmartContract) getHistoryForAsset(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	chickenName := args[0]

	resultsIterator, err := stub.GetHistoryForKey(chickenName)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the marble
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON marble)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoryForAsset returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
