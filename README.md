# Hotel reservation backend

# Project environment variables
```
HTTP_LISTEN_ADDRESS=:3000
JWT_SECRET=thisisourjwtsecret
MONGO_DB_NAME=hotel-reservation
MONGO_DB_URL=mongodb://localhost:27017/
TEST_MONGO_DB_URL=mongodb://localhost:27017/
```

## Dependency modules
### MongoDB
Docs
```
https://mongodb.com/docs/drivers/go/current/quick-start
```


Install
```
go get go.mongodb.org/mongo-driver/mongo
```


### gofiber
Docs
```
https://gofiber.io
```


Install
```
go get github.com/gofiber/fiber/v2
```
