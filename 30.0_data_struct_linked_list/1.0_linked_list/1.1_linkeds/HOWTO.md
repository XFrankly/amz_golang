# 从零开始， 以双向链表实现 运势计算(周易)
## 0.0 需要获得每天的运势吗？
    本文实现算法 主要参考 熊逸大师的说明，《当你的邻居为九五之尊》，这里做一个简单的汇总，如果有需要参考原文的请自行搜索。

    占卜算卦是一个古老的方法，人们无法把握未来时，经常通过这个方法获得一些确定性的心理效应。
    其社会影响深远而广泛，虽然它并不是一种完全准确的科学，但我们可以通过go实现 一个简易的占卜程序，用以了解我们的古人如何处理不确定性的。 并尝试了解其中的变化规则。


## 0.1 计算步骤初略
    步骤
        1， 产生49个数，为大衍之数为 54,实际参与计算的为 七七四十九
        2， 随机 分成 两组（随机的方法非常重要，可以影响甚至决定结果），
            为了可以产生正常的 三变 运算，最小12
            也就是说范围是：天 12 ~ 37， 地 37 ~ 12

        3， 从 天，地 随机选择一侧 取出一个 作为 人（一般在实际操作时，可以让寻求运势的人操作，以加强参与感）
        4，  一变 把象征天的那组棋子数数有多少颗，数清楚之后把这个数字除以4（4象征一年四季），余数是几。任何数字除以4，余数都只有四种可能：1、2、3、整除。如果遇到整除的情况，我们就当做余数是4。好了，现在把余数拿开。
        把象征地的那组棋子照猫画虎，和“3”的做法一致。
        把“2”里用来象征人的那一颗棋子，加上“3”中作为余数被拿掉的棋子，还有“4”里同样作为余数被拿掉的棋子归在一起。得出的数字只有两种可能：不是9就是5。如果错了，你就从头再来吧。 　　 　好了，从1-5完成动作，叫做“一变”。

        5， 二变，重复 2~4 步
        6， 三变，重复 2~4 步，最后剩下的卡 除以 4 得出 第一爻

        7， 重复 1~6 步，6次，产生 6个爻 即为 一卦
        8， 解卦象

## 0.2 双向链表的实现
    我们可以通过双向链表 存储卦象的6个爻，这样我们可以知道前后顺序，并且可以在每个节点存储变爻的值

    //定义链表结构体
    type node struct {
            number  int    //爻值
            yaobian [][]int   //三次爻变的 具体算子
            prev    *node   //前一个爻 节点
            next    *node   //后一个爻 节点
        }

    //定义双向链表
    type dlist struct {
            lens int
            head *node
            tail *node
        }
    
    //链表构造函数
    func makeDlist() *dlist {
            return &dlist{}
        }

    //判断是否空链表 
    func (this *dlist) newNodeList(n *node) bool { 

        if this.lens == 0 {
            this.head = n
            this.tail = n
            n.prev = nil
            n.next = nil
            this.lens += 1
            return true
        } else {
            Logg.Panic("not empty node list.")
        }
        return false
    }

    
    // 头部添加 节点
    func (this *dlist) pushHead(n *node) bool {

        if this.lens == 0 {
            return this.newNodeList(n)
        } else {
            this.head.prev = n
            n.prev = nil
            n.next = this.head
            this.head = n
            this.lens += 1
            return true
        }
    }

    
    //  添加尾部节点，我们主要使用此方法，用以保持爻的相对位置
    func (this *dlist) append(n *node) bool {

        if this.lens == 0 {
            return this.newNodeList(n)
        } else {
            this.tail.next = n
            n.prev = this.tail
            n.next = nil
            this.tail = n
            this.lens += 1
            return true
        }
    }

    
    /// 显示并返回链表的值
    func (this *dlist) display() []int {
        
        numbs := []int{}
        node := this.head
        t := 0 
        for node != nil {

            Logg.Println(node.number, node.yaobian)
            numbs = append(numbs, node.number)
            t += 1
            if t >= this.lens {
                break
            }

            node = node.next
        }

        fmt.Println("length:", this.lens)
        return numbs
    }

## 1.0 准备开始
    我们先制定一个计算结构体，一切开始于此，我们首先需要知道参与计算的 算子有多少，我们就按最常用的数49来制定

    type DataTao struct {
        DefMean    map[int]string // 卦象坐标名称
        DefValue   map[int]string // 卦象坐标值 4个
        GuaData    []int          // 卦象的 6爻初始值
        Indent     int            // 格式化输出间隔
        Env        int            // 影响偏差，Env 把一年分为 3个部分，每个部分4个月，这个值可以影响天地分开的比例，自己随意
        Total      int            // 大衍之数 个数
        PrintLevel int            // 计算过程显示控制，如果为 1 则显示
        Coordinate map[int]string //卦象坐标的 爻，古称
    }

    制作该计算对象的构造函数
    
        /*
        @param env 环境参数
        @param indent 格式化输出的缩进
        @param total 天地 大衍数 49
        */
        func MakeNewDataTao(env string, indent, total int) *DataTao {
            envNum := map[string]int{"A": 1, "B": 0, "C": 2}
            guaCoor := map[int]string{6: "六", 7: "七", 8: "八", 9: "九"}

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

## 1.1 开天辟地，创建爻变
    创建天地大衍之数
    // 创建一个新的天地列表 大衍之数，使用49
    func (dt *DataTao) NewCircle() []int {
        var newCir []int
        for i := 0; i < dt.Total; i++ {
            newCir = append(newCir, i)
        }
        return newCir
    }


    由于每一变 除以4的余数可能为 0，此时需要拿走4个算子，所以，3变最多需要12个算子，而
    // 分天地，默认随机算法   计算后 将其值 限定为下限 12，上限 37
    func (dt *DataTao) SplitNum(topNum int) int {
        /*
            Seed 使用提供的种子值将生成器初始化为确定性状态。
            Seed 不应与任何其他 Rand 方法同时调用。
        */
        rand.Seed(time.Now().UnixNano())
        top := topNum - 4
        sn := rand.Intn(top) + 4 // A random int in [4, 45]

        if dt.Env != 0 {
            if dt.Env == 1 {
                sn = sn + (dt.Env * 3)
            } else {
                sn = sn + (dt.Env * 2)
            }

        }
        if sn < 4 {
            sn = 4
        }

        if sn > topNum - 4 {
            sn = topNum - 4
        }

        logger.Printf("outcome sn:%v with topNum:%v\n", sn, topNum)
        return sn
    }

### 1.1.1  创建爻变,得到三才
    从天，地随机选择一方，拿出一个算子 作为天地人的 人才
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


    //做一变
    /*
    @param: lis, 天地混沌，需要分开的算子集合
    @param: sn, 从算子集合的哪一个位置执行分开
    // 做一变，为得一爻
    */
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


    // 每一个爻都需要对 天地大衍之数做三变，最后剩下的算子数 除以4 就是 一爻，算6次爻 那么就是一卦
    // 连续三变 得一爻，返回新天地
    func (dt *DataTao) TianDiYao(lis []int) ([]int, int) {
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

 ### 1.1.2  创建爻变, 得到6爻 即一卦
    计算6个爻 则需要 3变 6次
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

 ### 1.1.3    格式化输出， 让我们漂亮地查看爻的卦象
    // 格式化输出
    func (dt *DataTao) FormatShow(cont string) string {
        
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



## 1.2 准备解释材料 
    世界的变化都在天地之间，人类关注的内容主要就是 天，地，人 三才
    https://github.com/XFrankly/amz_golang/tree/exerice/30.0_data_struct_linked_list/1.0_linked_list/1.1_linkeds/taodata/taodata.go

    //解释文本,找到运势卦象后 通过此函数获得对应的 卦象解释和对应数字
    func (ic *ICh) CommonText(arrays []int) (string, int) {
        if len(arrays) != ic.Num {
            return ic.Warn + fmt.Sprintf("%v", arrays), 0
        }

        if ic.CompSlice(arrays, []int{1, 1, 1, 1, 1, 1}) {
            return `
            第一卦 乾 乾为天 乾上乾下 

            乾：元，亨，利，贞。 
            ....
    
    在显示卦象时，提示用户爻变的位置和解释方式
        AnySis = map[int]string{
            0: "  六爻全不变,以本 卦卦辞占   ----- 不变 只看本卦 ",
            1: "一爻变,以本卦变爻的爻辞占          ----- 变卦前 的那一 变爻 含义",
            2: " 如果两个爻发生变动,使用本卦的两个变爻占辞判断吉凶,以位置靠上的为主   ----- 变卦前 的 第一变爻 含义",
            3: " 三爻变,以本卦及之卦的卦辞占,以本卦的卦辞为主    ----- 变卦前后都看， 本卦及之卦 含义  本卦为主",
            4: "  四爻变,以之卦中二不变的爻辞占,以下爻的爻辞为主   ----- 变卦 后 看，  之卦 含义  二不变的爻辞占,以下爻的爻辞为主",
            5: "  五爻表,以之卦中不变的爻辞占卜   ----- 变卦 后 看，  之卦 含义  不变的爻辞占 ",
            6: "  六爻全变,乾坤两卦以用九和用六的辞占卜,并参考之卦卦辞,余六十二卦占以之卦卦辞  ----- 变卦 后 看，  之卦 含义，乾坤 用九和用六 ，其他卦全看",
        }
    



## 1.3 解释我们计算出来的卦象
    
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

    我们可以计算一万次，并统计64卦 分别出现多少次，从而知道中有多少人与自己的运势相同

    计算10000次，并统计次数。

# 结语  
    本文介绍了 双向链表的实现和在周易算卦中的简单实现，我们可以理解其具有广泛民众基础的算卦如何运作
    也可以对双向链表加深理解，本项目是本人众多项目中的一个，希望能对看官有所益

    本文的完整项目链接地址，
    https://github.com/XFrankly/amz_golang/tree/exerice/30.0_data_struct_linked_list/1.0_linked_list/1.1_linkeds/taodata












