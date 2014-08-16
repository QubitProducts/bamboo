var ngModule = angular.module("bamboo.DomainItem", [])
  .directive("domainItem", require("./domain-item-directive"));

module.exports = ngModule;