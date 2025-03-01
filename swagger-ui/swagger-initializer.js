const urlParams = new URLSearchParams(window.location.search);
const serverHost = urlParams.get('serverHost') || 'http://192.168.1.15:1388'; // Default to server IP
const url = serverHost + '/api-docs/swagger.json';
window.ui = SwaggerUIBundle({
  url: url,
  dom_id: '#swagger-ui',
  presets: [
    SwaggerUIBundle.presets.apis,
    SwaggerUIStandalonePreset
  ],
  layout: "BaseLayout",
  deepLinking: true,
  requestInterceptor: (req) => {
    // Ensure all requests use the server host
    req.url = serverHost + req.url;
    return req;
  }
});