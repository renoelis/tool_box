package utils

import (
	"regexp"
	"strings"
	"time"
	"unicode"
)

// 常量定义
var (
	// 中文数字映射（大写）
	ChineseNumbersUpper = []string{"零", "壹", "贰", "叁", "肆", "伍", "陆", "柒", "捌", "玖"}
	
	// 中文数字映射（普通）
	ChineseNumbersSimple = []string{"零", "一", "二", "三", "四", "五", "六", "七", "八", "九"}
	
	// 中文单位（大写）
	ChineseUnitsUpper = []string{"", "拾", "佰", "仟", "万", "拾", "佰", "仟", "亿"}
	
	// 中文单位（普通）
	ChineseUnitsSimple = []string{"", "十", "百", "千", "万", "十", "百", "千", "亿"}
	
	// 星期名称
	WeekdayNames = [7]string{"周日", "周一", "周二", "周三", "周四", "周五", "周六"}
)

// 字符串驼峰转换
func ToCamelCase(s string) string {
	// 将字符串先转为小写并分割
	words := strings.Fields(strings.ToLower(s))
	if len(words) == 0 {
		return ""
	}

	// 第一个单词保持小写
	result := words[0]

	// 后续单词首字母大写
	for _, word := range words[1:] {
		if word == "" {
			continue
		}
		result += strings.Title(word)
	}

	return result
}

// 字符串转帕斯卡命名
func ToPascalCase(s string) string {
	words := strings.Fields(strings.ToLower(s))
	if len(words) == 0 {
		return ""
	}

	// 所有单词首字母大写
	for i, word := range words {
		if word == "" {
			continue
		}
		words[i] = strings.Title(word)
	}

	return strings.Join(words, "")
}

// 字符串转蛇形命名
func ToSnakeCase(s string) string {
	// 首先处理空格、连字符等
	s = strings.Join(strings.Fields(s), "_")
	
	// 处理驼峰命名的情况
	var result strings.Builder
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			// 前一个字符不是下划线，则添加下划线
			if s[i-1] != '_' {
				result.WriteRune('_')
			}
		}
		result.WriteRune(unicode.ToLower(r))
	}
	
	return result.String()
}

// 字符串转kebab命名
func ToKebabCase(s string) string {
	// 同样处理空格等
	s = strings.Join(strings.Fields(s), "-")
	
	// 处理驼峰命名
	var result strings.Builder
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			// 前一个字符不是连字符，则添加连字符
			if s[i-1] != '-' {
				result.WriteRune('-')
			}
		}
		result.WriteRune(unicode.ToLower(r))
	}
	
	return result.String()
}

// 是否为周末
func IsWeekend(date time.Time) bool {
	weekday := date.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

// 验证日期格式是否有效
func IsValidDateFormat(dateStr string) bool {
	pattern := `^\d{4}-\d{2}-\d{2}$`
	match, _ := regexp.MatchString(pattern, dateStr)
	return match
}

// 解析日期字符串
func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

// 解析日期时间字符串
func ParseDateTime(datetimeStr string) (time.Time, error) {
	// 尝试多种常见日期时间格式
	formats := []string{
		"2006-01-02 15:04:05",
		"2006/01/02 15:04:05",
		time.RFC3339,
		"02-01-2006 15:04:05",
		"02/01/2006 15:04:05",
	}
	
	var parseErr error
	for _, format := range formats {
		t, err := time.Parse(format, datetimeStr)
		if err == nil {
			return t, nil
		}
		parseErr = err
	}
	
	return time.Time{}, parseErr
}

// 获取月内周数
func GetWeekNumberInMonth(date time.Time, startWithMonday bool) int {
	// 获取当月第一天
	firstDay := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	
	// 计算第一周的偏移
	var firstWeekOffset int
	if startWithMonday {
		// 周一为每周第一天
		firstWeekday := int(firstDay.Weekday())
		if firstWeekday == 0 { // 周日
			firstWeekday = 7
		}
		firstWeekOffset = firstWeekday - 1
	} else {
		// 周日为每周第一天
		firstWeekOffset = int(firstDay.Weekday())
	}
	
	// 计算当前日期是当月第几周
	dayOfMonth := date.Day()
	return (dayOfMonth+firstWeekOffset-1)/7 + 1
}

// 获取年内周数
func GetWeekNumberInYear(date time.Time, startWithMonday bool) int {
	// 获取当年第一天
	firstDay := time.Date(date.Year(), 1, 1, 0, 0, 0, 0, date.Location())
	
	// ISO周数计算
	isoYear, isoWeek := date.ISOWeek()
	if isoYear == date.Year() && startWithMonday {
		return isoWeek
	}
	
	// 自定义周数计算（周日为周起始日）
	daysSinceFirstDay := int(date.Sub(firstDay).Hours() / 24)
	firstDayWeekday := int(firstDay.Weekday())
	if !startWithMonday { // 周日为第一天
		return (daysSinceFirstDay + firstDayWeekday) / 7 + 1
	} else { // 周一为第一天
		if firstDayWeekday == 0 { // 第一天是周日
			firstDayWeekday = 7
		}
		return (daysSinceFirstDay + (firstDayWeekday - 1)) / 7 + 1
	}
} 