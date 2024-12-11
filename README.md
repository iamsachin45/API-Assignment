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
