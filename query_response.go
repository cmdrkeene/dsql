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

// inspect all items because dynamo is schemaless
func (q *QueryResponse) Columns() (cols []string) {
	ids := map[string]bool{}

	for _, item := range q.Items {
		for k, _ := range item {
			if !ids[k] {
				cols = append(cols, k)
				ids[k] = true
			}
		}
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
