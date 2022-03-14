package ergo

type Err struct {
	Error  int    `json:"error"`
	Reason string `json:"reason"`
	Detail string `json:"detail"`
}
