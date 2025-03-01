window.onload = function() {
    //<editor-fold desc="Changeable Configuration Block">

    // the following lines will be replaced by docker/configurator, when it runs in a docker-container
    //window.ui = SwaggerUIBundle({
    //  url: "https://petstore.swagger.io/v2/swagger.json",
    //  dom_id: '#swagger-ui',
    //  deepLinking: true,
    //  presets: [
    //    SwaggerUIBundle.presets.apis,
    //    SwaggerUIStandalonePreset
    //  ],
    //  plugins: [
    //    SwaggerUIBundle.plugins.DownloadUrl
    //  ],
    //  layout: "StandaloneLayout"
    //});
    window.ui = SwaggerUIBundle({
      url: "http://192.168.1.15:1388/api-docs/swagger.json", // Hardcode for testing
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
