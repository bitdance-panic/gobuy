package utils

import (
	"regexp"
	"strings"
)

// 移除代码块标记
func CleanBlock(rawSQL string) string {
	// 正则表达式匹配代码块标记（如 ```sql 或 ```）
	re := regexp.MustCompile(`(?i)^[[:space:]]*` + "`{3}" + `[[:alnum:]_]*\n`)
	cleaned := re.ReplaceAllString(rawSQL, "")
	// 移除末尾的代码块标记（如果有）
	cleaned = strings.TrimSuffix(cleaned, "```")
	// 返回清理后的 SQL
	return strings.TrimSpace(cleaned)
}
