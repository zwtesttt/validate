package validator

import (
	"fmt"
	"strconv"
	"strings"
)

type ValidationFunc func(ValidationStruct) error

type ValidationStruct struct {
	MessageSep    string
	fieldName     string
	validatorName string
	value         interface{}
	args          []string
	msgMap        map[string]map[string]string
}

var validationFunctions = map[string]ValidationFunc{
	"required": validateRequired,
	"email":    validateEmail,
	"min":      validateMin,
	// 添加其他验证函数
}

func validateRequired(vs ValidationStruct) error {
	// 验证逻辑
	if vs.value == "" {
		errMsg := vs.msgMap[vs.fieldName][vs.validatorName]
		if errMsg == "" {
			errMsg = "The value must"
		}
		return fmt.Errorf(errMsg)
	}
	return nil
}

func validateEmail(vs ValidationStruct) error {
	// 验证逻辑
	return nil
}

func validateMin(vs ValidationStruct) error {
	// 验证逻辑
	if len(vs.args) == 0 {
		return fmt.Errorf("min validation requires an argument")
	}
	ops := strings.Split(vs.args[0], vs.MessageSep)
	min, err := strconv.Atoi(ops[0])
	if err != nil {
		return fmt.Errorf("%s Parameter error", vs.validatorName)
	}
	if intValue, ok := vs.value.(int); ok {
		if intValue < min {
			errMsg := vs.msgMap[vs.fieldName][vs.validatorName]
			if !strings.Contains(errMsg, "%d") {
				return fmt.Errorf(errMsg)
			}

			if errMsg == "" {
				errMsg = "value min is %d"
			}
			return fmt.Errorf(errMsg, min)
		}
	}
	return nil
}
