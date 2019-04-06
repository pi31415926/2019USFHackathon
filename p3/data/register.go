package data

import (
	"encoding/json"
)

type RegisterData struct {
	AssignedId  int32  `json:"assignedId"`
	PeerMapJson string `json:"peerMapJson"`
}

func NewRegisterData(id int32, peerMapJson string) (RegisteredData RegisterData) {
	RegisteredData.AssignedId = id
	RegisteredData.PeerMapJson = peerMapJson
	return RegisteredData
}

func (data *RegisterData) EncodeToJson() (string, error) {
	byteJson, err := json.Marshal(data)
	return string(byteJson), err
}
