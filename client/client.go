package main

import (
	"context"
	"log"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

func main() {

	mcpCli, err := client.NewSSEMCPClient("http://localhost:8081/sse")
	if err != nil {

	}
	defer mcpCli.Close()

	err = mcpCli.Start(context.Background())
	if err != nil {
		log.Fatalf("Failed to start MCP client: %v", err)
	}

	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "Incidents client",
		Version: "0.0.1",
	}

	_, err = mcpCli.Initialize(context.Background(), initRequest)
	if err != nil {
		log.Fatalf("Failed to initialize MCP client: %v", err)
	}

	d := 1 * time.Minute
	incidentsRequest := mcp.CallToolRequest{}

	incidentsRequest.Params.Name = "incidents"
	incidentsRequest.Params.Arguments = map[string]interface{}{
		//	"component": "monitoring",
		"time": d.String(),
		// namespace ??
		// severity ??
		// alertname ??
	}

	incidentsData, err := mcpCli.CallTool(context.Background(), incidentsRequest)
	if err != nil {
		log.Fatalf("Failed to call the Incidents tool: %v\n", err)
	}
	log.Print(incidentsData.Content)
}
