package exp

import (
	"fmt"
	"testing"

	goExpr "gitee.com/rachel_os/fastsearch/goexpr"
)

type ExpEngine struct {
	expr   string
	Engine goExpr.Engine
}

func NewEvaluableExpression(exprs string) (ExpEngine, error) {

	eg := goExpr.NewEngine()
	eg.AddFunc("ADD", func(v ...interface{}) interface{} {
		return goExpr.FloatVal(v[0]) + goExpr.FloatVal(v[1])
	})
	eg.AddFunc("FloatVal", func(v ...interface{}) interface{} {
		return goExpr.FloatVal(v[0])
	})
	eg.AddPrefix("#", func(v interface{}) interface{} {
		return goExpr.FloatVal(v) * goExpr.FloatVal(v)
	})
	eg.AddInfix("Contain", 30, func(v1, v2 interface{}) interface{} {
		return goExpr.Contain(v1, v2)
	})

	return ExpEngine{
		expr:   exprs,
		Engine: *eg,
	}, nil
}

func (eng ExpEngine) Evaluate(params map[string]interface{}) (interface{}, error) {
	result := eng.Engine.Execute(eng.expr, params)
	return result, nil
}

func TestEngine(t *testing.T) {
	exprs := `-4+3>(-9)&&5<4+5&&3NotIN[1,2,4]&&ADD(1,2)<4&&-(#-3-4)<=30&&4>1&&[1,2,4] Contain 4 && ADD(1,2)!=1 && user.name=='kiteee' && user_count>20`
	//exprs = `-------1`
	//exprs=`user.name=='kiteee' && user_count>20`
	//exprs=`user_count>20 && user_count>20`
	//exprs=`#--3*-4-#2`
	//exprs=`-4-#2`
	//exprs = `3NotIN([1,2,3])&&ADD(1,2)<4`
	//exprs = `[1,2,4] Contain 4 IN [true] NotIN [false]`
	eg := goExpr.NewEngine()
	eg.AddFunc("ADD", func(v ...interface{}) interface{} {
		return goExpr.FloatVal(v[0]) + goExpr.FloatVal(v[1])
	})
	eg.AddPrefix("#", func(v interface{}) interface{} {
		return goExpr.FloatVal(v) * goExpr.FloatVal(v)
	})
	eg.AddInfix("Contain", 30, func(v1, v2 interface{}) interface{} {
		return goExpr.Contain(v1, v2)
	})
	var params = map[string]interface{}{
		"user": map[string]interface{}{
			"name": "kiteee",
			"age":  50,
		},
		"user_count": 30,
	}
	//eg.SetPriority("NotIN", 30)
	result := eg.Execute(exprs, params)
	fmt.Println(result)
}
