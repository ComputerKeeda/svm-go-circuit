package prover

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/ComputerKeeda/svm-go-circuit/types"
	"github.com/consensys/gnark-crypto/ecc"
	// tedwards "github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/btcsuite/btcutil/base58"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	// "github.com/consensys/gnark/std/algebra/native/twistededwards"
)

const BatchSize = 25

type MyCircuit struct {
	To              [BatchSize]frontend.Variable `gnark:",public"`
	From            [BatchSize]frontend.Variable `gnark:",public"`
	Amount          [BatchSize]frontend.Variable `gnark:",public"`
	TransactionHash [BatchSize]frontend.Variable `gnark:",public"`
	SenderBalance   [BatchSize]frontend.Variable `gnark:",public"`
	ReceiverBalance [BatchSize]frontend.Variable `gnark:",public"`
	Messages        [BatchSize]frontend.Variable `gnark:",public"`
}

func (circuit *MyCircuit) Define(api frontend.API) error {
	for i := 0; i < BatchSize; i++ {

		api.AssertIsLessOrEqual(circuit.Amount[i], circuit.SenderBalance[i])

		updatedFromBalance := api.Sub(circuit.SenderBalance[i], circuit.Amount[i])
		updatedToBalance := api.Add(circuit.ReceiverBalance[i], circuit.Amount[i])

		api.AssertIsEqual(updatedFromBalance, api.Sub(circuit.SenderBalance[i], circuit.Amount[i]))
		api.AssertIsEqual(updatedToBalance, api.Add(circuit.ReceiverBalance[i], circuit.Amount[i]))
	}

	return nil
}

func ComputeCCS() constraint.ConstraintSystem {
	var circuit MyCircuit
	ccs, _ := frontend.Compile(ecc.BLS12_381.ScalarField(), r1cs.NewBuilder, &circuit)

	return ccs
}

func GenerateVerificationKey() (groth16.ProvingKey, groth16.VerifyingKey, error) {
	ccs := ComputeCCS()
	pk, vk, error := groth16.Setup(ccs)
	return pk, vk, error
}

func getTransactionHash(tx types.SVMTransaction) string {
	record := tx.To + tx.From + tx.Amount + tx.SenderBalance + tx.ReceiverBalance + tx.TransactionHash + tx.Message
	h := sha256.New()
	h.Write([]byte(record))
	return hex.EncodeToString(h.Sum(nil))
}

func GetMerkleRoot(transactions []types.SVMTransaction) string {
	var merkleTree []string

	for _, tx := range transactions {
		merkleTree = append(merkleTree, getTransactionHash(tx))
	}

	for len(merkleTree) > 1 {
		var tempTree []string
		for i := 0; i < len(merkleTree); i += 2 {
			if i+1 == len(merkleTree) {
				tempTree = append(tempTree, merkleTree[i])
			} else {
				combinedHash := merkleTree[i] + merkleTree[i+1]
				h := sha256.New()
				h.Write([]byte(combinedHash))
				tempTree = append(tempTree, hex.EncodeToString(h.Sum(nil)))
			}
		}
		merkleTree = tempTree
	}

	return merkleTree[0]
}

func GenerateProof(values types.SVMBatchStruct, podNumber int) (any, string, []byte, error) {
	ccs := ComputeCCS()
	// snarkField, err := twistededwards.GetSnarkField(tedwards.BLS12_381)
	// if err != nil {
	// 	fmt.Println("Error getting snark field")
	// 	return nil, "", nil, err
	// }
	var transactions []types.SVMTransaction
	for i := 0; i < BatchSize; i++ {
		decodedFrom := base58.Decode(values.From[i])
		decodedTo := base58.Decode(values.To[i])
		decodedTransactionHash := base58.Decode(values.TransactionHash[i])
		// decodedMessage := base58.Decode(values.Messages[i])
		transaction := types.SVMTransaction{
			To:              string(decodedTo),
			From:            string(decodedFrom),
			Amount:          values.Amounts[i],
			SenderBalance:   values.SenderBalances[i],
			ReceiverBalance: values.ReceiverBalances[i],
			TransactionHash: string(decodedTransactionHash),
		}
		transactions = append(transactions, transaction)
	}

	currentStatusHash := GetMerkleRoot(transactions)
	pk, err := ReadProvingKeyFromFile("provingKey.txt")
	if err != nil {
		fmt.Println("Error reading proving key:", err)
		return nil, "", nil, err
	}

	var inputValueLength int

	fromLength := len(values.From)
	toLength := len(values.To)
	amountsLength := len(values.Amounts)
	txHashLength := len(values.TransactionHash)
	senderBalancesLength := len(values.SenderBalances)
	receiverBalancesLength := len(values.ReceiverBalances)
	messagesLength := len(values.Messages)

	if fromLength == toLength &&
		fromLength == amountsLength &&
		fromLength == txHashLength &&
		fromLength == senderBalancesLength &&
		fromLength == receiverBalancesLength &&
		fromLength == messagesLength {
		inputValueLength = fromLength
	} else {
		fmt.Println("Error: Input data is not correct")
		return nil, "", nil, fmt.Errorf("input data is not correct")
	}

	if inputValueLength < BatchSize {
		leftOver := BatchSize - inputValueLength
		for i := 0; i < leftOver; i++ {
			values.From = append(values.From, "0")
			values.To = append(values.To, "0")
			values.Amounts = append(values.Amounts, "0")
			values.TransactionHash = append(values.TransactionHash, "0")
			values.SenderBalances = append(values.SenderBalances, "0")
			values.ReceiverBalances = append(values.ReceiverBalances, "0")
			values.Messages = append(values.Messages, "0")
		}
	}

	inputs := MyCircuit{
		To:              [BatchSize]frontend.Variable{},
		From:            [BatchSize]frontend.Variable{},
		Amount:          [BatchSize]frontend.Variable{},
		TransactionHash: [BatchSize]frontend.Variable{},
		SenderBalance:   [BatchSize]frontend.Variable{},
		ReceiverBalance: [BatchSize]frontend.Variable{},
		Messages:        [BatchSize]frontend.Variable{},
	}

	for i := 0; i < BatchSize; i++ {

		decodedFrom := base58.Decode(values.From[i])
		decodedTo := base58.Decode(values.To[i])
		decodedTransactionHash := base58.Decode(values.TransactionHash[i])
		decodedMessage := base58.Decode(values.Messages[i])

		inputs.To[i] = frontend.Variable(decodedTo)
		inputs.From[i] = frontend.Variable(decodedFrom)
		inputs.Amount[i] = frontend.Variable(values.Amounts[i])
		inputs.TransactionHash[i] = frontend.Variable(decodedTransactionHash)
		inputs.SenderBalance[i] = frontend.Variable(values.SenderBalances[i])
		inputs.ReceiverBalance[i] = frontend.Variable(values.ReceiverBalances[i])
		msg := decodedMessage
		// msg := []byte(values.Messages[i])
		// msg := make([]byte, len(snarkField.Bytes()))
		inputs.Messages[i] = msg
	}

	witness, err := frontend.NewWitness(&inputs, ecc.BLS12_381.ScalarField())
	if err != nil {
		fmt.Printf("Error creating a witness: %v\n", err)
		return nil, "", nil, err
	}

	witnessVector := witness.Vector()

	proof, err := groth16.Prove(ccs, pk, witness)
	if err != nil {
		fmt.Printf("Error generating proof: %v\n", err)
		return nil, "", nil, err
	}

	proofByte, err := json.Marshal(proof)
	if err != nil {
		fmt.Println("Error marshalling proof:", err)
		return nil, "", nil, err
	}

	return witnessVector, currentStatusHash, proofByte, nil
}

func ReadProvingKeyFromFile(filename string) (groth16.ProvingKey, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	pk := groth16.NewProvingKey(ecc.BLS12_381)
	_, err = pk.ReadFrom(file)
	if err != nil {
		return nil, err
	}

	return pk, nil
}
