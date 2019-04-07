package data

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
