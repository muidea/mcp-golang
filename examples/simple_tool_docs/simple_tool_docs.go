package main

import (
	"fmt"

	"github.com/metoro-io/mcp-golang/transport/stdio"
	mcp_golang "github.com/muidea/mcp-golang"
)

type HelloArguments struct {
	Submitter string `json:"submitter" jsonschema:"required,description=The name of the thing calling this tool (openai or google or claude etc)'"`
}

// This is explained in the docs at https://mcpgolang.com/tools
func main() {
	done := make(chan struct{})
	server := mcp_golang.NewServer(stdio.NewStdioServerTransport())
	err := server.RegisterTool("hello", "Say hello to a person", func(arguments HelloArguments) (*mcp_golang.ToolResponse, error) {
		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(fmt.Sprintf("Hello, %s!", arguments.Submitter))), nil
	})
	err = server.Serve()
	if err != nil {
		panic(err)
	}
	<-done
}
