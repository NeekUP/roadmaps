# roadmaps
## Build
Repositories for real DB engine will be added later, when functionality will be determined

Use `go build -tags DEV` when developing.
This tag enable repositories with in-memory db, what not persistence. 

While starts dev env:
1. Creating default users (name: "Neek", email: "neek@neek.com", pass "123456") and (name:"Alen", email:"alen@alen.com", pass: "123456") 
2. Server read file `dev_db.json` and store its content into inmemory database. For UsersId uses usernames created in previous step. See data example.

## Deploy
1. On the first start, app should be able to create default users. For allow this, we should need to add environment variables called `adminname{n}`,`adminemail{n}`,`adminpass{n}` where `{n}` is positive integer up to 10. It will be created only if table `users` is empty.

## User Api
Successfull login will return atoken - authorization token and rtoken - refresh token.
atoken must be passed in `Authorization` header for every request. Example: `Authorization: Bearer A7QSDW3...`
rtoken used with `/api/user/refresh` method when atoken is expired.
fp - fingerprint (https://github.com/Valve/fingerprintjs2). Auth token can be refreshed only if request contains the same fingerprint like a login when token pair has been created.

In future Login and Registration request will be protected with Invisible Recaptcha.

### Registration
#### /api/user/registration

Request
```javascript
{
	"name": "string",
	"email": "string",
	"pass": "string"
}
```

Response
### 200 - OK,
No Body



### 400 - BadRequest
```javascript
{
    "error": "INVALID_REQUEST | INTERNAL_ERROR",
    "validation": {
        "email": "INVALID_FORMAT | ALREADY_EXISTS",
		"name": "INVALID_FORMAT | ALREADY_EXISTS",
		"pass": "INVALID_FORMAT"
    }
}
```
| Parameter  | Description                                | Value                                                                                 |
| ---------- | ------------------------------------------ | ------------------------------------------------------------------------------------- |
| error      | One or more request paratemers is inavalid | "INVALID_REQUEST"                                                                     |
| validation | Parameter validation description           | "INVALID_PASSWORD", "INVALID_EMAIL", "INVALID_USERNAME", "ALREADY_EXISTS"(name,email) |


### 500 - Internal Error
No Body

-----

### Login
#### /api/user/login
Request
```javascript
{
	"email": "string",
	"pass": "string",
	"fp": "string (fingerprint (hash))"
}
```
Response
### 200 - OK
```javascript
{
	"atoken": "string",
	"rtoken": "string",
	"user": {
		"id":"string",
		"name":"string",
		"img":"relative url"
	}
}
```

### 400 - BadRequest
```javascript
{
    "error": "INVALID_REQUEST | AUTHENTICATION_ERROR",
    "validation": {
        "email": "INVALID_FORMAT",
		"pass": "INVALID_FORMAT"
    }
}
```

### 500 - Internal Error
No Body

---

### Refresh roken
#### /api/user/refresh
Request
```javascript
{
	"atoken": "string",
	"rtoken": "string",
	"fp": "string (fingerprint (hash))"
}
```
Response
### 200 - OK
```javascript
{
	"atoken": "string",
	"rtoken": "string"
}
```

### 400 - BadRequest
```javascript
{
    "error": "INVALID_REQUEST | INTERNAL_ERROR",
    "validation": {
        "atoken": "INVALID_FORMAT",
		"rtoken": "INVALID_FORMAT",
		"fp": "INVALID_FORMAT",
		"useragent": "INVALID_FORMAT"
    }
}
```
### 500 - Internal Error
No Body

---

## Plans

Plans identifies by id(string). This is a short interpretation of int. 
0 -> "a"
1 -> "b"
9999 -> "cLr"

### Add plan
#### /api/plan/add
Request
```javascript
{
	"topic": "string (topic name)",
	"title": "string",
	"steps": [{
		"referenceId": "int",
		"referenceType":"Resource | Topic | Test"
	}]
}
```
### 400 - BadRequest
```javascript
{
    "error": "INVALID_REQUEST | INTERNAL_ERROR ",
    "validation": {
		"topic": "INVALID_FORMAT"
        "title": "INVALID_FORMAT",
		"steps": "INVALID_COUNT"
    }
}
```
### 500 - Internal Error
No Body

---

### Choose plan as favorite within topic
#### /api/user/plan/favorite
Request
```javascript
{
	"planId": "string"
}
```

Response
### 200 - OK
```javascript
{
    "success": bool
}
```

### 400 - BadRequest
```javascript
{
    "error": "INVALID_REQUEST"
}
```

### 500 - Internal Error
No Body

---

### Remove plan from favorite within topic
#### /api/user/plan/unfavorite
Request
```javascript
{
	"planId": "string"
}
```

Response
### 200 - OK
```javascript
{
    "success": bool
}
```

### 500 - Internal Error
No Body

---

### Get plan with steps
#### /api/plan/get
Request
```javascript
{
	"id": "string"
}
```
Responses
### 200 - OK
```javascript
{
    "id": "string",
    "title": "string",
    "topicName": "string",
    "owner": {
        "id": "string",
        "name": "string",
        "img": "string(relative url)"
    },
    "points": int,
    "inFavorites": bool,
    "steps": [
        {
            "id": int,
            "type": "Resource | Test | Topic",
            "position": int,
            "source": {
                "id": int,
                "title": "string1",
                "type": "Article | Video | Audio | Book",
                "props": "string json",
                "img": "string url",
                "desc": "string"
            }
        },
        {
            "id": 0,
            "type": "Topic",
            "position": 2,
            "source": {
                "id": "javascript",
                "title": "Javascript",
                "desc": "Most popular programming language"
            }
        }
    ]
}
```
inFavorites: if this plan user select as favorite within topic
For steps.type == Topic fields:
    steps.source.type, steps.source.props, steps.source.img will be omitted. But img in future should be used.  
    id - string

### 400 - BadRequest
```javascript
{
    "error": "INVALID_REQUEST | INTERNAL_ERROR",
    "validation": {
        "id": "INVALID_FORMAT",
    }
}
```

### 500 - Internal Error
No Body

---

### Plan List
#### /api/plan/list
List of plan by topic name
Request
```javascript
{
	"topicName": "string"
}
```
Response
### 200 - OK
```javascript
[
    {
        "id": "string",
        "title": "string",
        "topicName": "string",
        "owner": {
            "id": "string",
            "name": "string",
            "img": "string"
        },
        "points": int,
        "inFavorites": bool
    }
]
```
inFavorites: if this plan user select as favorite within topic

### 400 - BadRequest
```javascript
{
    "error": "INVALID_REQUEST | INTERNAL_ERROR",
    "validation": {
        "topicName": "INVALID_FORMAT"
    }
}
```

### 500 - Internal Error
No Body

---
### Plan tree
#### /api/plan/tree
Request
```javascript
{

	"id": "string"
}
```

Response
### 200 - OK
```javascript
{
    "nodes": [
        {
            "topicName": "string",
            "topicTitle": "string",
            "planId": "string",
            "planTitle": "string",
            "child": [
                {
                    "topicName": "string",
                    "topicTitle": "string",
                    "planId": "string",
                    "planTitle": "string"
                }
            ]
        }
    ]
}
```

### 400 - BadRequest
```javascript
{
    "error": "INVALID_REQUEST | INTERNAL_ERROR ",
    "validation": {
		"id": "INVALID_FORMAT"
    }
}
```
### 500 - Internal Error
No Body

---

## Resources
### App resource
#### /api/source/add
```javascript
{
	"identifier":"string (absolute url / isbn)",
	"type":"Article | Vidoe | Audio | Books"
}
```

Response
### 200 - OK
```javascript
{
	"id": "int",
	"title": "string",
	"identifier": "URL | ISBN-13 | ISBN-10",
	"type": "Article | Book | Video | Audio",
	"img":"URL",
	"desc":"string"
}
```
### 400 - BadRequest
```javascript
{
    "error": "INVALID_REQUEST | INTERNAL_ERROR",
    "validation": {
        "identifier": "INVALID_URL | INVALID_ISBN | INVALID_FORMAT | SOURCE_NOT_FOUND",
		"type": "INVALID_SOURCE_TYPE"
    }
}
```
### 500 - Internal Error
No Body

---


## Topics
### Add topic
#### /api/topic/add
Request
```javascript
{
	"title": "string",
	"desc": "string"
}
```

Response
### 200 - OK
```javascript
{
	"name": "string (Id)",
	"title": "string",
	"desc": "string"
}
```

### 400 - BadRequest
```javascript
{
    "error": "INVALID_REQUEST | ALREADY_EXISTS | INTERNAL_ERROR ",
    "validation": {
        "title": "INVALID_FORMAT"
    }
}
```
### 500 - Internal Error
No Body

---


### Topic get
#### /api/topic/get
```javascript
{
	"name":"string"
}
```

Response
### 200 - OK
```javascript
{
    "topic": {
        "name": "string",
        "title": "string",
        "desc": "string"
    }
}
```

### 400 - BadRequest
```javascript
{
    "error": "INVALID_REQUEST | INTERNAL_ERROR",
    "validation": {
        "name": "INVALID_FORMAT"
    }
}
```
### 500 - Internal Error
No Body

---

### Topic tree
#### /api/topic/tree
Request
```javascript
{
	"name": "string"
}
```

Response
### 200 - OK
```javascript
{
    "nodes": [
        {
            "topicName": "string",
            "topicTitle": "string",
            "planId": "string",
            "planTitle": "string",
            "child": [
                {
                    "topicName": "string",
                    "topicTitle": "string",
                    "planId": "string",
                    "planTitle": "string"
                }
            ]
        }
    ]
}
```

### 400 - BadRequest
```javascript
{
    "error": "INVALID_REQUEST | INTERNAL_ERROR",
    "validation": {
		"name": "INVALID_FORMAT"
    }
}
```
### 500 - Internal Error
No Body

---

## Errors
| Error                    | Description |
| ------------------------ | ----------- |
| "NONE"                   |
| "INVALID_PASSWORD"       |
| "INVALID_EMAIL"          |
| "NONEXISTENT_EMAIL"      |
| "INVALID_USERNAME"       |
| "ALREADY_EXISTS"         |
| "INTERNAL_ERROR"         |
| "AUTHENTICATION_ERROR"   |
| "AUTHENTICATION_EXPIRED" |
| "INVALID_REQUEST"        |
| "INVALID_URL"            |
| "INVALID_ISBN"           |
| "INVALID_TITLE"          |
| "INVALID_PROPS"          |
| "INVALID_SOURCE_TYPE"    |
| "INACCESSIBLE_WEBPAGE"   |
| "INVALID_FORMAT"         |
| "SOURCE_NOT_FOUND"       |
| "INVALID_COUNT"          |
| "NOT_EXISTS"             |
```