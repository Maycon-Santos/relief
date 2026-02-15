const http = require('http');

const PORT = process.env.PORT || 3000;

const server = http.createServer((req, res) => {
  const response = {
    message: 'Hello from Relief Orchestrator!',
    project: 'hello-world',
    timestamp: new Date().toISOString(),
    environment: process.env.NODE_ENV || 'production',
    path: req.url,
    method: req.method,
  };

  res.writeHead(200, { 'Content-Type': 'application/json' });
  res.end(JSON.stringify(response, null, 2));
});

server.listen(PORT, () => {
  console.log(`🚀 Hello World server running on http://localhost:${PORT}`);
  console.log(`📝 Try accessing: http://hello.local.dev`);
});

// Graceful shutdown
process.on('SIGTERM', () => {
  console.log('SIGTERM received, shutting down gracefully...');
  server.close(() => {
    console.log('Server closed');
    process.exit(0);
  });
});

process.on('SIGINT', () => {
  console.log('SIGINT received, shutting down gracefully...');
  server.close(() => {
    console.log('Server closed');
    process.exit(0);
  });
});
