package counter

// SortBy is the type of a "less" function that defines the ordering.
// Based on: https://golang.org/pkg/sort/#example__sortKeys
type SortBy func(d1, d2 *Dict) bool

// SortAsc is a SortBy function for sorting in descending order based on values
var SortAsc = SortBy(func(d1, d2 *Dict) bool {
	return d1.Value < d2.Value
})

// SortDesc is a SortBy function for sorting in descending order based on values
var SortDesc = SortBy(func(d1, d2 *Dict) bool {
	return d1.Value > d2.Value
})

// Counter is an interface that allows for counting the number of unique keys added
// to the map so it can be sorted
type Counter struct {
	counterMap map[string]int
	Dict       []Dict
	SortByFunc SortBy
}

// Dict is a "Dictionary" entry that is used to sort based on the value
type Dict struct {
	Key   string
	Value int
}

// New returns an initialized Counter
func New() *Counter {
	return &Counter{
		counterMap: make(map[string]int),
		SortByFunc: SortAsc,
	}
}

// Increment increments the counter for the supplied key
func (m *Counter) Increment(key string) {
	m.counterMap[key]++

	newValue := m.counterMap[key]
	for i := range m.Dict {
		if m.Dict[i].Key == key {
			// Found one, update the value
			m.Dict[i].Value = newValue
			return
		}
	}
	// Key wasn't found in m.Dict, it needs to be initialized
	m.Dict = append(m.Dict, Dict{Key: key, Value: newValue})
}

// Len is part of sort.Interface.
func (m *Counter) Len() int {
	return len(m.Dict)
}

// Swap is part of sort.Interface.
func (m *Counter) Swap(i, j int) {
	m.Dict[i], m.Dict[j] = m.Dict[j], m.Dict[i]
}

// Less is part of sort.Interface. We use value as the value to sort by
func (m *Counter) Less(i, j int) bool {
	return m.SortByFunc(&m.Dict[i], &m.Dict[j])
}
