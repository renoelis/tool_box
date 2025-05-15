package service

import (
	"github.com/renoz/toolbox-api/model"
	"github.com/renoz/toolbox-api/utils"
	"strconv"
	"strings"
	"time"
)

// 支持的时区
var SUPPORTED_TIMEZONES = map[string]string{
	"asia_shanghai":   "Asia/Shanghai",   // 东8区
	"america_new_york": "America/New_York", // 西5区
	"europe_london":    "Europe/London",    // 0区
	"asia_tokyo":      "Asia/Tokyo",       // 东9区
	"australia_sydney": "Australia/Sydney", // 东10区
	"europe_paris":     "Europe/Paris",     // 东1区
	"america_los_angeles": "America/Los_Angeles", // 西8区
}

// 时区描述映射
var TIMEZONE_DESCRIPTIONS = map[string]string{
	"Asia/Shanghai":   "东8区",
	"America/New_York": "西5区",
	"Europe/London":    "0区",
	"Asia/Tokyo":       "东9区",
	"Australia/Sydney": "东10区",
	"Europe/Paris":     "东1区",
	"America/Los_Angeles": "西8区",
}

// 获取时区信息
func getTimezoneInfo(timezoneName string, tzOffset *int) (*model.TimezoneInfo, error) {
	var loc *time.Location
	var err error
	var offset int
	
	// 优先使用偏移量
	if tzOffset != nil {
		// 使用指定的小时偏移量
		offset = *tzOffset * 3600
		loc = time.FixedZone("FixedZone", offset)
	} else if timezoneName != "" {
		// 查找映射的时区名称
		if tz, ok := SUPPORTED_TIMEZONES[timezoneName]; ok {
			timezoneName = tz
		}
		
		// 尝试加载时区
		loc, err = time.LoadLocation(timezoneName)
		if err != nil {
			return nil, err
		}
		
		// 计算与UTC的偏移量
		now := time.Now().In(loc)
		_, offset = now.Zone()
	} else {
		// 默认使用亚洲/上海
		loc, _ = time.LoadLocation("Asia/Shanghai")
		now := time.Now().In(loc)
		_, offset = now.Zone()
		timezoneName = "Asia/Shanghai"
	}
	
	// 计算小时和分钟偏移
	offsetHours := offset / 3600
	offsetMinutes := (offset % 3600) / 60
	
	// 获取描述
	description := TIMEZONE_DESCRIPTIONS[timezoneName]
	if description == "" {
		if offsetHours > 0 {
			description = "东" + strconv.Itoa(offsetHours) + "区"
		} else if offsetHours < 0 {
			description = "西" + strconv.Itoa(-offsetHours) + "区"
		} else {
			description = "0区"
		}
	}
	
	return &model.TimezoneInfo{
		Name:          timezoneName,
		OffsetHours:   offsetHours,
		OffsetMinutes: offsetMinutes,
		OffsetSeconds: offset,
		Description:   description,
	}, nil
}

// 将Python风格的时间格式转换为Go风格的格式
func convertPythonFormatToGo(pythonFormat string) string {
	// Python格式到Go格式的映射表
	formatMap := map[string]string{
		"%Y": "2006",     // 年份（4位数字）
		"%y": "06",       // 年份（2位数字）
		"%m": "01",       // 月份（01-12）
		"%d": "02",       // 日期（01-31）
		"%H": "15",       // 小时（24小时制，00-23）
		"%I": "03",       // 小时（12小时制，01-12）
		"%M": "04",       // 分钟（00-59）
		"%S": "05",       // 秒（00-59）
		"%f": "000000",   // 微秒
		"%a": "Mon",      // 星期几缩写
		"%A": "Monday",   // 星期几全称
		"%b": "Jan",      // 月份缩写
		"%B": "January",  // 月份全称
		"%p": "PM",       // AM/PM
		"%z": "-0700",    // UTC偏移量
		"%Z": "MST",      // 时区名称
		"%j": "002",      // 一年中的第几天（001-366）
		"%U": "01",       // 一年中的第几周（00-53），以周日为一周的开始
		"%W": "01",       // 一年中的第几周（00-53），以周一为一周的开始
		"%%": "%",        // 百分号
		// %w 留给外部处理，这里暂不转换
	}
	
	// 特殊格式的完全匹配
	specialFormats := map[string]string{
		// %c 不再硬编码，留给外部处理
		"%x": "2006年01月02日",              // 本地日期表示
		"%X": "15:04:05",                 // 本地时间表示
	}
	
	// 先检查是否是完全匹配的特殊格式
	if specialFormat, ok := specialFormats[pythonFormat]; ok {
		return specialFormat
	}
	
	// 使用替换处理所有格式代码
	result := pythonFormat
	for pyFmt, goFmt := range formatMap {
		result = strings.Replace(result, pyFmt, goFmt, -1)
	}
	
	return result
}

// 格式化时间
func formatTime(t time.Time, format string, customFormat string) string {
	switch format {
	case "date_only":
		return t.Format("2006-01-02")
	case "time_only":
		return t.Format("15:04:05")
	case "datetime":
		return t.Format("2006-01-02 15:04:05")
	case "iso":
		return t.Format(time.RFC3339)
	case "rfc3339":
		return t.Format(time.RFC3339)
	case "human":
		return t.Format("2006年01月02日 15时04分05秒")
	case "timestamp":
		return strconv.FormatInt(t.Unix(), 10)
	case "timestamp_ms":
		return strconv.FormatInt(t.UnixNano()/1e6, 10)
	case "custom":
		if customFormat != "" {
			// 普通情况直接转换
			goFormat := convertPythonFormatToGo(customFormat)
			return t.Format(goFormat)
		}
		return t.Format("2006-01-02 15:04:05")
	default:
		return t.Format("2006-01-02 15:04:05") // 默认为datetime格式
	}
}

// 获取当前时间
func GetCurrentTime(format, timezone string, tzOffset *int, customFormat string) (*model.CurrentTimeResponse, error) {
	// 获取时区信息
	timezoneInfo, err := getTimezoneInfo(timezone, tzOffset)
	if err != nil {
		return nil, &model.ErrorResponse{Code: 9000, Message: "获取时区信息失败: " + err.Error()}
	}
	
	// 加载时区
	var loc *time.Location
	if tzOffset != nil {
		loc = time.FixedZone("FixedZone", *tzOffset*3600)
	} else {
		loc, err = time.LoadLocation(timezoneInfo.Name)
		if err != nil {
			return nil, &model.ErrorResponse{Code: 9000, Message: "加载时区失败: " + err.Error()}
		}
	}
	
	// 获取当前时间
	now := time.Now().In(loc)
	
	// 格式化输出
	if format == "" {
		format = "datetime"
	}
	
	// 记录原始自定义格式（用于返回）
	rawCustomFormat := customFormat
	
	// 特殊处理各种格式
	var formattedTime string
	
	if format == "custom" && customFormat != "" {
		// 处理特殊格式%c
		if customFormat == "%c" {
			// 标准C语言格式: "Sun May 18 14:26:36 2025"
			formattedTime = now.Format("Mon Jan 2 15:04:05 2006")
		} else if customFormat == "%U" {
			// %U: 一年中的第几周，以周日为一周的开始（00-53）
			weekNumber := utils.GetWeekNumberInYear(now, false) // false表示以周日为开始
			
			// 确保是两位数格式
			weekNumberStr := strconv.Itoa(weekNumber)
			if weekNumber < 10 {
				weekNumberStr = "0" + weekNumberStr
			}
			
			formattedTime = weekNumberStr
		} else if strings.Contains(customFormat, "%U") {
			// 包含%U的复杂格式
			weekNumber := utils.GetWeekNumberInYear(now, false) // false表示以周日为开始
			
			// 确保是两位数格式
			weekNumberStr := strconv.Itoa(weekNumber)
			if weekNumber < 10 {
				weekNumberStr = "0" + weekNumberStr
			}
			
			// 使用特殊占位符替换%U
			const placeholder = "WEEKNUMBER_PLACEHOLDER_U"
			tempFormat := strings.Replace(customFormat, "%U", placeholder, -1)
			
			// 转换其他Python格式为Go格式
			goFormat := convertPythonFormatToGo(tempFormat)
			
			// 使用Go的时间格式化
			tempResult := now.Format(goFormat)
			
			// 最后将占位符替换为实际的周数
			formattedTime = strings.Replace(tempResult, placeholder, weekNumberStr, -1)
		} else if strings.Contains(customFormat, "%w") {
			// 使用手动替换方法处理包含%w的格式
			// 首先替换所有格式标记
			result := customFormat
			
			// 处理年份
			if strings.Contains(result, "%Y") {
				result = strings.Replace(result, "%Y", strconv.Itoa(now.Year()), -1)
			}
			
			// 处理月份
			if strings.Contains(result, "%m") {
				month := int(now.Month())
				monthStr := strconv.Itoa(month)
				if month < 10 {
					monthStr = "0" + monthStr
				}
				result = strings.Replace(result, "%m", monthStr, -1)
			}
			
			// 处理日期
			if strings.Contains(result, "%d") {
				day := now.Day()
				dayStr := strconv.Itoa(day)
				if day < 10 {
					dayStr = "0" + dayStr
				}
				result = strings.Replace(result, "%d", dayStr, -1)
			}
			
			// 处理小时（24小时制）
			if strings.Contains(result, "%H") {
				hour := now.Hour()
				hourStr := strconv.Itoa(hour)
				if hour < 10 {
					hourStr = "0" + hourStr
				}
				result = strings.Replace(result, "%H", hourStr, -1)
			}
			
			// 处理分钟
			if strings.Contains(result, "%M") {
				minute := now.Minute()
				minuteStr := strconv.Itoa(minute)
				if minute < 10 {
					minuteStr = "0" + minuteStr
				}
				result = strings.Replace(result, "%M", minuteStr, -1)
			}
			
			// 处理秒
			if strings.Contains(result, "%S") {
				second := now.Second()
				secondStr := strconv.Itoa(second)
				if second < 10 {
					secondStr = "0" + secondStr
				}
				result = strings.Replace(result, "%S", secondStr, -1)
			}
			
			// 最后处理星期几
			weekday := int(now.Weekday())
			weekdayStr := strconv.Itoa(weekday)
			result = strings.Replace(result, "%w", weekdayStr, -1)
			
			formattedTime = result
		} else {
			// 不包含%w的自定义格式
			goFormat := convertPythonFormatToGo(customFormat)
			formattedTime = now.Format(goFormat)
		}
	} else {
		// 非自定义格式
		formattedTime = formatTime(now, format, customFormat)
	}
	
	// 如果是自定义格式，记录实际使用的格式
	if format == "custom" && customFormat != "" {
		customFormat = rawCustomFormat
	}
	
	// 添加额外的响应字段
	response := &model.CurrentTimeResponse{
		CurrentTime:  formattedTime,
		Timezone:     timezoneInfo.Name,
		TimezoneInfo: *timezoneInfo,
		Format:       format,
		Year:         now.Year(),
		Month:        int(now.Month()),
		Day:          now.Day(),
		Hour:         now.Hour(),
		Minute:       now.Minute(),
		Second:       now.Second(),
		Timestamp:    now.Unix(),
		Weekday:      int(now.Weekday()),  // 正确的星期几 (0-6)
		CustomFormat: customFormat, 
	}
	
	return response, nil
}

// 检查是否为周末
func CheckIsWeekend(dateStr string) (*model.IsWeekendResponse, error) {
	var date time.Time
	var err error
	
	if dateStr == "" {
		// 使用当前日期
		date = time.Now()
	} else {
		// 验证并解析日期
		if !utils.IsValidDateFormat(dateStr) {
			return nil, &model.ErrorResponse{Code: 3001, Message: "日期格式错误，应为YYYY-MM-DD"}
		}
		
		date, err = utils.ParseDate(dateStr)
		if err != nil {
			return nil, &model.ErrorResponse{Code: 3001, Message: "日期解析失败: " + err.Error()}
		}
	}
	
	// 判断是否为周末
	weekday := date.Weekday()
	isWeekend := utils.IsWeekend(date)
	
	return &model.IsWeekendResponse{
		Date:        date.Format("2006-01-02"),
		IsWeekend:   isWeekend,
		Weekday:     int(weekday),
		WeekdayName: utils.WeekdayNames[weekday],
	}, nil
}

// 获取周数信息
func GetWeekNumber(dateStr string, inMonth bool, startWithMonday bool) (*model.WeekNumberResponse, error) {
	var date time.Time
	var err error
	
	if dateStr == "" {
		// 使用当前日期
		date = time.Now()
	} else {
		// 验证并解析日期
		if !utils.IsValidDateFormat(dateStr) {
			return nil, &model.ErrorResponse{Code: 3002, Message: "日期格式错误，应为YYYY-MM-DD"}
		}
		
		date, err = utils.ParseDate(dateStr)
		if err != nil {
			return nil, &model.ErrorResponse{Code: 3002, Message: "日期解析失败: " + err.Error()}
		}
	}
	
	// 计算周数
	var weekNumber int
	if inMonth {
		weekNumber = utils.GetWeekNumberInMonth(date, startWithMonday)
	} else {
		weekNumber = utils.GetWeekNumberInYear(date, startWithMonday)
	}
	
	// 获取星期几
	weekday := date.Weekday()
	
	// 设置周起始日
	var startDay string
	if startWithMonday {
		startDay = "周一"
	} else {
		startDay = "周日"
	}
	
	// 计算总周数
	var totalWeeks int
	if inMonth {
		// 获取月份的总天数
		year, month, _ := date.Date()
		lastDay := time.Date(year, month+1, 0, 0, 0, 0, 0, date.Location())
		daysInMonth := lastDay.Day()
		
		// 获取月份第一天和最后一天的weekday
		firstDay := time.Date(year, month, 1, 0, 0, 0, 0, date.Location())
		
		// 计算月份总周数
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
		
		// 计算总周数
		totalWeeks = (daysInMonth + firstWeekOffset - 1) / 7 + 1
	} else {
		// 年份总周数
		if startWithMonday {
			// ISO周计算法
			// 获取当年12月31日所在的周
			lastDay := time.Date(date.Year(), 12, 31, 0, 0, 0, 0, date.Location())
			_, lastWeek := lastDay.ISOWeek()
			
			// 处理特殊情况：年末几天可能属于下一年的第一周
			if lastWeek == 1 {
				// 检查12月30日
				prevDay := lastDay.AddDate(0, 0, -1)
				_, prevWeek := prevDay.ISOWeek()
				totalWeeks = prevWeek
			} else {
				totalWeeks = lastWeek
			}
		} else {
			// 基于周日的计算法
			firstDay := time.Date(date.Year(), 1, 1, 0, 0, 0, 0, date.Location())
			
			// 计算当年第一天的偏移
			firstDayOffset := int(firstDay.Weekday())
			
			// 获取当年总天数
			daysInYear := 365
			if date.Year()%4 == 0 && (date.Year()%100 != 0 || date.Year()%400 == 0) {
				daysInYear = 366 // 闰年
			}
			
			// 计算总周数
			totalWeeks = (daysInYear + firstDayOffset) / 7
			if (daysInYear + firstDayOffset) % 7 > 0 {
				totalWeeks++
			}
		}
	}
	
	return &model.WeekNumberResponse{
		Date:        date.Format("2006-01-02"),
		WeekNumber:  weekNumber,
		Weekday:     int(weekday),
		WeekdayName: utils.WeekdayNames[weekday],
		InMonth:     inMonth,
		StartDay:    startDay,
		TotalWeeks:  totalWeeks,
	}, nil
}

// 获取时区信息
func GetTimezoneInfo(timezone string, tzOffset *int) (*model.TimezoneInfoResponse, error) {
	// 获取时区信息
	timezoneInfo, err := getTimezoneInfo(timezone, tzOffset)
	if err != nil {
		return nil, &model.ErrorResponse{Code: 9000, Message: "获取时区信息失败: " + err.Error()}
	}
	
	// 获取当前时间
	var loc *time.Location
	if tzOffset != nil {
		loc = time.FixedZone("FixedZone", *tzOffset*3600)
	} else {
		loc, err = time.LoadLocation(timezoneInfo.Name)
		if err != nil {
			return nil, &model.ErrorResponse{Code: 9000, Message: "加载时区失败: " + err.Error()}
		}
	}
	
	now := time.Now().In(loc)
	
	// 可用的时区偏移
	availableOffsets := []int{-12, -11, -10, -9, -8, -7, -6, -5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}
	
	return &model.TimezoneInfoResponse{
		Timezone:         timezoneInfo.Name,
		TimezoneInfo:     *timezoneInfo,
		CurrentTime:      now.Format("2006-01-02 15:04:05"),
		AvailableOffsets: availableOffsets,
	}, nil
}

// 时间格式转换
func ConvertTime(req model.TimeConvertRequest) (*model.TimeConvertResponse, error) {
	var inputTime time.Time
	var err error
	var originalStr string
	
	// 解析输入时间
	switch v := req.TimeInput.(type) {
	case string:
		// 字符串输入，尝试解析
		originalStr = v
		inputTime, err = utils.ParseDateTime(v)
		if err != nil {
			return nil, &model.ErrorResponse{Code: 9000, Message: "时间格式解析失败: " + err.Error()}
		}
	case float64:
		// 数字输入，当作时间戳处理
		originalStr = strconv.FormatInt(int64(v), 10)
		if v > 9999999999 {
			// 毫秒时间戳
			inputTime = time.Unix(0, int64(v)*1000000)
		} else {
			// 秒时间戳
			inputTime = time.Unix(int64(v), 0)
		}
	case int:
		originalStr = strconv.Itoa(v)
		if v > 9999999999 {
			inputTime = time.Unix(0, int64(v)*1000000)
		} else {
			inputTime = time.Unix(int64(v), 0)
		}
	case int64:
		originalStr = strconv.FormatInt(v, 10)
		if v > 9999999999 {
			inputTime = time.Unix(0, v*1000000)
		} else {
			inputTime = time.Unix(v, 0)
		}
	default:
		return nil, &model.ErrorResponse{Code: 9000, Message: "不支持的时间输入格式"}
	}
	
	// 设置默认输出格式
	if req.OutputFormat == "" {
		req.OutputFormat = "default"
	}
	
	// 获取时区信息
	timezoneInfo, err := getTimezoneInfo(req.Timezone, req.TzOffset)
	if err != nil {
		return nil, &model.ErrorResponse{Code: 9000, Message: "获取时区信息失败: " + err.Error()}
	}
	
	// 加载时区
	var loc *time.Location
	if req.TzOffset != nil {
		loc = time.FixedZone("FixedZone", *req.TzOffset*3600)
	} else {
		loc, err = time.LoadLocation(timezoneInfo.Name)
		if err != nil {
			return nil, &model.ErrorResponse{Code: 9000, Message: "加载时区失败: " + err.Error()}
		}
	}
	
	// 转换时区
	inputTime = inputTime.In(loc)
	
	// 格式化输出，与GetCurrentTime保持一致的处理逻辑
	var formattedTime string
	
	if req.OutputFormat == "custom" && req.CustomFormat != "" {
		// 处理特殊格式%c
		if req.CustomFormat == "%c" {
			// 标准C语言格式: "Sun May 18 14:26:36 2025"
			formattedTime = inputTime.Format("Mon Jan 2 15:04:05 2006")
		} else if req.CustomFormat == "%U" {
			// %U: 一年中的第几周，以周日为一周的开始（00-53）
			weekNumber := utils.GetWeekNumberInYear(inputTime, false) // false表示以周日为开始
			
			// 确保是两位数格式
			weekNumberStr := strconv.Itoa(weekNumber)
			if weekNumber < 10 {
				weekNumberStr = "0" + weekNumberStr
			}
			
			formattedTime = weekNumberStr
		} else if strings.Contains(req.CustomFormat, "%U") {
			// 包含%U的复杂格式
			weekNumber := utils.GetWeekNumberInYear(inputTime, false) // false表示以周日为开始
			
			// 确保是两位数格式
			weekNumberStr := strconv.Itoa(weekNumber)
			if weekNumber < 10 {
				weekNumberStr = "0" + weekNumberStr
			}
			
			// 使用特殊占位符替换%U
			const placeholder = "WEEKNUMBER_PLACEHOLDER_U"
			tempFormat := strings.Replace(req.CustomFormat, "%U", placeholder, -1)
			
			// 转换其他Python格式为Go格式
			goFormat := convertPythonFormatToGo(tempFormat)
			
			// 使用Go的时间格式化
			tempResult := inputTime.Format(goFormat)
			
			// 最后将占位符替换为实际的周数
			formattedTime = strings.Replace(tempResult, placeholder, weekNumberStr, -1)
		} else if strings.Contains(req.CustomFormat, "%w") {
			// 使用手动替换方法处理包含%w的格式
			// 首先替换所有格式标记
			result := req.CustomFormat
			
			// 处理年份
			if strings.Contains(result, "%Y") {
				result = strings.Replace(result, "%Y", strconv.Itoa(inputTime.Year()), -1)
			}
			
			// 处理月份
			if strings.Contains(result, "%m") {
				month := int(inputTime.Month())
				monthStr := strconv.Itoa(month)
				if month < 10 {
					monthStr = "0" + monthStr
				}
				result = strings.Replace(result, "%m", monthStr, -1)
			}
			
			// 处理日期
			if strings.Contains(result, "%d") {
				day := inputTime.Day()
				dayStr := strconv.Itoa(day)
				if day < 10 {
					dayStr = "0" + dayStr
				}
				result = strings.Replace(result, "%d", dayStr, -1)
			}
			
			// 处理小时（24小时制）
			if strings.Contains(result, "%H") {
				hour := inputTime.Hour()
				hourStr := strconv.Itoa(hour)
				if hour < 10 {
					hourStr = "0" + hourStr
				}
				result = strings.Replace(result, "%H", hourStr, -1)
			}
			
			// 处理分钟
			if strings.Contains(result, "%M") {
				minute := inputTime.Minute()
				minuteStr := strconv.Itoa(minute)
				if minute < 10 {
					minuteStr = "0" + minuteStr
				}
				result = strings.Replace(result, "%M", minuteStr, -1)
			}
			
			// 处理秒
			if strings.Contains(result, "%S") {
				second := inputTime.Second()
				secondStr := strconv.Itoa(second)
				if second < 10 {
					secondStr = "0" + secondStr
				}
				result = strings.Replace(result, "%S", secondStr, -1)
			}
			
			// 最后处理星期几
			weekday := int(inputTime.Weekday())
			weekdayStr := strconv.Itoa(weekday)
			result = strings.Replace(result, "%w", weekdayStr, -1)
			
			formattedTime = result
		} else {
			// 不包含%w的自定义格式
			goFormat := convertPythonFormatToGo(req.CustomFormat)
			formattedTime = inputTime.Format(goFormat)
		}
	} else {
		// 非自定义格式
		formattedTime = formatTime(inputTime, req.OutputFormat, req.CustomFormat)
	}
	
	return &model.TimeConvertResponse{
		Original:     originalStr,
		Converted:    formattedTime,
		Timezone:     timezoneInfo.Name,
		TimezoneInfo: *timezoneInfo,
	}, nil
}

// 工作日计算
func CalculateWorkdays(req model.WorkdayRangeRequest) (*model.WorkdayRangeResponse, error) {
	// 默认休息日模式: 周六日休息
	if req.RestDayPattern == "" {
		req.RestDayPattern = "0000011" // 周一到周五工作，周六日休息
	}
	
	// 验证日期
	startDate, err := utils.ParseDate(req.StartDate)
	if err != nil {
		return nil, &model.ErrorResponse{Code: 9000, Message: "开始日期格式错误: " + err.Error()}
	}
	
	endDate, err := utils.ParseDate(req.EndDate)
	if err != nil {
		return nil, &model.ErrorResponse{Code: 9000, Message: "结束日期格式错误: " + err.Error()}
	}
	
	if endDate.Before(startDate) {
		return nil, &model.ErrorResponse{Code: 9000, Message: "结束日期不能早于开始日期"}
	}
	
	// 创建排除日期集合
	excludeDates := make(map[string]bool)
	for _, dateStr := range req.ExcludeDates {
		excludeDates[dateStr] = true
	}
	
	// 创建额外工作日集合
	addDates := make(map[string]bool)
	for _, dateStr := range req.AddDates {
		addDates[dateStr] = true
	}
	
	// 计算总天数
	totalDays := int(endDate.Sub(startDate).Hours()/24) + 1
	
	// 遍历日期范围，统计工作日
	var workdayList []string
	var restdayList []string
	
	current := startDate
	for i := 0; i < totalDays; i++ {
		dateStr := current.Format("2006-01-02")
		
		// 判断是否为工作日
		weekday := int(current.Weekday())
		if weekday == 0 {
			weekday = 7 // 调整周日为7
		}
		weekday-- // 调整为0-6的索引
		
		// 检查是否为休息日
		isRestDay := false
		if req.RestDayPattern[weekday] == '1' {
			isRestDay = true
		}
		
		// 应用自定义规则
		if excludeDates[dateStr] {
			// 排除日期，标记为休息日
			isRestDay = true
		}
		
		if addDates[dateStr] {
			// 添加日期，标记为工作日
			isRestDay = false
		}
		
		// 添加到相应列表
		if isRestDay {
			restdayList = append(restdayList, dateStr)
		} else {
			workdayList = append(workdayList, dateStr)
		}
		
		// 移动到下一天
		current = current.AddDate(0, 0, 1)
	}
	
	return &model.WorkdayRangeResponse{
		TotalDays:   totalDays,
		Workdays:    len(workdayList),
		Restdays:    len(restdayList),
		WorkdayList: workdayList,
		RestdayList: restdayList,
	}, nil
} 