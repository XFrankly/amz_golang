package main

// "github.com/jcelliott/lumber"  /// 创建记录器

/*
初始化创建 jsondatabase
go mod init gojsondatabase
*/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
)

const Version = "1.0.1"

var (
	Logger = log.New(os.Stderr, "INFO -", 18)
	dir    = "./"
)

type (
	Driver struct {
		mutex   sync.Mutex
		mutexes map[string]*sync.Mutex
		dir     string
		lg      *log.Logger
	}
)

type Options struct {
	Logg *log.Logger
}

//// 创建一个db 链接
func New(dir string, options *Options) (*Driver, error) {
	////
	dir = filepath.Clean(dir)
	opts := Options{}
	if options != nil {
		opts = *options
	}

	if opts.Logg == nil {
		opts.Logg = log.New(os.Stderr, "INFO -", 18)
	}

	dirver := Driver{
		dir:     dir,
		mutexes: make(map[string]*sync.Mutex),
		lg:      opts.Logg,
	}

	if _, err := os.Stat(dir); err == nil {
		//// 是否已经存在数据库文件
		opts.Logg.Printf("Using '%s' (database already exists)\n", dir)
		return &dirver, nil
	}

	//// 数据库文件不存在，创建一个
	opts.Logg.Printf("Creating database at '%s'...\n ", dir)
	return &dirver, os.MkdirAll(dir, 0755)
}

//// 写入db
func (d *Driver) Write(collection, resource string, v interface{}) error {
	if collection == "" {
		return fmt.Errorf("Missing collection - no place to save record!(no name)")
	}
	mutex := d.getOrCreateMutex(collection)
	mutex.Lock() //// 使用数据库时，加互斥锁
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, collection)
	d.lg.Println("write to file path:", d.dir, collection, dir)
	fnlPath := filepath.Join(dir, resource+".json")
	tmpPath := fnlPath + ".tmp"

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	//// []byte
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}
	b = append(b, byte('\n'))
	if err := ioutil.WriteFile(tmpPath, b, 0644); err != nil {
		return err
	}
	return os.Rename(tmpPath, fnlPath)
}

//// 读取一个db
func (d *Driver) Read(collection, resource string, v interface{}) error {
	if collection == "" {
		return fmt.Errorf("Missing collection - no place to read!")
	}
	if resource == "" {
		return fmt.Errorf("Missing resource - unable to save record(no name)")
	}

	record := filepath.Join(d.dir, collection, resource)
	if _, err := stat(record); err != nil {
		return err
	}
	b, err := ioutil.ReadFile(record + ".json")
	if err != nil {
		return err
	}

	return json.Unmarshal(b, &v)
}

/// 读取全部
func (d *Driver) ReadAll(collection string) ([]string, error) {
	if collection == "" {
		return nil, fmt.Errorf("Missiong collection - unable to read.")
	}
	dir := filepath.Join(d.dir, collection)
	if _, err := stat(dir); err != nil {
		return nil, err
	}
	files, _ := ioutil.ReadDir(dir)
	var records []string
	for _, f := range files {
		b, err := ioutil.ReadFile(filepath.Join(dir, f.Name()))
		if err != nil {
			return nil, err
		}
		records = append(records, string(b))
	}
	return records, nil
}

/// 删除一个数据
func (d *Driver) Delete(collection string, resource string) error {
	path := filepath.Join(collection, resource)
	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, path)
	switch fi, err := stat(dir); {
	//// 根据 dir 的属性来判断 和操作
	case fi == nil, err != nil:
		return fmt.Errorf("unable to find file or directory naned %v \n", path)
	case fi.Mode().IsDir():
		return os.RemoveAll(dir)
	case fi.Mode().IsRegular():
		return os.RemoveAll(dir + ".json")
	}
	return nil
}

//// 获取或创建一个互斥锁
func (d *Driver) getOrCreateMutex(collection string) *sync.Mutex {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	m, ok := d.mutexes[collection]
	if !ok {
		m = &sync.Mutex{}
		d.mutexes[collection] = m
	}
	return m
}

////检查path 路径
func stat(path string) (fi os.FileInfo, err error) {
	////
	if fi, err = os.Stat(path); os.IsNotExist(err) {
		fi, err = os.Stat(path + ".json")
	}
	return
}

type Address struct {
	City    string
	State   string
	Country string
	Pincode json.Number /// 邮政编码
}

type User struct {
	Name    string
	Age     json.Number
	Contact string
	Company string
	Address Address
}

func putAndShow() {
	db, err := New(dir, nil)
	if err != nil {
		Logger.Println("Error", err)
	}
	Logger.Println("this is a json database base on golang.")

	//// 结构数组
	employees := []User{
		{"John", "23", "9832912", "New Tech", Address{
			"U.S.A", "Online", "EU", "230001",
		}},
		{"Paul", "39", "9832913", "New Tech", Address{
			"U.S.A", "Online", "Asia", "230039",
		}},
		{"Frank", "29", "9832914", "New Tech", Address{
			"U.S.A", "Online", "EU", "230029",
		}},
		{"Tony", "31", "9832915", "New Tech", Address{
			"U.S.A", "Online", "EU", "230031",
		}},
		{"Dave", "28", "9832916", "New Tech", Address{
			"Asia", "Online", "EU", "230082",
		}},
		{"Neo", "24", "9832917", "New Tech", Address{
			"Asia", "Remote", "Asia", "230004",
		}},
	}

	for _, emp := range employees {
		Logger.Printf("emp struct:%+v\n", emp)
		//// 链接到 user表
		db.Write("users", emp.Name, User{
			Name:    emp.Name,
			Age:     emp.Age,
			Contact: emp.Contact,
			Company: emp.Company,
			Address: emp.Address,
		})
	}

	records, err := db.ReadAll("users")
	if err != nil {
		Logger.Println("error", err)
	}
	Logger.Println(records)

	allusers := []User{}
	for _, f := range records {
		empFound := User{}
		if err := json.Unmarshal([]byte(f), &empFound); err != nil {
			Logger.Println("error", err)
		}
		allusers = append(allusers, empFound)
	}
	Logger.Println(allusers)

	////// 删除某个信息
	if err := db.Delete("users", "John"); err != nil {
		Logger.Println("Error", err)
	}

	///// 删除 users 全部信息
	// if err := db.Delete("users", ""); err != nil {
	// 	Logger.Println("Error", err)
	// }
}

func main() {
	putAndShow()
}
