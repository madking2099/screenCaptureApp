const urlParams = new URLSearchParams(window.location.search);
const serverHost = urlParams.get('serverHost') || window.location.origin;
const baseUrl = serverHost;

window.ui = SwaggerUIBundle({
  url: baseUrl + '/api-docs/swagger.json',
  dom_id: '#swagger-ui',
  presets: [
    SwaggerUIBundle.presets.apis,
    SwaggerUIStandalonePreset
  ],
  layout: "BaseLayout",
  deepLinking: true,
  requestInterceptor: (req) => {
    // Ensure all requests use the server host
    req.url = baseUrl + req.url;
    return req;
  }
});