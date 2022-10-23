package main

/*
gorm + nats + mysql

write to db

translation lostions.
*/

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	// _ "github.com/go-sql-driver/mysql" //并不需要把整个包都导入进来，仅仅是是希望它执行init()函数而已。这个时候就可以使用 import _ 引用该包
)

const (
	host     = "192.168.30.131"
	port     = 3306
	username = "admin"
	password = "admin2022.post"
	dbname   = "mystate"
)

var (
	DSN  = "admin:admin2022.post@tcp(192.168.30.131:3306)/mystate?multiStatements=true&allowNativePasswords=false&checkConnLiveness=true&maxAllowedPacket=0&charset=utf8&parseTime=True&loc=Local&timeout=10s&readTimeout=30s&writeTimeout=60s"
	Logg = log.New(os.Stderr, "INFO -", 18)
	db   *gorm.DB
	err  error
)

func init() {
	db, err = gorm.Open(mysql.Open(DSN), &gorm.Config{
		// DefaultStringSize:         256,
		// DisableDatetimePrecision:  true,  // 禁用datetime精度， mysql>5.6
		// DontSupportRenameIndex:    true,  //重命名索引时采用删除并新建的方式 Mysql>5.7才支持重命名索引
		// DontSupportRenameColumn:   true,  // 用 change重命名 列，mysql > 8
		// SkipInitializeWithVersion: false, // 根据当前Mysql版本自动配置
	})
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}
}

// func tools() *grom.DB {
// 	// 定义一个工具函数
// 	var _db *gorm.DB

// 	var err error
// 	_db, err = gorm.Open("mysql", DSN)
// 	if err != nil {
// 		panic("连接失败")
// 	}
// 	/// 数据库连接池
// 	_db.DB().SetMaxOpenConns(100) // 设置数据库连接池最大连接数
// 	_db.DB().SetMaxIdleConns(20)  /// 连接池最大允许空闲连接数，如果没有sql任务需要执行的连接数大于20，超过的连接将被连接池关闭
// 	return _db
// }
type User struct {
	// gorm.Model
	Id    int
	Name  string
	Price float64
	// CreditCard CreditCard
	// CreatedAt  time.Time
	// UpdatedAt  time.Time
	// DeletedAt  time.Time
}

/*
CREATE TABLE IF NOT EXISTS `users` (`id` INT UNSIGNED, `name` VARCHAR(40) NOT NULL, `price` DECIMAL(15,2) NOT NULL  DEFAULT 0, PRIMARY KEY (`id`))ENGINE=InnoDB DEFAULT CHARSET=utf8;
*/
type CreditCard struct {
	gorm.Model
	Number string
	UserID uint
}
type Participant struct {
	Participantid int
	Name          string
	Email         string
	Cash          float64
}

// 迁移工具
/*
GORM 的 AutoMigrate 在大多数情况下都工作得很好，但如果您正在寻找更严格的迁移工具，GORM 提供一个通用数据库接口，可能对您有帮助。
// returns `*sql.DB`
db.DB()
*/

//// 执行数据库操作前执行 的钩子，检查如果有错误，则不执行 任何插入
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	/*
			// 开始事务
		BeforeSave
		BeforeCreate
		// 关联前的 save
		// 插入记录至 db
		// 关联后的 save
		AfterCreate
		AfterSave
		// 提交或回滚事务

	*/
	if u.Id > 200 {
		Logg.Println("have bigger user:", u)
		return errors.New("invalid id number.")
	}
	// else {
	// 	tx.Create(u)
	// }
	return
}

func BasicFunc() {
	//1 连接MYSQL, 获得DB类型实例，用于后面的数据库读写操作。

	t := time.Now()
	user := User{Id: 112, Name: "Lucy22", Price: 1300.00} //CreditCard: CreditCard{Number: "212831212"}}
	///2  插入一个数据
	// result := db.Create(&user)

	Logg.Printf("create result:%+v, time:%v\n", user, t)

	/// 3 插入数据时指定字段
	// db.Select("Id", "Name", "Price").Create(&User{3, "Bob", 2810.33})
	// gorm  连接池
	// dbs := tools()

	//  4 创建记录并更新未给出的字段 ？？？报错
	// db.Omit("Id", "Name", "Price").Create(&User{4, "Bob1", 2810.33})

	///  5 批量插入 将一个slice 传递给 create 方法，将切片数据传递给Create方法，GORM将生成一个单一SQL插入所有数据
	//// 并回填主键的值，钩子方法也会被调用
	// var users []User

	// for i := 0; i < 10; i++ {
	// 	name := fmt.Sprintf("Lily-%v", i)
	// 	users = append(users, User{i + 10, name, 1002.00})
	// }
	// db.Create(&users)
	// Logg.Printf("%+v\n", users)
	// for i := 0; i < 28; i++ {
	// 	name := fmt.Sprintf("Lily1-%v", i)
	// 	users = append(users, User{i + 80, name, 2002.00})
	// }
	/// 插入时指定 数量， 无效。。。
	// db.CreateInBatches(users, 8)

	//6  跳过钩子的使用 ，false 表示不跳过钩子的检查
	// db.Session(&gorm.Session{SkipHooks: false}).Create(&users)
	/// 创建钩子 GORM 允许用户定义的钩子有 beforeSave， beforeCreate，after Save，AfterCreate 创建记录时将调用这些钩子方法，参考 Hooks https://learnku.com/docs/gorm/v2/hooks

	//// 7 根据Map创建 Gorm 支持
	// db.Model(&User{}).Create(map[string]interface{}{
	// 	"Id": 111, "Name": "NameBob", "Price": 291.21,
	// })

	//8  从 map插入 map 创建记录时，association 不会被调用，且主键也不会自动填充
	db.Model(&User{}).Create([]map[string]interface{}{
		{"Id": 112, "Name": "Bobs2", "Price": 1281.2},
		{"Id": 113, "Name": "Bobs3", "Price": 1281.2},
	})

	//// 9 关联插入 必须要有 (`created_at`,`updated_at`,`deleted_at`
	// db.Create(&User{
	// 	Name: "BobLucy",

	////10 upsert 及冲突
	/// 有冲突时 不做任何操作, 不会插入Id为 112 的lucy
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&user)

	////11 当id有冲突时，将库中有冲突的数据，更新指定列为默认值
	/*
			// MERGE INTO "users" USING *** WHEN NOT MATCHED THEN INSERT *** WHEN MATCHED THEN UPDATE SET ***; SQL Server
		// INSERT INTO `users` *** ON DUPLICATE KEY UPDATE ***; MySQL

	*/
	// db.Clauses(clause.OnConflict{
	// 	Columns:   []clause.Column{{Name: "id"}},
	// 	DoUpdates: clause.Assignments(map[string]interface{}{"id": 0}),
	// }).Create(&user)

	///12 当 id 冲突，更新指定列为新值
	/*
			// MERGE INTO "users" USING *** WHEN NOT MATCHED THEN INSERT *** WHEN MATCHED THEN UPDATE SET "name"="excluded"."name"; SQL Server
		// INSERT INTO "users" *** ON CONFLICT ("id") DO UPDATE SET "name"="excluded"."name", "age"="excluded"."age"; PostgreSQL
		// INSERT INTO `users` *** ON DUPLICATE KEY UPDATE `name`=VALUES(name),`age=VALUES(age); MySQL
	*/
	// db.Clauses(clause.OnConflict{
	// 	Columns:   []clause.Column{{Name: "id"}},
	// 	DoUpdates: clause.AssignmentColumns([]string{"name", "price"}),
	// }).Create(&user)

	////13 强制更新全部列为新数据
	/// 这将强制更新冲突的内容为 新的数据
	db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&user)

	////// 14 链式调用 链式方法，Finisher 方法，新建会话方法
	///链式方法是将 Clauses 修改或添加到当前 Statement 的方法，例如：
	// Where, Select, Omit, Joins, Scopes, Preload, Raw (Raw
	/*
		Finishers 是会立即执行注册回调的方法，然后生成并执行 SQL，比如这些方法：
		Create, First, Find, Take, Save, Update, Delete, Scan, Row, Rows…
	*/
	// db.Where("name = 'jack' ").First(&user)
	db.Find(&user)

	// 安全的使用新初始化的 *gorm.DB
	for i := 0; i < 100; i++ {
		go db.Where("id > 1").First(&user)
	}

	tx := db.Where("name = ?", "jack")
	// 不安全的复用 Statement
	for i := 0; i < 100; i++ {
		go tx.Where("id > 10").First(&user)
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	ctxDB := db.WithContext(ctx)
	// 在 `新建会话方法` 之后是安全的
	for i := 0; i < 100; i++ {
		go ctxDB.Where("price > 100").First(&user)
	}

	ctx2, _ := context.WithTimeout(context.Background(), time.Second)
	ctxDB2 := db.Where("name = ?", "jack").WithContext(ctx2)
	// 在 `新建会话方法` 之后是安全的
	for i := 0; i < 100; i++ {
		go ctxDB2.Where("id > 100").First(&user) // `name = 'jinzhu'` 会应用到查询中
	}

	tx2 := db.Where("name = ?", "lucy").Session(&gorm.Session{})
	// 在 `新建会话方法` 之后是安全的
	for i := 0; i < 100; i++ {
		go tx2.Where("id > 19").First(&user) // `name = 'jinzhu'` 会应用到查询中
	}

	// 	//CreditCard: CreditCard{Number: "87312312333"},
	// })
	// 执行查询
	// u := User{}
	// 指定生成 sql，SELECT * FROM `users` WHERE (username = 'tizi365') LIMIT 1
	// dbs.Where("username = ?", "tizi365").First(&u)
}

func AmountGreaterThan1000(db *gorm.DB) *gorm.DB {
	return db.Where("id > ?", 10)
}

func PaidWithCreditCard(db *gorm.DB) *gorm.DB {
	return db.Where("name = ?", "jack")
}

func PaidWithCod(db *gorm.DB) *gorm.DB {
	return db.Where("price > ?", 1000)
}

func OrderStatus(status []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("name IN (?)", status)
	}
}

///// 分页
// func Paginate(r *http.Request) func(db *gorm.DB) *gorm.DB {
// 	return func(db *gorm.DB) *gorm.DB {
// 		page, _ := strconv.Atoi(r.Query("page"))
// 		if page == 0 {
// 			page = 1
// 		}

// 		pageSize, _ := strconv.Atoi(r.Query("page_size"))
// 		switch {
// 		case pageSize > 100:
// 			pageSize = 100
// 		case pageSize <= 0:
// 			pageSize = 10
// 		}

// 		offset := (page - 1) * pageSize
// 		return db.Offset(offset).Limit(pageSize)
// 	}
// 	///// 分页
// 	db.Scopes(Paginate(r)).Find(&user)
// 	db.Scopes(Paginate(r)).Find(&articles)
// }

func main() {
	a := db.Scopes(AmountGreaterThan1000, PaidWithCreditCard).Find(&User{})
	// 查找所有金额大于 1000 的信用卡订单

	b := db.Scopes(AmountGreaterThan1000, PaidWithCod).Find(&User{})
	// 查找所有金额大于 1000 的 COD 订单

	c := db.Scopes(AmountGreaterThan1000, OrderStatus([]string{"jack", "Lucy"})).Find(&User{})
	// 查找所有金额大于1000 的已付款或已发货订单
	Logg.Printf("%+v\n, %+v\n, %+v\n", a, b.Config, c.Statement)

}
