
zhconvgo
===============

**zhconvgo 是 [zhconv](https://github.com/gumblex/zhconv) 实现的 golang 版本** 目前仅实现了 Convert 函数

### Usage
```go
package main

import (
	"fmt"

	"github.com/DCRcoder/zhconvgo"
)

func main() {
	s := zhconvgo.Convert("我幹什麼不干你事", zhconvgo.ZHCN)
	// 我干什么不干你事
	s := zhconvgo.Convert("秦川雄帝宅，函谷壯皇居。 綺殿千尋起，離宮百雉餘。 連甍遙接漢，飛觀迥凌虛。 雲日隱層闕，風煙出綺疎。", zhconvgo.ZHCN)
	// 秦川雄帝宅，函谷壮皇居。 绮殿千寻起，离宫百雉余。 连甍遥接汉，飞观迥凌虚。 云日隐层阙，风烟出绮疏。
}
```
# License
Licensed under the [MIT license](./LICENSE)
