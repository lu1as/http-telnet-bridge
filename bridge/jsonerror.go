package bridge

import (
	"encoding/json"
)

type jsonError struct {
	Err string `json:"error"`
}

func JsonError(err string) string {
	e := &jsonError{
		Err: err,
	}
	j, _ := json.Marshal(e)
	return string(j)
}
