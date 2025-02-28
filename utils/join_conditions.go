package utils

import (
	"strings"
)

// 条件式のスライスを指定の演算子で結合
func JoinConditions(conditions []string, op string) string {
	return strings.Join(conditions, op)
}
