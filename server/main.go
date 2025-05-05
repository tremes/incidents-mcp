package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

func main() {
	parentCtx := context.Background()
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()
	addr := "localhost:8081"

	mcpServer := server.NewMCPServer("incidents-mcp", "0.0.1", server.WithToolCapabilities(true))
	sseServer := server.NewSSEServer(mcpServer, server.WithBaseURL(fmt.Sprintf("http://%s", addr)))

	mcpServer.AddTool(mcp.NewTool("incidents",
		mcp.WithDescription("Incidents provides list of active incidents in the cluster.")),
		handleIncidentTool)

	go func() {
		err := sseServer.Start(addr)
		if err != nil {
			log.Fatalf("Failed to run the MCP server: %v", err)
		}

	}()
	log.Printf("MCP server is running at %s\n", addr)
	<-ctx.Done()
}

func handleIncidentTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	arguments := request.Params.Arguments
	component := arguments["component"]
	qTime, ok := arguments["time"].(string)
	d := 24 * time.Hour
	if ok {
		durationArg, err := time.ParseDuration(qTime)
		if err != nil {
			return nil, err
		}
		d = durationArg
	}

	promQuery := "cluster:health:components:map{}"
	if component != nil {
		promQuery = fmt.Sprintf(`cluster:health:components:map{"component"="%s"}`, component)
	}
	fmt.Println("Component: ", component)
	fmt.Println("Time: ", d)

	api_config := api.Config{
		Address: "http://localhost:8080",
	}
	/* 	certs := x509.NewCertPool()
	   certs.AppendCertsFromPEM([]byte(promCert))
	   defaultRt := api.DefaultRoundTripper.(*http.Transport)
	   defaultRt.TLSClientConfig = &tls.Config{RootCAs: certs}

	   api_config.RoundTripper = promConfig.NewAuthorizationCredentialsRoundTripper(
		   "Bearer", promConfig.NewInlineSecret(promToken), defaultRt)
	*/
	promClient, err := api.NewClient(api_config)
	if err != nil {
		return nil, err
	}

	promAPI := v1.NewAPI(promClient)

	incidentsData, _, err := promAPI.Query(ctx, promQuery, time.Now().Add(d))
	if err != nil {
		return nil, err
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("The incidents data is %s", incidentsData.String()),
			},
		},
	}, nil
}
