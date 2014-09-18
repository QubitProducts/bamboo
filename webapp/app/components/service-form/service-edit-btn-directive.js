module.exports = ["Service", function (Service) {
  return {
    restrict: "AE",

    template: '<button class="btn btn-default" title="Edit" ng-click="new()"><i class="icon ion-compose"></i></button>',

    scope: {
      serviceModel: "="
    },

    controller: require("./service-form-ctrl.js"),

    link: function (scope) {
      scope.actionName = "Update";
      scope.disableMarathonIdChange = true;

      scope.service = {
        id: scope.serviceModel.id,
        acl: scope.serviceModel.service.Acl
      };

      var modalOptions = {
        title: "Edit service configuration",
        template: "bamboo/modal-confirm",
        contentTemplate: "bamboo/service-form",
        scope: scope,
        show: false,
        html: true
      };


      scope.new = function () {
        scope.showModal(modalOptions);
      };

      scope.makeRequest = function (model) {
        return Service.update(model);
      };
    }
  };
}];