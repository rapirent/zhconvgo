package main

import (
	"fmt"

	"github.com/rapirent/zhconvgo"
)

func main() {
	s := zhconvgo.Convert("秦川雄帝宅，函谷壯皇居。綺殿千尋起，離宮百雉餘。連甍遙接漢，飛觀迥凌虛。雲日隱層闕，風煙出綺疎。", zhconvgo.ZHCN)
	fmt.Println(s)
	// s = zhconvgo.Convert("秦川雄帝宅，函谷壯皇居。綺殿千尋起，離宮百雉餘。連甍遙接漢，飛觀迥凌虛。雲日隱層闕，風煙出綺疎。", zhconvgo.ZHCN)
}
