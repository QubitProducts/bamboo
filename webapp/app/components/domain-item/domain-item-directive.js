module.exports = function () {
  return {
    restrict: "AE",
    replace: true,
    scope: {
      domainId: "=",
      domainValue: "=",
      domainApp: "=",
      domainActionType: "=?"
    },
    controller: ["$scope", function ($scope) {

      $scope.instancesCount = function () {
        if ($scope.domainApp) {
          return $scope.domainApp.Tasks.length;
        }

        return "-";
      };


    }],
    template: require("./domain-item.html")
  };
};