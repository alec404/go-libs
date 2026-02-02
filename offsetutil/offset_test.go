package offsetutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Offset(t *testing.T) {
	str := Md5ToDecimalMod("8436fdbdc599b4fd9eeeb63b27139232", 10000)
	strVal := "1602"
	assert.Equal(t, strVal, str)

	threeStr := Md5ToDecimalMod("0ac5b8dca43069438c05a154b59b3dc5", 10000)
	threeStrVal := "0629"
	assert.Equal(t, threeStrVal, threeStr)

	zeroStr := Md5ToDecimalMod("", 10000)
	zeroStrVal := "0000"
	assert.Equal(t, zeroStrVal, zeroStr)

	offset := GenerateOffset("2024-12-30 22:59:48", "8436fdbdc599b4fd9eeeb63b27139232")
	var offsetVal int64 = 17355707881602
	assert.Equal(t, offsetVal, offset)

	threeOffset := GenerateOffset("2024-12-30 22:59:48", "0ac5b8dca43069438c05a154b59b3dc5")
	var threeOffsetVal int64 = 17355707880629
	assert.Equal(t, threeOffsetVal, threeOffset)
}
