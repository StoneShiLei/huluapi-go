package model

type OpenComputerRequest struct {
	MacAddr   string `json:"macAddr"`
	BcastAddr string `json:"bcastAddr"`
}
