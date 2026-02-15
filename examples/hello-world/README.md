# Hello World - Relief Orchestrator Example

This is a simple example project to test **Relief Orchestrator**.

## 📋 What it does

A basic Node.js HTTP server that:
- Listens on the configured port (default: 3000)
- Returns JSON with request information
- Demonstrates how to use `relief.yaml` to configure a project

## 🚀 How to use with the Orchestrator

### 1. Add the project

In Relief Orchestrator, click "Add Local Project" and select this folder.

Or add manually to the configuration file:

```yaml
projects:
  - name: "hello-world"
    path: "./examples/hello-world"
    domain: "hello.local.dev"
    type: "node"
```

### 2. Start the project

In the Orchestrator panel:
1. Find the "hello-world" project
2. Click the "Start" button
3. Wait for status to change to "Running"

### 3. Test

Access in browser:
```
http://hello.local.dev
```

Or via curl:
```bash
curl http://hello.local.dev
```

You should see a JSON response like:
```json
{
  "message": "Hello from Relief Orchestrator!",
  "project": "hello-world",
  "timestamp": "2026-02-15T10:30:00.000Z",
  "port": 3000,
  "env": "development"
}
```

## 🔍 Structure

- `relief.yaml` - Project manifest (configuration)
- `index.js` - Simple HTTP server
- `package.json` - Node.js project metadata

## ⚙️ Requirements

- Node.js >= 18.0.0 (automatically checked by Orchestrator)

## 📝 Modifying

Try modifying:

1. **Port:** Change `PORT` in `relief.yaml`
2. **Domain:** Change `domain` to `test.local.dev`
3. **Response:** Edit the `response` object in `index.js`

After modifying, restart the project in the Orchestrator.

## ❓ Common issues

### Port already in use
If port 3000 is occupied, change in `relief.yaml`:
```yaml
env:
  PORT: "3001"
```

### Domain doesn't resolve
Check if Orchestrator added the entry to `/etc/hosts`:
```bash
cat /etc/hosts | grep relief
```

Should have:
```
127.0.0.1 hello.local.dev # RELIEF
```

## 🎓 Next steps

Now that you've tested the basic example:

1. Create your own `relief.yaml` in other projects
2. Explore managed dependencies
3. Configure multiple services
4. Use Docker for more complex projects
