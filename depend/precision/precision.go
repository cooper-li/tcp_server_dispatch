package precision

import (
	"github.com/shopspring/decimal"
	"gitee.com/johng/gf/g/util/gconv"
)

// 高精度运算

// 加
func Add(x, y interface{}) float64 {
	tx, ty := decimal.NewFromFloat(gconv.Float64(x)), decimal.NewFromFloat(gconv.Float64(y))
	z, _ := tx.Add(ty).Float64()
	return z
}

// 减
func Sub(x, y interface{}) float64 {
	tx, ty := decimal.NewFromFloat(gconv.Float64(x)), decimal.NewFromFloat(gconv.Float64(y))
	z, _ := tx.Sub(ty).Float64()
	return z
}

// 乘
func Mul(x, y interface{}) float64 {
	tx, ty := decimal.NewFromFloat(gconv.Float64(x)), decimal.NewFromFloat(gconv.Float64(y))
	z, _ := tx.Mul(ty).Float64()
	return z
}

// 除
func Div(x, y interface{}) float64 {
	tx, ty := decimal.NewFromFloat(gconv.Float64(x)), decimal.NewFromFloat(gconv.Float64(y))
	z, _ := tx.Div(ty).Float64()
	return z
}

// 对比 x < y : -1  ||  x == y : 0   ||  x > y : 1
func Compare(x, y interface{}) int {
	tx, ty := decimal.NewFromFloat(gconv.Float64(x)), decimal.NewFromFloat(gconv.Float64(y))
	return tx.Cmp(ty)
}
