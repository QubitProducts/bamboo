module.exports = ["$resource", function ($resource) {
  var index = $resource("/api/state/domains", {},
    {
      get: { method: "GET" },
      create: { method: "POST" }
    });
  var entity = $resource("/api/state/domains/:id", { id: "@id" }, {
    update: { method: "PUT", params: { id: "@id" } },
    destroy: { method: "DELETE", params: { id: "@id"} }
  });

  return {
    all: function () {
      return index.get().$promise;
    },

    create: function (params) {
      return index.create(params).$promise;
    },

    update: function (params) {
      return entity.update(params).$promise;
    },

    destroy: function (params) {
      return entity.destroy(params).$promise;
    }
  }
}];