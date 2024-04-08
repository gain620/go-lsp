package rpc

import (
	"encoding/json"
	"fmt"
)

func EncodeMessage(msg any) string {
	content, _ := json.Marshal(msg)
	return fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content)
}
