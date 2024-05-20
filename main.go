package main

import (
	"bufio"
	"encoding/json"
	"golsp/analysis"
	"golsp/lsp"
	"golsp/rpc"
	"io"
	"log"
	"os"
)

func main() {
	logger := getLogger("./go-lsp.log")
	logger.Println("Starting go-lsp")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	state := analysis.NewState()
	writer := os.Stdout

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("Failed to decode message: %s\n", err)
			continue
		}

		handleMessage(method, contents, state, writer, logger)
	}
}

func handleMessage(method string, contents []byte, state analysis.State, writer io.Writer, logger *log.Logger) {
	logger.Printf("Received message with method : %s\n", method)

	switch method {
	case "initialize":
		logger.Println("Client initialize request")
		var request lsp.InitializeRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Failed to unmarshal request: %s\n", err)
			return
		}

		logger.Printf("Connected to client : %s %s\n", request.Params.ClientInfo.Name, request.Params.ClientInfo.Version)

		msg := lsp.NewInitializeResponse(request.ID)

		writeResponse(writer, msg)

	case "textDocument/didOpen":
		var request lsp.TextDocumentDidOpenNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Failed to unmarshal request: %s\n", err)
			return
		}

		logger.Printf("Text document opened: %s\n", request.Params.TextDocument.URI)

		state.OpenDocument(request.Params.TextDocument.URI, request.Params.TextDocument.Text)

	case "textDocument/didChange":
		var request lsp.TextDocumentDidChangeNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Failed to unmarshal request: %s\n", err)
			return
		}

		logger.Printf("Text document changed: %s\n", request.Params.TextDocument.URI)

		for _, change := range request.Params.ContentChanges {
			state.UpdateDocument(request.Params.TextDocument.URI, change.Text)
		}

	case "textDocument/hover":
		var request lsp.HoverRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Failed to unmarshal request: %s\n", err)
			return
		}

		response := state.Hover(request.ID, request.Params.TextDocument.URI, request.Params.Position)

		writeResponse(writer, response)

	case "textDocument/definition":
		var request lsp.DefinitionRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Failed to unmarshal request: %s\n", err)
			return
		}

		response := state.Definition(request.ID, request.Params.TextDocument.URI, request.Params.Position)

		writeResponse(writer, response)

	case "textDocument/codeAction":
		var request lsp.CodeActionRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Failed to unmarshal request: %s\n", err)
			return
		}

		response := state.TextDocumentCodeAction(request.ID, request.Params.TextDocument.URI)

		writeResponse(writer, response)
	}
}

func writeResponse(writer io.Writer, msg any) {
	reply := rpc.EncodeMessage(msg)
	writer.Write([]byte(reply))
}

func getLogger(filename string) *log.Logger {
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %s\n", err)
	}

	return log.New(logFile, "[go-lsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
