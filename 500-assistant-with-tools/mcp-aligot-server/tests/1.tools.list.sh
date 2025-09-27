#!/bin/bash
: <<'COMMENT'
# Tools list
COMMENT

# STEP 1: Load the session ID from the environment file
source mcp.env
source mcp.server.env

MCP_SERVER=${MCP_SERVER:-"http://localhost:${MCP_HTTP_PORT}"}
echo "Using MCP Server: ${MCP_SERVER}"
read -r -d '' DATA <<- EOM
{
  "jsonrpc": "2.0",
  "id": "test",
  "method": "tools/list",
  "params": {}
}
EOM

# NOTE: always use the session ID from the mcp.env file
curl ${MCP_SERVER}/mcp \
  -H "Content-Type: application/json" \
  -H "Mcp-Session-Id: $SESSION_ID" \
  -d "${DATA}" | jq
