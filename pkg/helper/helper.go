// Package helper global helper, todo remove
package helper

import (
	"encoding/json"
	"log"
)

func ToJson(v interface{}) string {
	b, err := json.Marshal(v)
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

	b, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(string(b))
}
