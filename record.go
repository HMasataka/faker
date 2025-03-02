package faker

type DB map[string]Records

func (d DB) Has(key string) bool {
	_, has := d[key]
	return has
}

func (d DB) HasAll(keys []string) bool {
	for i := range keys {
		if has := d.Has(keys[i]); !has {
			return false
		}
	}

	return true
}

type Records []Record
type Record map[string]any
