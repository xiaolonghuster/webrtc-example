package srs

type SdpRequest struct {
	Api       string `json:"api"`
	ClientIp  string `json:"clentip"`
	StreamURL string `json:"streamurl"`
	Sdp       string `json:"sdp"`
}

type SdpResponse struct {
	Code      int    `json:"code"`
	Server    string `json:"server"`
	SessionId string `json:"sessionid"`
	Sdp       string `json:"sdp"`
}
