package middleware

import (
	"reflect"

	"github.com/captainlee1024/fast-gin/internal/pkg/public"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

// TranslationMiddleware 设置 Translation
func TranslationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置支持语言
		enT := en.New() // 英文翻译器
		zhT := zh.New() // 中文翻译器

		// 设置国际化翻译器
		// 第一个参数是备用 (fallback) 的语言环境，后面的参数是支持的语言环境，可以是多个
		uni := ut.New(zhT, zhT, enT)
		val := validator.New()

		// 根据参数取翻译器实例
		// 也可以使用 uni.FindTranslator(...) 传入多个locale进行查找
		local := c.DefaultQuery("local", "zh")
		trans, _ := uni.GetTranslator(local)

		// 把翻译器注册到验证器（validator）
		// 校验器是真正做校验的，这里注册到在 gin 中拿到的校验其中
		switch local {
		case "en":
			enTranslations.RegisterDefaultTranslations(val, trans)
			val.RegisterTagNameFunc(func(fld reflect.StructField) string {
				return fld.Tag.Get("en_comment")
			})
			break
		default:
			zhTranslations.RegisterDefaultTranslations(val, trans)
			val.RegisterTagNameFunc(func(fld reflect.StructField) string {
				return fld.Tag.Get("comment")
			})

			// 自定义验证方法
			val.RegisterValidation("tag string", func(fl validator.FieldLevel) bool {
				return fl.Field().String() == "admin"
			})

			// 自定义验证器
			val.RegisterTranslation("is-validuer", trans,
				func(ut ut.Translator) error {
					return ut.Add("is-validuser", "{0}填写不正确", true)
				},
				func(ut ut.Translator, fe validator.FieldError) string {
					t, _ := ut.T("is-validuser", fe.Field())
					return t
				})
			break
		}
		c.Set(public.CtxTranslatorKey, trans)
		c.Set(public.CtxValidatorKey, val)
		c.Next()
	}
}
