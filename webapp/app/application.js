var DomainListModule = require("./components/domain-list/domain-list");
var DomainFormModule = require("./components/domain-form/domain-form");

var bambooApp = angular.module("bamboo", [
    DomainListModule.name,
    DomainFormModule.name
  ])
  .factory("State", require("./components/resources/state-resource"))
  .factory("Domain", require("./components/resources/domain-resource"))
  .run(["$templateCache", function ($templateCache) {
    $templateCache.put("bamboo/modal-confirm", require("./components/modal/modal-confirm.html"));
  }]);

module.exports = bambooApp;
