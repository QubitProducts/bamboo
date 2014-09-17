module.exports = function () {
  return {
    restrict: "AE",
    template:  '<button class="btn btn-danger" ng-click="showModal()" title="Delete"><i class="icon ion-android-trash"></i></button>',
    scope: {
      serviceId: "="
    },
    controller: ["$scope", "Service", "$modal", "$rootScope", function ($scope, Service, $modal, $rootScope) {
      $scope.actionName = "Delete It!";

      $scope.showModal = function () {
        $scope.modal = $modal({
          title: "Are you sure?",
          template: "bamboo/modal-confirm",
          content: "Delete Marathon ID " + $scope.serviceId,
          scope: $scope,
          show: true
        });
      };

      $scope.doAction = function () {
        Service.destroy({
            id: $scope.serviceId
          })
          .then(function () {
            $scope.modal.hide();
            $scope.modal = null;
            $rootScope.$broadcast("services.reset");
          });
      };
    }]
  }
};