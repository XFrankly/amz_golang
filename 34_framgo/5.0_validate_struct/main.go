package main

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/asaskevich/govalidator"
)

type User struct {
	Id        int    `validata:"number,min=1,max=1000" json:"id"`
	Name      string `validata:"string, min=2,max=10" json:"name"`
	Bio       string `validata:"string" json:"bio"`
	Email     string `validata:"email" json:"email"`
	Active    string `validata:"active" json:"active"`
	Admin     string `validata:"admin" json:"admin"`
	CreatedAt string `validata:"created_at" json:"created_at"`
}

/*
Id 值在某一个范围
Name 长度在某一个范围
Email格式校验
//// 不使用接口，校验结构字段
if tagIsOfNumber(){
	validator := NumberValidator{}
} else if tagIsOfString() {
		validator := StringValidator{}
...
}

使用接口，所有validator都去实现这个接口
*/

const (
	tagName = "validate"
)

//邮箱正则校验
var (
	mailRe = regexp.MustCompile(`A[w+-.]+@[a-zd-]+(.[a-z]+)*.[a-z]+z`)
)

//验证接口
type Validator interface {
	Validate(interface{}) (bool, error)
}

//默认验证
type DefaultValidator struct{}

func (v DefaultValidator) Validate(val interface{}) (bool, error) {
	return true, nil
}

type StringValidator struct {
	Min int
	Max int
}

func (v StringValidator) Validate(val interface{}) (bool, error) {
	le := len(val.(string))

	if le == 0 {
		return false, fmt.Errorf("cannot be blank.")
	}
	if le < v.Min {
		return false, fmt.Errorf("should be at least %v chars long", v.Min)
	}
	if v.Max >= v.Min && le > v.Max {
		return false, fmt.Errorf("should be less than %v chars long", v.Max)
	}
	return true, nil
}

type NumberValidator struct {
	Min int
	Max int
}

func (v NumberValidator) Validate(val interface{}) (bool, error) {
	le := val.(int)

	if le < v.Min {
		return false, fmt.Errorf("should be at bigger than %v ", v.Min)
	}
	if v.Max >= v.Min && le > v.Max {
		return false, fmt.Errorf("should be less than %v  ", v.Max)
	}
	return true, nil
}

type EmailValidator struct{}

func (v EmailValidator) Validate(val interface{}) (bool, error) {

	if !mailRe.MatchString(val.(string)) {
		return false, fmt.Errorf("not valid email addr: %v ", val.(string))
	}
	return true, nil
}

func getValidatorFromTag(tag string) Validator {
	args := strings.Split(tag, ",")

	switch args[0] {
	case "number":
		validator := NumberValidator{}
		//将structTag 中min和max解析到结构体
		fmt.Sscanf(strings.Join(args[1:], ","), "min=%d,max=%d", &validator.Min, &validator.Max)
		return validator
	case "string":
		validator := StringValidator{}
		fmt.Sscanf(strings.Join(args[1:], ","), "min=%d,max=%d", &validator.Min, &validator.Max)
		return validator
	case "email":
		return EmailValidator{}
	}
	return DefaultValidator{}
}

// 使用结构体字段校验
func validateStruct(s interface{}) []error {
	errs := []error{}

	v := reflect.ValueOf(s)
	for i := 0; i < v.NumField(); i++ {
		//利用反射获取structTag
		tag := v.Type().Field(i).Tag.Get(tagName)

		if tag == "" || tag == "-" {
			continue
		}

		validator := getValidatorFromTag(tag)

		valid, err := validator.Validate(v.Field(i).Interface())
		if !valid && err != nil {
			errs = append(errs, fmt.Errorf("%s %s", v.Type().Field(i).Name, err.Error()))
		}
	}
	return errs
}

// 使用自定义方法校验
func ClustomCheck() {
	user := User{
		Id:    0,
		Name:  "superlongstring",
		Bio:   "",
		Email: "foobar",
	}
	fmt.Println("check Errors:")
	for i, err := range validateStruct(user) {
		fmt.Printf("t%d. %sn", i+1, err.Error())
	}
}

// 使用第三方包校验，支持内置tag 和 自定义tag
func ThirdCheckPackage() {
	user := &User{
		Id:    1,
		Name:  "IX01",
		Bio:   "bio its me.",
		Email: "foobar@",
	}
	//自定义tag验证函数
	govalidator.TagMap["email"] = govalidator.Validator(func(str string) bool {
		// return strings.HasPrefix(str, "IX")
		return strings.Contains(str, "@")
	})

	if ok, err := govalidator.ValidateStruct(user); err != nil {
		panic(err)
	} else {
		fmt.Printf("Validate OK:%v\n", ok)
	}

}

func main() {
	ClustomCheck()
	// ThirdCheckPackage()
}
