package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric/common/flogging"
)

type SmartContract struct {
	contractapi.Contract
}

var logger = flogging.MustGetLogger("fabcar_cc")

type Car struct {
	ID      string `json:"id"`
	Make    string `json:"make"`
	Model   string `json:"model"`
	Color   string `json:"color"`
	Owner   string `json:"owner"`
	AddedAt uint64 `json:"addedAt"`
}

//custom functions by Omkar Dabir

type Transaction struct {
	ID       string `json:"id"`
	FromBank string `json:"fromBank"`
	ToBank   string `json:"toBank"`
	Amount   string `json:"amount"`
	Status   string `json:"status"`
	AddedAt  string `json:"addedAt"`
}

type carPrivateDetails struct {
	Owner string `json:"owner"`
	Price string `json:"price"`
}

func (s *SmartContract) SetTransaction(ctx contractapi.TransactionContextInterface, transactionData string) (string, error) {

	if len(transactionData) == 0 {
		return "", fmt.Errorf("Please pass the correct transaction data")
	}

	var transaction Transaction
	err := json.Unmarshal([]byte(transactionData), &transaction)
	if err != nil {
		return "", fmt.Errorf("Failed while unmarshling transaction. %s", err.Error())
	}

	transactionAsBytes, err := json.Marshal(transaction)
	if err != nil {
		return "", fmt.Errorf("Failed while marshling transaction. %s", err.Error())
	}

	ctx.GetStub().SetEvent("CreateAsset", transactionAsBytes)

	return ctx.GetStub().GetTxID(), ctx.GetStub().PutState(transaction.ID, transactionAsBytes)
}

func (s *SmartContract) GetTransactionById(ctx contractapi.TransactionContextInterface, transactionID string) (*Transaction, error) {
	if len(transactionID) == 0 {
		return nil, fmt.Errorf("Please provide correct contract Id")
		// return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	transactionAsBytes, err := ctx.GetStub().GetState(transactionID)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if transactionAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", transactionID)
	}

	transaction := new(Transaction)
	_ = json.Unmarshal(transactionAsBytes, transaction)

	return transaction, nil

}

//args[0] => privateCollectionName
func (s *SmartContract) SetPrivateTransaction(ctx contractapi.TransactionContextInterface, args []string) (string, error) {

	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect number of arguments. Expecting 1")
	}

	transMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return "", fmt.Errorf("Error getting transient: " + err.Error())
	}

	transactionDataAsBytes, ok := transMap["transaction"]
	if !ok {
		return "", fmt.Errorf("transaction must be a key in the transient map")
	}
	logger.Infof(string(transactionDataAsBytes))

	if len(transactionDataAsBytes) == 0 {
		return "", fmt.Errorf("Transaction value in the transient map must be a non-empty JSON string")
	}

	var transactionInput Transaction
	err = json.Unmarshal(transactionDataAsBytes, &transactionInput)
	if err != nil {
		return "", fmt.Errorf("Failed to decode JSON of: " + string(transactionDataAsBytes) + "Error is : " + err.Error())
	}

	if len(transactionInput.ID) == 0 {
		return "", fmt.Errorf("ID field must be a non-empty string")
	}
	if len(transactionInput.FromBank) == 0 {
		return "", fmt.Errorf("FromBank field must be a non-empty string")
	}
	if len(transactionInput.ToBank) == 0 {
		return "", fmt.Errorf("ToBank field must be a non-empty string")
	}
	if len(transactionInput.Amount) == 0 {
		return "", fmt.Errorf("Amount field must be a non-empty string")
	}
	if len(transactionInput.Status) == 0 {
		return "", fmt.Errorf("Status field must be a non-empty string")
	}
	if len(transactionInput.AddedAt) == 0 {
		return "", fmt.Errorf("AddedAt field must be a non-empty string")
	}

	// ==== Check if transaction already exists ====
	transactionAsBytes, err := ctx.GetStub().GetPrivateData(args[0], transactionInput.ID)
	if err != nil {
		return "", fmt.Errorf("Failed to get transaction: " + err.Error())
	} else if transactionAsBytes != nil {
		fmt.Println("This transaction already exists: " + transactionInput.ID)
		return "", fmt.Errorf("This transaction already exists: " + transactionInput.ID)
	}

	var transaction = Transaction{ID: transactionInput.ID,
		FromBank: transactionInput.FromBank,
		ToBank:   transactionInput.ToBank,
		Amount:   transactionInput.Amount,
		Status:   transactionInput.Status,
		AddedAt:  transactionInput.AddedAt}

	transactionAsBytes, err = json.Marshal(transaction)
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}
	return ctx.GetStub().GetTxID(), ctx.GetStub().PutPrivateData(args[0], transactionInput.ID, transactionAsBytes)
}

func (s *SmartContract) ReadPrivateTransaction(ctx contractapi.TransactionContextInterface, args []string) (string, error) {

	if len(args) != 2 {
		return "", fmt.Errorf("Incorrect number of arguments. Expecting 2")
	}
	// collectionTransactions12
	transactionAsBytes, err := ctx.GetStub().GetPrivateData(args[0], args[1])
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get private details for " + args[1] + ": " + err.Error() + "\"}"
		return "", fmt.Errorf(jsonResp)

	} else if transactionAsBytes == nil {
		jsonResp := "{\"Error\":\"Transaction private details does not exist: " + args[1] + "\"}"
		return "", fmt.Errorf(jsonResp)

	}
	return string(transactionAsBytes), nil
}

// custom functions end

func (s *SmartContract) CreateCar(ctx contractapi.TransactionContextInterface, carData string) (string, error) {

	if len(carData) == 0 {
		return "", fmt.Errorf("Please pass the correct car data")
	}

	var car Car
	err := json.Unmarshal([]byte(carData), &car)
	if err != nil {
		return "", fmt.Errorf("Failed while unmarshling car. %s", err.Error())
	}

	carAsBytes, err := json.Marshal(car)
	if err != nil {
		return "", fmt.Errorf("Failed while marshling car. %s", err.Error())
	}

	ctx.GetStub().SetEvent("CreateAsset", carAsBytes)

	return ctx.GetStub().GetTxID(), ctx.GetStub().PutState(car.ID, carAsBytes)
}

//
func (s *SmartContract) UpdateCarOwner(ctx contractapi.TransactionContextInterface, carID string, newOwner string) (string, error) {

	if len(carID) == 0 {
		return "", fmt.Errorf("Please pass the correct car id")
	}

	carAsBytes, err := ctx.GetStub().GetState(carID)

	if err != nil {
		return "", fmt.Errorf("Failed to get car data. %s", err.Error())
	}

	if carAsBytes == nil {
		return "", fmt.Errorf("%s does not exist", carID)
	}

	car := new(Car)
	_ = json.Unmarshal(carAsBytes, car)

	car.Owner = newOwner

	carAsBytes, err = json.Marshal(car)
	if err != nil {
		return "", fmt.Errorf("Failed while marshling car. %s", err.Error())
	}

	//  txId := ctx.GetStub().GetTxID()

	return ctx.GetStub().GetTxID(), ctx.GetStub().PutState(car.ID, carAsBytes)

}

func (s *SmartContract) GetHistoryForAsset(ctx contractapi.TransactionContextInterface, carID string) (string, error) {

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(carID)
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return "", fmt.Errorf(err.Error())
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
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

	return string(buffer.Bytes()), nil
}

func (s *SmartContract) GetCarById(ctx contractapi.TransactionContextInterface, carID string) (*Car, error) {
	if len(carID) == 0 {
		return nil, fmt.Errorf("Please provide correct contract Id")
		// return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	carAsBytes, err := ctx.GetStub().GetState(carID)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if carAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", carID)
	}

	car := new(Car)
	_ = json.Unmarshal(carAsBytes, car)

	return car, nil

}

func (s *SmartContract) DeleteCarById(ctx contractapi.TransactionContextInterface, carID string) (string, error) {
	if len(carID) == 0 {
		return "", fmt.Errorf("Please provide correct contract Id")
	}

	return ctx.GetStub().GetTxID(), ctx.GetStub().DelState(carID)
}

func (s *SmartContract) GetContractsForQuery(ctx contractapi.TransactionContextInterface, queryString string) ([]Car, error) {

	queryResults, err := s.getQueryResultForQueryString(ctx, queryString)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from ----world state. %s", err.Error())
	}

	return queryResults, nil

}

func (s *SmartContract) getQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]Car, error) {

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []Car{}

	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		newCar := new(Car)

		err = json.Unmarshal(response.Value, newCar)
		if err != nil {
			return nil, err
		}

		results = append(results, *newCar)
	}
	return results, nil
}

func (s *SmartContract) GetDocumentUsingCarContract(ctx contractapi.TransactionContextInterface, documentID string) (string, error) {
	if len(documentID) == 0 {
		return "", fmt.Errorf("Please provide correct contract Id")
	}

	params := []string{"GetDocumentById", documentID}
	queryArgs := make([][]byte, len(params))
	for i, arg := range params {
		queryArgs[i] = []byte(arg)
	}

	response := ctx.GetStub().InvokeChaincode("document_cc", queryArgs, "mychannel")

	return string(response.Payload), nil

}

func (s *SmartContract) CreateDocumentUsingCarContract(ctx contractapi.TransactionContextInterface, functionName string, documentData string) (string, error) {
	if len(documentData) == 0 {
		return "", fmt.Errorf("Please provide correct document data")
	}

	params := []string{functionName, documentData}
	queryArgs := make([][]byte, len(params))
	for i, arg := range params {
		queryArgs[i] = []byte(arg)
	}

	response := ctx.GetStub().InvokeChaincode("document_cc", queryArgs, "mychannel")

	return string(response.Payload), nil

}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		return
	}
	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting chaincodes: %s", err.Error())
	}

}
