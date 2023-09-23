package models

type Block struct {
	Number uint64 `gorm:"primaryKey;autoIncrement:false"`
	Hash string `gorm:"type:VARCHAR(66);index"`
	Time uint64
	Parent_Hash string `gorm:"type:VARCHAR(66)"`
	Transactions []Transaction `gorm:"foreignKey:Block_Number;references:Number"`
}

type Transaction struct {
	ID uint64 `gorm:"primaryKey"`
	Hash string `gorm:"type:VARCHAR(66);index"`
	From string `gorm:"type:VARCHAR(42)"`
	To string `gorm:"type:VARCHAR(42)"`
	Nonce uint64
	Value string
	Input_Data string `gorm:"type:LONGTEXT"`
	Block_Number uint64
	Logs []Log `gorm:"foreignKey:Transaction_ID;references:ID"`
}

type Log struct {
	Address string `gorm:"type:VARCHAR(42);index"`
	Topic_0 string `gorm:"type:VARCHAR(66)"`
	Topic_1 string `gorm:"type:VARCHAR(66)"`
	Topic_2 string `gorm:"type:VARCHAR(66)"`
	Topic_3 string `gorm:"type:VARCHAR(66)"`
	Data string `gorm:"type:LONGTEXT"`
	Transaction_ID uint64 
}