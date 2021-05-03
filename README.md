# A GraphQL API developed using GoLang,MySQL,Redis,JWT & OAuth2

This is a Ready to deploy GraphQL API developed using GoLang,MySQL,Redis,JWT & OAuth2. You can use this if you want to quick start developing your own custom Micro service by skipping 95% of your scratch works. Hopefully this will save lot of your time as this API includes all the basic stuffs you need to get started.


## Set up

```
 git clone gitlab.com/sirinibin/go-mysql-graphql
 cd go-mysql-graphql
 export GO111MODULE="on"
 go mod tidy
 go run server.go
 ```
## Set up MySQL & Redis Db Credentials
   - Use the mysql schema golang_mysql_graphql.sql to create your new db
   - Update the DB server credentials of both MySQL & Redis on config/db.go & config.redis.go
## Run

```
 go run server.go
 ```


## Playground API Documentation

https://graphqlbin.com/v2/wRxKsR

