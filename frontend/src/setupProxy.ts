import { createProxyMiddleware } from "http-proxy-middleware";
module.exports = function (app: any) {
  app.use(
    ["/store/*", "/item/*", "order/*"],
    createProxyMiddleware({
      target: "http://localhost:8086",
      changeOrigin: true,
    })
  );
};
