package dsql

import "database/sql/driver"

type QueryResponse struct {
	ConsumedCapacity struct {
		CapcityUnits int
		TableName    string
	}
	Count            int
	Items            []Item
	LastEvaluatedKey Item
}

// problematic due to the schemaless nature of dynamo
// driver wants fixed-width scan operations but items are variable-width
// consider preprocessing and properly accepting nil gaps
func (q *QueryResponse) Columns() (cols []string) {
	for id, _ := range q.Items[0] {
		cols = append(cols, id)
	}
	return cols
}

func (q *QueryResponse) Values() (values [][]driver.Value) {
	cols := q.Columns()

	for _, item := range q.Items {
		row := []driver.Value{}
		for _, col := range cols {
			if v := item[col].Value(); v != nil {
				row = append(row, v)
			}
		}
		values = append(values, row)
	}
	return values
}