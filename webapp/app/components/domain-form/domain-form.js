var ngModule = angular.module("bamboo.DomainForm", [])
  .directive("domainNewBtn", require("./domain-new-btn-directive"))
  .directive("domainEditBtn", require("./domain-edit-btn-directive"))
  .directive("domainDeleteBtn", require("./domain-delete-btn-directive"))

  .run(["$templateCache", function ($templateCache) {
    $templateCache.put("bamboo/domain-form", require("./domain-form.html"));
  }]);

module.exports = ngModule;