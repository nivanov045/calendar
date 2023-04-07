# Calendar System

## Features
* create a user
* create a meeting in a user's calendar with a list of invited users
* get meeting details
* accept or decline another user's invitation
* find all user meetings for a given time range
* for a given list of users and a minimum meeting duration, find the nearest time interval in which all these users are free

Meetings in the calendar can have the following recurrence settings:

* no repeat
* every day
* every week
* every year
* Monday through Friday

## Server Settings
To specify the server address, you can use the command line flag `a` or the environment variable `ADDRESS`. By default, `127.0.0.1:8080`.

## Usage
The server accepts `POST` and `GET` requests with `content-type application/json`.

### Create a user
#### Request
`POST` to `/create-user` in the format

    {
        "info" : {
            "name" : "Ivan"
        }
    }

#### Responses
* `200 OK` upon successful user addition
* `400 Bad Request` upon request error
* `404 Not Found` upon other errors

#### Successful response format
id of the created user

    {
        "id": "8c487d7a-a734-4c08-82f2-162c854ce827"
    }


### Create a meeting in the calendar
#### Request
`POST` to `/create-event-with-users` in the format

    {
        "candidates" : ["8c487d7a-a734-4c08-82f2-162c854ce827"],
        "participants" : ["c10ab64d-3860-46ef-bed6-46b8d3759928"],
        "start" : "2022-09-02T10:00:05Z",
        "finish" : "2022-09-02T11:00:05Z",
        "repeat_type" : 0,
        "info": {
            "name": "Some meeting name"
        }
    }
where `repeat_type` takes one of the following values depending on the type of repetition:
* `0` - no repeat
* `1` - every day
* `2` - every week
* `3` - every year
* `4` - Monday through Friday

#### Responses
* `200 OK` upon successful event addition to the calendar
* `400 Bad Request` upon request error
* `404 Not Found` upon other errors

#### Successful response format
id of the created event:

    {
        "id": "788dfa05-0f5d-4799-899a-c3b0e9eb3044"
    }


### Get meeting details
#### Request
`GET` to `/event-details` with the event id in the format

    {
        "event": "788dfa05-0f5d-4799-899a-c3b0e9eb3044"
    }

#### Responses
* `200 OK` upon successful event existence and successful detail retrieval
* `400 Bad Request` upon request error
* `404 Not Found` upon other errors, including event absence

#### Successful response format
All information about the event:

    {
        "candidates" : ["8c487d7a-a734-4c08-82f2-162c854ce827"],
        "participants" : ["c10ab64d-3860-46ef-bed6-46b8d3759928"],
        "start" : "2022-09-02T10:00:05Z",
        "finish" : "2022-09-02T11:00:05Z",
        "repeat_type" : 0,
        "info": {
            "name": "Some meeting name"
        }
    }


### Accepting an invitation to a meeting
#### Request
`POST` to `/accept-invitation` with the user id and event id in the format

    {
        "user": "8c487d7a-a734-4c08-82f2-162c854ce827"
        "event": "788dfa05-0f5d-4799-899a-c3b0e9eb3044"
    }

#### Responses
* `200 OK` upon successful invitation acceptance
* `400 Bad Request` upon request error
* `404 Not Found` upon other errors

### Declining an invitation to a meeting
#### Request
`POST` to `/reject-invitation` with the user id and event id in the format

    {
        "user": "8c487d7a-a734-4c08-82f2-162c854ce827"
        "event": "788dfa05-0f5d-4799-899a-c3b0e9eb3044"
    }

#### Responses
* `200 OK` upon successful invitation rejection
* `400 Bad Request` upon request error
* `404 Not Found` upon other errors

### Getting all user meetings within a specified interval
#### Request
`GET` to `/events`  with the user id and time interval in the format

    {
        "user" : "8c487d7a-a734-4c08-82f2-162c854ce827",
        "from" : "2022-09-02T10:00:05Z",
        "to"   : "2022-09-02T11:00:05Z",
    }

#### Responses
* `200 OK` upon successful event retrieval
* `400 Bad Request` upon request error
* `404 Not Found` upon other errors

#### Successful response format
All information about events intersecting with the specified interval:

    {
        [
            {
                "id" : "375d9831-592c-4373-8398-e22a54eaff2c"
                "candidates" : [],
                "participants" : ["c10ab64d-3860-46ef-bed6-46b8d3759928"],
                "start" : "2022-09-02T10:00:05Z",
                "finish" : "2022-09-02T11:00:05Z",
                "repeat_type" : 0,
                "info": {
                    "name": "Some meeting name 1"
                }
            },
            {
                "id" : "36492a24-24ba-478b-90b4-583486123291"
                "candidates" : ["87bbd840-a07a-4ef2-9407-09523695f690"],
                "participants" : ["ce48717b-ec68-4dfd-8917-fe59c8724b7e"],
                "start" : "2022-09-02T10:00:05Z",
                "finish" : "2022-09-02T11:00:05Z",
                "repeat_type" : 1,
                "info": {
                    "name": "Some meeting name 2"
                }
            }
        ]
    }


### Finding a free slot for a group of users
#### Request
`GET` to `/find-slot` with user ids, meeting duration in nanoseconds, and time after which the search for a meeting is no longer needed

    {
        "users" : [
            "8c487d7a-a734-4c08-82f2-162c854ce827",
            "375d9831-592c-4373-8398-e22a54eaff2c"
        ],
        "duration" : 1800000000000,
        "valid_until" : "2022-10-02T11:00:00Z",
    }

#### Responses
* `200 OK` upon successful slot identification
* `400 Bad Request` upon request error
* `404 Not Found` upon other errors

#### Successful response format
Successful response format:

    {
        "begin": "2022-09-05T11:00:00Z"
    }

## Planned improvements
* Add tests
* Add database support
