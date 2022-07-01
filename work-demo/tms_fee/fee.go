package tms_fee

import (
	"fmt"

	"github.com/spf13/cast"
)

func fee_price() {
	countNum := 5
	skuQuantity := 4
	countType := 1
	var goodsQuantity float64
	// 6期优化，计件数/10，客户端传值乘以10
	countNums := cast.ToFloat64(countNum) / 10
	// 一件计多件
	if countType == 1 {
		goodsQuantity = cast.ToFloat64(skuQuantity) * countNums
	} else {
		// 多件计一件,
		goodsQuantityInt := int(cast.ToFloat64(skuQuantity) / countNums)

		skuQuantity := skuQuantity * 10
		remainder := skuQuantity % countNum
		if cast.ToFloat64(remainder)/cast.ToFloat64(countNum) >= 0.5 {
			goodsQuantityInt++
		}
		goodsQuantity = cast.ToFloat64(goodsQuantityInt)
	}
	fmt.Println("goodsQuantity: ", goodsQuantity)
}
