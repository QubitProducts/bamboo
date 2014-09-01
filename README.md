== Bamboo


=== REST APIs for managing services

POST /api/state/domains

```
curl -i -X POST -d '{"id":"app-1","value":"app1.example.com"}' http://localhost:8000/api/state/domains
```

PUT /api/state/domains/:id

```
curl -i -X PUT -d '{"id":"app-1","value":"app1-beta.example.com"}' http://localhost:8000/api/state/domains/app-1
```

DELETE /api/state/domains/:id

```
curl -i -X DELETE http://localhost:8000/api/state/domains/app-1
```
