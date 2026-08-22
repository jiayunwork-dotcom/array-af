package geometry

import (
	"errors"
	"fmt"
	"math"
)

// FieldError 描述校验失败的具体字段与原因。
type FieldError struct {
	Field string
	Value any
	Rule  string
}

func (e FieldError) Error() string {
	return fmt.Sprintf("%s: %v violates %s", e.Field, e.Value, e.Rule)
}

// ValidationError 聚合一次输入的全部字段错误。
type ValidationError struct {
	Fields []FieldError
}

func (e *ValidationError) Error() string {
	if e == nil || len(e.Fields) == 0 {
		return "validation failed"
	}
	msg := "invalid input:"
	for i, f := range e.Fields {
		if i > 0 {
			msg += ";"
		}
		msg += " " + f.Error()
	}
	return msg
}

// FieldCount 返回聚合的错误字段数。
func (e *ValidationError) FieldCount() int {
	if e == nil {
		return 0
	}
	return len(e.Fields)
}

// Add 追加一个字段错误。
func (e *ValidationError) Add(f FieldError) {
	e.Fields = append(e.Fields, f)
}

// Validate 校验输入参数，返回 nil 表示全部合法。
//
// 规则：
//   - N 必须是整数且 >= 2（N=1 退化为单阵元，没有可分析的阵因子峰）
//   - d 必须 > 0（零间距时阵元重叠，ψ 恒等于 β，栅瓣判定失真）
//   - lambda 必须 > 0（非正波长导致波数 k 无意义）
func Validate(p Params) error {
	ve := &ValidationError{}
	if p.N < 2 {
		ve.Add(FieldError{
			Field: "N",
			Value: p.N,
			Rule:  "must be an integer >= 2",
		})
	}
	if p.D <= 0 {
		ve.Add(FieldError{
			Field: "d",
			Value: p.D,
			Rule:  "must be > 0",
		})
	}
	if p.Lambda <= 0 {
		ve.Add(FieldError{
			Field: "lambda",
			Value: p.Lambda,
			Rule:  "must be > 0",
		})
	}
	if ve.FieldCount() == 0 {
		return nil
	}
	return ve
}

// ErrNotFinite 表示输入中存在 NaN 或无穷大数值。
var ErrNotFinite = errors.New("input contains NaN or Inf")

// CheckFinite 校验输入浮点数都是有限值。
func CheckFinite(p Params) error {
	if p.D != p.D || p.Lambda != p.Lambda || p.Beta != p.Beta {
		return ErrNotFinite
	}
	if !finite(p.D) || !finite(p.Lambda) || !finite(p.Beta) {
		return ErrNotFinite
	}
	return nil
}

func finite(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0)
}
