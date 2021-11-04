package postgres

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	domain "lostpets"
)

var errorEmptyIn = errors.New("A filter with IN operation has no value")

func getFilters(fieldMap map[string]string, filterMap domain.FilterMap) (string, map[string]interface{}, error) {
	filterArr := make([]string, 0)
	params := make(map[string]interface{}, 0)

	for key, filters := range filterMap {
		for i, filter := range filters {
			dbField := key
			if str, ok := fieldMap[strings.ToLower(dbField)]; ok {
				dbField = str
			}

			index := dbField + strconv.Itoa(i)
			f, err := getFilterStr(dbField, index, &filter)
			if err != nil {
				return "", nil, err
			}
			if f != "" {
				filterArr = append(filterArr, f)
				params[index] = filter.Value
			}
		}
	}

	if len(filterArr) != 0 {
		return strings.Join(filterArr, " and "), params, nil
	}

	return "", nil, nil
}

func getFilterStr(key string, index string, filter *domain.Filter) (string, error) {
	switch filter.Comparator {
	case "=":
		//if interface is string
		if filter.Value == nil {
			return fmt.Sprintf("%s IS NULL", key), nil
		}

		if str, ok := filter.Value.(string); ok {
			filter.Value = strings.ToLower(str)
			return fmt.Sprintf("lower(%s) %s :%s", key, "like", index), nil
		}
		return fmt.Sprintf("%s %s :%s", key, filter.Comparator, index), nil
	case "in":
		//if interface is string array
		if strArr, ok := filter.Value.([]string); ok {
			if len(strArr) == 0 {
				return "", errorEmptyIn
			}

			for i := range strArr {
				strArr[i] = strings.ToLower(strArr[i])
			}
			return fmt.Sprintf("lower(%s) IN (:%s)", key, index), nil
		}

		if intArr, ok := filter.Value.([]int); ok {
			if len(intArr) == 0 {
				return "", errorEmptyIn
			}
			return fmt.Sprintf("%s IN (:%s)", key, index), nil
		}

		return "", fmt.Errorf("unsupported type for IN filter: %T", filter.Value)
	default:
		return fmt.Sprintf("%s %s :%s", key, filter.Comparator, index), nil
	}
}
