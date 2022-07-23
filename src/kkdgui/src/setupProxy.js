const { createProxyMiddleware } = require('http-proxy-middleware');

module.exports = function(app) {
    app.use(
        '/api/v1',
        createProxyMiddleware({
          target: 'http://192.168.1.100:30000',
          changeOrigin: true,
        })
      );
      app.use(
        '/apinode',
        createProxyMiddleware({
          target: 'http://192.168.1.100:30013',
          changeOrigin: true,
        })
      );
};
