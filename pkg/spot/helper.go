package spot

import (
	"encoding/json"
	"strconv"
)

func toJson(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func str2Float64(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		// todo
		panic(err)
	}
	return f
}
