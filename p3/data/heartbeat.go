package data

import (
	"../../p1"
	"math/rand"
	"time"
)

type HeartBeatData struct {
	IfNewBlock  bool   `json:"ifNewBlock"`
	Id          int32  `json:"id"`
	BlockJson   string `json:"blockJson"`
	PeerMapJson string `json:"peerMapJson"`
	Addr        string `json:"addr"`
	Hops        int32  `json:"hops"`
}

func NewHeartBeatData(ifNewBlock bool, id int32, blockJson string, peerMapJson string, addr string) HeartBeatData {
	var heartBeatData HeartBeatData

	heartBeatData.PeerMapJson = peerMapJson
	heartBeatData.Addr = addr
	heartBeatData.BlockJson = blockJson
	heartBeatData.Id = id
	heartBeatData.IfNewBlock = ifNewBlock
	return heartBeatData
}

func PrepareHeartBeatData(sbc *SyncBlockChain, selfId int32, peerMapBase64 string, addr string) HeartBeatData {
	var heartBeatData HeartBeatData
	var err error
	//TODO: check how sbc link to BeatData
	if rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100)%2 == 0 {
		mpt := p1.MerklePatriciaTrie{}
		mpt.Initial()
		mpt.Insert("0", "root")
		newBlock := sbc.GenBlock(mpt)

		sbc.bc.Insert(newBlock)
		heartBeatData.IfNewBlock = true
	}

	heartBeatData.BlockJson, err = sbc.bc.EncodeToJson()
	if err != nil {
		panic(err)
	}
	heartBeatData.Id = selfId
	heartBeatData.PeerMapJson = peerMapBase64
	heartBeatData.Addr = addr
	return heartBeatData
}
