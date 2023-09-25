package models

type EventLog struct {
	Index uint8 `json:"index"`
	Data string `json:"data"`
}