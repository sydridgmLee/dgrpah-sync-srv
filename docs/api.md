# API

## User

### create new user

```bash
curl  http://localhost:8888/dgraph_sync/user -X POST -H 'Content-Type: application/json' -d '{"name":"laura", "email":"laura@ab.cd"}' 
```



### get user info

```bash
curl http://localhost:8888/dgraph_sync/user -X GET -H 'Content-Type: application/json' -d '{"email":"laura@ab.cd"}' 
```



### update user info

```bash
curl http://localhost:8888/dgraph_sync/user -X PUT -H 'Content-Type: application/json' -d '{"email":"laura@ab.cd","name":"laura Q"}'
```



###  delete user

```bash
curl http://localhost:8888/dgraph_sync/user -X DELETE -H 'Content-Type: application/json' -d '{"email":"laura@ab.cd"}' 
```



## Task

### create new task

```bash
curl  http://localhost:8888/dgraph_sync/task -X POST -H 'Content-Type: application/json' -d '
	{	
  			"title":"run 10000m", 
  			"description":"it is hard",
  			"priority": "midum",
  			"status": "todo",
  			"reporter": {
  				"email": "laura@ab.cd"
  			},
  			"assignee": {
  				"email": "yang@ab.cd"
  			}
  	}' 
```



### update task

```bash
curl  http://localhost:8888/dgraph_sync/task -X PUT -H 'Content-Type: application/json' -d '
{
	"uid": "0x9c49",
  "title":"run 100m",
  "description":"try you best",
  "priority": "midum",
  "status": "todo"
}'
```



### get task

```bash
curl  http://localhost:8888/dgraph_sync/task -X GET -H 'Content-Type: application/json' -d '{"uid":"0x9c49"}' 
```



### delete task

```bash
curl  http://localhost:8888/dgraph_sync/task -X DELETE -H 'Content-Type: application/json' -d '{"uid":"0x9c49"}' 
```



