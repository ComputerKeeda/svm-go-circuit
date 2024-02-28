package modules

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ComputerKeeda/svm-go-circuit/types"
)

func FetchJsonData() types.SVMPodStruct {
	file, err := os.ReadFile("data/person.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(3)
	}

	var dv types.SVMPodStruct
	err = json.Unmarshal(file, &dv)
	if err != nil {
		fmt.Println("Error unmarshalling:", err)
		os.Exit(3)
	}
	return dv
}
