package models

type TxAndLog struct {
	Hash string 
	From string 
	To string 
	Nonce uint64
	Value string
	Input_Data string 
	Data string
}

type TxAndLogs struct {
	Hash string `json:"tx_hash"`
	From string `json:"from"`
	To string `json:"to"`
	Nonce uint64 `json:"nonce"`
	Value string `json:"value"`
	Input_Data string `json:"data"`
	Logs []EventLog `json:"logs"`
}

