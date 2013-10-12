# sql-to-dynamo

 ¯\_(ツ)_/¯

## Mapping

BatchGetItem    SELECT
BatchWriteItem  INSERT
CreateTable     CREATE
DeleteItem      DELETE
DeleteTable     DROP
DescribeTable   DESCRIBE
GetItem         SELECT
ListTables      SHOW TABLES
PutItem         INSERT
Query           SELECT
Scan            -- see below
UpdateItem      UPDATE
UpdateTable     ALTER

## Scan

Map Scan operations onto DECLARE/OPEN/FETCH/CLOSE in some kind of cursor model.
Since scans require special care and offer a simpler interface than the cursor
model, this is not worth it. This is mainly for completeness.
