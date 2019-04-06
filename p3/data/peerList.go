package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"sync"
)

type PeerList struct {
	selfId    int32
	peerMap   map[string]int32
	maxLength int32
	mux       sync.Mutex
}

type PeerPair struct {
	Addr string `json:"addr"`
	Id   int32  `json:"id"`
}

func NewPeerList(id int32, maxLength int32) PeerList {
	var peers PeerList
	peers.maxLength = maxLength
	peers.selfId = id
	peers.peerMap = make(map[string]int32)
	return peers
}

func (peers *PeerList) Add(addr string, id int32) {
	peers.mux.Lock()
	peers.peerMap[addr] = id
	peers.mux.Unlock()
}

func (peers *PeerList) Delete(addr string) {
	peers.mux.Lock()
	if _, ok := peers.peerMap[addr]; ok {
		delete(peers.peerMap, addr)
	}
	peers.mux.Unlock()
}

func (peers *PeerList) Rebalance() {
	peers.mux.Lock()
	idArray := make([]int32, len(peers.peerMap))
	peerReverseMap := make(map[int32]string)

	i := 0
	for k, v := range peers.peerMap {
		idArray[i] = v
		peerReverseMap[v] = k
		i++
	}
	sort.Sort(Int32(idArray))

	newIdArray := findClosestArray(peers.selfId, idArray, int(peers.maxLength/2))
	newPeerMap := make(map[string]int32)

	for i := range newIdArray {
		newId := newIdArray[i]
		newPeerMap[peerReverseMap[newId]] = newId
	}
	peers.peerMap = newPeerMap
	peers.mux.Unlock()
}

func findClosestArray(target int32, originalArray []int32, rangeLen int) []int32 {
	pivot := len(originalArray)
	for i := range originalArray {
		if originalArray[i] > target {
			pivot = i
			break
		} else if originalArray[i] == target {
			originalArray = append(originalArray[:i], originalArray[i+1:]...)
			pivot = i
			break
		}
	}

	frontIndex := (((pivot - rangeLen) % len(originalArray)) + len(originalArray)) % len(originalArray)
	backIndex := (((pivot + rangeLen) - 1) % len(originalArray)) + len(originalArray)%len(originalArray)

	if frontIndex < backIndex {
		return originalArray[frontIndex : backIndex+1]
	} else if frontIndex > backIndex {
		return append(originalArray[:backIndex+1], originalArray[frontIndex:]...)
	} else {
		if rangeLen == 0 {
			return []int32{}
		} else {
			return originalArray
		}

	}
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

func TestPeerListRebalance() {
	peers := NewPeerList(6, 4)
	peers.Add("1111", 1)
	peers.Add("4444", 4)
	peers.Add("-1-1", -1)
	peers.Add("0000", 0)
	peers.Add("2121", 21)
	peers.Rebalance()
	expected := NewPeerList(6, 4)
	expected.Add("1111", 1)
	expected.Add("4444", 4)
	expected.Add("2121", 21)
	expected.Add("-1-1", -1)
	fmt.Println(reflect.DeepEqual(peers, expected))

	peers = NewPeerList(5, 2)
	peers.Add("1111", 1)
	peers.Add("4444", 4)
	peers.Add("-1-1", -1)
	peers.Add("0000", 0)
	peers.Add("2121", 21)
	peers.Rebalance()
	expected = NewPeerList(5, 2)
	expected.Add("4444", 4)
	expected.Add("2121", 21)
	fmt.Println(reflect.DeepEqual(peers, expected))

	peers = NewPeerList(5, 4)
	peers.Add("1111", 1)
	peers.Add("7777", 7)
	peers.Add("9999", 9)
	peers.Add("11111111", 11)
	peers.Add("2020", 20)
	peers.Rebalance()
	expected = NewPeerList(5, 4)
	expected.Add("1111", 1)
	expected.Add("7777", 7)
	expected.Add("9999", 9)
	expected.Add("2020", 20)
	fmt.Println(reflect.DeepEqual(peers, expected))
}

type Int32 []int32

func (u Int32) Len() int { return len(u) }

func (u Int32) Less(i, j int) bool { return u[i] < u[j] }

func (u Int32) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}
