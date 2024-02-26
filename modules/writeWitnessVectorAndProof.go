package modules

import (
	"fmt"
	"os"
)

func WriteWitnessVector(witnessVector any) {
	// write witness vector to witnessVector.json
	witnessVectorFile, err := os.Create("data/witnessVector.json")
	if err != nil {
		fmt.Println("Error creating witnessVector.json: ", err)
		os.Exit(0)
	}

	_, err = witnessVectorFile.WriteString(fmt.Sprintf("%v", witnessVector))
	if err != nil {
		fmt.Println("Error writing witnessVector.json: ", err)
		os.Exit(0)
	}
}

func WriteProof(proofByte []byte) {
	// write proof to proof.json
	proofFile, err := os.Create("data/proof.json")
	if err != nil {
		fmt.Println("Error creating proof.json: ", err)
		os.Exit(0)
	}

	_, err = proofFile.Write(proofByte)
	if err != nil {
		fmt.Println("Error writing proof.json: ", err)
		os.Exit(0)
	}
}