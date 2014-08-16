module.exports = ["Domain", function (Domain) {
  return {
    restrict: "AE",
    controller: require("./domain-form-ctrl"),


    template: function (element, attrs) {
      var cta = attrs.hasOwnProperty('text') ? attrs['text'] : '<i class="icon ion-plus"></i> New';
      return '<button class="btn btn-primary btn-create-domain" ng-click="new()">' +
        cta + '</button>';
    },

    scope: {
      domainId: "=?"
    },

    link: function (scope) {

      scope.actionName = "Create";
      scope.domain = {
        id: scope.domainId || "",
        value: ""
      };

      var modalOptions = {
        title: "Create new mapping",
        template: "bamboo/modal-confirm",
        contentTemplate: "bamboo/domain-form",
        scope: scope,
        animation: "am-fade-and-scale",
        show: false,
        html: true
      };

      scope.new = function () {
        scope.showModal(modalOptions);
      };

      scope.makeRequest = function (domainModel) {
        return Domain.create(domainModel);
      };
    }
  };
}];