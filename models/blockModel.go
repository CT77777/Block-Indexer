package models

type BlockHeader struct {
	Number    uint64   `json:"block_num"`
	Hash   string   `json:"block_hash"`
	Time   uint64   `json:"block_time"`
	Parent_Hash  string   `json:"parent_hash"`
}

type BlockAndTx struct {
	BlockHeader
	Transaction string
}

type BlockAndTxs struct {
	BlockHeader
	Transactions []string `json:"transactions"`
}