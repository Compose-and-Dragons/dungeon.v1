# MCP server demo

1. `docker compose up`
2. Use the url displayed in the logs to open MCP inspector.
3. Choose `Streamable HTTP` as transport type
4. Use `http://host.docker.internal:9011` as URL, this is the internal url of MCP Gateway
5. Click on `Connect`
6. Then you can test the MCP server by clicking on `Tools`, `List Tools` and `choose_character_by_species
` or `detect_real_topic_in_user_message`


## Alternative testing
Testing the MCP without MCP Gateway
1. build the Docker image: `docker build -t mcp-roll-dices:demo .`
2. run the Docker container: `docker run --rm -i -p 9095:9090 mcp-roll-dices:demo`
3. test with the bash scripts.