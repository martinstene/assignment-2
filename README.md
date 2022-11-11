# Assignment - 2 & 3
_PROG2005 - Cloud Technologies_

**Table of contents**

[[_TOC_]]

# Description
This is a project from the PROG2005 course at NTNU.
This project is about creating a RESTful API using Golang.
The assignment is to use other APIs and retrieve information in form
of a json body structure. 
A challenge is also to use these APIs and
all their features to use these same features on your own API. 
The assignment also tests the use of graphql and hwo to get that
information back to a struct in the form of a JSON.
The biggest challenge of this assignment is the webhooks and how to 
invoke them. And also caching using the firestore cloud DB.

The API will use the https://covid19-graphql.vercel.app/ API for
gathering covid cases _(GRAPHQL)_ and the 
https://covidtracker.bsg.ox.ac.uk/about-api API for
gathering stringency and policies.

This project also take advantage of Firebase cloud services called 'Firestore.'
That is used to store webhooks and for the use of caching results of all
endpoints to rate limit the amount of calls needed to the different APIs.

# Cloning
To use this software you will have to clone this from GitLab. You will need to enter a terminal and enter the following:

    1. Cloning the Remote Repo to your Local host. example: git clone https://github.com/user-name/repository.git.
    2. Pulling the Remote Repo to your Local host. First you have to create a git local repo by, example: git init or git init repo-name then, git pull https://github.com/user-name/repository.git.

# Useage

To use this API you will need four URLs.

1. /corona/v1/cases/
2. /corona/v1/policy/
3. /corona/v1/status/
4. /corona/v1/notifications/


## 1 /corona/v1/cases/
To use the endpoint you will have to enter a country-name. 
After this you will get a json body returned to you with information about
covid cases in that specified country.

Advanced tasks are done, caching should work and is deleted after a
certain amount of time from the collection.

    Method: GET
    Path: /corona/v1/cases/{:country_name}
    Path: /corona/v1/cases/{:cca3_code}

Body (Example)

    {
        "country": "Norway",
        "date": "2022-03-05",
        "confirmed": 1305006,
        "recovered": 0,
        "deaths": 1664,
        "growth_rate": 0.004199149089414866
    }


## 2 /corona/v1/policy/
To use the endpoint you will have to enter a country-code and a date.
The date parameter is optional and if you don't enter any you will get today'notification
date which may or may not have any dataFromDb. If it doesn't an error will be returned
saying that there is no dataFromDb available.

This is a known design choice seeing as the API the info is gotten from
chooses to say that it has no dataFromDb about the country and then proceeding to tell you
that this country, which it has no information about, has 20 policies in place. Which 
in least some countries, is just blatantly wrong.

Advanced tasks are done, caching should work and is deleted after a
certain amount of time from the collection.

    Method: GET
    Path: /corona/v1/policy/{:cca3_code}{?scope=YYYY-MM-DD}
    Path: /corona/v1/policy/{:country_name}{?scope=YYYY-MM-DD}

Body (Example)

    {
        "country_code": "FRA",
        "scope": "2022-03-01",
        "stringency": 63.89,
        "policies": 0
    }

## 3 /corona/v1/status/
This is a status interface which checks if services are up and running and showing
you the status codes for the given services. Also telling you how many webhooks are 
stored in the firestore database.

    Method: GET
    Path: corona/v1/status/

Body

    {
        "cases_api": "<http status code for *Covid 19 Cases API*>",
        "policy_api": "<http status code for *Corona Policy Stringency API*>",
        ...
        "webhooks": <number of registered webhooks>,
        "version": "v1",
        "uptime": <time in seconds from the last service restart>
    }

## 4 /corona/v1/notifications/

This is a notification endpoint that uses webhooks and holds the 
amount of calls. It has some methods.

1. Register a webhook and get an ID as response
2. Delete a webhook based on ID
3. View a webhook based on ID or view all if no ID is provided
4. Webhook invocation based on the amount of country calls

Here are all the methods with explanations below.

**Add a webhook**

    Method: POST
    Path: /corona/v1/notifications/

Where you post a body which is in the example below using postman, this
will be stored, and you will get a response which will be the generated
webhook ID.

Body(Example)

    {
        "url": "https://localhost:8080/client/",
        "country": "France",
        "calls": 5
    }

Response Body (Example)

    {
        "webhook_id": "OIdksUDwveiwe"
    }

**Delete a webhook**

    Method: DELETE
    Path: /corona/v1/notifications/{id}

Gets a webhook based on the ID given and then deletes this webhook from
the DB.


**View a webhook**

    Method: GET
    Path: /corona/v1/notifications/{id}

Gets the webhook based on the ID given and then decodes and gives this
completed webhook as a response. With the body below.

Body(Example)

    {
        "webhook_id": "OIdksUDwveiwe",
        "url": "https://localhost:8080/client/",
        "country": "France",
        "calls": 5
    }

**View all webhooks**

    Method: GET
    Path: /corona/v1/notifications/

If no ID is given it will give the completed webhooks from the DB.
It will give this in the format as the example below.


Body(Example)

    [{
        "webhook_id": "OIdksUDwveiwe",
        "url": "https://localhost:8080/client/",
        "country": "France",
        "calls": 5
    },
    {
        "webhook_id": "DiSoisivucios",
        "url": "https://localhost:8080/client/",
        "country": "Norway",
        "calls": 2
    },
    ...
    ]

**Webhook invocation**

    Method: POST
    Path: <url specified in the corresponding webhook registration>

Sends a post request to the specified URL in the webhook which the
user gave when sending the first post request when adding a webhook.
You can use the [webhook.site](https://webhook.site/) to test this.

Body(Example)

    {
        "webhook_id": "OIdksUDwveiwe",
        "country": "Norway",
        "calls": 3
    }


### Firebase

This project uses firebase as a database for the project. Every response except
the status endpoint writes a collection to the DB and is stored. The DB
caches the results which is checked by the program after getting specific
field information.

# Deployment

**Openstack:**
The service is deployed using SkyHigh from NTNU, which means that you will
have to be on the NTNU network if you want to test that version. 
To access this you will need to go to the IP and port number which is linked
[here](http://10.212.136.201:8080/)

_note: this is the old version and the one on google is updated._

**Google cloud:**
The API is deployed to the google cloud service and will 
stay there for a month after uploading this API to it.
You do not have to be on any specific network to use it, just
make sure that when you're accessing it you use HTTP and **not** HTTPS.
It can be accessed on this IP and port number [here](http://34.88.164.157:8080/) 

It is deployed using a docker file which is in the [repo](https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2022-workspace/mmstene/assignment-2/-/tree/master).

# Known issues

_code issue:_ The tests for the firestore functionality are lackluster, 
they include some relevant testing but since most of the methods 
written were written for printing, these are really hard to test 
without breaking them all down.

Issues with the restcountries API; some countries don't match up with 
the Cases and stringency APIs.

# Examples
These are a few examples per endpoint which makes testing it easier and faster. The examples 
show one easy example for each of teh deployed platforms.

- `Cases endpoint`


    http://localhost:8080/corona/v1/cases/norway
    https://10.212.136.201:8080/corona/v1/cases/norway
    http://34.88.164.157:8080/corona/v1/cases/norway


- `Policy endpoint`

**GET**


    http://localhost:8080/corona/v1/policy/NOR?scope=2022-01-01
    https://10.212.136.201:8080/corona/v1/policy/NOR?scope=2022-01-01
    http://34.88.164.157:8080/corona/v1/policy/NOR?scope=2022-01-01

    http://localhost:8080/corona/v1/policy/NOR
    https://10.212.136.201:8080/corona/v1/policy/NOR
    http://34.88.164.157:8080/corona/v1/policy/NOR

- **`Notification endpoint`**


**GET**

    GET every webhook:
    http://localhost:8080/corona/v1/notifications/
    https://10.212.136.201:8080/corona/v1/notifications/
    http://34.88.164.157:8080/corona/v1/notifications/

    GET one webhook
    http://localhost:8080/corona/v1/notifications/{ID}
    https://10.212.136.201:8080/corona/v1/notifications/{ID}
    http://34.88.164.157:8080/corona/v1/notifications/{ID}

**POST**

    {
        "url": "webhook.site/{your_unique_URL}",
        "country": "Norway",
        "calls": 2
    }

**DELETE**

    http://localhost:8080/corona/v1/notifications/{ID}
    https://10.212.136.201:8080/corona/v1/notifications/{ID}
    http://34.88.164.157:8080/corona/v1/notifications/{ID}