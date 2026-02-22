package sip

type CallInfo struct {
	Id        int    `json:"id"`
	RemoteUri string `json:"remoteUri"`
	State     string `json:"state"`
}
