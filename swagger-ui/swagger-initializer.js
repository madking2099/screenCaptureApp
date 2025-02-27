const urlParams = new URLSearchParams(window.location.search);
const serverHost = urlParams.get('serverHost') || window.location.origin;
const url = serverHost + '/api-docs/swagger.json';
window.ui = SwaggerUIBundle({
  url: url,
  dom_id: '#swagger-ui',
  presets: [
    SwaggerUIBundle.presets.apis,
    SwaggerUIStandalonePreset
  ],
  layout: "BaseLayout",
  deepLinking: true
});