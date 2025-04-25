# Incidents as a MCP server prototype

This is a very basic prototype of a MCP server (using [mcp-golang](https://github.com/metoro-io/mcp-golang) library) 
with one registered tool called "incidents"

## Run locally
1. Connect to an OpenShift cluster (with installed Incidents)
2. Run `./forward.sh` to find the `thanos-querier` Pod and locally port-forward Thanos
3. Run the MCP server with:
    ```bash
    go run server/main.go
    ```
    Server is running at `localhost:8081` by default
4. Run the MCP client with:
    ```bash
    go run client/client.go
    ```
