package compiler

import (
	"encoding/base64"
	"encoding/json"
)

type (
	Auth struct {
		Address  string
		Username string
		Password string
	}
)

func Header(username, password string) string {
	v := struct {
		Username string `json:"username,omitempty"`
		Password string `json:"password,omitempty"`
	}{
		Username: username,
		Password: password,
	}
	buf, _ := json.Marshal(&v)
	return base64.URLEncoding.EncodeToString(buf)
}
