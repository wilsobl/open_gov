Open_Gov 

app for connecting people to their government


TODO:
- check for duplicates on userRep edit function
- set up kafka stream for adding/removing reps
- set up external data sources (S3?) for users, reps, etc.
- set up CI/CD pipeline
- deploy to AWS EC2


To set up kafka comsumer/producer:  
`$ su <pc account>`  
start zookeeper server  
`$ zookeeper-server-start /usr/local/etc/kafka/zookeeper.properties`  
start kafka producer  
`$ kafka-console-producer --broker-list lohost:9092 --topic test`  
start kafka consumer  
`$ kafka-console-consumer --bootstrap-server localhost:9092 --topic test --from-beginning`

