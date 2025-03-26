package convert

import (
	"fmt"
	"testing"
)

func TestZltcMustToAddress(t *testing.T) {
	zltc := "zltc_nFgGmfSks6uQPT5hqMRQ8fdkKMprSUdbN"
	address := ZltcMustToAddress(zltc)
	fmt.Println(address.String())
}
