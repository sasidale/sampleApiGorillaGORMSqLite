implemented 
- routing using Gorilla mux
- GORM
- in file db sqlite


Install Gorilla mux
go get -u github.com/gorilla/mux

Install GORM
go get -u github.com/jinzhu/gorm

Install sqlite3
GORM sqlite dilect uses go-sqlite3
go get -u github.com/mattn/go-sqlite3

Supported Endpoints
POST /users
GET /users
GET /users/{id}
PUT /users/{id}
DELETE /users/{id}

Validation : 
basic API validation on id is provided

