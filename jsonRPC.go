package portal

import (
	"encoding/json"
	"fmt"
)

type jsonRPCRequest struct {
	Method string          `json:"method"`
	ID     json.RawMessage `json:"id"`
	Params json.RawMessage `json:"params"`

	// Non-Standard authorization token
	Auth string `json:"auth"`
}

// type jsonRPCRequest struct {
// 	Method string        `json:"method"`
// 	Params []interface{} `json:"params"`
// 	ID     string        `json:"id"`
// }

type jsonRPCRersult struct {
	RawResult json.RawMessage `json:"result"`
	// RawError  json.RawMessage `json:"error"`
	Error *jsonRPCError `json:"error"`
	ID    string        `json:"id"`
}

type jsonRPCError struct {
	Code    int
	Message string
}

func (err *jsonRPCError) Error() string {
	return fmt.Sprintf("[code: %d] %s", err.Code, err.Message)
}
