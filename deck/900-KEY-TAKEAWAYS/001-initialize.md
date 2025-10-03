---
marp: true
theme: default
paginate: true
---
# 🗝️ Key Takeaways & Perspectives

---
## Small Language Models (SLMs)
- SLMs (and even "Tiny" ones) are largely sufficient for this type of functionality
  - You need to provide them with context at the right moment and not too much
- And they don't need many resources
- SLMs, if well selected and for simple use cases, handle function calling very well

---
## MCP & Streamable HTTP
- The MCP protocol is very simple to implement
- Streamable HTTP is a real plus for user experience
  - Also allows deploying components on different machines
- Provides functionality to compensate for SLM limitations

---
## A2A Protocol

- Interesting, but too new to be used in production
- Official Go implementation in progress
- Specification might be a bit heavy for simple use cases

---
## Perspectives
- Replace A2A with MCP for the Boss?
- Generate the dungeon in advance?
- A Python version?
- A web interface?

[← Previous](../800-DUNGEON/03-dungeon-end-of-level-boss/001-initialize.md) | [Next →](002-thank-you.md)
