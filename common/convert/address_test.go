package convert

import (
	"fmt"
	"testing"
)

func TestZltcMustToAddress(t *testing.T) {
	zltc := "zltc_gLpU8MFUgdECP5wJdhoZXVr4kmYH5xuir"
	address := ZltcMustToAddress(zltc)
	fmt.Println(address.String())
}
