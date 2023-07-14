window.onload = function() {
  //<editor-fold desc="Changeable Configuration Block">

  // the following lines will be replaced by docker/configurator, when it runs in a docker-container
  window.ui = SwaggerUIBundle({
    urls: [{url:"protos/login.swagger.json", name:"login service"}, {url:"protos/user.swagger.json", name:"user service"}],
    dom_id: '#swagger-ui',
    deepLinking: true,
    presets: [
      SwaggerUIBundle.presets.apis,
      SwaggerUIStandalonePreset
    ],
    plugins: [
      SwaggerUIBundle.plugins.DownloadUrl,
      SwaggerUIBundle.plugins.Topbar
    ],
    layout: "StandaloneLayout"
  });

  //</editor-fold>
};
