package member

import (
	"github.com/go-playground/validator/v10"
)

func NameValid(fl validator.FieldLevel)bool{
		if s,ok := fl.Field().Interface().(string);ok{
			if s == "admin"{
				return false
			}
		}
		return true
}
func AgeValid(f1 validator.FieldLevel)bool{
	if s,ok := f1.Field().Interface().(int);ok{
		if s < 0 {
			return false
		}
	}
	return true
}