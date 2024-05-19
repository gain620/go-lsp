package main

import (
	"bufio"
	"encoding/json"
	"golsp/lsp"
	"golsp/rpc"
	"log"
	"os"
)

func main() {
	logger := getLogger("./go-lsp.log")
	logger.Println("Starting go-lsp")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("Failed to decode message: %v\n", err)
			continue
		}

		handleMessage(method, contents, logger)
	}
}

func handleMessage(method string, contents []byte, logger *log.Logger) {
	logger.Printf("Received message with method : %v\n", method)

	switch method {
	case "initialize":
		logger.Println("Handling initialize request")
		var request lsp.InitializeRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Failed to unmarshal initialize request: %v\n", err)
		}
		logger.Printf("Connected to : %v %v\n", request.Params.ClientInfo.Name, request.Params.ClientInfo.Version)

		msg := lsp.NewInitializeResponse(request.ID)
		reply := rpc.EncodeMessage(msg)

		writer := os.Stdout
		writer.Write([]byte(reply))

		logger.Printf("Sent the response: %v\n", reply)
	}
}

func getLogger(filename string) *log.Logger {
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v\n", err)
	}

	return log.New(logFile, "[go-lsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
