# Elasticsearch Config

## Prerequisite

- Download Elasticsearch and run it locally
- Elasticsearch service will run on port :9200 by default



## Data 

The task data structure will be stored in ES:

| attributes  | type   | description        |
| ----------- | ------ | ------------------ |
| uid         | string | task uid in Dgraph |
| title       | string | task title         |
| description | string | task description   |



## Mapping

```
PUT /tasks
{
    "mappings": {
      "properties": {
        "uid": {
            "type": "keyword"
          }
      }
    }
}
```



## Check tasks in ES

```
GET /_search
{
    "query": {
        "match_all": {}
    }
}
```



# Dgraph

## Prerequisite

- install docker

- run Dgraph container with docker:

  ```bash
  docker run --rm -it -p 9001:8080 -p 9080:9080 -p 9000:8000 -v ~/dgraph:/dgraph dgraph/standalone:v20.03.0
  ```

  rpc port: 9080

  http port: 9001

  console port: 9000



## schema

```
curl "localhost:9001/alter" -XPOST -d $'
	email: string @index(exact) .
  name: string .
  work_on: [uid] .
  manage: [uid] .

  notification: [uid] .
  err: [uid] .

  title: string .
  description: string .
  assignee: [uid] @reverse .
  status: string .
  priority: string .
  reporter: [uid] @reverse .

  message: string .  
  n_task: [uid] @reverse .
  send_to: [uid] @reverse .

  sync_action: string .
  try_times: int .
  err_task: [uid] @reverse .
  
  create_date: int .

  type User {
    email
    name
  }

  type Task {
    title
    description
    assignee
    status
    priority
    reporter
  }

  type Notification {
    n_task
    message 
    send_to
    create_date
  }

  type SynError {
    sync_action
    err_task
    try_times
    create_date
  }
'
```



# go-micro

## Prerequisite

- install go-micro

  ```bash
  go install github.com/micro/micro/v2
  ```

  

## Run micro service

- run the server

  ```bash
  micro server
  ```

- run micro api

  ```bash
  micro api --address=0.0.0.0:8888 --handler=api --enable_rpc
  ```



# run the code

```bash
$ cd path/to/dgrpah-sync-srv
$ go get -d -v ./...
$ go run api/api.go
$ go run srv/main.go
```

