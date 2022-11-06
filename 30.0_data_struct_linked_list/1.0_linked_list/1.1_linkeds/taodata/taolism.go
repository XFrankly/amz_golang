package taodata

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

/*
步骤
1， 产生49个数
2， 随机 分成 两组， 为了可以产生正常的 三变 运算，最小12
	也就是说范围是：天 12 ~ 37， 地 37 ~ 12

3， 从 天，地 随机选择一侧 取出一个 作为 人
4，  一变 把象征天的那组棋子数数有多少颗，数清楚之后把这个数字除以4（4象征一年四季），看看余数是几。任何数字除以4，余数都只有四种可能：1、2、3、整除。如果遇到整除的情况，我们就当做余数是4。好了，现在把余数拿开。
把象征地的那组棋子照猫画虎，和“3”的做法一致。
把“2”里用来象征人的那一颗棋子，加上“3”中作为余数被拿掉的棋子，还有“4”里同样作为余数被拿掉的棋子归在一起。得出的数字只有两种可能：不是9就是5。如果错了，你就从头再来吧。 　　 　好了，从1-5完成动作，叫做“一变”。

5， 二变，重复 2~4 步
6， 三变，重复 2~4 步，最后剩下的卡 除以 4 得出 第一爻

7， 重复 1~6 步，6次，产生 6个爻 即为 一卦
8， 解卦象

*/

var (
	defMean     = map[int]string{6: "老阴", 7: "少阳", 8: "少阴", 9: "老阳"}
	Coordinates = map[int]string{1: "初", 2: "二", 3: "三", 4: "四", 5: "五", 6: "上"}
	logger      = log.New(os.Stderr, "INFO -", 4)
	ResultNodes = makeDlist()
	SixYao      = [][]int{}
	start       = `
           _______
           _______
           _______
           .......
___ ___....元亨利贞...._______
___ ___....运势占卜...._______
_______....大吉大利....___ ___ 
           .......
           ___ ___
           ___ ___
           ___ ___

`
)

type DataTao struct {
	DefMean  map[int]string
	DefValue map[int]string //
	GuaData  []int
	//format print
	Indent     int
	Env        int
	Total      int
	PrintLevel int
	Coordinate map[int]string
}

/*
@param env 环境参数
@param indent 格式化输出的缩进
@param total 天地 大衍数 49
*/
func MakeNewDataTao(env string, indent, total int) *DataTao {
	envNum := map[string]int{"A": 1, "B": 0, "C": 2}
	guaCoor := map[int]string{6: "六", 7: "七", 8: "八", 9: "九"} //# 卦象坐标值 4个

	var envLevel = 0

	if indent <= 0 {
		indent = 8
	}

	if env != "" {
		envLevel = envNum[env]
	}

	if total < 49 {
		total = 49
	}

	newData := &DataTao{
		DefMean:    defMean,
		DefValue:   guaCoor,
		GuaData:    []int{},
		Indent:     indent,
		Env:        envLevel,
		Total:      total,
		PrintLevel: envNum["A"],
		Coordinate: Coordinates,
	}
	return newData
}

// 创建一个新的天地列表 大衍之数，使用49
func (dt *DataTao) NewCircle() []int {
	var newCir []int
	for i := 0; i < dt.Total; i++ {
		newCir = append(newCir, i)
	}
	return newCir
}

// 分天地，默认随机算法   计算后 将其值 限定为下限 12，上限 37
func (dt *DataTao) SplitNum(topNum int) int {
	/*
		Seed 使用提供的种子值将生成器初始化为确定性状态。
		Seed 不应与任何其他 Rand 方法同时调用。
	*/
	rand.Seed(time.Now().UnixNano())
	top := topNum - 12
	sn := rand.Intn(top) + 12 // A random int in [12, 37]

	if dt.Env != 0 {
		if dt.Env == 1 {
			sn = sn + (dt.Env * 3)
		} else {
			sn = sn + (dt.Env * 2)
		}

	}
	if sn < 12 {
		sn = 12
	}

	if sn > topNum-12 {
		sn = topNum - 12
	}

	logger.Printf("outcome sn:%v with topNum:%v\n", sn, topNum)
	return sn
}

// 分三才 从选定的一方 得一人
func (dt *DataTao) RandLis(lis []int) ([]int, int) {

	rand.Seed(time.Now().UnixNano())
	if len(lis) < 2 {
		return lis, 0
	}
	t := rand.Intn(len(lis) - 1)
	ren := lis[t]
	lis = append(lis[:t], lis[t+1:]...)
	return lis, ren

}

// 做一变，为得一爻
func (dt *DataTao) DoOneYao(lis []int, sn int) []int {

	//分天地
	tian := lis[:sn]
	di := lis[sn:]

	var ran int
	//从天地 取一人
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(2)
	if r == 0 {
		tian, ran = dt.RandLis(tian)
	} else {
		di, ran = dt.RandLis(di)
	}

	//一变
	var tia, d, newOrigin []int
	xt := len(tian) % 4
	if xt == 0 {
		//取 4个
		tia = tian[:4]
		tian = tian[4:]
	} else {
		tia = tian[:xt]
		tian = tian[xt:]
	}

	dtv := len(di) % 4
	if dtv == 0 {
		d = di[:4]
		di = di[4:]
	} else {
		d = di[:dtv]
		di = di[dtv:]
	}

	//返回一变
	newOrigin = append(tia, d...)
	newOrigin = append(newOrigin, ran)

	if dt.Env == dt.PrintLevel {
		SixYao = append(SixYao, newOrigin)
		logger.Println("余天:", tian, "\n余地:", di, "\n变数:", newOrigin, "人才:", ran, "\n新天地:", len(tian)+len(di))
		logger.Println("已变总数:", len(SixYao))
	}
	newDT := append(tian, di...)
	return newDT
}

// 连续三变 得一爻，返回新天地
func (dt *DataTao) TianDi(lis []int) ([]int, int) {
	/*
			分天地 进行 变一
		        :param lis:
		        :param sn:
		        :return:  最后的 天，地，归档人列表， 和 爻 值
	*/
	finalTiandi := lis
	for _, k := range []int{1, 2, 3} {
		if dt.Env == dt.PrintLevel {
			logger.Println("第", k, "变")
		}
		finalTiandi = dt.DoOneYao(finalTiandi, dt.SplitNum(len(finalTiandi)))
	}

	yao := len(finalTiandi) / 4
	if dt.Total > 49 {
		db := dt.Total / 49
		yao := yao / db
		if dt.Env == dt.PrintLevel {
			logger.Println("大衍之数:", dt.Total, "三变后的爻:", yao, "距离 倍数 db:", db, "新天地:", len(finalTiandi))
		}

	}
	if dt.Env == dt.PrintLevel {
		logger.Println("三变后的爻:", yao)
	}
	return finalTiandi, yao
}

// 入口函数，开天辟地，计算6爻
func (dt *DataTao) SuanGua() *dlist {
	// gua := []int{}
	logger.Printf("开始占卜运势...")
	for i := 1; i < 7; i++ {
		nc := dt.NewCircle()
		if dt.Env == dt.PrintLevel {
			logger.Println(len(nc), nc, dt.SplitNum(len(nc)))
		}
		///分天地人
		ft, y := dt.TianDi(nc)
		if dt.Env == dt.PrintLevel {
			logger.Println("新天地:", ft)
			logger.Println(dt.Coordinate[i]+dt.DefValue[y], "\n第", i, "爻", y, dt.DefMean[y])
		}
		time.Sleep(time.Millisecond * 100)
		// gua = append(gua, y)
		ResultNodes.append(&node{number: y})
	}
	return ResultNodes
}

// 格式化输出
func (dt *DataTao) FormatShow(cont string) string {
	/*

		:param cont:  需要显示的内容
		:return:
	*/
	spaceStr := []string{} //{" ", dt.Indent}
	for i := 0; i < dt.Indent; i++ {
		spaceStr = append(spaceStr, " ")
	}
	msg := strings.Join(spaceStr, "") + cont
	print(msg)
	return msg + "\n"
}

// 转变 6 爻为 卦象后 显示卦象
func (dt *DataTao) Common(gua []int) string {
	/*
	   gua 转变后的列表，只应该有 0，1
	          :param gua: example [1,0,1,0,1,0]
	          :return: string
	*/
	guas := ""
	for _, g := range gua {
		if g == 0 { // 0 表示阴爻
			guas += dt.FormatShow("__ __ (阴)")
		} else if g == 1 { // 1 表示阳爻
			guas += dt.FormatShow("_____ (阳)")
		} else {
			guas += dt.FormatShow("")
		}
	}

	return guas
}

// # 6爻卦象 变卦后，解释其含义
func (dt *DataTao) KanGua(gua []int, n int) (string, int) {
	/*
	   解释卦 的含义
	   :param gua: 转 7, 8 为 1, 0
	   :return:
	*/
	newG := []int{}
	for _, g := range gua {
		if g == 7 {
			newG = append(newG, 1)
		}
		if g == 8 {
			newG = append(newG, 0)
		}
	}

	nc := ICApp
	means, nob := nc.CommonText(newG)
	nc.HowAnysis(n)
	guas := dt.Common(newG)
	print(means)
	result := guas + "\n" + means
	return result, nob
}

// # 6爻卦象 显示原始卦象
func (dt *DataTao) KanGuaOrigin() ([]int, int) {
	c := "n" //# input("是否查看本卦 需要 键入 n:")
	guas := dt.SuanGua()
	gua := guas.display()
	v := 0 //# 变爻次数
	fmt.Println("卦象已出:")
	if c == "n" {
		for _, g := range gua {
			if g == 6 || g == 8 {
				if g == 6 {
					dt.FormatShow("__ __ (6 " + dt.DefMean[6] + ")\n")
					dt.GuaData = append(dt.GuaData, 7) //# 只写入变卦后的 少阳 少阴
					v += 1
				}
				if g == 8 {
					dt.FormatShow("__ __ (8 " + dt.DefMean[8] + ")\n")
					dt.GuaData = append(dt.GuaData, g) //# 只写入变卦后的 少阳 少阴
				}

			} else if g == 7 || g == 9 {
				if g == 7 {
					dt.FormatShow("_____ (7 " + dt.DefMean[7] + ")\n")
					dt.GuaData = append(dt.GuaData, g) //# 只写入变卦后的 少阳 少阴
				}

				if g == 9 {
					dt.FormatShow("_____ (9 " + dt.DefMean[9] + ")\n")
					dt.GuaData = append(dt.GuaData, 8) //# 只写入变卦后的 少阳 少阴
					v += 1
				}

			} else {
				dt.FormatShow("")
			}
		}
	} else {
		for _, g := range gua {
			if g == 6 || g == 8 {
				if g == 6 {
					dt.FormatShow("_____ (6 " + dt.DefMean[6] + " -> 7 " + dt.DefMean[7] + ")\n")
					dt.GuaData = append(dt.GuaData, 7) //# 只写入变卦后的 少阳 少阴
					v += 1
				}

				if g == 8 {
					dt.FormatShow("__ __ (8 " + dt.DefMean[8] + ")\n")
					dt.GuaData = append(dt.GuaData, g) //# 只写入变卦后的 少阳 少阴
				}

			} else if g == 7 || g == 9 {
				if g == 7 {
					dt.FormatShow("_____ (7 " + dt.DefMean[7] + ")\n")
					dt.GuaData = append(dt.GuaData, g) //# 只写入变卦后的 少阳 少阴
				}

				if g == 9 {
					dt.FormatShow("__ __ (9 " + dt.DefMean[9] + "-> 8 " + dt.DefMean[8] + ")\n")
					dt.GuaData = append(dt.GuaData, 8) //# 只写入变卦后的 少阳 少阴
					v += 1
				} else {
					dt.FormatShow("")
				}

			}
		}

	}

	return dt.GuaData, v

}