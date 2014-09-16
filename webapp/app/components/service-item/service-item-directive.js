module.exports = function () {
  return {
    restrict: "AE",
    replace: true,
    scope: {
      serviceModel: "="
    },
    controller: ["$scope", function ($scope) {

      $scope.instancesCount = function () {
        if ($scope.serviceModel.app) {
          return $scope.serviceModel.app.Tasks.length;
        }

        return "-";
      };
    }],
    template: require("./service-item.html")
  };
};