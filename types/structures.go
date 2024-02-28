package types

// type SVMBatchStruct struct {
// 	From             []string `json:"from_txn"`
// 	To               []string `json:"to_txn"`
// 	Amounts          []string `json:"amounts_txn"`
// 	TransactionHash  []string `json:"transaction_hash_txn"`
// 	SenderBalances   []string `json:"sender_balances_txn"`
// 	ReceiverBalances []string `json:"receiver_balances_txn"`
// 	Messages         []string `json:"messages_txn"`
// }

type SVMPodStruct struct {
	TxSignature  []string        `json:"tx_signature"`
	Fee          []string        `json:"fee"`
	PreBalance   [][]interface{} `json:"pre_balance"`
	PostBalance  [][]interface{} `json:"post_balance"`
	AccountKeys  [][]string      `json:"account_keys"`
	Instructions [][]struct {
		Accounts       []int       `json:"accounts"`
		Data           string      `json:"data"`
		ProgramIDIndex int         `json:"programIdIndex"`
		StackHeight    interface{} `json:"stackHeight"`
	} `json:"instructions"`
	RecentBlockhash []string `json:"recent_blockhash"`
}
type SVMTransaction struct {
	TxSignature  string
	Fee          string
	PreBalance   []interface{}
	PostBalance  []interface{}
	AccountKeys  []string
	Instructions []SVMInstruction
	RecentBlockhash string
}

type SVMInstruction struct {
	Accounts       []int
	Data           string
	ProgramIDIndex int
	StackHeight    interface{}
}