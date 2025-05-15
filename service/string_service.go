package service

import (
	"github.com/mozillazg/go-pinyin"
	"github.com/renoz/toolbox-api/model"
	"github.com/renoz/toolbox-api/utils"
	"regexp"
	"strconv"
	"strings"
)

// 字符串分割
func SplitString(req model.SplitRequest) (*model.SplitResponse, error) {
	// 分割字符串
	parts := strings.Split(req.Input, req.Delimiter)
	
	// 构建响应
	response := &model.SplitResponse{
		Count: len(parts),
	}
	
	// 检查是否需要返回键值对映射
	if req.MapFormat && req.KeyValueDelimiter != "" {
		result := make(map[string]string)
		for _, part := range parts {
			kvParts := strings.SplitN(part, req.KeyValueDelimiter, 2)
			if len(kvParts) == 2 {
				result[kvParts[0]] = kvParts[1]
			}
		}
		response.Result = result
	} else {
		response.Result = parts
	}
	
	return response, nil
}

// 索引切分字符串
func SplitIndexedString(req model.SplitIndexedRequest) (*model.SplitResponse, error) {
	// 分割字符串
	parts := strings.Split(req.Content, req.Delimiter)
	
	// 构建索引映射
	result := make(map[string]string)
	for i, part := range parts {
		result[strconv.Itoa(i)] = part
	}
	
	return &model.SplitResponse{
		Result: result,
		Count:  len(parts),
	}, nil
}

// 字符串替换
func ReplaceString(req model.ReplaceRequest) (*model.ReplaceResponse, error) {
	var result string
	
	if req.UseRegex {
		// 使用正则表达式替换
		re, err := regexp.Compile(req.Target)
		if err != nil {
			return nil, err
		}
		result = re.ReplaceAllString(req.Input, req.Replacement)
	} else {
		// 普通替换
		result = strings.ReplaceAll(req.Input, req.Target, req.Replacement)
	}
	
	return &model.ReplaceResponse{
		Result: result,
		Length: len(result),
	}, nil
}

// 驼峰命名转换
func ToCamelCase(req model.CaseConversionRequest) (*model.CaseConversionResponse, error) {
	result := utils.ToCamelCase(req.Text)
	
	return &model.CaseConversionResponse{
		Result:   result,
		Original: req.Text,
	}, nil
}

// 帕斯卡命名转换
func ToPascalCase(req model.CaseConversionRequest) (*model.CaseConversionResponse, error) {
	result := utils.ToPascalCase(req.Text)
	
	return &model.CaseConversionResponse{
		Result:   result,
		Original: req.Text,
	}, nil
}

// 蛇形命名转换
func ToSnakeCase(req model.CaseConversionRequest) (*model.CaseConversionResponse, error) {
	result := utils.ToSnakeCase(req.Text)
	
	return &model.CaseConversionResponse{
		Result:   result,
		Original: req.Text,
	}, nil
}

// Kebab命名转换
func ToKebabCase(req model.CaseConversionRequest) (*model.CaseConversionResponse, error) {
	result := utils.ToKebabCase(req.Text)
	
	return &model.CaseConversionResponse{
		Result:   result,
		Original: req.Text,
	}, nil
}

// 提取中文拼音首字母
func ExtractInitials(req model.ExtractInitialsRequest) (*model.ExtractInitialsResponse, error) {
	// 配置拼音转换器
	args := pinyin.NewArgs()
	args.Fallback = func(r rune, a pinyin.Args) []string {
		return []string{string(r)}
	}
	
	// 转换为拼音数组
	pinyinArray := pinyin.Pinyin(req.Text, args)
	
	// 提取首字母
	var initials strings.Builder
	for _, py := range pinyinArray {
		if len(py) > 0 && len(py[0]) > 0 {
			initial := py[0][0:1]
			if req.Uppercase {
				initial = strings.ToUpper(initial)
			} else {
				initial = strings.ToLower(initial)
			}
			initials.WriteString(initial)
		}
	}
	
	return &model.ExtractInitialsResponse{
		OriginalText: req.Text,
		Initials:     initials.String(),
	}, nil
}

// 日期转换为中文大写格式
func ConvertDateToChinese(req model.ConvertDateRequest) (*model.ConvertDateResponse, error) {
	// 验证日期格式
	if !utils.IsValidDateFormat(req.DateStr) {
		return nil, &model.ErrorResponse{Code: 4001, Message: "日期格式错误，应为YYYY-MM-DD"}
	}
	
	// 解析日期
	date, err := utils.ParseDate(req.DateStr)
	if err != nil {
		return nil, err
	}
	
	// 获取年、月、日
	year := date.Year()
	month := int(date.Month())
	day := date.Day()
	
	// 转换年份 - 年份必须用中文四位数完整表述
	yearStr := strconv.Itoa(year)
	var yearChineseBuilder strings.Builder
	for _, digit := range yearStr {
		digitInt, _ := strconv.Atoi(string(digit))
		yearChineseBuilder.WriteString(utils.ChineseNumbersUpper[digitInt])
	}
	
	// 转换月份
	var monthChineseBuilder strings.Builder
	if month >= 1 && month <= 2 {
		// 1月、2月前需加"零"
		monthChineseBuilder.WriteString("零")
		monthChineseBuilder.WriteString(utils.ChineseNumbersUpper[month])
	} else if month >= 3 && month <= 9 {
		// 3月到9月前不加"零"
		monthChineseBuilder.WriteString(utils.ChineseNumbersUpper[month])
	} else if month == 10 {
		// 10月前需加"零"
		monthChineseBuilder.WriteString("零壹拾")
	} else if month == 11 || month == 12 {
		// 11月、12月前需加"壹"
		monthChineseBuilder.WriteString("壹拾")
		if month == 11 {
			monthChineseBuilder.WriteString(utils.ChineseNumbersUpper[1])
		} else {
			monthChineseBuilder.WriteString(utils.ChineseNumbersUpper[2])
		}
	}
	monthChineseBuilder.WriteString("月")
	
	// 转换日期
	var dayChineseBuilder strings.Builder
	if day >= 1 && day <= 9 {
		// 1日到9日前需加"零"
		dayChineseBuilder.WriteString("零")
		dayChineseBuilder.WriteString(utils.ChineseNumbersUpper[day])
	} else if day == 10 {
		// 10日前需加"零"
		dayChineseBuilder.WriteString("零壹拾")
	} else if day >= 11 && day <= 19 {
		// 11日到19日前需加"壹"
		dayChineseBuilder.WriteString("壹拾")
		dayChineseBuilder.WriteString(utils.ChineseNumbersUpper[day-10])
	} else if day == 20 {
		// 20日前需加"零"
		dayChineseBuilder.WriteString("零贰拾")
	} else if day >= 21 && day <= 29 {
		// 21日至29日前不加"零"
		dayChineseBuilder.WriteString("贰拾")
		dayChineseBuilder.WriteString(utils.ChineseNumbersUpper[day-20])
	} else if day == 30 {
		// 30日前需加"零"
		dayChineseBuilder.WriteString("零叁拾")
	} else if day == 31 {
		// 31日
		dayChineseBuilder.WriteString("叁拾壹")
	}
	dayChineseBuilder.WriteString("日")
	
	return &model.ConvertDateResponse{
		OriginalDate:  req.DateStr,
		ConvertedDate: yearChineseBuilder.String() + "年" + monthChineseBuilder.String() + dayChineseBuilder.String(),
	}, nil
}

// 日期转换为中文普通格式
func ConvertDateToChineseSimple(req model.ConvertDateSimpleRequest) (*model.ConvertDateResponse, error) {
	// 验证日期格式
	if !utils.IsValidDateFormat(req.DateStr) {
		return nil, &model.ErrorResponse{Code: 4001, Message: "日期格式错误，应为YYYY-MM-DD"}
	}
	
	// 解析日期
	date, err := utils.ParseDate(req.DateStr)
	if err != nil {
		return nil, err
	}
	
	// 获取年、月、日
	year := date.Year()
	month := int(date.Month())
	day := date.Day()
	
	// 转换年份
	yearStr := strconv.Itoa(year)
	var yearChineseBuilder strings.Builder
	for _, digit := range yearStr {
		digitInt, _ := strconv.Atoi(string(digit))
		yearChineseBuilder.WriteString(utils.ChineseNumbersSimple[digitInt])
	}
	
	// 转换月份
	var monthChineseBuilder strings.Builder
	if month < 10 {
		monthChineseBuilder.WriteString(utils.ChineseNumbersSimple[month])
	} else if month == 10 {
		monthChineseBuilder.WriteString("十")
	} else if month == 11 {
		monthChineseBuilder.WriteString("十一")
	} else if month == 12 {
		monthChineseBuilder.WriteString("十二")
	}
	monthChineseBuilder.WriteString("月")
	
	// 转换日期
	var dayChineseBuilder strings.Builder
	if day < 10 {
		dayChineseBuilder.WriteString(utils.ChineseNumbersSimple[day])
	} else if day < 20 {
		if day == 10 {
			dayChineseBuilder.WriteString("十")
		} else {
			dayChineseBuilder.WriteString("十")
			dayChineseBuilder.WriteString(utils.ChineseNumbersSimple[day-10])
		}
	} else {
		dayChineseBuilder.WriteString(utils.ChineseNumbersSimple[day/10])
		dayChineseBuilder.WriteString("十")
		if day%10 > 0 {
			dayChineseBuilder.WriteString(utils.ChineseNumbersSimple[day%10])
		}
	}
	dayChineseBuilder.WriteString("日")
	
	return &model.ConvertDateResponse{
		OriginalDate:  req.DateStr,
		ConvertedDate: yearChineseBuilder.String() + "年" + monthChineseBuilder.String() + dayChineseBuilder.String(),
	}, nil
} 