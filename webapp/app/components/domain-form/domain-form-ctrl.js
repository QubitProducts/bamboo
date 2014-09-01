module.exports = ["$scope", "$modal", "$rootScope", function ($scope, $modal, $rootScope) {

  $scope.showModal = function (modalOptions) {
    var modal;
    $scope.modal = modal = $modal(modalOptions);
    modal.$promise.then(modal.show);
  };

  $scope.loading = false;


  var resetError = function () {
    $scope.errors = null;
  };

  var handleSuccess = function () {
    $scope.loading = false;
    $scope.modal.hide();
    $scope.modal = null;
    $rootScope.$broadcast("domains.reset");
  };

  var handleError = function (payload) {
    $scope.loading = false;
    $scope.errors = payload.data;
  };

  $scope.doAction = function () {
    resetError();
    $scope.loading = true;
    $scope.makeRequest({
        id: $scope.domain.id,
        value: $scope.domain.value
      })
     .then(handleSuccess, handleError);
  };

}];