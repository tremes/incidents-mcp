package main

import (
	"context"
	"log"

	mcpGolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/http"
)

type IncidentArgs struct {
	// TODO we probably need to pass time - how far back in history we want go
	// additionally component name, alert name etc.
	Name string `json:"name" jsonschema:"required,description=The name to say hello to"`
}

func main() {
	transport := http.NewHTTPClientTransport("/mcp")
	transport.WithBaseURL("http://localhost:8081")

	ctx := context.Background()
	client := mcpGolang.NewClient(transport)

	_, err := client.Initialize(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize client: %v", err)
	}

	resp, err := client.CallTool(ctx, "incidents", IncidentArgs{
		Name: "test",
	})
	if err != nil {
		log.Fatalf("Failed to call incidents tool: %v", err)
	}
	log.Println(resp.Content[0].TextContent)
}
