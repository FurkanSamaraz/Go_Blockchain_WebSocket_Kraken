package convertertw

import (
	"encoding/json"
	"fmt"
)

func Converter(con json.Number) (s string) {
	f, _ := con.Float64()

	askPrice_3 := f * 20
	s = fmt.Sprintf("%f", askPrice_3)
	return s
}
