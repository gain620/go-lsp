package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

func EncodeMessage(msg any) string {
	content, _ := json.Marshal(msg)
	return fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content)
}

type BaseMessage struct {
	Method string `json:"method"`
}

func DecodeMessage(msg []byte) (string, []byte, error) {
	header, content, found := bytes.Cut(msg, []byte{'\r', '\n', '\r', '\n'})
	if !found {
		return "", nil, errors.New("separator not found")
	}

	// Content-Length: num
	contentLengthBytes := bytes.TrimPrefix(header, []byte("Content-Length: "))
	contentLength, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return "", nil, err
	}

	var baseMsg BaseMessage
	if err = json.Unmarshal(content[:contentLength], &baseMsg); err != nil {
		return "", nil, err
	}

	return baseMsg.Method, content[:contentLength], nil
}
