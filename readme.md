# DSQL - Dynamo Structured Query Language

 SQL dialect for interacting with DynamoDB.

 Writing Dynamo queries is cumbersome, SQL is easy.

## TODO

* database/sql interface
* test against fake dynamo
* test against actual dynamo
* insert multiple items
* create table with local secondary index

## NICE TO HAVE

* `Expected` conditions (how does this map to SQL?)
* `BatchWriteItem` is smarter than 25 item cap (batch in driver)
* `Scan` operations (maybe a cursor?)
* `RETURNING` keyword for `ReturnValues`

