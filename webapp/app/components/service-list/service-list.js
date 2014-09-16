var DomainItemModule = require("../service-item/service-item.js");

var ngModule = angular.module("bamboo.ServiceList", [
    DomainItemModule.name
  ])
  .directive("serviceList", require("./service-list-directive.js"));

module.exports = ngModule;