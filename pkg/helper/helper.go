// Package helper global helper
package helper

import (
	"encoding/json"
	"log"
	"math"
	"strconv"
)

func ToJson(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func ToJsonIndent(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(b)
}

func JsonLog(data interface{}, err error) {
	if err != nil {
		log.Println(err)
		return
	}

	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(string(b))
}

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

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func Round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return math.Floor(num*output) / output
}
