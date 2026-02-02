package offsetutil

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"
)

// Md5ToDecimalMod 对MD5哈希进行转换并进行模运算
func Md5ToDecimalMod(md5Hash string, mod int64) string {
	// 确保 mod 是合法值
	if mod <= 0 {
		return "0"
	}

	// 计算结果的最小位数
	minDigits := len(strconv.FormatInt(mod-1, 10))

	// 将16进制字符串转换为大整数
	bigIntHash, ok := new(big.Int).SetString(md5Hash, 16)
	if !ok {
		return strings.Repeat("0", minDigits) // 返回全零字符串
	}
	// 创建mod的大整数表示
	bigMod := big.NewInt(mod)
	// 进行模运算
	result := new(big.Int).Mod(bigIntHash, bigMod)
	// 转换为字符串
	resultStr := result.String()

	// 根据 mod 的位数进行补零
	if len(resultStr) < minDigits {
		resultStr = strings.Repeat("0", minDigits-len(resultStr)) + resultStr
	}
	return resultStr
}

// GenerateOffset 根据时间和md5模运算生成offset
func GenerateOffset(publishTime string, md5 string) int64 {
	// 解析时间字符串
	t, err := time.ParseInLocation("2006-01-02 15:04:05", publishTime, time.Local)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return 0
	}
	// 获取Unix时间戳（秒）
	strTimestamp := strconv.FormatInt(t.Unix(), 10)
	// 获取MD5模运算后的结果
	strMd5ToNumber := Md5ToDecimalMod(md5, 10000)
	// 拼接字符串并转换为整数
	offset, err := strconv.ParseInt(strTimestamp+strMd5ToNumber, 10, 64)
	if err != nil {
		fmt.Println("Error converting offset:", err)
		return 0
	}
	return offset
}
