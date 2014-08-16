module.exports = function () {
  return {
    restrict: "AE",
    template:  '<button class="btn btn-danger" ng-click="showModal()" title="Delete"><i class="icon ion-android-trash"></i></button>',
    scope: {
      domainId: "=",
      domainValue: "="
    },
    controller: ["$scope", "Domain", "$modal", "$rootScope", function ($scope, Domain, $modal, $rootScope) {
      $scope.actionName = "Delete It!";

      $scope.showModal = function () {
        $scope.modal = $modal({
          title: "Are you sure?",
          template: "bamboo/modal-confirm",
          content: "Delete Marathon ID " +
            $scope.domainId + " mapping to " + $scope.domainValue,
          scope: $scope,
          show: true
        });
      };

      $scope.doAction = function () {
        Domain.destroy({
            id: $scope.domainId,
            value: $scope.domainValue
          })
          .then(function () {
            $scope.modal.hide();
            $scope.modal = null;
            $rootScope.$broadcast("domains.reset");
          });
      };
    }]
  }
};