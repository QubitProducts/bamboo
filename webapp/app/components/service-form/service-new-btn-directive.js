module.exports = ["Service", function (Service) {
  return {
    restrict: "AE",
    controller: require("./service-form-ctrl.js"),


    template: function (element, attrs) {
      var cta = attrs.hasOwnProperty('text') ? attrs['text'] : '<i class="icon ion-plus"></i> New';
      return '<button class="btn btn-primary btn-create-service" ng-click="new()">' +
        cta + '</button>';
    },

    scope: {
      serviceModel: "=?"
    },

    link: function (scope) {

      scope.actionName = "Create";
      scope.service = {
        id: scope.serviceModel ? (scope.serviceModel.id || "") : "",
        acl: ""
      };

      var modalOptions = {
        title: "Create new service configuration",
        template: "bamboo/modal-confirm",
        contentTemplate: "bamboo/service-form",
        scope: scope,
        animation: "am-fade-and-scale",
        show: false,
        html: true
      };

      scope.new = function () {
        scope.showModal(modalOptions);
      };

      scope.makeRequest = function (serviceModel) {
        return Service.create(serviceModel);
      };
    }
  };
}];