var ngModule = angular.module("bamboo.ServiceItem", [])
  .directive("serviceItem", require("./service-item-directive.js"));

module.exports = ngModule;