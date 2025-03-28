package main

// TODO 使用golang的http server实现一个简单的服务器
// 要求该服务器提供一个Post的接口
// 不额外引入非golang官方的包，以保持demo最小化运行

import (
	"context"
	"fmt"
	"net/http"
	"time"

	mcp_golang "github.com/muidea/mcp-golang"
	mcp_http "github.com/muidea/mcp-golang/transport/http"
)

// TimeArgs defines the arguments for the time tool
type TimeArgs struct {
	Format string `json:"format" jsonschema:"description=The time format to use"`
}

func main() {
	transport := mcp_http.NewSimpleTransport()
	// Create a new server with the transport
	server := mcp_golang.NewServer(transport, mcp_golang.WithName("mcp-golang-simple-example"), mcp_golang.WithVersion("0.0.1"))

	// Register a simple tool
	err := server.RegisterTool("time", "Returns the current time in the specified format", func(ctx context.Context, args TimeArgs) (*mcp_golang.ToolResponse, error) {
		format := args.Format
		if format == "" {
			format = time.RFC3339
		}

		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(time.Now().Format(format))), nil
	})
	if err != nil {
		panic(err)
	}

	go server.Serve()

	http.HandleFunc("/mcp", transport.Handler())
	fmt.Println("Server starting on :8081...")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}
