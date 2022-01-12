package fur

import "strconv"

func Str2Float64(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		// todo
		panic(err)
	}
	return f
}
