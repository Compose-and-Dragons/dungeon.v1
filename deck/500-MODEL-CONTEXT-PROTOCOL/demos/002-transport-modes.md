---
marp: true
theme: default
paginate: true
---
<style>
.dodgerblue {
  color: dodgerblue;
}
.firebrick {
  color: firebrick;
}
</style>
## How MCP Works - Transport Modes

- **Stdio Transport**: Direct process communication via standard input/output
- **SSE Transport**: Server-Sent Events for web-based applications (<span class="firebrick">`deprecated`</span>)
- <span class="dodgerblue">**Streamable HTTP Transport**</span>: Built upon HTTP and can support various communication patterns, including request-response and streaming
- **`<Your Transport>`**: ...


Each transport mode enables different deployment scenarios and use cases.