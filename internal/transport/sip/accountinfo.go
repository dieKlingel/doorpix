package sip

type AccountInfo struct {
	Uri        string `json:"uri"`
	IsActive   bool   `json:"isActive"`
	StatusText string `json:"statusText"`
}
