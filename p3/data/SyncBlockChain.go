package data

import (
	"../../p2"
	"fmt"
	"sync"
)

type SyncBlockChain struct {
	bc  p2.BlockChain
	mux sync.Mutex
}

func NewBlockChain() SyncBlockChain {
	return SyncBlockChain{bc: p2.NewBlockChain()}
}

func (sbc *SyncBlockChain) Get(height int32) ([]p2.Block, bool) {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return sbc.bc.Get(height), true
}

func (sbc *SyncBlockChain) GetBlock(height int32, hash string) (p2.Block, bool) {
	blocks := sbc.bc.Get(height)

	for blockIndex := range blocks {
		if blocks[blockIndex].Header.Hash == hash {
			return blocks[blockIndex], true
		}
	}
	return p2.Block{}, false
}

func (sbc *SyncBlockChain) Insert(block p2.Block) {
	sbc.mux.Lock()
	sbc.bc.Insert(block)
	sbc.mux.Unlock()
}

func (sbc *SyncBlockChain) CheckParentHash(insertBlock p2.Block) bool {
	return insertBlock.Header.ParentHash == sbc.bc.Get(sbc.bc.Length - 1)[0].Header.Hash
}

func (sbc *SyncBlockChain) UpdateEntireBlockChain(blockChainJson string) {
	blockChain, err := p2.DecodeFromJSON(blockChainJson)
	if err == nil {
		for k, v := range blockChain.Chain {
			sbc.bc.Chain[k] = v
		}

	}

}

func (sbc *SyncBlockChain) BlockChainToJson() (string, error) {
	return sbc.bc.EncodeToJson()
}

func (sbc *SyncBlockChain) Show() string {
	return sbc.bc.Show()
}

func PrintError(err error, msg string) {
	fmt.Println(msg, err)
}
