package main

import (
	"encoding/base64"
	"encoding/hex"
	"github.com/bitbeliever/binance-api/pkg/fapi"
	"log"
)

func main() {
	//b, err := base64.StdEncoding.DecodeString("TGBktDVykhChFj3AJE33j2nMX67skQdzCuX")
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//println(hex.EncodeToString(b))

	b2 := "4c6064b435729210a1163dc0244df78f69cc5faeec9107730ae5"
	_ = b2
	a2 := "60b32561A880eD2A1ee42AC72796455B86B7B1c8"
	bb, err := hex.DecodeString(a2)
	if err != nil {
		log.Println(err)
		return
		println(base64.StdEncoding.EncodeToString(bb))
		fapi.GetLogger().Println("afaasfasf")

	}

	fapi.PositionMode()
}
