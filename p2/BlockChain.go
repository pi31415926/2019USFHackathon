package p2

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/sha3"
	"sort"
)

/*
BlockChain structure
*/
type BlockChain struct {
	Chain  map[int32][]Block
	Length int32
}

func NewBlockChain() (blockChain BlockChain) {
	return blockChain
}

/*
Description: This function takes a Height as the argument, returns the list of blocks stored in that Height or None if the Height doesn't exist.
Argument: int32
Return type: []Block
*/
func (blockChain *BlockChain) Get(height int32) (blocks []Block) {
	if blockChain.Length < height {
		return nil
	} else {
		return blockChain.Chain[height]
	}

}

/*
Description: This function takes a block as the argument, use its Height to find the corresponding list in blockchain's Chain map. If the list has already contained that block's Hash, ignore it because we don't store duplicate blocks; if not, insert the block into the list.
Argument: block
*/
func (blockChain *BlockChain) Insert(block Block) {
	if blockChain.Chain == nil {
		blockChain.Chain = make(map[int32][]Block)
	}
	targetBlocks := blockChain.Chain[block.Header.Height]
	hasHash := false
	for i := range targetBlocks {
		if targetBlocks[i].Header.Hash == block.Header.Hash {
			hasHash = true
		}
	}
	if hasHash {
		return
	} else {
		if targetBlocks == nil {
			targetBlocks = []Block{}
		}
		blockChain.Chain[block.Header.Height] = append(targetBlocks, block)
		if block.Header.Height > blockChain.Length {
			blockChain.Length += block.Header.Height
		}

	}
}

/*
Description: This function iterates over all the blocks, generate blocks' JsonString by the function you implemented previously, and return the list of those JsonStritgns.
Return type: string
*/
func (blockChain *BlockChain) EncodeToJson() (jsonString string, err error) {
	resultString := "["
	position := 0
	for _, v := range blockChain.Chain {
		for i := range v {
			block := v[i]
			if position == 0 {
				resultString = resultString + EncodeToJSON(&block)
			} else {
				resultString = resultString + ", " + EncodeToJSON(&block)
			}
		}
		position += 1
	}
	resultString = resultString + "]"
	return resultString, nil
}

/*
Description: This function is called upon a blockchain instance. It takes a blockchain JSON string as input, decodes the JSON string back to a list of block JSON strings, decodes each block JSON string back to a block instance, and inserts every block into the blockchain.
Argument: self, string
*/
func DecodeFromJSON(jsonString string) (blockChain BlockChain, err error) {
	var strArray []interface{}
	err = json.Unmarshal([]byte(jsonString), &strArray)
	if err == nil {
		for i := range strArray {
			bytesArray, err := json.Marshal(strArray[i])
			newBlock := DecodeFromJson(string(bytesArray))
			if err == nil {
				blockChain.Insert(newBlock)
			}
		}

	}

	return blockChain, err
}

func (bc *BlockChain) Show() string {
	rs := ""
	var idList []int
	for id := range bc.Chain {
		idList = append(idList, int(id))
	}
	sort.Ints(idList)
	for _, id := range idList {
		var hashs []string
		for _, block := range bc.Chain[int32(id)] {
			hashs = append(hashs, block.Header.Hash+"<="+block.Header.ParentHash)
		}
		sort.Strings(hashs)
		rs += fmt.Sprintf("%v: ", id)
		for _, h := range hashs {
			rs += fmt.Sprintf("%s, ", h)
		}
		rs += "\n"
	}
	sum := sha3.Sum256([]byte(rs))
	rs = fmt.Sprintf("This is the BlockChain: %s\n", hex.EncodeToString(sum[:])) + rs
	return rs
}
