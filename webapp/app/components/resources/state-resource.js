module.exports = ["$resource", function ($resource) {
  var index = $resource("/api/state", {});
  return {
    get: function () {
      return index.get().$promise;
    }
  }
}];