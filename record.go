package faker

type TableName string

type DB map[TableName]Records

func (d DB) Has(key TableName) bool {
	_, has := d[key]
	return has
}

func (d DB) HasAll(keys []TableName) bool {
	for i := range keys {
		if has := d.Has(keys[i]); !has {
			return false
		}
	}

	return true
}

type Records []Record
type Record map[string]any
