package tool

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"math/rand"
	"strconv"

	"github.com/bwmarrin/snowflake"
)

func GenerateUUID() string {
	node, _ := snowflake.NewNode(1)
	snowflakeId := node.Generate().String()
	return snowflakeId
}

func GenerateUUID36() string {
	node, _ := snowflake.NewNode(1)
	snowflakeId := node.Generate()
	return strconv.FormatInt(snowflakeId.Int64(), 36)
}

func JsonToMap(jsonStr string) map[string]string {
	m := make(map[string]string)
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		return nil
	}
	return m
}

func MapToJson(m map[string]string) string {
	jsonByte, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(jsonByte)
}

func StrToMd5(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}

func RandomStr(strLen int64) string {
	letters := []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	randBytes := make([]rune, strLen)
	for i := range randBytes {
		randBytes[i] = letters[rand.Intn(len(letters))]
	}

	return string(randBytes)
}
