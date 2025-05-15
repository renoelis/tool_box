package service

import (
	"github.com/renoz/toolbox-api/model"
	"math/rand"
	"time"
)

// 初始化随机数种子
func init() {
	rand.Seed(time.Now().UnixNano())
}

// 生成随机整数
func GenerateRandomIntegers(min, max, count int, allowDuplicates bool) (*model.RandomIntegerResponse, error) {
	if min > max {
		return nil, &model.ErrorResponse{Code: 400, Message: "最小值不能大于最大值"}
	}
	
	if count < 1 {
		return nil, &model.ErrorResponse{Code: 400, Message: "生成数量必须大于0"}
	}
	
	// 如果不允许重复，检查可能值范围是否足够
	possibleValues := max - min + 1
	if !allowDuplicates && count > possibleValues {
		return nil, &model.ErrorResponse{Code: 400, Message: "不允许重复时，生成数量不能大于可能值范围"}
	}
	
	// 生成随机数
	var numbers []int
	
	if !allowDuplicates {
		// 不允许重复时使用Fisher-Yates洗牌算法
		// 创建有序数组
		pool := make([]int, possibleValues)
		for i := 0; i < possibleValues; i++ {
			pool[i] = min + i
		}
		
		// 随机选择元素
		for i := 0; i < count; i++ {
			j := rand.Intn(possibleValues - i)
			numbers = append(numbers, pool[j])
			// 将已选元素移到末尾
			pool[j], pool[possibleValues-i-1] = pool[possibleValues-i-1], pool[j]
		}
	} else {
		// 允许重复时直接生成随机数
		for i := 0; i < count; i++ {
			numbers = append(numbers, rand.Intn(max-min+1)+min)
		}
	}
	
	return &model.RandomIntegerResponse{
		Numbers: numbers,
		Count:   count,
		Min:     min,
		Max:     max,
	}, nil
}