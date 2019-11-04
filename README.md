⚙ Leadstore - basic lead/customer api and database.

The application is run from the leadstore file.

Composed of 3 components: A controler - leadstore.go, A simple CRUD sqlite package (sqldb) and a REST http API server package (apis).

A test suit (golang test) is provided to test database operations.

A Postman script is provided to test against API server endpoints. (point to localhost:3000)

A Basic .js script is provided to test CORS against API and database operations. (run on a server and point to localhost:3000)

The API server is currently set to port 3000.