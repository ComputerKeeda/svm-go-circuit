package main

import (
	"fmt"
	"os"

	"github.com/ComputerKeeda/svm-go-circuit/modules"
	"github.com/ComputerKeeda/svm-go-circuit/prover"
)

func main() {
	prover.CreateVkPk()
	values := modules.FetchJsonData()

	witnessVector, _, proofByte, pkErr := prover.GenerateProof(values, 1)
	if pkErr != nil {
		fmt.Println("Error in GenerateProof: ", pkErr)
		os.Exit(0)
	}

	// print witness vector and proof
	fmt.Println("Witness Vector: ", witnessVector)
	fmt.Println("Proof: ", proofByte)

	// create files in data folder and then write witness vector and proof
	modules.WriteWitnessVector(witnessVector)
	modules.WriteProof(proofByte)
}
