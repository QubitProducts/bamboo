module.exports = ["Domain", function (Domain) {
  return {
    restrict: "AE",

    template: '<button class="btn btn-default" title="Edit" ng-click="new()"><i class="icon ion-compose"></i></button>',

    scope: {
      domainId: "=",
      domainValue: "="
    },

    controller: require("./domain-form-ctrl"),

    link: function (scope) {
      scope.actionName = "Update";
      scope.disableMarathonIdChange = true;

      scope.domain = {
        id: scope.domainId,
        value: scope.domainValue
      };

      var modalOptions = {
        title: "Update domain mapping",
        template: "bamboo/modal-confirm",
        contentTemplate: "bamboo/domain-form",
        scope: scope,
        show: false,
        html: true
      };


      scope.new = function () {
        scope.showModal(modalOptions);
      };

      scope.makeRequest = function (domainModel) {
        return Domain.update(domainModel);
      };
    }
  };
}];