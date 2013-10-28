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

	cols []string
}

// inspect all items because dynamo is schemaless
func (q *QueryResponse) Columns() []string {
	if q.cols != nil {
		return q.cols
	}

	ids := map[string]bool{}

	for _, item := range q.Items {
		for k, _ := range item {
			if !ids[k] {
				q.cols = append(q.cols, k)
				ids[k] = true
			}
		}
	}

	return q.cols
}

func (q *QueryResponse) Values() (values [][]driver.Value) {
	for _, item := range q.Items {
		row := make([]driver.Value, len(q.cols))
		for j, col := range q.cols {
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
