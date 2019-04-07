package p3

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Show",
		"GET",
		"/show",
		Show,
	},
	Route{
		"Upload",
		"GET",
		"/upload",
		Upload,
	},
	Route{
		"UploadBlock",
		"GET",
		"/block/{height}/{hash}",
		UploadBlock,
	},
	Route{
		"HeartBeatReceive",
		"POST",
		"/heartbeat/receive",
		HeartBeatReceive,
	},
	Route{
		"Start",
		"GET",
		"/start",
		Start,
	},
	Route{
		"DownloadToUsers",
		"GET",
		"/downloadtousers",
		DownloadToUsers,
	},
	Route{
		"SaveNode",
		"POST",
		"/savenode",
		SaveNode,
	},
	Route{
		"AddVote",
		"POST",
		"/addvote",
		AddVote,
	},
	Route{
		"SynchronizeAll",
		"POST",
		"/synchronizeall",
		SynchronizeAll,
	},
}
