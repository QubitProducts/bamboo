var ServiceListModule = require("./components/service-list/service-list.js");
var ServiceFormModule = require("./components/service-form/service-form.js");

var bambooApp = angular.module("bamboo", [
    ServiceListModule.name,
    ServiceFormModule.name
  ])
  .factory("State", require("./components/resources/state-resource"))
  .factory("Service", require("./components/resources/service-resource"))
  .run(["$templateCache", function ($templateCache) {
    $templateCache.put("bamboo/modal-confirm", require("./components/modal/modal-confirm.html"));
  }]);

module.exports = bambooApp;