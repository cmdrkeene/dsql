# DSQL - Dynamo Structured Query Language

 SQL dialect for interacting with DynamoDB.

 Package implements a `database/sql` driver.

## TODO

* test against fake dynamo
* test against actual dynamo
* insert multiple items
* create table with local secondary index

## NICE TO HAVE

* `Expected` conditions (how does this map to SQL?)
* `BatchWriteItem` is smarter than 25 item cap (batch in driver)
* `Scan` operations (maybe a cursor?)
* `RETURNING` keyword for `ReturnValues`

