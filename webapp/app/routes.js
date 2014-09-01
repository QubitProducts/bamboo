module.exports = function ($stateProvider, $urlRouterProvider) {

  $stateProvider
    .state("main", {
      url: "/main",
      template: require("./layouts/application.html")
    });

  $urlRouterProvider.otherwise("/main");
}