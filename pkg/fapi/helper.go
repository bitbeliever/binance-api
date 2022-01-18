package fapi

import (
	"encoding/json"
	"math"
	"strconv"
)

func Str2Float64(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		// todo
		panic(err)
	}
	// todo using math round
	//return math.Round(f*100) / 100
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

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return math.Floor(num*output) / output
}
