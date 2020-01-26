# roadmaps

## Run
On the first start, app should be able to create default users. For allow this, we should need to add environment variables called `adminname{n}`,`adminemail{n}`,`adminpass{n}` where `{n}` is positive integer up to 10. It will be created only if table `users` is empty.

# API
## [Users](#user-api)
- [Registration](#registration)
- [Login](#login)
- [Refresh token](#refresh-token)

## [Plans](#plans)
- [Add](#add-plan)
- [Get](#get-plan-with-steps)
- [Get list](#plan-list)
- [Edit](#edit-plan)
- [Remove](#remove-plan)
- [Add to favorite](#choose-plan-as-favorite-within-topic) 
- [Remove from favorite](#remove-plan-from-favorite-within-topic)
- [Plan tree](#plan-tree)

## [Resources](#resources)
- [Add](#add-resource)

## [Topics](#topics)
- [Add](#add-topic)
- [Get](#topic-get)
- [Topic tree](#topic-tree)
- [Search](#topic-search)
- [Add tag](#topic-add-tag)
- [Remove tag](#topic-remove-tag)
- [Edit](#edit-topic-as-admin)

## [Comments](#comments)
- [Add](#add-comment)
- [Remove](#remove-comment)
- [Threads list](#thread-list)
- [Thread comments](#thread-comments)

---

## Types
- [EntityType](#EntityType)
- [Errors](#Errors)

---

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

### Refresh token
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

### Edit plan
#### /api/plan/edit
Request
```javascript
{
    "id": "string",
    "title":"string",
    "topic":"string",
    "steps":[{
        "referenceId": "int",
        "referenceType":"Resource | Topic | Test"
    }]
}
```

Response
### 200 - OK
No Body

### 400 - BadRequest
```javascript
{
    "error": "INVALID_REQUEST | INTERNAL_ERROR ",
    "validation": {
        "topic": "INVALID_FORMAT"
        "title": "INVALID_FORMAT",
        "steps": "INVALID_COUNT",
        "id":"NOT_EXISTS | ACCESS_DENIED | INVALID_VALUE"
    }
}
```
### 500 - Internal Error
No Body

---

### Remove plan
```javascript
{
    "id": "string"
}
```
Response
### 200 - OK
No Body

### 400 - BadRequest
```javascript
{
    "error": "INVALID_REQUEST | INTERNAL_ERROR",
    "validation": {
        "id":"NOT_EXISTS | ACCESS_DENIED | INVALID_VALUE"
    }
}
```
### 500 - Internal Error
No Body
 
---

## Resources
### Add resource
#### /api/source/add
Request
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
	"desc": "string",
    "tags": [{
        "name": "string",
        "title": "string"
    }],
    "istag": bool
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
        "desc": "string",
        "tags": [{
            "name": "string",
            "title": "string"
        }],
        "istag": bool
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

### Topic search
#### /api/topic/search
Request
```javascript
{
	"query": "string"
}
```
Response
### 200 - OK
```javascript
{
    "query": "string"
    "topic": {
        "name": "string",
        "title": "string",
        "desc": "string",
        "tags": [string]
    }
}
```

### Topic add tag
#### /api/topic/tag/add
Request
```javascript
{
	"topicname":"strint",
	"tagname":"string"
}
```
Response
### 200 - OK
```javascript
{
    "added": bool
}
```

### Topic remove tag
#### /api/topic/tag/remove
Request
```javascript
{
	"topicname":"strint",
	"tagname":"string"
}
```
Response
### 200 - OK
```javascript
{
    "removed": bool
}
```

### Edit topic as admin
#### /api/topic/edit
Required **M** rights

Request
```javascript
{
	"id": int,
	"title": "string",
	"desc": "string",
	"istag": bool
}
```
Response
### 200 - OK
NoBody

### 403 - Forbidden
NoBody

### 400 - BadRequest
```javascript
{
    "error": "INTERNAL_ERROR | INVALID_REQUEST",
    "validation":{
        "title": "INVALID_FORMAT",
        "id": "INVALID_VALUE"
    }
}
```

## Comments
### Add comment
#### /api/comment/add
Request
```javascript
{
	"entityType": int, // see EntityType
	"entityId": "string", // string for planId and int for other
	"parentId": int,
	"text": "string",
	"title": "string" // null if parentId == 0
}
```
Response
### 200 - OK
```javascript
{
    "id": int
}
```

### 403 - Forbidden
NoBody

### 400 - BadRequest
```javascript
{
    "error": "INTERNAL_ERROR | INVALID_REQUEST",
    "validation":{
        "entityType": "INVALID_VALUE",
        "entityId": "INVALID_VALUE",
        "text": "INVALID_VALUE",
        "title": "INVALID_VALUE",
        "parentId"  "INVALID_VALUE",
    }
}
```

### Remove comment
#### /api/comment/delete
Request
```javascript
{
    "id": int
}
```

### 200 - OK
NoBody

### 403 - Forbidden
NoBody

### 400 - BadRequest
```javascript
{
    "error": "INTERNAL_ERROR | INVALID_REQUEST",
    "validation":{
        "id": "INVALID_VALUE"
    }
}
```
### Thread list
#### /api/comment/threads
Request
```javascript
{
	"entityType": int,  // see EntityType
	"entityId": "string", // string for planId and int for other
	"count": int, // count per page
	"page": int // start from zero
}
```

### 200 - OK
```javascript
{
    "hasMore": false,
    "page": 0,
    "comments": [
        {
            "Id": 1,
            "EntityType": 1,
            "EntityId": "e",
            "ThreadId": 0,
            "ParentId": 0,
            "Date": "2020-01-26T11:54:55.3433Z",
            "User": {
                "id": "e45bdc37-6a74-4871-bbfe-0e03e1347920",
                "name": "Neek",
                "img": ""
            },
            "Text": "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum",
            "Title": "wft man!",
            "Deleted": false,
            "Points": 0
        }
    ]
}
```

### 400 - BadRequest
```javascript
{
    "error": "INTERNAL_ERROR | INVALID_REQUEST",
    "validation":{
        "entityType": "INVALID_VALUE",
        "entityId": "INVALID_VALUE",
        "count": "INVALID_VALUE",
        "page": "INVALID_VALUE"
    }
}
```

### Thread comments
#### /api/comment/thread
Request
```javascript
{
	"entityType": int,  // see EntityType
	"entityId": "string", // string for planId and int for other
	"threadId": int, 
}
```

### 200 - OK
```javascript
[
    {
        "Id": 3,
        "EntityType": 1,
        "EntityId": "e",
        "ThreadId": 1,
        "ParentId": 1,
        "Date": "2020-01-26T11:55:32.941379Z",
        "User": {
            "id": "e45bdc37-6a74-4871-bbfe-0e03e1347920",
            "name": "Neek",
            "img": ""
        },
        "Text": "anim id est laborum",
        "Title": "",
        "Deleted": false,
        "Points": 0,
        "Childs": [
            {
                "Id": 5,
                "EntityType": 1,
                "EntityId": "e",
                "ThreadId": 1,
                "ParentId": 3,
                "Date": "2020-01-26T15:14:28.766358Z",
                "User": {
                    "id": "e45bdc37-6a74-4871-bbfe-0e03e1347920",
                    "name": "Neek",
                    "img": ""
                },
                "Text": "!!!!! 3-",
                "Title": "",
                "Deleted": false,
                "Points": 0,
                "Childs": []
            }
        ]
    }
]
```
### 400 - BadRequest
```javascript
{
    "error": "INTERNAL_ERROR | INVALID_REQUEST",
    "validation":{
        "entityType": "INVALID_VALUE",
        "entityId": "INVALID_VALUE",
        "threadId": "INVALID_VALUE"
    }
}
```


## EntityType
| Name                      | Value |
| ------------------------  | ----------- |
| Plan                      | 1
| Topic                     | 2
| Project                   | 3
| Resource                  | 4

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
| "INVALID_VALUE"          |
| "INACCESSIBLE_WEBPAGE"   |
| "INVALID_FORMAT"         |
| "SOURCE_NOT_FOUND"       |
| "INVALID_COUNT"          |
| "NOT_EXISTS"             |
```