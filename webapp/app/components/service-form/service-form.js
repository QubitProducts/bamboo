var ngModule = angular.module("bamboo.ServiceForm", [])
  .directive("serviceNewBtn", require("./service-new-btn-directive.js"))
  .directive("serviceEditBtn", require("./service-edit-btn-directive.js"))
  .directive("serviceDeleteBtn", require("./service-delete-btn-directive.js"))

  .run(["$templateCache", function ($templateCache) {
    $templateCache.put("bamboo/service-form", require("./service-form.html"));
  }]);

module.exports = ngModule;