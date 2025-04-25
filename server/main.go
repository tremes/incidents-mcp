package main

import (
	"context"
	"log"
	"time"

	mcpGolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/http"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

type IncidentArgs struct {
	// TODO we probably need to pass time - how far back in history we want go
	Name string `json:"name" jsonschema:"required,description=The name to say hello to"`
}

func main() {

	parentCtx := context.Background()
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()
	addr := "localhost:8081"

	transport := http.NewHTTPTransport("/mcp").WithAddr(addr)
	server := mcpGolang.NewServer(transport, mcpGolang.WithName("mcp-incidents-prototype"), mcpGolang.WithVersion("0.0.1"))

	err := server.RegisterTool("incidents", "Prints list of active incidents in the cluster", incidentTool)
	if err != nil {
		log.Fatalf("Failed to register incidents tool: %v ", err)
	}

	go func() {
		err = server.Serve()
		if err != nil {
			log.Fatalf("Failed to run the MCP server: %v", err)
		}

	}()
	log.Printf("MCP server is running at %s\n", addr)
	<-ctx.Done()
}

func incidentTool(args IncidentArgs) (*mcpGolang.ToolResponse, error) {
	log.Printf("Incidents tool called with arguments: %s\n", args)
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
		log.Fatalf("Failed to create a new Prometheus client: %v", err)
	}

	promAPI := v1.NewAPI(promClient)

	incidentsData, _, err := promAPI.Query(context.Background(), `cluster:health:components:map{}`, time.Now().Add(-1*time.Minute))
	if err != nil {
		return nil, err
	}
	incidentsAsText := mcpGolang.NewTextContent(incidentsData.String())
	return mcpGolang.NewToolResponse(incidentsAsText), nil
}
