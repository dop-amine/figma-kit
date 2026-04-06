# Figma Desktop MCP Server Setup

Alternative to the remote MCP server (`mcp.figma.com`). Runs locally via the Figma Desktop app for lower latency and no rate limits.

## Prerequisites

- Figma Desktop app (latest version)
- A Full seat on your Figma plan (Dev seat is read-only)
- Edit permission on the file you want to modify

## Setup Steps

### 1. Enable the Desktop MCP Server

1. Open **Figma Desktop** (not the browser version)
2. Open any design file
3. Enter **Dev Mode** with `Shift+D`
4. In the inspect panel, click **"Enable desktop MCP server"**
5. The server starts at `http://127.0.0.1:3845/mcp`

### 2. Configure Cursor

Update your project's `.cursor/mcp.json` to use the local server:

```json
{
  "mcpServers": {
    "figma": {
      "url": "http://127.0.0.1:3845/mcp"
    }
  }
}
```

Then reload Cursor (Cmd+Shift+P -> "Reload Window").

### 3. Verify Connection

Ask the AI to run `whoami` via the Figma MCP tools. It should return your Figma account info.

## Switching Between Remote and Desktop

To switch back to the remote server:

```json
{
  "mcpServers": {
    "figma": {
      "url": "https://mcp.figma.com/mcp"
    }
  }
}
```

You can keep both configurations and comment/uncomment as needed.

## Comparison

| Feature | Remote (`mcp.figma.com`) | Desktop (`localhost:3845`) |
|---------|-------------------------|---------------------------|
| Latency | Higher (network) | Low (local) |
| Rate limits | Yes | No |
| Requires Figma Desktop open | No | Yes |
| Requires internet | Yes | No (for open files) |
| File must be open | No | Yes |
| Authentication | OAuth via Figma | Inherits desktop session |
| Available tools | All 17 tools | All 17 tools |
| Max code size | 50KB | 50KB |

## When to Use Desktop Server

- Working on a single file for extended periods
- Iterating rapidly (many `use_figma` calls)
- On slower internet connections
- When hitting rate limits on the remote server

## When to Use Remote Server

- Working across multiple files
- Files not currently open in Figma Desktop
- Need to create new files (`create_new_file`)
- Collaborating and need cloud-synced state

## Troubleshooting

**Server not responding:**
- Ensure Figma Desktop is running and a file is open
- Re-enable Dev Mode (`Shift+D`)
- Check the port isn't blocked by a firewall

**Changes not appearing:**
- Figma Desktop may need focus to process Plugin API calls
- Try clicking into the Figma window, then retry

**"Permission denied" errors:**
- Verify you have edit access to the file
- Desktop server inherits your Figma session permissions
