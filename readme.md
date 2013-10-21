# DSQL - DynamoDB Structured Query Language

EXPERIMENTAL SQL dialect for interacting with DynamoDB.

Package implements a `database/sql` driver:

    import (
      _ "github.com/cmdrkeene/dsql"
      "database/sql"
    )

    func main() {
      db, _ := sql.Open("dynamodb", "dynamodb://access:secret@us-east-1")
      rows, _ := db.Query("SELECT name FROM users WHERE id=$1", 123)
      // ...
    }

## Statments

The package does not prevent you from writing incorrect or insane (e.g.
unbounded scans) queries. It is merely a translator.

Here are some statements the parser understands:

    SELECT * FROM users;

    SELECT * FROM users LIMIT 10;

    SELECT * FROM users LIMIT 10 ORDER BY id ASC;

    SELECT id, name FROM users;

    SELECT id, name FROM users WHERE id = "1" AND name > 2;

    SELECT id, name FROM users WHERE id = 1 AND name BETWEEN("a", "z");

    INSERT INTO users (id, name) VALUES (1, "A");

    UPDATE users SET name = "B" WHERE name = "A";

    DELETE FROM users WHERE name = "A";

    CREATE TABLE messages (
      group string HASH,
      id number RANGE,
      created string,
      updated string,
      INDEX created WITH (HASH=group, RANGE=created, PROJECTION=(id, created)),
      INDEX updated WITH (HASH=group, RANGE=updated, PROJECTION=ALL)
    )
    WITH (READ=10, WRITE=10);

## TODO

* interpolate args in query
* test against actual dynamo
* insert multiple items
* create table with local secondary index

## NICE TO HAVE

* `Expected` conditions (how does this map to SQL?)
* `BatchWriteItem` is smarter than 25 item cap (batch in driver)
* `Scan` operations (maybe a cursor?)
* `RETURNING` keyword for `ReturnValues`

