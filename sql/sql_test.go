package sql

import "testing"

func TestCreate(t *testing.T) {
	t.Skip()
	var in, out, trans string
	// CREATE [EXTERNAL] TABLE [IF NOT EXISTS] [db_name.]table_name
	//   [(col_name data_type [COMMENT col_comment], ...)]
	//   [COMMENT table_comment]
	//   [PARTITIONED BY (col_name data_type [COMMENT col_comment], ...)]
	//   [CLUSTERED BY (col_name, col_name, ...) [SORTED BY (col_name [ASC|DESC], ...)] INTO num_buckets BUCKETS]
	//   [SKEWED BY (col_name, col_name, ...) ON ([(col_value, col_value, ...), ...|col_value, col_value, ...]) (Note: only available starting with 0.10.0)]
	//   [
	//    [ROW FORMAT row_format] [STORED AS file_format]
	//    | STORED BY 'storage.handler.class.name' [WITH SERDEPROPERTIES (...)]  (Note: only available starting with 0.6.0)
	//   ]
	//   [LOCATION hdfs_path]
	//   [TBLPROPERTIES (property_name=property_value, ...)]  (Note: only available starting with 0.6.0)
	//   [AS select_statement]  (Note: this feature is only available starting with 0.5.0, and is not supported when creating external tables.)
	in = `
	CREATE TABLE messages (
		chat_id text hash key, 
		id int range key,
		name text, 
		body text
	) WITH read_units=10, write_units=10
	`
	out = `{
		"TableName": "messages",
		"AttributeDefinitions": [
			{
				"AttributeName": "chat_id",
				"AttributeType": "S"
			},
			{
				"AttributeName": "id",
				"AttributeType": "N"
			}
		],
		"KeySchema": [
			{
				"AttributeName": "chat_id",
				"KeyType": "HASH"
			},
			{
				"AttributeName": "id",
				"KeyType": "RANGE"
			}
		],
		"ProvisionedThroughput": {
			"ReadCapacityUnits": 10,
			"WriteCapacityUnits": 10
		}
	}`
	trans, err := Translate(in)
	if err != nil {
		t.Error(err)
	}
	if trans != out {
		t.Error(in, out, trans)
	}
}

func TestSelect(t *testing.T) {
	var in, out, trans string
	// SELECT [ALL | DISTINCT] select_expr, select_expr, ...
	// FROM table_reference
	// [WHERE where_condition]
	// [GROUP BY col_list]
	// [CLUSTER BY col_list
	//   | [DISTRIBUTE BY col_list] [SORT BY col_list]
	// ]
	// [LIMIT number]
	in = "SELECT chat_id, name FROM messages WHERE chat_id=$1"
	out = `{
		"TableName":"messages",
		"KeyConditions": {
	        "chat_id": {
	            "AttributeValueList": [{"S": "1"}],
	            "ComparisonOperator": "EQ"
	        }
	   	}
	}`
	trans, err := Translate(in)
	if err != nil {
		t.Error(err)
	}
	if trans != out {
		t.Error(in, out, trans)
	}
}
