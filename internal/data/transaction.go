package data

type OutMessage struct {
	Message string `json:"message,omitempty"`
	Value   string `json:"value,omitempty"`
}

type Transaction struct {
	OutMessage []OutMessage `json:"out_msgs,omitempty"`
}

type Response struct {
	Ok     bool          `json:"ok"`
	Result []Transaction `json:"result,omitempty"`
	Error  string        `json:"error,omitempty"`
	Code   int32         `json:"code,omitempty"`
}
