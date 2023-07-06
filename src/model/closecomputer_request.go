package model

type CloseComputerRequest struct {
	IpAddr   string `json:"ipaddr"`
	Username string `json:"username"`
	Password string `json:"password"`
}
