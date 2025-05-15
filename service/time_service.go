package service

import (
	"github.com/renoz/toolbox-api/model"
	"github.com/renoz/toolbox-api/utils"
	"fmt"
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
	// 如果是空格式，返回默认格式
	if pythonFormat == "" {
		return "2006-01-02 15:04:05"
	}
	
	// 处理直接对应的Python到Go格式映射
	result := pythonFormat
	
	// 检查特殊情况，整个格式是否为特定格式
	if pythonFormat == "%c" {
		return "Mon Jan 2 15:04:05 2006" // ANSIC 格式
	} else if pythonFormat == "%x" {
		return "01/02/06" // 日期格式
	} else if pythonFormat == "%X" {
		return "15:04:05" // 时间格式
	}
	
	// 包含特殊处理的格式
	// 这些指令需要特殊处理，不能直接映射到Go时间格式
	var specialDirectives = []string{"%U", "%W", "%w", "%j", "%V", "%u"}
	for _, directive := range specialDirectives {
		if strings.Contains(pythonFormat, directive) {
			// 这些特殊格式需要自定义处理，不使用直接替换
			return pythonFormat
		}
	}
	
	// 进行常规替换
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
		"%%": "%",        // 百分号
		// 不添加 %V、%u、%w、%j、%U、%W 因为它们需要特殊处理
	}
	
	// 使用替换处理所有格式代码
	for pyFmt, goFmt := range formatMap {
		result = strings.Replace(result, pyFmt, goFmt, -1)
	}
	
	return result
}

// 格式化指定日期时间（处理包含特殊Python格式的情况）
func formatWithPythonFormat(t time.Time, pythonFormat string) string {
	// 获取1月1日，用于计算周数
	jan1 := time.Date(t.Year(), time.January, 1, 0, 0, 0, 0, t.Location())
	jan1Weekday := int(jan1.Weekday()) // Sunday=0
	
	// 特殊格式直接处理
	if pythonFormat == "%c" {
		return t.Format("Mon Jan 2 15:04:05 2006")
	} else if pythonFormat == "%x" {
		return t.Format("01/02/06")
	} else if pythonFormat == "%X" {
		return t.Format("15:04:05")
	}
	
	// 识别格式字符串中的所有Python指令
	var directives = []string{
		"%Y", "%y", "%m", "%B", "%b", "%d", "%j", "%U", "%W", "%V", 
		"%A", "%a", "%w", "%u", "%H", "%I", "%p", "%M", "%S", "%f",
		"%z", "%Z", "%c", "%x", "%X", "%%",
	}
	
	// 特殊计算函数：计算周数
	weekNumber := func(doy, jan1Weekday int, firstDay int) int {
		// firstDay: 0=Sunday, 1=Monday
		offset := (7 - ((jan1Weekday - firstDay + 7) % 7)) % 7
		if doy-1 < offset {
			return 0
		}
		return (doy-1-offset)/7 + 1
	}
	
	// 检查是否需要特殊处理
	needsSpecialHandling := false
	for _, directive := range []string{"%U", "%W", "%w", "%j", "%V", "%u"} {
		if strings.Contains(pythonFormat, directive) {
			needsSpecialHandling = true
			break
		}
	}
	
	// 如果不需要特殊处理，使用常规转换
	if !needsSpecialHandling {
		goFormat := convertPythonFormatToGo(pythonFormat)
		return t.Format(goFormat)
	}
	
	// 需要特殊处理，逐个替换指令
	result := pythonFormat
	
	// 依次处理每个指令
	for _, directive := range directives {
		// 如果格式中不包含该指令，跳过
		if !strings.Contains(result, directive) {
			continue
		}
		
		var replacement string
		
		switch directive {
		case "%Y":
			replacement = fmt.Sprintf("%04d", t.Year())
		case "%y":
			replacement = fmt.Sprintf("%02d", t.Year()%100)
		case "%m":
			replacement = fmt.Sprintf("%02d", int(t.Month()))
		case "%B":
			replacement = t.Month().String()
		case "%b":
			replacement = t.Month().String()[:3]
		case "%d":
			replacement = fmt.Sprintf("%02d", t.Day())
		case "%j":
			replacement = fmt.Sprintf("%03d", t.YearDay())
		case "%U": // 以周日为一周的第一天
			week := weekNumber(t.YearDay(), jan1Weekday, 0)
			replacement = fmt.Sprintf("%02d", week)
		case "%W": // 以周一为一周的第一天
			week := weekNumber(t.YearDay(), jan1Weekday, 1)
			replacement = fmt.Sprintf("%02d", week)
		case "%V": // ISO标准周数
			_, isoWeek := t.ISOWeek()
			replacement = fmt.Sprintf("%02d", isoWeek)
		case "%A":
			replacement = t.Weekday().String()
		case "%a":
			replacement = t.Weekday().String()[:3]
		case "%w":
			replacement = fmt.Sprintf("%d", int(t.Weekday()))
		case "%u": // ISO周几（1-7，周一到周日）
			weekday := int(t.Weekday())
			if weekday == 0 { // 周日
				weekday = 7
			}
			replacement = fmt.Sprintf("%d", weekday)
		case "%H":
			replacement = fmt.Sprintf("%02d", t.Hour())
		case "%I":
			hour12 := t.Hour() % 12
			if hour12 == 0 {
				hour12 = 12
			}
			replacement = fmt.Sprintf("%02d", hour12)
		case "%p":
			if t.Hour() < 12 {
				replacement = "AM"
			} else {
				replacement = "PM"
			}
		case "%M":
			replacement = fmt.Sprintf("%02d", t.Minute())
		case "%S":
			replacement = fmt.Sprintf("%02d", t.Second())
		case "%f":
			micro := t.Nanosecond() / 1000
			replacement = fmt.Sprintf("%06d", micro)
		case "%z":
			replacement = t.Format("-0700")
		case "%Z":
			replacement = t.Format("MST")
		case "%c":
			replacement = t.Format(time.ANSIC)
		case "%x":
			replacement = t.Format("01/02/06")
		case "%X":
			replacement = t.Format("15:04:05")
		case "%%":
			replacement = "%"
		}
		
		// 替换指令
		result = strings.Replace(result, directive, replacement, -1)
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
			// 使用改进的格式处理函数
			return formatWithPythonFormat(t, customFormat)
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
	
	// 格式化输出
	var formattedTime string
	
	if format == "custom" && customFormat != "" {
		// 使用新的格式处理函数
		formattedTime = formatWithPythonFormat(now, customFormat)
	} else {
		// 非自定义格式
		formattedTime = formatTime(now, format, customFormat)
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
		CustomFormat: rawCustomFormat, 
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
	
	// 格式化输出
	var formattedTime string
	
	if req.OutputFormat == "custom" && req.CustomFormat != "" {
		// 使用新的格式处理函数
		formattedTime = formatWithPythonFormat(inputTime, req.CustomFormat)
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