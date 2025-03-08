package faker

type TableName string

type ColumnName string
type ColumnNames []ColumnName

func (c ColumnNames) ToStrings() []string {
	s := make([]string, len(c))

	for i := range c {
		s[i] = string(c[i])
	}

	return s
}

func (c ColumnNames) IndexOf(columnName ColumnName) int {
	for i := range c {
		if c[i] == columnName {
			return i
		}
	}

	return -1
}

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

type Records struct {
	ColumnNames ColumnNames
	Values      [][]any
}

func (r Records) Len() int {
	return len(r.Values)
}

func (r Records) GetByColumnName(idx int, columnName ColumnName) any {
	return r.Values[idx][r.ColumnNames.IndexOf(columnName)]
}

func (r Records) GetLast() ([]any, error) {
	length := r.Len()

	if length == 0 {
		return nil, ErrNoRecord
	}

	return r.Values[length-1], nil
}
