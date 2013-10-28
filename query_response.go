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
		row := make([]driver.Value, len(cols))
		for j, col := range cols {
			v := item[col].Value()
			if v != nil {
				row[j] = v
			} else {
				row[j] = []byte{}
			}
		}
		values = append(values, row)
	}
	return values
}
