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
    const urlParams = new URLSearchParams(window.location.search);
    const serverHost = urlParams.get('serverHost') || 'http://192.168.1.15:1388'; // Default to server IP
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
        // Ensure all requests use the server host, even if UI defaults to localhost
        if (!req.url.startsWith('http')) {
          req.url = baseUrl + req.url;
        }
        console.log('Intercepting request, new URL:', req.url); // Debug in browser
        return req;
      }
    });
};
