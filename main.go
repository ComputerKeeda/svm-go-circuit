package main

import (
	"fmt"
	"strconv"

	"github.com/ComputerKeeda/svm-go-circuit/modules"
	"github.com/ComputerKeeda/svm-go-circuit/types"
)

func main() {
	// prover.CreateVkPk()
	values := modules.FetchJsonData()
	// now we will get the first transaction data from the json file
	var firstTransaction types.SVMTransaction
	firstTransaction.TxSignature = values.TxSignature[0]
	firstTransaction.Fee = values.Fee[0]
	firstTransaction.PreBalance = values.PreBalance[0]
	firstTransaction.PostBalance = values.PostBalance[0]
	firstTransaction.AccountKeys = values.AccountKeys[0]
	var svmInstructions []types.SVMInstruction
	for _, instruction := range values.Instructions[0] {
		svmInstruction := types.SVMInstruction{
			Accounts:       instruction.Accounts,
			Data:           instruction.Data,
			ProgramIDIndex: instruction.ProgramIDIndex,
		}
		svmInstructions = append(svmInstructions, svmInstruction)
	}
	firstTransaction.RecentBlockhash = values.RecentBlockhash[0]

	// now we will print the first transaction data with indentation

	fmt.Printf("Fee: %s\n", firstTransaction.Fee)
	fmt.Printf("PreBalance: %.f\n", firstTransaction.PreBalance[0])
	preBalance, ok := firstTransaction.PreBalance[0].(float64)
	if !ok {
		fmt.Println("PreBalance is not a string")
		fmt.Printf("not ok : %v\n", ok)
	}
	feesToBeSub, error := strconv.ParseFloat(firstTransaction.Fee, 64)
	if error != nil {
		fmt.Println("Error in parsing float: ", error)
	}

	newBalance := preBalance - feesToBeSub
	fmt.Printf("newBalance: %.f\n", newBalance)
	fmt.Printf("PostBalance: %.f\n", firstTransaction.PostBalance[0])
	fmt.Println(firstTransaction.PostBalance[0] == newBalance)


	fmt.Printf("TxSignature: %s\n", firstTransaction.TxSignature)
	fmt.Printf("PreBalance: %v\n", firstTransaction.PreBalance)
	fmt.Printf("PostBalance: %v\n", firstTransaction.PostBalance)
	fmt.Printf("AccountKeys: %v\n", firstTransaction.AccountKeys)
	fmt.Printf("Instructions[0] Account : %v\n", svmInstructions[0].Accounts)
	fmt.Printf("Instructions[0] Data : %s\n", svmInstructions[0].Data)
	fmt.Printf("Instructions[0] ProgramIDIndex : %d\n", svmInstructions[0].ProgramIDIndex)
	fmt.Printf("RecentBlockhash: %s\n", firstTransaction.RecentBlockhash)
	// witnessVector, _, proofByte, pkErr := prover.GenerateProof(values, 1)
	// if pkErr != nil {
	// 	fmt.Println("Error in GenerateProof: ", pkErr)
	// 	os.Exit(0)
	// }

	// // print witness vector and proof
	// fmt.Println("Witness Vector: ", witnessVector)
	// fmt.Println("Proof: ", proofByte)

	// // create files in data folder and then write witness vector and proof
	// modules.WriteWitnessVector(witnessVector)
	// modules.WriteProof(proofByte)
}
