package postgres

import (
	"errors"
	"lostpets"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetFilterStr(t *testing.T) {
	type testInput struct {
		key    string
		index  string
		filter *lostpets.Filter
	}
	type testResult struct {
		filterStr string
		value     interface{}
		err       error
	}

	type test struct {
		name     string
		input    testInput
		expected testResult
	}

	tests := []test{
		test{
			name: "Should return case insensitive equals",
			input: testInput{
				key:   "dbField",
				index: "mapIndex",
				filter: &lostpets.Filter{
					Comparator: "=",
					Value:      "testValue",
				},
			},
			expected: testResult{
				filterStr: "lower(dbField) like :mapIndex",
				value:     "testvalue",
			},
		},
		test{
			name: "Should return equals",
			input: testInput{
				key:   "dbField",
				index: "mapIndex",
				filter: &lostpets.Filter{
					Comparator: "=",
					Value:      35,
				},
			},
			expected: testResult{
				filterStr: "dbField = :mapIndex",
				value:     35,
			},
		},
		test{
			name: "Should return default",
			input: testInput{
				key:   "dbField",
				index: "mapIndex",
				filter: &lostpets.Filter{
					Comparator: "<=",
					Value:      35,
				},
			},
			expected: testResult{
				filterStr: "dbField <= :mapIndex",
				value:     35,
			},
		},
		test{
			name: "Should return lower in",
			input: testInput{
				key:   "dbField",
				index: "mapIndex",
				filter: &lostpets.Filter{
					Comparator: "in",
					Value:      []string{"ONE", "tWo"},
				},
			},
			expected: testResult{
				filterStr: "lower(dbField) IN (:mapIndex)",
				value:     []string{"one", "two"},
			},
		},
		test{
			name: "Should return in",
			input: testInput{
				key:   "dbField",
				index: "mapIndex",
				filter: &lostpets.Filter{
					Comparator: "in",
					Value:      []int{34, 56},
				},
			},
			expected: testResult{
				filterStr: "dbField IN (:mapIndex)",
				value:     []int{34, 56},
			},
		},
		test{
			name: "Should return in",
			input: testInput{
				key:   "dbField",
				index: "mapIndex",
				filter: &lostpets.Filter{
					Comparator: "in",
					Value:      23,
				},
			},
			expected: testResult{
				err:   errors.New("unsupported type for IN filter: int"),
				value: 23,
			},
		},
		test{
			name: "Should return In Error if in filter is empty for type int",
			input: testInput{
				key:   "dbField",
				index: "mapIndex",
				filter: &lostpets.Filter{
					Comparator: "in",
					Value:      []int{},
				},
			},
			expected: testResult{
				err:   errorEmptyIn,
				value: []int{},
			},
		},
		test{
			name: "Should return errorEmptyIn if in filter is empty for type string",
			input: testInput{
				key:   "dbField",
				index: "mapIndex",
				filter: &lostpets.Filter{
					Comparator: "in",
					Value:      []string{},
				},
			},
			expected: testResult{
				err:   errorEmptyIn,
				value: []string{},
			},
		},
	}

	for _, test := range tests {
		str, err := getFilterStr(test.input.key, test.input.index, test.input.filter)
		assert.Equal(t, test.expected.filterStr, str)
		assert.Equal(t, test.expected.value, test.input.filter.Value)
		assert.Equal(t, test.expected.err, err)
	}
}

func TestGetFilters(t *testing.T) {
	lowerTime := time.Date(2014, 4, 12, 0, 0, 0, 0, time.UTC)
	upperTime := time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)
	filters := lostpets.FilterMap{
		"createdat": { //with mapping and multi-value
			{
				Comparator: "<=",
				Value:      upperTime,
			},
			{
				Comparator: ">=",
				Value:      lowerTime,
			},
		},
		"id": { //no ,mapping
			{
				Comparator: "in",
				Value:      []string{"ABC", "oNE", "two"},
			},
		},
		"age": { //no ,mapping
			{
				Comparator: "=",
				Value:      56,
			},
		},
	}

	fieldMap := map[string]string{
		"createdat": "created_at",
	}

	filterStr, params, err := getFilters(fieldMap, filters)
	assert.NoError(t, err)
	assert.Contains(t, filterStr, "created_at <= :created_at0")
	assert.Contains(t, filterStr, "created_at >= :created_at1")
	assert.Contains(t, filterStr, "lower(id) IN (:id0)")
	assert.Contains(t, filterStr, "age = :age0")
	assert.Equal(t, map[string]interface{}{
		"age0":        56,
		"created_at0": upperTime,
		"created_at1": lowerTime,
		"id0":         []string{"abc", "one", "two"},
	}, params)
}
