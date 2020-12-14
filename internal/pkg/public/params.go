package public

import (
	"errors"
	"strings"

	ut "github.com/go-playground/universal-translator"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// DefaultGetValidParams 获取
func DefaultGetValidParams(c *gin.Context, params interface{}) error {
	if err := c.ShouldBind(params); err != nil {
		return err
	}
	// 获取验证器
	valid, err := GetValidator(c)
	if err != nil {
		return err
	}
	// 获取翻译器
	trans, err := GetTranslation(c)
	if err != nil {
		return err
	}

	err = valid.Struct(params)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		sliceErrs := []string{}
		for _, e := range errs {
			sliceErrs = append(sliceErrs, e.Translate(trans))
		}
		return errors.New(strings.Join(sliceErrs, "."))
	}
	return nil
}

// GetValidator 获取验证器
func GetValidator(c *gin.Context) (*validator.Validate, error) {
	val, ok := c.Get(CtxValidatorKey)
	if !ok {
		return nil, errors.New("未设置验证器")
	}
	validator, ok := val.(*validator.Validate)
	if !ok {
		return nil, errors.New("获取验证器失败")
	}
	return validator, nil
}

// GetTranslation 获取翻译器
func GetTranslation(c *gin.Context) (ut.Translator, error) {
	trans, ok := c.Get(CtxTranslatorKey)
	if !ok {
		return nil, errors.New("未设置翻译器")
	}
	translator, ok := trans.(ut.Translator)
	if !ok {
		return nil, errors.New("获取翻译器失败")
	}
	return translator, nil
}
