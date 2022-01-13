package fapi

import (
	"encoding/json"
	"strconv"
)

func Str2Float64(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		// todo
		panic(err)
	}
	return f
}

func toJson(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)

		//log.Println(err)
		//return ""
	}
	return string(b)
}

func toJsonIndent(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)

		//log.Println(err)
		//return ""
	}
	return string(b)
}
