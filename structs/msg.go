package structs

type Msg struct {
	Data []MsgData `json:"data"`
	Type string    `json:"type"`
}

type MsgData struct {
	P float64 `json:"p"` //price
	T float64 `json:"t"` //timestamp
	S string  `json:"s"` //symbol
}
