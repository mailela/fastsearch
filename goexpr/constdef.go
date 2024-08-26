package goexpr

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	Variable string = "Variable"
	Func     string = "function"
	PreOp    string = "prefixOp"
	InfxOp   string = "infixOp"
	Value    string = "value"
	Array    string = "array"
	Args     string = "funcArgs"
)

const (
	Trim       string = "Trim"
	IN         string = "IN"
	NotIn      string = "NotIN"
	Like       string = "Like"
	Equal      string = "=="
	Less       string = "<"
	LessEqual  string = "<="
	AboveEqual string = ">="
	Above      string = ">"
	NotEqual   string = "!="
	Add        string = "+"
	Sub        string = "-"
	Mult       string = "*"
	Div        string = "/"
	Rest       string = "%"
	And        string = "&&"
	Or         string = "||"
	Not        string = "!"
	BraktLeft  string = "("
	BraktRight string = ")"
	ArrayLeft  string = "["
	ArrayRight string = "]"
	ItemSpit   string = ","
)

var OpPriority = map[string]int32{
	Mult:       60,
	Div:        60,
	Rest:       60,
	Add:        50,
	Sub:        50,
	Above:      40,
	AboveEqual: 40,
	Less:       40,
	LessEqual:  40,
	Equal:      40,
	NotEqual:   40,
	IN:         30,
	NotIn:      30,
	Like:       30,
	And:        20,
	Or:         10,
}

// 函数运算
type FunctionOp func(v ...interface{}) interface{}

// 前缀运算
type PrefixOp func(v interface{}) interface{}

// 中缀运算
type InfixOp func(v1, v2 interface{}) interface{}

var OpFunctionSet = map[string]FunctionOp{
	Trim: func(v ...interface{}) interface{} {
		v1 := fmt.Sprint(v)
		v1 = strings.Trim(v1, "[")
		v1 = strings.Trim(v1, "]")
		v1 = strings.Replace(v1, " ", "", -1)
		return strings.TrimSpace(v1)
	},
}

var PrefixOpSet = map[string]PrefixOp{
	Not: func(v1 interface{}) interface{} {
		return !v1.(bool)
	},
	Sub: func(v interface{}) interface{} {
		return -FloatVal(v)
	},
}

var InfixOpSet = map[string]InfixOp{
	Equal: func(v1, v2 interface{}) interface{} {
		return fmt.Sprint(v1) == fmt.Sprint(v2)
	},
	Add: func(v1, v2 interface{}) interface{} {
		return FloatVal(v1) + FloatVal(v2)
	},
	Sub: func(v1, v2 interface{}) interface{} {
		return FloatVal(v1) - FloatVal(v2)
	},
	Mult: func(v1, v2 interface{}) interface{} {
		return FloatVal(v1) * FloatVal(v2)
	},
	Div: func(v1, v2 interface{}) interface{} {
		return FloatVal(v1) / FloatVal(v2)
	},
	Rest: func(v1, v2 interface{}) interface{} {
		return int64(FloatVal(v1)) % int64(FloatVal(v2))
	},
	And: func(v1, v2 interface{}) interface{} {
		return v1.(bool) && v2.(bool)
	},
	Or: func(v1, v2 interface{}) interface{} {
		return v1.(bool) || v2.(bool)
	},
	Less: func(v1, v2 interface{}) interface{} {
		return FloatVal(v1) < FloatVal(v2)
	},
	LessEqual: func(v1, v2 interface{}) interface{} {
		return FloatVal(v1) <= FloatVal(v2)
	},
	AboveEqual: func(v1, v2 interface{}) interface{} {
		return FloatVal(v1) >= FloatVal(v2)
	},
	Above: func(v1, v2 interface{}) interface{} {
		return FloatVal(v1) > FloatVal(v2)
	},
	NotEqual: func(v1, v2 interface{}) interface{} {
		return fmt.Sprint(v1) != fmt.Sprint(v2)
	},
	IN: func(v1 interface{}, v2 interface{}) interface{} {
		return in(v1, v2)
	},
	NotIn: func(v1 interface{}, v2 interface{}) interface{} {
		return notIn(v1, v2)
	},
	Like: func(v1 interface{}, v2 interface{}) interface{} {
		return like(v1, v2)
	},
}

func notIn(a, b interface{}) interface{} {
	if b == nil {
		return true
	}
	aStr := fmt.Sprint(a)
	array := reflect.ValueOf(b)
	length := array.Len()
	for i := 0; i < length; i++ {
		bStr := fmt.Sprint(array.Index(i).Interface())
		if bStr == aStr {
			return false
		}
	}
	return true
}

func like(a, b interface{}) interface{} {
	if b == nil {
		return false
	}
	if a == nil {
		return false
	}
	aStr := fmt.Sprint(a)
	str := strings.TrimSpace(fmt.Sprint(b))
	reg := regexp.MustCompile(str)
	result := reg.FindString(aStr)
	if result == "" {
		return false
	}
	return true
}

func in(a, b interface{}) interface{} {
	if b == nil {
		return false
	}
	aStr := fmt.Sprint(a)
	array := reflect.ValueOf(b)
	length := array.Len()
	for i := 0; i < length; i++ {
		bStr := fmt.Sprint(array.Index(i).Interface())
		if bStr == aStr {
			return true
		}
	}
	return false
}

func FloatVal(v interface{}) float64 {
	if v == nil {
		v = ""
	}
	s := fmt.Sprint(v)
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func Contain(a, b interface{}) interface{} {
	bStr := fmt.Sprint(b)
	array := reflect.ValueOf(a)
	length := array.Len()
	for i := 0; i < length; i++ {
		aStr := fmt.Sprint(array.Index(i).Interface())
		if bStr == aStr {
			return true
		}
	}
	return false
}
