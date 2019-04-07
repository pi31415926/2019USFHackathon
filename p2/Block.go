package p2

import (
	"encoding/json"
	"fmt"
)

/*
Block structure
*/
type Block struct {
	Header *Header

	InsertedMap map[string]string
}

/*
 Header structure
*/
type Header struct {
	Height     int32
	Timestamp  int64
	Hash       string
	ParentHash string
}

/*
convert json to BlockJson structure
*/
type BlockJson struct {
	Hash       string            `json:"Hash"`
	Height     int32             `json:"Height"`
	ParentHash string            `json:"ParentHash"`
	Mpt        map[string]string `json:"mpt"`
	TimeStamp  int64             `json:"timeStamp"`
}

/*
Description: This function takes arguments(such as Height, ParentHash, and value of MPT type) and forms a block. This is a method of the block struct.
*/
func (block *Block) Initial(Hash string, Height int32, ParentHash string, Mpt map[string]string, TimeStamp int64) {
	block.Header = new(Header)
	block.Header.Hash = Hash
	block.Header.Height = Height
	block.Header.ParentHash = ParentHash
	block.Header.Timestamp = TimeStamp
	block.InsertedMap = Mpt
}

/*
Description: This function takes a string that represents the JSON value of a block as an input, and decodes the input string back to a block instance. Note that you have to reconstruct an MPT from the JSON string, and use that MPT as the block's value.
Argument: a string of JSON format
Return value: a block instance
*/
func DecodeFromJson(jsonString string) (block Block) {
	var blockJson BlockJson

	err := json.Unmarshal([]byte(jsonString), &blockJson)
	if err == nil {
		block.Initial(blockJson.Hash, blockJson.Height, blockJson.ParentHash, blockJson.Mpt, blockJson.TimeStamp)
		return block
	} else {
		fmt.Println(err)
		return block
	}

}

/*
Description: This function encodes a block instance into a JSON format string. Note that the block's value is an MPT, and you have to record all of the (key, value) pairs that have been inserted into the MPT in your JSON string.
*/
func EncodeToJSON(block *Block) (jsonString string) {
	var blockJson BlockJson
	blockJson.Hash = block.Header.Hash
	blockJson.Height = block.Header.Height
	blockJson.ParentHash = block.Header.ParentHash
	blockJson.TimeStamp = block.Header.Timestamp
	blockJson.Mpt = block.InsertedMap

	byteJson, err := json.Marshal(blockJson)
	if err == nil {
		return string(byteJson)
	} else {
		fmt.Println(err)
		return ""
	}
}
