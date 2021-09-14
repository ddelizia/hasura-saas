package gqlreq

type InsertBaseResponse struct {
	Returning    []InsertBaseResponseID `json:"returning"`
	AffectedRows int                    `json:"affected_rows"`
}

type InsertBaseResponseID struct {
	ID string `json:"id"`
}

type RequestHeader struct {
	Key   string
	Value string
}

type RequestVar struct {
	Key   string
	Value interface{}
}

type HeaderInfo struct {
	AccountId string
	UserId    string
	Role      string
}
