package types

type SVMBatchStruct struct {
	From             []string `json:"from_txn"`
	To               []string `json:"to_txn"`
	Amounts          []string `json:"amounts_txn"`
	TransactionHash  []string `json:"transaction_hash_txn"`
	SenderBalances   []string `json:"sender_balances_txn"`
	ReceiverBalances []string `json:"receiver_balances_txn"`
	Messages         []string `json:"messages_txn"`
}

type SVMTransaction struct {
	From            string
	To              string
	Amount          string
	TransactionHash string
	SenderBalance   string
	ReceiverBalance string
	Message         string
}