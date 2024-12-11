Steps to setup 
1. Clone the repo.
2. go run main.go
3. Now try with the following curl commands for sign, signup, token generator, token revocation and token renew

For sign up 
```
curl -X POST http://localhost:8080/signup \
-H "Content-Type: application/json" \
-d '{"email": "user@example.com", "password": "securepassword"}'
```

For Sign in
```
curl -X POST http://localhost:8080/signin \
-H "Content-Type: application/json" \
-d '{"email": "user@example.com", "password": "securepassword"}'
```
Response will be 
```
{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzM5MzI2MTksInVzZXIiOnsiZW1haWwiOiJzYWNoaW5AZXhhbXBsZS5jb20iLCJwYXNzd29yZCI6InNhY2hpbjEyMyJ9fQ.AXFa57YFEi3iG-imvWGDEhwNIcklGKz_nv-mb7bJb9Y"}

```

For Protected area
```
curl -X GET http://localhost:8080/protected \
-H "Authorization: Bearer <token>"
```

For Token revocation
```
curl -X POST http://localhost:8080/revoke \
-H "Authorization: Bearer <token>"
```

For Token Renew
```
curl -X POST http://localhost:8080/renew \
-H "Authorization: Bearer <token>"
```
