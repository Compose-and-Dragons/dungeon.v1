---
marp: true
theme: default
paginate: true
---

## Why Use Docker MCP Gateway?

- **Managing MCP server lifecycle**: Each local MCP sever in the catalog runs in an isolated Docker container. npx and uvx servers are granted minimal host privileges.
- **Providing a unified interface**: AI models access MCP servers through a single Gateway.
- **Handling authentication and security**: Keep secrets out of environment variables using Docker Desktop's secrets management.
- **Supports dynamic tool discovery** and configuration. Each MCP client (eg VS Code, Cursor, Claude Desktop, etc.) connects to the same Gateway configuration, ensuring consistency across different clients.
- **Enables OAuth flows** for MCPs that require OAuth access token service connections.