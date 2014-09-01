var DomainItemModule = require("../domain-item/domain-item");

var ngModule = angular.module("bamboo.DomainList", [
    DomainItemModule.name
  ])
  .directive("domainList", require("./domain-list-directive"));

module.exports = ngModule;