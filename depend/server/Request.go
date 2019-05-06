package server

type (

	RequestMsgBody struct {
		S  string `json:"s"`            // type struct
		M  string `json:"m"`            // method
		V  string `json:"v"`            // version
		P  string `json:"p"`            // Body
		Md ProcessUid `json:"md,omitempty"` // some_param
	}

	ProcessUid struct {
		Uid    int `json:"uid"`
		BakUid int `json:"bak_uid,omitempty"`
	}

)