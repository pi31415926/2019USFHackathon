package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
)

type PeerList struct {
	selfId  int32
	peerMap map[string]int32
	Num     int
	mux     sync.Mutex
}

type PeerPair struct {
	Addr string `json:"addr"`
	Id   int32  `json:"id"`
}

func NewPeerList(id int32) PeerList {
	var peers PeerList
	peers.selfId = id
	peers.Num = 0
	peers.peerMap = make(map[string]int32)
	return peers
}

func (peers *PeerList) Add(addr string, id int32) {
	peers.mux.Lock()
	peers.peerMap[addr] = id
	peers.Num = peers.Num + 1
	peers.mux.Unlock()
}

func (peers *PeerList) IsEmpty() bool {
	peers.mux.Lock()
	isEmpty := (peers.Num == 1)
	peers.mux.Unlock()
	return isEmpty
}

func (peers *PeerList) Delete(addr string) {
	peers.mux.Lock()
	if _, ok := peers.peerMap[addr]; ok {
		delete(peers.peerMap, addr)
	}
	peers.Num = peers.Num - 1
	peers.mux.Unlock()
}

func (peers *PeerList) Show() string {
	strBuffer := new(bytes.Buffer)
	for key, value := range peers.peerMap {
		fmt.Fprintf(strBuffer, "%s=\"%s\"\n", key, value)
	}
	return strBuffer.String()
}

func (peers *PeerList) Register(id int32) {
	peers.selfId = id
	fmt.Printf("SelfId=%v\n", id)
}

func (peers *PeerList) Copy() map[string]int32 {
	peers.mux.Lock()
	var peerMapCopy map[string]int32
	for k, v := range peers.peerMap {
		peerMapCopy[k] = v
	}
	peers.mux.Unlock()
	return peerMapCopy

}

func (peers *PeerList) GetSelfId() int32 {
	return peers.selfId
}

func (peers *PeerList) PeerMapToJson() (string, error) {
	peerMapCopy := peers.Copy()
	byteJson, err := json.Marshal(peerMapCopy)
	if err == nil {
		return string(byteJson), err
	} else {
		return "", err
	}
}

func (peers *PeerList) InjectPeerMapJson(peerMapJsonStr string, senderAddr string, id int32) {
	var strArray []interface{}

	err1 := json.Unmarshal([]byte(peerMapJsonStr), &strArray)
	if err1 == nil {
		for i := range strArray {
			bytesArray, err2 := json.Marshal(strArray[i])
			if err2 == nil {
				var peerPair PeerPair
				err3 := json.Unmarshal(bytesArray, &peerPair)
				if err3 == nil {
					peers.Add(peerPair.Addr, peerPair.Id)
				}
			}

		}
		peers.Add(senderAddr, id)
	} else {
		fmt.Println(err1)
	}
}

func TestIsEmpty() bool {
	peers := NewPeerList(2)
	return peers.IsEmpty()
}
