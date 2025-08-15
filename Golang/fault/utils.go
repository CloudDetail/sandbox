package fault

import "fmt"

// ParseStringToInt 解析字符串为整数
func ParseStringToInt(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}

// ParseStringToFloat 解析字符串为浮点数
func ParseStringToFloat(s string) (float64, error) {
	var result float64
	_, err := fmt.Sscanf(s, "%f", &result)
	return result, err
}

// ParseStringToBool 解析字符串为布尔值
func ParseStringToBool(s string) (bool, error) {
	switch s {
	case "true", "1", "on", "yes":
		return true, nil
	case "false", "0", "off", "no":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean value: %s", s)
	}
}
