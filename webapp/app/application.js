var DomainListModule = require("./components/service-list/service-list.js");
var DomainFormModule = require("./components/service-form/service-form.js");

var bambooApp = angular.module("bamboo", [
    DomainListModule.name,
    DomainFormModule.name
  ])
  .factory("State", require("./components/resources/state-resource"))
  .factory("Service", require("./components/resources/service-resource"))
  .run(["$templateCache", function ($templateCache) {
    $templateCache.put("bamboo/modal-confirm", require("./components/modal/modal-confirm.html"));
  }]);

module.exports = bambooApp;