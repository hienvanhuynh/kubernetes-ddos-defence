const { createProxyMiddleware } = require('http-proxy-middleware');

const host = process.env.NODE_IP || "localhost";

console.log('host', host);

module.exports = function(app) {
    app.use(
        '/api/v1',
        createProxyMiddleware({
          target: `http://${host}:30000`,
          changeOrigin: true,
        })
      );
      app.use(
        '/apinode',
        createProxyMiddleware({
          target: `http://${host}:30013`,
          changeOrigin: true,
        })
      );
};
