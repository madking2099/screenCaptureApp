window.onload = function() {
  //<editor-fold desc="Changeable Configuration Block">
  // the following lines will be replaced by docker/configurator, when it runs in a docker-container
  // Get server host from query parameter or default to window.location.origin
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
  //</editor-fold>
};
