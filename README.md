Open_Gov 

app for connecting people to their government


TODO:
- in init: add mysql queries to pull from representatives and user_favorite_reps tables
    - switch user_rep map to have it's own loading function and store repGUID has string, not int
- create a kafka consumer go script that reads in messages and writes to mysql db
- check for duplicates on userRep edit function
- set up external data sources (S3?) for users, reps, etc.
- set up CI/CD pipeline
- deploy to AWS EC2

DONE:
- set up kafka stream for adding/removing reps


To set up kafka comsumer/producer:  
`$ su <pc account>`  
start zookeeper server  
`$ zookeeper-server-start /usr/local/etc/kafka/zookeeper.properties`  
start kafka producer  
`$ kafka-console-producer --broker-list lohost:9092 --topic test`  
start kafka consumer  
`$ kafka-console-consumer --bootstrap-server localhost:9092 --topic test --from-beginning`

