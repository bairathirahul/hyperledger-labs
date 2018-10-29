package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strings"
	"time"
)

// Chaincode01 : the chaincode struct
type Chaincode01 struct {
}

// Asset : that will be stored in our ledger
type Asset struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Init : Chaincode initialization function
func (c *Chaincode01) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke : Invoke call of chaincode
func (c *Chaincode01) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()

	fmt.Println("Inside invoke method")
	fmt.Println("Received function name - " + function)
	fmt.Print("Received arguments - ")
	fmt.Println(args)

	if function == "get" {
		return c.get(stub, args)
	} else if function == "set" {
		return c.set(stub, args)
	} else if function == "history" {
		return c.history(stub, args)
	} else if function == "delete" {
		return c.delete(stub, args)
	} else {
		fmt.Println("Invalid function " + function)
		return shim.Error("Invalid function " + function)
	}
}

func (c *Chaincode01) set(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2 {
		fmt.Println("Key and value pair required")
	}

	asset := Asset{args[0], args[1]}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		fmt.Println("Failed to marshal asset with key " + args[0])
		return shim.Error(err.Error())
	}

	// Save asset state
	err = stub.PutState(args[0], assetJSON)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(assetJSON)
}

func (c *Chaincode01) get(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		fmt.Println("Key is required")
	}

	// Read asset
	assetJSON, err := stub.GetState(args[0])
	if err != nil {
		fmt.Println("Failed to get asset with key " + args[0])
		return shim.Error(err.Error())
	} else if assetJSON == nil {
		fmt.Println("Asset with key " + args[0] + " not found")
		return shim.Error("Asset with key " + args[0] + " not found")
	}

	return shim.Success(assetJSON)
}

func (c *Chaincode01) history(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		fmt.Println("Key is required")
	}

	// Read transaction history
	transactionIter, err := stub.GetHistoryForKey(args[0])
	if err != nil {
		return shim.Error("Error reading transaction history: " + err.Error())
	}

	// To close the iterator at the end
	defer transactionIter.Close()

	// Prepare asset history output
	output := fmt.Sprintf("+" + strings.Repeat("-", 146) + "+\n")
	output += fmt.Sprintf("| %64s | %32s | %32s | %7s |\n", "Transaction ID", "Value", "Time", "Deleted")
	output += fmt.Sprintf("+" + strings.Repeat("-", 146) + "+\n")
	for transactionIter.HasNext() {
		transaction, err := transactionIter.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		asset := Asset{}
		assetJSON := transaction.GetValue()
		if assetJSON != nil {
			err = json.Unmarshal(assetJSON, &asset)
			if err != nil {
				fmt.Println("Error in asset JSON " + string(assetJSON))
				return shim.Error("Asset unmarshalling error: " + err.Error())
			}
		}

		transactionTime := time.Unix(transaction.GetTimestamp().GetSeconds(), 0)
		output += fmt.Sprintf("| %64s | %32s | %32s | %7t |\n",
			transaction.GetTxId(), asset.Value, transactionTime.String(),
			transaction.GetIsDelete())
		output += fmt.Sprintf("+" + strings.Repeat("-", 146) + "+\n")
	}

	fmt.Print(output)

	return shim.Success([]byte(output))
}

func (c *Chaincode01) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		fmt.Println("Key is required")
	}

	// Delete State
	err := stub.DelState(args[0])
	if err != nil {
		return shim.Error("Failed to delete asset:" + err.Error())
	}

	return shim.Success([]byte("Asset with key " + args[0] + " deleted"))
}

func main() {
	err := shim.Start(new(Chaincode01))
	if err != nil {
		fmt.Printf("Error starting Chaincode: %s\n", err)
	}
}
