# tag-measurements-microservices
Temperature monitoring system  
# Architecture
Common work scheme:  

                                              
                                                 +-------------+                 
                                                 |  Database   |
                                                 +--+------+---+
                                                    |      ˄
                          HTTP Response + JSON      |      |
	                                    |           |      |
	 +-------------------------+        ˅     +-----˅------+------+
	 |     Wirelesstag API     +------------->+   tag-measurements-microservices APP   |
	 |     wirelesstag.com     +<-------------+   Web API (REST)  +<-------+
	 +-------------------------+        ˄     +-------------------+        |
	                                    |                                  |<--JSON 
	                               JSON + HTTP Requests                    ˅
	                                                              +--------+--------+
	                                                              | REST API Client |          
	                                                              +-----------------+
	  
Main purpose:
 - <b>System that collect data to our database.</b> (DONE) 
 - Fix 05.10.2020: Develop web application for warehouse clients, QA, system admins and other people that display temperature, humidity, signal and battery voltage of tags, measurement
  plot for certain period
 - Fix 14.05.2020: Collect <b>ALL DATA</b> from cloud. (?)
 
App architecture:

                                              +----------+
    +-----------------+                    +--+Notify SRV|
    |  Fetch Service  |                    |  +----------+
    | Real Time Spring+            +----+  |                 +------+
    +-----------------+----------->+    +<-+ +---------+     | User |
    |  Fetch Service  +----------->| DB +<---|Clean SRV|     |  DB  |
    +-----------------+   INSERT   +--+-+    +---------+     +---+--+
             ˄                        |                          ˄
             |JSON             FETCH  |                    FETCH ˅ MODIFY
       +-----+-----+            +-----˅-----+              +-----+----+
       | Data from |            |ResourceSRV|              | Auth SRV |
       |   cloud   |            +-----+-----+              +----+-----+
       +-----------+                  ˄                         ˄ 
                                      |                         |
                                      ˅       User credentials  |
                                   Client<----------------------+
                                                   JWT
                 
# Purpose of services

## Fetch service (Golang version)
Fetch service should fetch measurements every N time from cloud service.  
The main loop logic:

                                                                     +-----------------------------------+
                                                           mainLoop()|      get session ids              |
     +---------+   +-------------------+   +--------------------+    ˅  +--------------------+           |
     |  Start  +-->| Create app struct +-->+ Init db connection |------->  Init WST clients  +-------+   |
     +---------+   +-------------------+   +--------------------+       +--------------------+       |   |
                                                new   +------------------------+                     ˅   |
                               +----------------------+  Iterate WST clients   +<--------------------+   |
                               |              thread  +------------------------+                         |
                               |                                 /api/GetTagList                         |
                               |                              +------------------+                       |
                               +----------------------------->+ StoreMeasurement +-----------------------+
                                                              +------------------+


## Fetch service (Java version)
That fetch service should fetch data from N to M date and sync data.                                    
                                            
## Resource service
Resource service provides API that allow fetching data from datasource.
Default response body format - application/json.
Implemented 8 http endpoints:
 - /api/tagManagers  
 GET - return set of tag managers
 /api/tagManagers/{id:[0-9]+}
 - /api/tags  
 GET - return set of tags
 url parameters: 
    - mac - string value, MAC address of his tag manager (optional)
 - /api/temperatureZones  
 GET - return array of temperature zones (require ADMIN role)
 POST - create zone (require ADMIN role)
 PUT  - update zone (require ADMIN role)
 DELETE - delete zone (require ADMIN ROLE)
    
    
    {
        id: number,
        name: string,
        description: string,
        lower_temp_limit: number,
        higher_temp_limit: number
    }
    
    
 PUT - edit warehouse group -> /api/temperatureZones/{id:[0-9]+} (require privilege)
 - /api/measurements?uuidList=[... , ...]&startDate=...&endDate=...&epsilon=...
 GET - uuidList, startDate and endDate query parameters required. The epsilon parameter is optional.
 - /api/measurementsRT?uuidList=[... , ...]&startDate=...&endDate=...&epsilon=...
 GET - uuidList, startDate and endDate query parameters required. The epsilon parameter is optional.
 - /api/users - admin only - CRUD operations  
 
## Authentication service
It provides authentication API that issues JWT token by passing user credentials. 
The JWT token should be stored in X database for resource service usage.  

## Clean service
The clean service removes measurements rows filtered by a specified (> 40 days by default) date.

## Notify service
This service every 1 minute check rows in measurement table by comparing temperature column with lower_temp_limit and higher_temp_limit.
If records were found service send mail to emails that specified in models.TemperatureZone.NotifyEmails 

## Frontend
On the main page should be displayed group of tag managers and tags   
with select opportunity (checkbox), and display plot Temperature(timestamp), Humidity(timestamp), Voltage(timestamp), Signal(timestamp).  
Also, should be settings page that user can define fetch an interval of each tag manager.

## Project Structure

    .
    .
    ├── build                   # dockerfiles, deploy files
    ├── cmd                     # go main files
    │   ├── auth_service       
    │   ├── clean_service       
    │   ├── fetch_service       # reserve fetch service
    │   ├── notify_service       
    │   ├── resource_service    
    │   └── tgbot_service
    ├── configs                 # json configs for go applications
    ├── fetch_service           # main fetch service on Spring Boot
    ├── internal                # internal (not reuse) go files
    │   ├── auth_service
    │   │   ├── controllers
    │   │   └── structures
    │   ├── clean_service
    │   │   └── structures
    │   ├── fetch_service
    │   │   ├── api
    │   │   ├── structures
    │   │   └── types
    │   ├── notify_service
    │   │   └── structures
    │   └── resource_service
    │   │   ├── controllers
    │   │   └── structures
    │   ├── tgbot_service
    │       └── structures
    ├── pkg                     # common files (reusable)
    │   ├── datasource
    │   ├── dto
    │   ├── models
    │   ├── repository
    │   └── utils
    ├── sql                     # sql commands for deploy database
    └── web                     # web ui
        └── thermo-ng
            ├── e2e
            └── src



# How to build (docker)
 ${ROOT} - project directory
 Backend  - Launch shell-script docker_commands.sh in ${ROOT}/build path to build golang services and docker images.  
 Frontend - Launch shell-script build_save.sh in ${ROOT}/web/thermo-ng path to build Angular app and docker image.
 Fetch Service on Java (Spring Boot) - run in ${ROOT}/fetch_service docker build
 
# Code of conduct
The main programming languages to implement that system can be:
 - Python 2/3 (PEP8)
 - Go (Standard name convention)
 - Java (Java naming convention)
 - Pure JS (ES6) or JS framework (TypeScript) (Angular)
 
Also, contributor should add commentaries to functions, methods and obscure code.
