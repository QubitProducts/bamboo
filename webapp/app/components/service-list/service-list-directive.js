var _ = require("lodash");

module.exports = ["State", "$rootScope", function (State, $rootScope) {
  return {
    restrict: "AE",
    link: function (scope) {
      var fetch = function () {
        State.get().then(function (payload) {
          var appsMap = _.indexBy(payload.Apps, "Id");
          var appIds  = _.keys(appsMap);
          var servicesKeys = _.keys(payload.Services);

          var services = _.map(_.union(appIds, servicesKeys), function (id) {
            var actionType;
            var app = appsMap[id];
            var serviceModel = payload.Services[id];

            // Add required action
            if (app && serviceModel !== undefined) {
              actionType = "default";
            } else {
              if (!app) {
                actionType = "marathon";
              } else if (!serviceModel) {
                actionType = "service";
              }
            }

            return {
              id: id,
              service: serviceModel,
              app: app,
              actionType: actionType
            };
          });

          scope.services = _.sortBy(services, function (d) {
            if (d.actionType === "service") {
              return 0;
            } else if (d.actionType === "marathon") {
              return -1;
            }
            return 5;
          });
        });
      };

      fetch();

      $rootScope.$on("services.reset", fetch);
    },

    template: require("./service-list.html")
  };
}];