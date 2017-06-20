package uweb

//
// View data
//
type Map map[string]interface{}

func (m Map) Merge(data Map) {
	for k, v := range data {
		m[k] = v
	} 
}
