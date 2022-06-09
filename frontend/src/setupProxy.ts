import { createProxyMiddleware } from "http-proxy-middleware";
module.exports = function (app: any) {
  app.use(
    ["/store/*", "/item/*"],
    createProxyMiddleware({
      target: "http://localhost:8086",
      changeOrigin: true,
    })
  );
};
