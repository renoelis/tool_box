package model

// 通用响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// 字符串分割请求
type SplitRequest struct {
	Input             string `json:"input" binding:"required"`
	Delimiter         string `json:"delimiter" binding:"required"`
	MapFormat         bool   `json:"mapFormat"`
	KeyValueDelimiter string `json:"keyValueDelimiter"`
}

// 字符串分割响应
type SplitResponse struct {
	Result interface{} `json:"result"`
	Count  int         `json:"count"`
}

// 索引切分字符串请求
type SplitIndexedRequest struct {
	Content   string `json:"content" binding:"required"`
	Delimiter string `json:"delimiter" binding:"required"`
}

// 字符串替换请求
type ReplaceRequest struct {
	Input       string `json:"input" binding:"required"`
	Target      string `json:"target" binding:"required"`
	Replacement string `json:"replacement" binding:"required"`
	UseRegex    bool   `json:"use_regex"`
}

// 字符串替换响应
type ReplaceResponse struct {
	Result string `json:"result"`
	Length int    `json:"length"`
}

// 命名格式转换请求
type CaseConversionRequest struct {
	Text string `json:"text" binding:"required"`
}

// 命名格式转换响应
type CaseConversionResponse struct {
	Result   string `json:"result"`
	Original string `json:"original"`
}

// 中文拼音首字母提取请求
type ExtractInitialsRequest struct {
	Text      string `json:"text" binding:"required"`
	Uppercase bool   `json:"uppercase" default:"true"`
}

// 中文拼音首字母提取响应
type ExtractInitialsResponse struct {
	OriginalText string `json:"original_text"`
	Initials     string `json:"initials"`
}

// 日期转换为中文大写格式请求
type ConvertDateRequest struct {
	DateStr string `json:"date_str" binding:"required"`
}

// 日期转换为中文大写格式响应
type ConvertDateResponse struct {
	OriginalDate  string `json:"original_date"`
	ConvertedDate string `json:"converted_date"`
}

// 日期转换为中文普通格式请求
type ConvertDateSimpleRequest struct {
	DateStr string `json:"date_str" binding:"required"`
}

// 随机数生成响应
type RandomIntegerResponse struct {
	Numbers []int `json:"numbers"`
	Count   int   `json:"count"`
	Min     int   `json:"min"`
	Max     int   `json:"max"`
}

// 工作日计算请求
type WorkdayRangeRequest struct {
	StartDate      string   `json:"start_date" binding:"required"`
	EndDate        string   `json:"end_date" binding:"required"`
	RestDayPattern string   `json:"rest_day_pattern"`
	ExcludeDates   []string `json:"exclude_dates"`
	AddDates       []string `json:"add_dates"`
}

// 工作日计算响应
type WorkdayRangeResponse struct {
	TotalDays   int      `json:"total_days"`
	Workdays    int      `json:"workdays"`
	Restdays    int      `json:"restdays"`
	WorkdayList []string `json:"workday_list"`
	RestdayList []string `json:"restday_list"`
}

// 时区信息响应
type TimezoneInfo struct {
	Name          string `json:"name"`
	OffsetHours   int    `json:"offset_hours"`
	OffsetMinutes int    `json:"offset_minutes"`
	OffsetSeconds int    `json:"offset_seconds"`
	Description   string `json:"description"`
}

// 时间转换请求
type TimeConvertRequest struct {
	TimeInput     interface{} `json:"time_input" binding:"required"`
	OutputFormat  string      `json:"output_format"`
	Timezone      string      `json:"timezone"`
	TzOffset      *int        `json:"tz_offset"`
	CustomFormat  string      `json:"custom_format"`
}

// 时间转换响应
type TimeConvertResponse struct {
	Original     string      `json:"original"`
	Converted    string      `json:"converted"`
	Timezone     string      `json:"timezone"`
	TimezoneInfo TimezoneInfo `json:"timezone_info"`
}

// 当前时间响应
type CurrentTimeResponse struct {
	CurrentTime  string      `json:"current_time"`
	Timezone     string      `json:"timezone"`
	TimezoneInfo TimezoneInfo `json:"timezone_info"`
	Format       string      `json:"format"`
	Year         int         `json:"year"`
	Month        int         `json:"month"`
	Day          int         `json:"day"`
	Hour         int         `json:"hour"`
	Minute       int         `json:"minute"`
	Second       int         `json:"second"`
	Timestamp    int64       `json:"timestamp"`
	Weekday      int         `json:"weekday"`
	CustomFormat string      `json:"custom_format,omitempty"`
}

// 周末检查响应
type IsWeekendResponse struct {
	Date        string `json:"date"`
	IsWeekend   bool   `json:"is_weekend"`
	Weekday     int    `json:"weekday"`
	WeekdayName string `json:"weekday_name"`
}

// 周数信息响应
type WeekNumberResponse struct {
	Date        string `json:"date"`
	WeekNumber  int    `json:"week_number"`
	Weekday     int    `json:"weekday"`
	WeekdayName string `json:"weekday_name"`
	InMonth     bool   `json:"in_month"`
	StartDay    string `json:"start_day"`
	TotalWeeks  int    `json:"total_weeks"`
}

// 时区信息完整响应
type TimezoneInfoResponse struct {
	Timezone         string       `json:"timezone"`
	TimezoneInfo     TimezoneInfo `json:"timezone_info"`
	CurrentTime      string       `json:"current_time"`
	AvailableOffsets []int        `json:"available_offsets"`
}