"use strict";

require("jquery");
require("angular");

require("angular-ui-router");
require("angular-animate");
require("angular-strap-tpl");
require("angular-strap");
require("angular-resource");
require("angular-moment");

var routes = require("./routes");
var app = require("./application");

angular.element(document).ready(function() {
  angular.module("app", [
    "ui.router", "ngResource",
    "ngAnimate", "mgcrea.ngStrap",
    "angularMoment",
    app.name
  ])
  .config(routes);

  angular.bootstrap(document, ['app']);
});
