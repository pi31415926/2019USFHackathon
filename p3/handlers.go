package p3

import (
	"../p2"
	"./data"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var TA_SERVER = "http://localhost:6688"
var REGISTER_SERVER = TA_SERVER + "/peer"
var BC_DOWNLOAD_SERVER = TA_SERVER + "/upload"
var SELF_ADDR = "http://localhost:6686"

var SBC data.SyncBlockChain
var Peers data.PeerList
var ifStarted bool

func init() {
	// This function will be executed before everything else.
	// Do some initialization here.
	SBC = data.NewBlockChain()
	ifStarted = false
}

// Register ID, download BlockChain, start HeartBeat
func Start(w http.ResponseWriter, r *http.Request) {
	if ifStarted == false {
		ifStarted = true
		Register()
		Peers = data.NewPeerList(Peers.GetSelfId(), 32)
		Download()
		StartHeartBeat()
	}
}

// Display peerList and sbc
func Show(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s\n%s", Peers.Show(), SBC.Show())
}

// Register to TA's server, get an ID
func Register() {
	var client http.Client
	resp, err := client.Get(REGISTER_SERVER)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		if err == nil {
			id, err := strconv.ParseInt(bodyString, 10, 32)
			if err != nil {
				panic(err)
			}
			Peers.Register(int32(id)) //TODO: check if correct
		}
	}
}

// Download blockchain from TA server
func Download() {
	var client http.Client
	resp, err := client.Get(BC_DOWNLOAD_SERVER)
	if err != nil {
		log.Print(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		if err == nil {
			SBC.UpdateEntireBlockChain(bodyString) //TODO: check if correct

		}
	}
}

// Upload blockchain to whoever called this method, return jsonStr
func Upload(w http.ResponseWriter, r *http.Request) {
	blockChainJson, err := SBC.BlockChainToJson()
	if err != nil {
		data.PrintError(err, "Upload")
	}
	comingIp := r.Header.Get("X-FORWARDED-FOR")
	Peers.Add(comingIp, 0)
	fmt.Fprint(w, blockChainJson)
}

// Upload a block to whoever called this method, return jsonStr
func UploadBlock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	heightStr := vars["height"]
	hash := vars["hash"]
	heightIn64, err := strconv.ParseInt(heightStr, 10, 32)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, "InternalServerError")
	}
	height := int32(heightIn64)
	block, hasFound := SBC.GetBlock(height, hash)
	if hasFound {
		fmt.Fprint(w, p2.EncodeToJSON(&block))
	} else {
		w.WriteHeader(204)
		fmt.Fprint(w, "StatusNoContent")
	}

}

// Received a heartbeat
func HeartBeatReceive(w http.ResponseWriter, r *http.Request) {

	var heartBeat data.HeartBeatData
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
		return
	}
	err = json.Unmarshal(bodyBytes, &heartBeat)
	if err != nil {
		panic(err)
	}
	//TODO: check if correct
	comingIp := r.Header.Get("X-FORWARDED-FOR")
	Peers.InjectPeerMapJson(heartBeat.PeerMapJson, comingIp, heartBeat.Id)

	if heartBeat.IfNewBlock {
		blockChain, err := p2.DecodeFromJSON(heartBeat.BlockJson)
		_, hasFound := SBC.GetBlock(blockChain.Length-1, blockChain.Chain[blockChain.Length][0].Header.ParentHash)
		if !hasFound {
			AskForBlock(blockChain.Length-1, blockChain.Chain[blockChain.Length][0].Header.ParentHash)
		}
		SBC.UpdateEntireBlockChain(heartBeat.BlockJson)
		peerMapToJson, err := Peers.PeerMapToJson()
		if err != nil {
			panic(err)
		}
		ForwardHeartBeat(data.NewHeartBeatData(true, Peers.GetSelfId(), heartBeat.BlockJson, peerMapToJson, SELF_ADDR))
	}
}

// Ask another server to return a block of certain height and hash
func AskForBlock(height int32, hash string) {

	for address, _ := range Peers.Copy() {
		resp, err := http.Get(address + "/block/" + string(height) + "/" + hash)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}
			var resultBlock p2.Block
			err = json.Unmarshal(bodyBytes, &resultBlock)
			if err != nil {
				panic(err)
			}
			SBC.Insert(resultBlock)
		}
	}

}

func ForwardHeartBeat(heartBeatData data.HeartBeatData) {
	heartBeatDataJson, _ := json.Marshal(heartBeatData)
	Peers.Rebalance()
	peerMap := Peers.Copy()
	for address, _ := range peerMap {
		http.Post(address+"/heartbeat/receive", "application/json; charset=UTF-8", strings.NewReader(string(heartBeatDataJson)))
	}
}

func StartHeartBeat() {
	if len(Peers.Copy()) != 0 {

		Peers.Rebalance()
		peerMap := Peers.Copy()
		for address, id := range peerMap {
			StartToSendHeartBeat(address, id)
		}
		time.Sleep(10 * time.Second)
		for true {
			Peers.Rebalance()
			peerMap := Peers.Copy()
			peerMapJson, err := Peers.PeerMapToJson()
			if err != nil {
				panic(err)
			}
			for address, _ := range peerMap {

				SendHeartBeat(address, Peers.GetSelfId(), peerMapJson)
			}
			time.Sleep(10 * time.Second)
		}
	}else{
		for true {
			if len(Peers.Copy()) != 0 {
				StartHeartBeat()
			}
			time.Sleep(10 * time.Second)
		}
	}
}

func SendHeartBeat(address string, selfId int32, peerMapBase64 string) {
	heartBeatDataJson, _ := json.Marshal(data.PrepareHeartBeatData(&SBC, selfId, peerMapBase64, address))
	resp, err := http.Post(address+"/heartbeat/receive", "application/json; charset=UTF-8", strings.NewReader(string(heartBeatDataJson)))

	bytes, _ := ioutil.ReadAll(resp.Body)
	var rData data.RegisterData
	err = json.Unmarshal(bytes, &rData)
	if err != nil {
		data.PrintError(err, "HeartBeat")
		return
	}
	SBC.UpdateEntireBlockChain(rData.PeerMapJson)
}

func StartToSendHeartBeat(address string, id int32) {
	heartBeatDataJson, _ := json.Marshal(data.NewHeartBeatData(false, id, "", "", SELF_ADDR)) //TODO: check NewHeartBeatData
	_, err := http.Post(address+"/heartbeat/receive", "application/json; charset=UTF-8", strings.NewReader(string(heartBeatDataJson)))

	if err != nil {
		data.PrintError(err, "SendHeartBeat")
		return
	}
}
