Open_Gov 

app for connecting people to their government

Must start zookeeper/kafka before running app, see below


TODO:
- edit userReps DB to include propublica ids
- integrate Propublica's apis
- move all data sources to AWS
- Add universal ids to reps (from)
- improve google civic api consumption at a local level (AdminstrativeArea1&2)
- create a kafka consumer go script that reads in messages and writes to mysql db
- check for duplicates on userRep edit function
- set up external data sources (S3?) for users, reps, etc.
- set up CI/CD pipeline
- deploy to AWS EC2

DONE:
- in init: add mysql query to pull from representatives and switch from repDB to rep map
- in init: add mysql queries to pull from user_favorite_reps, switch repGUID to string
- set up kafka stream for adding/removing reps


Install instructions
- Install Go: https://golang.org/doc/install  
- Install Postman (for API testing): https://www.postman.com/downloads/  
- Install MySQL (for local DB): https://dev.mysql.com/downloads/mysql/  
- Install MySQL Workbench: https://downloads.mysql.com/archives/workbench/  
- Install Kafka/Zookeeper (for Mac):  
    - FIRST:   
    `brew tap adoptopenjdk/openjdkbrew`  
    `cask install adoptopenjdk8`  
    - Then: https://medium.com/@Ankitthakur/apache-kafka-installation-on-mac-using-homebrew-a367cdefd273  

- Install Kafka/Zookeeper (for Windows):
    - Refer to article for more information https://towardsdatascience.com/running-zookeeper-kafka-on-windows-10-14fc70dcc771
    - To download Zookeeper from Oracle, must make account with Oracle https://www.oracle.com/index.html
    - Alternatively- May download Chocolatey, https://chocolatey.org/install, and install SDK from powershell
     - To download Kafka- https://kafka.apache.org/downloads. Select the Scala tgz file under binaries and download. Extract to C drive. *Note May have to extract twice for bin        and config folders
     - Other helpful resources
            -https://www.geeksforgeeks.org/how-to-install-and-run-apache-kafka-on-windows/ 
            -https://www.learningjournal.guru/article/kafka/installing-kafka-on-windows/
ref for below: https://medium.com/@Ankitthakur/apache-kafka-installation-on-mac-using-homebrew-a367cdefd273  
To set up kafka comsumer/producer:  
`$ su <pc account>` (if needed)  
start zookeeper server  
`$ zookeeper-server-start /usr/local/etc/kafka/zookeeper.properties`  
start kafka server  
`$ kafka-server-start /usr/local/etc/kafka/server.properties`  
create kafka topic  
`$ kafka-topics --create --zookeeper localhost:2181 --replication-factor 1 --partitions 1 --topic test`
start kafka producer  
`$ kafka-console-producer --broker-list localhost:9092 --topic test`  
start kafka consumer  
`$ kafka-console-consumer --bootstrap-server localhost:9092 --topic test --from-beginning`

