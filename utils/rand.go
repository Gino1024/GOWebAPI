package utils

import (
	"math/rand"
	"strings"
	"time"
)

//產生min ~ max 亂數值
func RandomInt(min int64, max int64) int64{
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	result := r.Int63n(max-min+1) + min // (max - min +1) 假設 10~100 > 取得 91以下的值 0~90 -> +10  最少就是10~100
	return result
}

const alphabet = "abcefghijklmnopqrstuvwxyz"
func RandomString(n int) string{
	var sb strings.Builder
	k := len(alphabet)

	for i:=0;i<n;i++{
		c := alphabet[RandomInt(0,int64(k)-1)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomOwner() string{
	return RandomString(6)
}

func RandomMoney() int64{
	return RandomInt(0, 1000)
}

func RandomCurrency() string{
	currency := []string{"EUR","USD","CAD"}
	len := len(currency)

	return currency[RandomInt(0, int64(len)-1)]
}