package json

import "encoding/json"

func Marshal(v any) []byte {
	result, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return result
}

func Unmarshal[T any](data []byte) T {
	var result T
	err := json.Unmarshal(data, &result)
	if err != nil {
		panic(err)
	}
	return result
}
