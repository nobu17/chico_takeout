const { createProxyMiddleware } = require("http-proxy-middleware");

module.exports = function (app) {
  app.use(
    ["/store/*", "/item/*", "order/*"],
    createProxyMiddleware({
      target: "http://localhost:8086",
      changeOrigin: true,
    })
  );
};
