# roadmaps
## Build
Repositories for real DB engine will be added later, when functionality will be determined

Use `go build -tags DEV` when developing.
This tag enable repositories with in-memory db, what not persistence. 
## Deplay
1. On the first start, app should be able to create default users. For allow this, we should need to add environment variables called `adminname{n}`,`adminemail{n}`,`adminpass{n}` where `{n}` is positive integer up to 10. It will be created only if table `users` is empty.

## User Api
### Registration
#### /api/user/reqistration

Request
```javascript
{
	"name": "string",
	"email: "string",
	"pass": "string"
}
```

Response
### 200 - OK,
No Body

-------------

### 400 - BadRequst
```javascript
{
    "error": "INVALID_REQUEST",
    "validation": {
        "email": "ERROR",
		"name": "ERROR"
    }
}
```
Parameter | Description | Value
------------ | ------------- | -------------
error | One or more request paratemers is inavalid | "INVALID_REQUEST" | 
validation | Parameter validation description | "INVALID_PASSWORD", "INVALID_EMAIL", "INVALID_USERNAME", "ALREADY_EXISTS"(name,email)

-------------

### 500 - Internal Error
No Body

-----

## Errors
Error | Description
------------ | -------------
"NONE" |
"INVALID_PASSWORD" |
"INVALID_EMAIL" |
"NONEXISTENT_EMAIL" |
"INVALID_USERNAME" |
"ALREADY_EXISTS" |
"INTERNAL_ERROR" |
"AUTHENTICATION_ERROR" |
"AUTHENTICATION_EXPIRED" |
"INVALID_REQUEST" |
"INVALID_URL" |
"INVALID_ISBN" |
"INVALID_TITLE" |
"INVALID_PROPS" |
"INVALID_SOURCE_TYPE" |
"INACCESSIBLE_WEBPAGE" |
"INVALID_FORMAT" |
"SOURCE_NOT_FOUND" |
```