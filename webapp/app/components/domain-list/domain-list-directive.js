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

          var domains = _.map(_.union(appIds, servicesKeys), function (id) {
            var actionType;
            var app = appsMap[id];
            var domainValue = payload.Services[id];


            if (app && domainValue !== undefined) {
              actionType = "default";
            } else {
              if (!app) {
                actionType = "marathon";
              } else if (!domainValue) {
                actionType = "domain";
              }
            }

            return {
              id: id,
              value: domainValue,
              app: app,
              actionType: actionType
            };
          });

          scope.domains = _.sortBy(domains, function (d) {
            if (d.actionType === "domain") {
              return 0;
            } else if (d.actionType === "marathon") {
              return -1;
            }
            return 5;
          });
        });
      };

      fetch();

      $rootScope.$on("domains.reset", fetch);
    },

    template: require("./domain-list.html")
  };
}];