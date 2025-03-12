package trace

import (
	"math/rand"
	"time"
)

var (
	localRand   = rand.New(rand.NewSource(time.Now().UnixNano()))
	encodeTable = []rune{
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
		'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
		'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	}
	encodeTableLen = len(encodeTable)
)

// NewTraceId 返回 8 位随机字符串, 包含数字和字母, 区分大小写. 一般千万级别不会重复
func NewTraceId() string {
	return nextString(8)
}

// nextString 返回指定位随机字符串, 包含数字和字母, 区分大小写.
func nextString(length int) string {
	value := make([]rune, length)
	for i := 0; i < length; i++ {
		value[i] = encodeTable[nextInt(0, encodeTableLen)]
	}
	return string(value)
}

// nextInt 返回随机数, 范围 [startInclusive, endExclusive)
func nextInt(startInclusive int, endExclusive int) int {
	if startInclusive > endExclusive {
		panic("Start value must be smaller or equal to end value.")
	}
	if startInclusive < 0 {
		panic("Both range values must be non-negative.")
	}
	if startInclusive == endExclusive {
		return startInclusive
	}
	return startInclusive + localRand.Intn(endExclusive-startInclusive)
}
