package main

/*
摹刻 自ChainMap，源代码为

class ChainMaps():  #_collections_abc.MutableMapping
    # 更新链中的第一个映射，但lookup会搜索整个链。
    # 然而，如果需要深度写和删除，也可以很容易的通过定义一个子类来实现它
    '''
    # 集合 收藏品 容器类型， 类似一个字典类型，将多个映射集合到一个视图
    这样它们作为一个单元处理，通常比创建一个新字典和多次调用 update()要快
    可以用作模拟嵌套作用域，并且在模板化的时候比较有用。

    ChainMap 将多个字典（或其他映射）组合在一起创建一个单一的、可更新的视图。
    提供一个默认空字典，这样一个新链至少有一个映射。

    底层映射存储在列表中。该列表是公开的并且可以使用 *maps* 属性访问或更新。没有其他的状态。

    查找连续搜索底层映射，直到找到一个键。  相反，写入、更新和删除只对第一个映射操作

    一个ChainMap通过引用合并底层映射。 所以，如果一个底层映射更新了，这些更改将反映到ChainMap
    支持所有常用字典方法，另外还有一个maps 属性 attribute
    一个创建 子上下文的方法 method，一个存取它们首个映射的属性 property

    '''

    def __init__(self, *maps):
        '''
        初始化 ChainMap，如果没有字典 则初始化一个空的
        '''
        # maps 为一个可以更新的映射列表，这个列表是按照第一次搜索到最后一次搜索的顺序组织的。
        # 它是仅有的存储状态，可以被修改。 它至少包括一个 映射 即 {}
        # 链映射
        self.maps = list(maps) or [{}]  # always at least one map

    #####################
    def __missing__(self, key):
        #  __missing__(key) 如果 default_factory 属性为 None，则调用本方法会抛出 KeyError 异常，附带参数 key。
        #  如果 default_factory 不为 None，则它会被（不带参数地）调用来为
        raise KeyError(key)

    #####################
    def __getitem__(self, key):
    #  支持 self[key] ，否则如下报错
          ### return self[key] if key in self else default
        ##### TypeError: 'ChainMaps' object is not subscriptable
        # 这个异常会原封不动地向外层传递。 在无法找到所需键值时，本方法会被 dict 中的 __getitem__() 方法调用。
        # 无论本方法返回了值还是抛出了异常，都会被 __getitem__() 传递。 注意，__missing__() 不会 被 __getitem__() 以外的其他方法调用
        for mapping in self.maps:
            try:
                return mapping[key]  # can't use 'key in mapping' with defaultdict
            except KeyError:
                pass
        return self.__missing__(key)  # support subclasses that define __missing__

    ###############################
    def get(self, key, default=None):

        return self[key] if key in self else default

    def __iter__(self):
        # 返回一个 可迭代对象，此对象将返回 mapping的所有key
        # 以支持 return self[key] if key in self else default
        d = {}
        # 返回给定序列值的反向迭代器。
        for mapping in reversed(self.maps):
            d.update(dict.fromkeys(mapping))  # reuses stored hash values if possible
        return iter(d)


    #######################################
    def new_child(self, m=None):  # like Django's Context.push()
        '''
        带有新map的新 ChainMap，后跟所有以前的 map。
        如果没有提供 map，则使用空字典。
        '''
        if m is None:
            m = {}
        nChild = self.__class__(m, *self.maps)
        print("new child dict length:", len(nChild.__dict__))
        return nChild

    ###########################################
    @property
    def parents(self):  # like Django's Context.pop()
        '来自地图 [1:] 的新 ChainMap。'
        parents = self.__class__(*self.maps[1:])
        print(f"parents:{len(parents.__dict__)}")
        return parents

    #####################
    def __setitem__(self, key, value):
        # 更新 maps 的第一个元素 键值对，支持子类操作： self.env[name] = value
        print(f'update key:{key}, value:{type(value)} to first maps length:{len(self.maps)}')
        self.maps[0][key] = value
*/
import (
	"fmt"
	"log"
	"os"
	"sync"
)

var (
	Mutex  sync.RWMutex
	Logger = log.New(os.Stderr, "[INFO] --", 18)
)

type Tindex interface {
	int | ~string
}

type TValue interface {
	~string
}

type VaStrings string
type MapValue[T Tindex, V TValue] struct {
	Key   T
	Value V
}

type MyChanMap[T Tindex, V TValue] struct {
	// *MyChan //
	Read <-chan *MapValue[T, V] //interface{} // 只读通道  为 channel 通道创建一个 按索引查看的方法
	//all   chan map[int]interface{}   // 可读可写
	Input chan<- *MapValue[T, V] //interface{} // 只写通道
	// maxsize int
	Cache   chan *MapValue[T, V]
	MaxSize int
}

/// 从 chan map 获取 指定编号 的 MapValue
func (self *MyChanMap[T, V]) Get(key int) *MapValue[T, V] {
	/// 默认返回第一个
	if len(self.Read) > 0 {
		Mutex.Lock()
		defer Mutex.Unlock()
		var vs *MapValue[T, V]
		for i := 0; i < len(self.Read); i++ {
			vs = <-self.Read
			if i == key {
				return vs
			}
		}
		return vs
	}
	return nil
}

//// 返回chan的所有key
func (self *MyChanMap[T, V]) Keys() ([]T, []V) {
	var keys []T
	var values []V
	totalRead := len(self.Read)
	for i := 0; i < totalRead; i++ {
		iV := self.Get(0)
		keys = append(keys, iV.Key)
		values = append(values, iV.Value)
		self.Cache <- iV
	}
	self.PushBackCache()
	return keys, values
}

///存入一个 Mapvalue 结构体到 ChanMap 队列
func (self *MyChanMap[T, V]) Put(Value *MapValue[T, V]) bool {
	if len(self.Input) < self.MaxSize {
		self.Input <- Value
		return true
	}
	return false
}

func (self *MyChanMap[T, V]) PushBackCache() bool {
	//// 回填 cache 队列到 Input队列，如果cache队列为空 或 回填失败，
	//返回 false
	c := len(self.Cache)
	Logger.Println("length of cache:", c)
	if c > 0 {
		for i := 0; i < c; i++ {
			vm := <-self.Cache
			self.Input <- vm

			Logger.Printf("push back cache:%+v,length of read:%d\n", vm, len(self.Read))
		}
		return true
	}
	return false
}

//// 是否当前通道 已存满
func (self *MyChanMap[T, V]) IsChanFull() bool {
	return self.MaxSize == len(self.Read)
}

//// 返回一个新的 ChanMap，大小是当前大小的 2倍，并包含当前ChanMap的元素
func (self *MyChanMap[T, V]) NewChildChanMap() *MyChanMap[T, V] {
	mapVa := make(chan *MapValue[T, V], self.MaxSize*2)
	cacheVa := make(chan *MapValue[T, V], self.MaxSize*2)
	new_chanmaps := &MyChanMap[T, V]{
		Read:    mapVa,
		Input:   mapVa,
		Cache:   cacheVa,
		MaxSize: self.MaxSize * 2,
	}
	Logger.Println("Make New ChildChan of length  self read:", len(self.Read))
	if len(self.Read) > 0 {
		/// 将 当前Read 通道的数据 存入新建的 Input通道
		totalLen := len(self.Read)
		for i := 0; i < totalLen; i++ {
			item := self.Get(0)
			new_chanmaps.Input <- item
			self.Cache <- item
		}
		self.PushBackCache()
	}
	return new_chanmaps
}

////parents 取原 map 列表 元素 第一个以后的全部，返回一个新的MyChanMap 结构体通道
func (self *MyChanMap[T, V]) ParentsChan() *MyChanMap[T, V] {
	//// 来自 map[1:]的新 ChanMap, 最大大小与原map一致
	mapVa := make(chan *MapValue[T, V], self.MaxSize)
	cacheVa := make(chan *MapValue[T, V], self.MaxSize)
	parents_chanmaps := &MyChanMap[T, V]{
		Read:    mapVa,
		Input:   mapVa,
		Cache:   cacheVa,
		MaxSize: self.MaxSize,
	}

	// parents/
	popitem := self.Get(0)
	Logger.Println("pop item:", popitem)
	if len(self.Read) > 0 {
		/// 将 当前Read 通道的数据 存入新建的 Input通道
		totalLen := len(self.Read)
		for i := 0; i < totalLen; i++ {
			item := self.Get(0)
			parents_chanmaps.Input <- item
			self.Cache <- item
		}
		self.PushBackCache()
	}
	Logger.Println("Make parents of pop self read:", len(self.Read))
	return parents_chanmaps
}

/// 插入到 chan 通道的哪一个位置， index 为 chan中的位置，
func (self *MyChanMap[T, V]) Insert(index int, node *MapValue[T, V]) bool {
	/// index 从 0 开始，也就是 index 最大为 self.MaxSize - 1
	if len(self.Input) == 0 || index >= len(self.Input) {
		/// 队列为空，或 插入位置 大于等于 Input队列长度，则直接添加到 尾部
		self.Put(node)
		return true
	} else if self.IsChanFull() == true || index >= self.MaxSize {
		//// 该chan 已满，无法继续存入
		Logger.Fatal("self.input is full of maxsize:", len(self.Input), index)
		return false
	} else {
		totalRead := len(self.Read)
		//// 处理 index 插入位置 在 chan队列中 的场景
		for i := 0; i < totalRead; i++ {
			/// 遍历每个元素
			nr := self.Get(0)
			if index == i {
				self.Cache <- node
				self.Cache <- nr
			} else {
				self.Cache <- nr
			}
		}
		self.PushBackCache()
		return true
	}

}

//// Set 更新 maps的第一个元素，支持嵌套子类操作？
func (self *MyChanMap[T, V]) SetItem(key T, value V) *MyChanMap[T, V] {
	// self.Cache.
	self.Insert(0, &MapValue[T, V]{Key: key,
		Value: value})
	return self
}

//// Delete 删除 maps的 某个位置的元素
func (self *MyChanMap[T, V]) Delete(index int) bool {
	if len(self.Read) == 0 || index > len(self.Read) {
		return false
	} else {
		totalRead := len(self.Read)
		//// 处理 index 插入位置 在 chan队列中 的场景
		for i := 0; i < totalRead; i++ {
			/// 遍历每个元素
			nr := self.Get(0)
			if index == i {
				Logger.Println("this item will not pushback", nr)
			} else {
				self.Cache <- nr
			}
		}
		self.PushBackCache()
		return true
	}

}

////////////////////////////////////////////////
func InMapsValue[T Tindex, V TValue](self *MyChanMap[T, V], a V) bool {
	//// 查找某个值 是否在 字典元素中
	for {
		mps := self.Get(0)
		if mps == nil {
			return false
		} else {
			va := &MapValue[T, V]{
				Key:   mps.Key,
				Value: mps.Value,
			}
			if a == mps.Value {
				self.Cache <- va
				/// 返回前 回填数据
				// self.PushBackCache()
				// return true
			} else {
				self.Cache <- va
			}
			///  清理缓存队列
			Logger.Println("push back cache at InmapValue:", len(self.Cache))
			self.PushBackCache()
			return true
		}
	}
}

//// 存入一个 自定义字典结构体到 chan map
func RetMaps[T Tindex, V TValue](MyC MyChanMap[T, V], a T, v V) *MyChanMap[T, V] {
	var Maps1 = &MapValue[T, V]{
		Key:   a,
		Value: v,
	}
	// ti := reflect.ValueOf(v).String()
	MyC.Input <- Maps1
	return &MyC
}

// 构造MyChanMap的函数
func MakeMyChans[T Tindex, V TValue](maxsize int) *MyChanMap[T, V] {
	// 只读 只写 分开
	fmt.Println("make chan and receuve maxsize:", maxsize)
	// maxsize := 10
	var MyChan = make(chan *MapValue[T, V], maxsize)
	var MyCache = make(chan *MapValue[T, V], maxsize)
	// MyC.Input <- Maps1
	// var ret_mychan MyChanMap
	ret_mychan := MyChanMap[T, V]{
		Read:    MyChan,
		Input:   MyChan,
		Cache:   MyCache,
		MaxSize: maxsize,
	}
	////  必须存入一个初始化数据 到 chan
	ret_mychans := RetMaps(ret_mychan, T(0), V("Hello"))
	return ret_mychans
}

func main() {
	// 调用时指明类型
	rmc := MakeMyChans[int, string](10)
	// rmc.Put(&MapValue[int, string]{Key: 0, Value: "Hello"})
	Logger.Printf("rmc:%+v\n", rmc)
	Logger.Println(rmc.Get(1)) //// 已经取出了 预置的 数据
	Logger.Printf("cache read:%+v\n", rmc.Read)
	for i := 0; i < rmc.MaxSize; i++ {
		rmc.Put(&MapValue[int, string]{Key: i,
			Value: fmt.Sprintf("value-%v", i)})
		// RetMaps[int,string](rmc, i, "HelloIo")
		Logger.Println("read length:", len(rmc.Read))
	}
	Logger.Println("chan read length:", len(rmc.Read))
	Logger.Println(InMapsValue(rmc, "Hello")) //// 预置的数据 Hello 已经不存在
	Logger.Println(InMapsValue(rmc, "value-9"))
	Logger.Println(len(rmc.Read))
	rmc.PushBackCache()
	// for len(rmc.Read) > 0 {
	// 	read_info := <-rmc.Read
	// 	Logger.Printf("read:%+v\n", read_info)
	// }

	// Logger.Println(rmc.Keys())

	/// 返回一个 两边当前chan 通道大小的 chan
	new_child := rmc.NewChildChanMap()
	Logger.Println("new_child read length:", len(new_child.Read), "maxsize:", new_child.MaxSize, "origin length and maxsize:", len(rmc.Read), rmc.MaxSize)

	Logger.Println("New Chan Full?", new_child.IsChanFull(), "origin Chan full?", rmc.IsChanFull())

	Parentschan := rmc.ParentsChan()
	Logger.Println("Parents chans:", len(Parentschan.Read), Parentschan.MaxSize)
	Parentschan.Insert(0, &MapValue[int, string]{Key: 110,
		Value: fmt.Sprintf("value-101")})
	Logger.Println("Parents chans:", len(Parentschan.Read), Parentschan.MaxSize)
	Logger.Printf("first item:%+v\n", Parentschan.Get(0))
	Logger.Println(rmc.Delete(3))
	Logger.Println("after delete from chans:", len(Parentschan.Read), Parentschan.MaxSize)
	Logger.Println(Parentschan.Keys())
}
