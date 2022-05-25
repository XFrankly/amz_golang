package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	// 是database/sql/driver接口的实现 拥有database/sql全部API
	// "github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql" //并不需要把整个包都导入进来，仅仅是是希望它执行init()函数而已。这个时候就可以使用 import _ 引用该包
)

/*
轻巧快速
本机 Go 实现。没有 C 绑定，只有纯 Go
通过 TCP/IPv4、TCP/IPv6、Unix 域套接字或自定义协议的连接
自动处理断开的连接
自动连接池（通过 database/sql 包）
支持大于 16MB 的查询
全力sql.RawBytes支持。
LONG DATA准备好的语句中的智能处理
LOAD DATA LOCAL INFILE通过文件许可名单和io.Reader支持获得安全支持
可选time.Time解析
可选占位符插值
*/
const (
	host     = "192.168.30.131"
	port     = 3306
	username = "admin"
	password = "admin2022.post"
	dbname   = "mystate"
)

/*
连接池
db.Ping() 调用完毕后会马上把连接返回给连接池。
db.Exec() 调用完毕后会马上把连接返回给连接池，但是它返回的Result对象还保留这连接的引用，当后面的代码需要处理结果集的时候连接将会被重用。
db.Query() 调用完毕后会将连接传递给sql.Rows类型，当然后者迭代完毕或者显示的调用.Clonse()方法后，连接将会被释放回到连接池。
db.QueryRow()调用完毕后会将连接传递给sql.Row类型，当.Scan()方法调用之后把连接释放回到连接池。
db.Begin() 调用完毕后将连接传递给sql.Tx类型对象，当.Commit()或.Rollback()方法调用后释放连接。

*/
var (
	db, db2 *sql.DB
	// dsn := "user:password@tcp(127.0.0.1:3306)/sql_test?charset=utf8mb4&parseTime=True"
	DSN  = "admin:admin2022.post@tcp(192.168.30.131:3306)/mystate?multiStatements=true&allowNativePasswords=false&checkConnLiveness=true&maxAllowedPacket=0"
	DSN1 = "admin:admin2022.post@tcp(192.168.30.131:3306)/mystate?multiStatements=true&allowNativePasswords=false&checkConnLiveness=true&maxAllowedPacket=100"

	ctx  context.Context
	Logg = log.New(os.Stderr, "INFO -", 18)
	// Capture connection properties.
	// 使用 MySQL 驱动程序Config - 和类型FormatDSN -
	// 以收集连接属性并将它们格式化为连接字符串的 DSN。
	// 该Config结构使代码比连接字符串更容易阅读。
	// cfg = mysql.Config{
	// 	// User:   os.Getenv("DBUSER"),  /// 从OS环境获取
	// 	// Passwd: os.Getenv("DBPASS"), /// 从OS环境获取
	// 	User:   username,
	// 	Passwd: password,
	// 	Net:    "tcp",
	// 	Addr:   "192.168.30.131:3306",
	// 	DBName: dbname,
	// }
	// cfg.MultiStatements = true
	// Get a database handle.
	err, err2 error
)

type User struct {
	Name string
	Cash float64
}

func InitDb(conn *sql.DB) bool {
	/// 初始化  DROP ... CASCADE 强制删除依赖，将导致其他表的依赖 列 混乱 context.Background(),
	sql_strs := []string{
		"CREATE DATABASE IF NOT EXISTS mystate;",
		"DROP TABLE IF EXISTS trunk ;",
		"DROP TABLE IF EXISTS participant;",
		"CREATE TABLE IF NOT EXISTS `trunk`(" +
			"`trunkid` INT UNSIGNED AUTO_INCREMENT," +
			"`participantid` INT NOT NULL," +
			"`name` VARCHAR(40) NOT NULL," +
			"`price` DECIMAL(15,2) NOT NULL," +
			"`description` VARCHAR(200) NOT NULL," +
			"PRIMARY KEY ( `trunkid` )" +
			" )ENGINE=InnoDB DEFAULT CHARSET=utf8;",
		"CREATE TABLE  IF NOT EXISTS  `participant` (" +
			"`participantid` INT UNSIGNED AUTO_INCREMENT," +
			"`name` VARCHAR(40) NOT NULL," +
			"`email` VARCHAR(200) NOT NULL," +
			"`cash` DECIMAL(15,2) NOT NULL," +
			"PRIMARY KEY (participantid)" +
			")ENGINE=InnoDB DEFAULT CHARSET=utf8;",
		"INSERT INTO participant (name, email, cash) VALUES ('Tom', 'Admin@example.com', '1100.00');",
		"INSERT INTO participant (name, email, cash) VALUES ('Jack', 'User@example.com', '1150.00');",
		"INSERT INTO trunk (participantid, name,price, description) VALUES (1,'Linux CD', '1.00', 'Complete OS on a CD'); ",
		"INSERT INTO trunk (participantid, name,price, description) VALUES (2,'ComputerABC', '12.90', 'a book about OS computer!');",
		"INSERT INTO trunk (participantid, name,price, description) VALUES (2,'Magazines', '6.90', 'Stack of Computer Magezines computer!');",
	}
	sql := ""
	for _, strs := range sql_strs {
		sql = fmt.Sprintf("%v%v", sql, strs) //sql + strs
	}
	// Logg.Println("full sqls:", sql)

	rst1, err := conn.Exec(sql)
	Logg.Println("full sql exec result:", rst1, err)
	// result, err := conn.Exec("INSERT INTO User (title, artist, price) VALUES (?, ?)", alb.Name, alb.Cash)
	if err != nil {
		Logg.Println("failed to panic", err)
		panic(err)
	}
	return true
}

// 创建数据库句柄
// 制作db一个全局变量可以简化这个例子。
//在生产环境中，您会避免使用全局变量，
//例如将变量传递给需要它的函数或将其包装在结构中。
// 声明一个db类型的变量*sql.DB。这是您的数据库句柄。
func init() {

	// 调用sql.Open 初始化db变量，传递 FormatDSN.
	//检查来自 的错误sql.Open。
	//例如，如果您的数据库连接细节格式不正确，它可能会失败。
	//为了简化代码，您调用log.Fatal结束执行并将错误打印到控制台。
	//在生产代码中，您会希望以更优雅的方式处理错误。
	// DSN1 := cfg.FormatDSN() // user:passwd@tcp(host:port)/database?multiStatements=true
	// DSN1 := fmt.Sprintf("%v:%v@tcp(%v:3306)/%v?multiStatements=true", username, password, dbname)

}

type Computers struct {
	Trunkid       int
	Participantid int
	Name          string
	Price         float64
	Description   string
}

// Create a helper function for preparing failure results.
func fail(err error) error {
	return fmt.Errorf("CreateOrder: %v", err)
}

//// 脏读，读未提交和 读已提交 隔离级
func dirtyRead(db1, db2 *sql.DB, isolationLevel string, pointName string) {
	//// 此连接开始事务 对应两个级别
	/// SET TRANSACTION ISOLATION LEVEL READ COMMITTED
	/// SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED
	//// 进行脏读，则步骤2和3的读取结果将相同。
	/// 但是由于更改是在事务内部进行的，因此在提交之前在外部是不可用的，其他连接 总是读到旧的数据
	// tx, err := conn1.Begin(ctx)
	/*
				// /// 模拟Tom 从 Jack 购买了一个设备ComputerABC
				START TRANSACTION;
				UPDATE participant SET cash=cash-12.99 WHERE participantid=1;
				UPDATE participant SET cash=cash+12.99 WHERE participantid=2;
				UPDATE trunk SET participantid = 1 WHERE name = 'ComputerABC' AND participantid=2;
				SAVEPOINT %s;
				SELECT * FROM participant;
				COMMIT;
				`
			READ-UNCOMMITTED
		提交事务前:连接内部值Tom cash from main transaction after update: 1087.010000
		提交事务前:其他连接Tom participant from conn2 : 1100.000000
			READ-COMMITTED
		提交事务前:连接内部值Tom cash from main transaction after update: 1087.010000
		提交事务前:其他连接Tom participant from conn2 : 1087.010000
			REPEATABLE-READ
		提交事务前:连接内部值Tom cash from main transaction after update: 1087.010000
		提交事务前:其他连接Tom participant from conn2 : 1100.000000
			SERIALIZABLE
		提交事务前:连接内部值Tom cash from main transaction after update: 1087.010000
		提交事务前:其他连接Tom participant from conn2 : 1087.010000
	*/

	Logg.Println("create sql.DB translation with level:.", isolationLevel)
	levels, err1 := db1.Exec("SET GLOBAL transaction_isolation= ? ;", isolationLevel)
	if err1 != nil {
		Logg.Println("dirtyRead conn1 set level failed:", levels, err1)
		panic(fail(err1))
	}
	ptx, err := db1.Begin()
	if err != nil {
		msg := fmt.Sprintf("dirtyRead start translation failure wuth db1 connection.: %+v\n", err)
		panic(msg)
	}

	//// 执行交易 从Tom 扣除12.99
	Logg.Println("start exec.")
	pconn, perr := ptx.Exec("UPDATE participant SET cash=cash-12.99 WHERE participantid=1;")
	if perr != nil {
		Logg.Printf("dirtyRead conn1 Failed to update Tom cash in tx: %v\n", perr)
		panic(perr)
	}
	//// 执行交易 给Jack 增加12.99
	pconn2, perr2 := ptx.Exec("UPDATE participant SET cash=cash+12.99 WHERE participantid=2;")
	if perr2 != nil {
		Logg.Printf("dirtyRead conn1 Failed to update Tom cash in tx: %v\n", perr2)
		panic(perr2)
	}

	//// 执行交易 把物品从 Jack交还给Tom
	pconn3, perr3 := ptx.Exec("UPDATE trunk SET participantid = 1 WHERE name = 'ComputerABC' AND participantid=2;")
	if perr != nil {
		Logg.Printf("Failed to trunk SET  in tx: %v\n", perr3)
		panic(perr)
	}
	Logg.Println("success update result:", pconn, pconn2, pconn3)
	/////保存事务执行点
	sql_savepoint := fmt.Sprintf("SAVEPOINT  %v", pointName)
	ptx.Exec(sql_savepoint)
	/// 事务提交前 当前事务连接 和 其他连接 查询对比
	var balance float64
	row := ptx.QueryRow("SELECT cash FROM participant WHERE name='Tom'")
	row.Scan(&balance)
	Logg.Printf("提交事务前:连接内部值Tom cash from main transaction after update: %f\n", balance) /// 事务连接内部值 tom cash 1087.01

	var othercash float64
	row1 := db2.QueryRow("SELECT cash FROM participant WHERE name='Tom'")
	row1.Scan(&othercash)
	/// 其他连接的 脏读
	Logg.Printf("提交事务前:其他连接Tom participant from conn2 : %f\n", othercash) /// 其他连接tom cash 1100
	Logg.Printf("same eq?:%+v\n", balance+12.99 == othercash)

	if err := ptx.Commit(); err != nil {
		Logg.Printf("dirtyRead conn1  Failed to commit: %v\n", err)
		panic(err)
	}

	/////ComputerABC
	var participantOwn int
	row2 := db2.QueryRow("SELECT  participantid FROM trunk WHERE name = 'ComputerABC'")
	row2s, e_ := db2.Query("SELECT  * FROM trunk WHERE name = 'ComputerABC'")
	if e_ != nil {
		Logg.Println(e_)
		panic(e_)
	}
	row2.Scan(&participantOwn)
	Logg.Printf("提交事务后:Tom participantOwn from conn2 trunk: %+v\n", participantOwn)

	// var participantOwns string
	// 非常重要：关闭rows释放持有的数据库链接
	defer row2s.Close()
	// 查询归属 循环读取结果集中的数据
	for row2s.Next() {
		var cu Computers
		err := row2s.Scan(&cu.Trunkid, &cu.Participantid, &cu.Name, &cu.Price, &cu.Description)
		if err != nil {
			Logg.Printf("dirtyRead conn1  scan failed, err:%v\n", err)
			panic(err)
		}
		Logg.Printf("归属 :id:%d name:%s Price:%f\n", cu.Participantid, cu.Name, cu.Price)
	}
	// Logg.Println(len(*row2s))

}

///// 重复读 当一个事务重新读取前面读取过的数据时，发现该数据已经被另一个已提交事务修改了
////
func NonrepeatableRead(db1, db2 *sql.DB, isolationLevel string, pointName string) {
	//// 测试两个级别 都使 非主事务连接读取的数值为旧的。 {"READ COMMITTED", "REPEATABLE READ"}
	// Logg.Println(conn1, conn2, isolationLevel, pointName) //IsolationLevel = isolationLevel
	/*
		conn1 修改 Tom的cash 值时，conn2 也修改了 Tom的cash值
								// READ-UNCOMMITTED 隔离级 locked 可能会被主事务锁表 冲突
							conn2 read:
					from conn Final table state:
					1 |        Tom |                                  Admin@example.com | 1101
					ptx translations add Tom cash 1100 to 1110 from connection 1
					Conn1 Read:
					conn1:  1 |        Tom |                                  Admin@example.com | 1110
									// READ-COMMITTED 隔离级 locked 可能会被主事务锁表 冲突
					conn2 read:
					from conn Final table state:
					1 |        Tom |                                  Admin@example.com | 1101
					ptx translations add Tom cash 1100 to 1110 from connection 1
					Conn1 Read:
					conn1:  1 |        Tom |                                  Admin@example.com | 1110
					//REPEATABLE-READ
					 conn2 read:
			from conn Final table state:
			1 |        Tom |                                  Admin@example.com | 1101
			ptx translations add Tom cash 1100 to 1110 from connection 1
			Conn1 Read:
			conn1:  1 |        Tom |                                  Admin@example.com | 1110
				//SERIALIZABLE
			1 |        Tom |                                  Admin@example.com | 1101
			ptx translations add Tom cash 1100 to 1110 from connection 1
			Conn1 Read:
			conn1:  1 |        Tom |                                  Admin@example.com | 1110
	*/

	// ///启动事务  测试两次isolationLevel {"READ UNCOMMITTED", "READ COMMITTED"}
	Logg.Println("create sql.DB translation with level:.", isolationLevel)
	levels, err1 := db1.Exec("SET GLOBAL transaction_isolation= ? ;", isolationLevel)
	if err1 != nil {
		Logg.Println("conn1 set level failed:", levels, err1)
		fail(err1)
	}

	ptx, err := db1.Begin()
	if err != nil {
		Logg.Println("db1 begin err", err)
		panic(err)
	}
	////事务查 tom 的cash值
	row := ptx.QueryRow("SELECT cash FROM participant WHERE name='Tom'")
	var balance float64
	row.Scan(&balance)
	Logg.Printf("Tom cash at the beginning of transaction: %f\n", balance)

	///其他连接 更新tom 的cash值
	Logg.Printf("Updating Tom cash to 1101 from connection 2\n")
	//// READ-UNCOMMITTED 级别隔离将导致 非主事务被卡住
	_, err = db2.Exec("UPDATE participant SET cash = 1101 WHERE name='Tom';")
	if err != nil {
		Logg.Printf("Failed to update Tom cash from conn2  %e", err)
		// panic(err)
		// UNCOMMITED 隔离级 locked 可能会被主事务锁表 冲突
		// panic: Error 1205: Lock wait timeout exceeded; try restarting transaction
		return
	}
	Logg.Println("conn2 read:")
	printTable(db2)
	///主事务 修改 Tom的cash，增加10, 此时主事务 在RC 和 REPEATABLE 等全部 隔离级都可以修改成功
	Logg.Printf("ptx translations add Tom cash 1100 to 1110 from connection 1\n")
	_, err = ptx.Exec("UPDATE participant SET cash = ? WHERE name='Tom';", balance+10)
	if err != nil {
		Logg.Printf("Failed to update Tom cash in tx conn1: %v\n", err)
		panic(err)
	}
	Logg.Println("Conn1 Read:")
	rows, _ := ptx.Query("SELECT participantid, name, email, cash FROM participant ORDER BY participantid;")
	for rows.Next() {
		var name, email []byte
		var id int
		var cash float64
		rows.Scan(&id, &name, &email, &cash)
		Logg.Printf("conn1: %2d | %10s | %50s | %+v\n", id, name, email, cash)
	}

	/////保存事务执行点
	sql_savepoint := fmt.Sprintf("SAVEPOINT  %v", pointName)
	ptx.Exec(sql_savepoint)
	ptx.Exec("Commit;")
	if err := ptx.Commit(); err != nil {
		///提交事务失败
		Logg.Printf("Failed to commit: %v\n", err)
	}
}

////幻读， 当事务对 admin邮箱用户进行操作时，其他连接修改了 一个用户的邮箱为 admin
///
//// 当处于 重复读隔离级别，事务将 只更改 在事务启动的时间点的 admin邮箱用户，不更改其他连接在事务中新增的admin用户
func PhantomRead(conn1, conn2 *sql.DB, isolationLevel string, pointName string) {
	//// 测试两个隔离级别 {"READ COMMITTED", "REPEATABLE READ"}
	/*
			/// READ-UNCOMMITTED
			conn2 transaction moves Jack email to same tom and check:
		from conn Final table state:
		1 |        Tom |                                  Admin@example.com | 1100
		2 |       Jack |                                  Admin@example.com | 1150
		conn1 Users same email after cuncurrent transaction:
		[{Tom 1100} {Jack 1150}]
		conn1 Update selected users cash by +15
		save pointer sql: SAVEPOINT  READUNCOMMITTED0;
		from conn Final table state:
		1 |        Tom |                                  Admin@example.com | 1115
		2 |       Jack |                                  Admin@example.com | 1165

		/// READ-COMMITTED
		conn2 transaction moves Jack email to same tom and check:
		from conn Final table state:
		1 |        Tom |                                  Admin@example.com | 1100
		2 |       Jack |                                  Admin@example.com | 1150
		conn1 Users same email after cuncurrent transaction:
		[{Tom 1100} {Jack 1150}]
		conn1 Update selected users cash by +15
		save pointer sql: SAVEPOINT  READCOMMITTED1;
		from conn Final table state:
		1 |        Tom |                                  Admin@example.com | 1115
		2 |       Jack |                                  Admin@example.com | 1165

		// REPEATABLE-READ
		conn2 transaction moves Jack email to same tom and check:
		from conn Final table state:
		1 |        Tom |                                  Admin@example.com | 1100
		2 |       Jack |                                  Admin@example.com | 1150
		conn1 Users same email after cuncurrent transaction:
		[{Tom 1100} {Jack 1150}]
		conn1 Update selected users cash by +15
		save pointer sql: SAVEPOINT  REPEATABLEREAD2;
		from conn Final table state:
		1 |        Tom |                                  Admin@example.com | 1115
		2 |       Jack |                                  Admin@example.com | 1165
		// SERIALIZABLE
		conn2 transaction moves Jack email to same tom and check:
				from conn Final table state:
				1 |        Tom |                                  Admin@example.com | 1100
				2 |       Jack |                                  Admin@example.com | 1150
				conn1 Users same email after cuncurrent transaction:
				[{Tom 1100} {Jack 1150}]
				conn1 Update selected users cash by +15
				save pointer sql: SAVEPOINT  SERIALIZABLE3;
				from conn Final table state:
				1 |        Tom |                                  Admin@example.com | 1115
				2 |       Jack |                                  Admin@example.com | 1165

	*/
	// Create a helper function for preparing failure results.
	fail := func(err error) error {
		return fmt.Errorf("CreateOrder: %v", err)
	}
	// Logg.Println(conn1, conn2, isolationLevel, pointName) //IsolationLevel = isolationLevel
	// 另一种方式的 事务启动BeginTx

	// ///启动事务  测试两次isolationLevel {"READ UNCOMMITTED", "READ COMMITTED"}
	Logg.Println("create sql.DB translation with level:.", isolationLevel)
	levels, err1 := conn1.Exec("SET GLOBAL transaction_isolation= ? ;", isolationLevel)
	if err1 != nil {
		Logg.Println("conn1 set level failed:", levels, err1)
		panic(fail(err1))
	}
	Logg.Println("set iso level rst:", levels, err1)
	ptx, err := conn1.Begin()
	Logg.Println("begin rst:", ptx, err)
	if err != nil {
		// if ptx != nil {
		// 	_ = ptx.Rollback()
		// }
		Logg.Println(fail(err))
		panic(fail(err))
	}

	defer ptx.Rollback()
	// change_iso := fmt.Sprintf("SET GLOBAL transaction_isolation='%s'; " + isolationLevel)

	_, errbegin := ptx.Exec("BEGIN;")
	if errbegin != nil {
		Logg.Println(errbegin)
		panic(errbegin)
	}

	var users []User
	var user User
	rows, er := ptx.Query("SELECT name, cash FROM participant WHERE  email = 'Admin@example.com';")
	if er != nil {
		Logg.Println(fail(er))
		panic(fail(er))
	}
	for rows.Next() { /// 遍历返回的所有行
		var user User
		rows.Scan(&user.Name, &user.Cash)
		users = append(users, user)
	}
	Logg.Println("Begin %+v\n", user)

	Logg.Printf("Users email Admin@example.com at the beginning of transaction:\n%v\n", users)
	printTable(conn1)

	//// 其他连接 更改了 Jack的 邮箱，导致相同邮箱的用户变多
	Logg.Printf("conn2 transaction moves Jack email to same tom and check:\n")
	// READ-UNCOMMITTED 导致其他事务死锁
	editemail, err3 := conn2.Exec("UPDATE participant SET email = 'Admin@example.com' WHERE name='Jack'")
	if err3 != nil {
		/// Error 1205: Lock wait timeout exceeded; try restarting transaction
		Logg.Println("emial err:", editemail, err3)
		panic(err3)
	}
	printTable(conn2)

	/// 当前事务连接 admin邮箱查询, 将只会查到 conn1 连接的信息，只有tom一个是admin邮箱，因为事务隔离
	users = []User{}
	rows, _ = ptx.Query("SELECT name, cash FROM participant WHERE email = 'Admin@example.com'")
	for rows.Next() {
		rows.Scan(&user.Name, &user.Cash)
		users = append(users, user)
	}
	Logg.Printf("conn1 Users same email after cuncurrent transaction:\n%v\n", users)

	Logg.Printf("conn1 Update selected users cash by +15\n")
	for _, user := range users {
		_, err = ptx.Exec("UPDATE participant SET cash = ? WHERE name=?", user.Cash+15, user.Name)
		if err != nil {
			Logg.Printf("conn1 Failed to update in tx: %v\n", err)
		}
	}
	/////保存事务执行点, 保存点名称不能有 -
	sql_savepoint := fmt.Sprintf("SAVEPOINT  %v;", pointName)
	clear_sqlpoint := strings.ReplaceAll(sql_savepoint, "-", "")
	Logg.Println("save pointer sql:", clear_sqlpoint)
	_, err2 := ptx.Exec(clear_sqlpoint)
	if err2 != nil {
		Logg.Println("savepoint failure.", err2)
		panic(err)
	}
	/// 提交事务
	// ptx.Exec("Commit;")
	if err := ptx.Commit(); err != nil {
		Logg.Printf("conn1 Failed to commit: %v\n", err)
	}
}

//// 可序列化 隔离级别
//// A事务执行时，其他事务不能执行。当然 A 不能读取 B 事务的任何内容
//// 这时候的 隔离级别就是 序列化隔离级 Serialization anomaly, 这将阻止 对同一批数据的 序列化操作
func SerializationAnomaly(conn1, conn2 *sql.DB, isolationLevel string, pointName string) {
	//// 事务1 将读取所有 admin 邮箱的用户的 cash
	/// 事务2 将 jack 移动到 admin 邮箱组
	//// 事务1 再次读取 全部 admin 邮箱 用户的 cash
	//// 事务1 操作 对全部所选 admin 邮箱用户 的cash + 15
	////  提交事务 1
	////  提交事务 2   ##  在 序列化隔离级 事务2将提交失败，报错 由于 事务之间的读/写依赖 而无法序列化访问

	// Logg.Println(conn1, conn2, isolationLevel, pointName) //IsolationLevel = isolationLevel
	// 另一种方式的 事务启动BeginTx
	/*
			////READ-COMMITTED 隔离级别 成功提交
			from conn Final table state:
		1 |        Tom |                                  Admin@example.com | 10
		2 |       Jack |                                  Admin@example.com | 1150

		///// REPEATABLE-READ 死锁场景: 隔离级别 当 主事务执行 修改时，其他事务也修改了相同位置将导致主事务被锁定，超时回退到初始状态
		conn1 update participant user cash. <nil> Error 1205: Lock wait timeout exceeded; try restarting transaction
		conn1 SerializationAnomaly Failed to update in tx:<nil> with err: Error 1205: Lock wait timeout exceeded; try restarting transaction
		from conn Final table state:
		1 |        Tom |                                  Admin@example.com | 1100
		2 |       Jack |                                   User@example.com | 1150

		//// SERIALIZABLE  A事务未提交，B事务就等待。 直到失败
		死锁场景: 隔离级别，多个事务依次执行，表数据的状态取决于 最后一个事务的提交。如果总有不提交的，那就失败。
			与其他场景一起使用，无法初始化
			单独场景 报错如下：
			panic: Error conn2 SerializationAnomaly result: <nil> in tx2 err: Error 1205: Lock wait timeout exceeded; try restarting transaction
	*/
	// ///启动事务  测试两次isolationLevel {"READ UNCOMMITTED", "READ COMMITTED"}
	Logg.Println("create sql.DB translation with level:.", isolationLevel)
	levels, err1 := conn1.Exec("SET SESSION transaction_isolation= ? ;", isolationLevel)
	/// 注意 多事务连接时，设置session 会话的 事务级别 SET SESSION  transaction_isolation=
	if err1 != nil {
		Logg.Println("conn1 set SerializationAnomaly level failed:", levels, err1)
		panic(fail(err1))
	}
	Logg.Println("set SerializationAnomaly iso level rst:", levels, err1)

	/// 启动事务1 设置隔离级别
	ptx, err := conn1.Begin()
	if err != nil {
		panic(err)
	}
	defer ptx.Rollback()
	Logg.Println("create conn2 SerializationAnomaly sql.DB translation with level:.", isolationLevel)
	/// 注意 多事务连接时，设置session 会话的 事务级别  SET SESSION  transaction_isolation=
	Logg.Println("create conn2 sql.DB translation with level:.", isolationLevel)
	levels2, err2 := conn2.Exec("SET SESSION transaction_isolation= ? ;", isolationLevel)
	if err2 != nil {
		Logg.Println("conn2 set level failed:", levels2, err2)
		panic(fail(err2))
	}
	Logg.Println("conn2 SerializationAnomaly set iso level rst:", levels2, err2)
	///启动事务 2 设置隔离级别与 事务 1相同
	tx2, err2 := conn2.Begin()
	Logg.Println("conn2 begin translations:", tx2, err2)
	if err2 != nil {
		panic(err2)
	}
	defer tx2.Rollback()
	// rst, err2 := tx2.Exec("SET TRANSACTION ISOLATION LEVEL " + isolationLevel)
	// /*
	//  Error 1568: Transaction characteristics can't be changed while a transaction is in progress
	// */
	// Logg.Println("tx2 set level rst, err", rst, err2)
	// if err2 != nil {
	// 	panic(err2)
	// }
	// rstb, err2b := tx2.Exec("BEGIN;")
	// Logg.Println("tx2 begin rst, err", rstb, err2b)
	/// 事务1 选择 admin 邮箱的所有用户
	var sum int
	Logg.Println("conn1 select cash who is admin.")
	row := ptx.QueryRow("SELECT cash FROM participant WHERE email = 'Admin@example.com'")
	row.Scan(&sum)

	/// 事务 2 修改 jack 到 admin 邮箱组
	Logg.Println("conn2 edit jack to admin.")

	uprst, errtx2 := tx2.Exec("UPDATE participant SET email = 'Admin@example.com' WHERE name='Jack'")
	if errtx2 != nil {
		msg := fmt.Sprintf("Error conn2 SerializationAnomaly result: %+v in tx2 err: %v\n", uprst, errtx2)
		Logg.Printf(msg)
		// tx2.Rollback()
		// tx2.Commit()
		panic(msg)
	}

	/// 事务1 再次选择 admin 邮箱组
	rows, _ := ptx.Query("SELECT name, cash FROM participant WHERE email = 'Admin@example.com'")
	type User struct {
		Name string
		Cash int
	}
	var users []User
	for rows.Next() {
		var user User
		rows.Scan(&user.Name, &user.Cash)
		users = append(users, user)
	}

	/// 事务 1更新 全部 admin 邮箱组用户的 cash
	for _, user := range users {
		rst2, err2 := ptx.Exec("UPDATE participant SET cash = ? WHERE name=?", user.Cash+10, user.Name)
		Logg.Println("conn1 update participant user cash.", rst2, err2)
		//REPEATABLE-READ 隔离级别 Error 1205: Lock wait timeout exceeded; try restarting transaction
		if err2 != nil {
			Logg.Printf("conn1 SerializationAnomaly Failed to update in tx:%v with err: %v\n", rst2, err2)
			// ptx.Rollback()
			// ptx.Commit()
			// panic(err2)
			// READ-UNCOMMITTED Lock wait timeout exceeded; try restarting transactinon
			return
		}
	}

	/////保存事务执行点
	sql_savepoint := fmt.Sprintf("SAVEPOINT  %v;", pointName)
	clear_sqlpoint := strings.ReplaceAll(sql_savepoint, "-", "")
	Logg.Println("save pointer sql:", clear_sqlpoint)
	saves, err2s := ptx.Exec(clear_sqlpoint)
	if err2s != nil {
		Logg.Println("conn1 try to save point fail:", saves, err2s)
		// ptx.Rollback()

		// ptx.Commit()
		return
	}
	sql_savepoint2 := fmt.Sprintf("SAVEPOINT  %v;", pointName)
	clear_sqlpoint2 := strings.ReplaceAll(sql_savepoint2, "-", "2")
	Logg.Println("save pointer2 sql:", clear_sqlpoint2)
	savets, errt2s := tx2.Exec(clear_sqlpoint2)
	if errt2s != nil {
		Logg.Println("conn1 try to save point fail:", savets, errt2s)
		// tx2.Rollback()

		// tx2.Commit()
		return
	}
	/// 提交事务 1
	if err := ptx.Commit(); err != nil {
		msg := fmt.Sprintf("conn1 SerializationAnomaly Failed to commit tx: %v\n", err)
		Logg.Printf(msg)
		// ptx.Rollback()

		// ptx.Commit()
		panic(msg)
	}
	/// 提交事务 2  在序列化 隔离级别 将被阻止提交
	if err := tx2.Commit(); err != nil {
		msg := fmt.Sprintf("conn2 SerializationAnomaly Failed to commit tx: %v\n", err)
		Logg.Printf(msg)
		// tx2.Rollback()

		// tx2.Commit()
		panic(msg)
	}
}

func IsolationDo() {
	////

	/*
		dbpool.QueryRow(ctx, "select 'Hello, world!'").Scan(&greeting)
	*/
	db, err = sql.Open("mysql", DSN) /// 创建连接池，默认无限制
	if err != nil {
		log.Fatal(err)

	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)

	Logg.Println("conn mysql success with", DSN)
	defer db.Close()

	// DSN := cfg.FormatDSN()
	db2, err2 = sql.Open("mysql", DSN1) /// 创建连接池
	if err2 != nil {
		log.Fatal(err2)
	}
	db2.SetConnMaxLifetime(time.Minute * 3)
	db2.SetMaxOpenConns(2000)
	db2.SetMaxIdleConns(1000)
	Logg.Println("conn mysql success with", DSN1)
	defer db2.Close()

	// 调用DB.Ping以确认连接到数据库有效。
	//在运行时，sql.Open可能不会立即连接，具体取决于驱动程序。
	//您在Ping此处使用以确认 database/sql包可以在需要时连接。
	pingErr := db.Ping()
	if pingErr != nil { // 检查来自 的错误Ping，以防连接失败。
		log.Fatal(pingErr)
	}
	fmt.Println("db Connected!") //Ping如果连接成功，则打印一条消息。

	pingErr2 := db2.Ping()
	if pingErr != nil { // 检查来自 的错误Ping，以防连接失败。
		log.Fatal(pingErr2)
	}
	fmt.Println("db2 connection Connected!") //Ping如果连接成功，则打印一条消息。

	/// 初始化数据
	// InitDb(db)
	type ReadCash struct {
		name            string
		isolationLevels []string
		testFunction    func(db, db2 *sql.DB, isolationLevel string, pointName string)
	}

	phenomenas := []ReadCash{
		{
			/* "READ-UNCOMMITTED", A 可以 读取 B 未提交的更改
			 "READ-COMMITTED", A 只能读取 B 已提交的更改
			"REPEATABLE-READ", 事务只能读取事务内的数据
			 "SERIALIZABLE", 事务 按顺序执行，不能相互读取任何内容，A 执行时，B 不能进行操作。
			*/
			/// 脏读 测试 读未提交 和 读已提交对比，脏读结果一致
			//// 其他连接 总是读取到 事务未提交 的 cash 值，所以是脏读，
			name:            "Dirty read",
			isolationLevels: []string{"READ-UNCOMMITTED", "READ-COMMITTED", "REPEATABLE-READ", "SERIALIZABLE"},
			testFunction:    dirtyRead,
		},
		{
			/// 重复读隔离对比
			//// 读提交 隔离，可以完成 该事务， 重复读 隔离，无法完成该事务，两个会话同时写，导致conn1不可完成事务
			/// 对比隔离 级别，读已提交 READ COMMITTED 时，事务连接 即使在有其他连接修改cash时，也可以读写成功，修改成功Tom的cash，为事务的修改值
			//// 在可重复读 REPEATABLE READ 隔离级别， 事务连接 在有其他连接修改cash时，修改失败， tom的cash 将为 其他连接的修改值
			//// 在 重复读 隔离级 REPEATABLE READ，如果事务执行过程中，有其他连接修改了 事务需要修改的数据，将导致事务执行失败报错，事务回滚，最后结果为其他连接的修改结果。
			//// REPEATABLE READ 则是不可重复的读取情况，不可重复修改隔离级
			name:            "Nonrepeatable read",
			isolationLevels: []string{"READ-UNCOMMITTED", "READ-COMMITTED", "REPEATABLE-READ", "SERIALIZABLE"},
			testFunction:    NonrepeatableRead,
		},
		{
			/////幻读 Mysql
			//  幻读 类似于 不可重复读，在事务中选择 一组 行。 如果外部更改了一行，则一组行都更改，这就是幻读的场景
			//// 可重复读 Repeatable read 隔离级别，可以防止这种场景。 此隔离级别将保存从事务开始的数据，而 1 和 3的读取将返回同一行集，与并发更改隔离
			//// 可重复读 隔离级，不包括在 事务操作中，新加入的 同一email地址的 Jack，
			//// 而在 读提交 隔离级 将包括 在事务中新增相同email的Jack
			name:            "Phantom read",
			isolationLevels: []string{"READ-UNCOMMITTED", "READ-COMMITTED", "REPEATABLE-READ", "SERIALIZABLE"},
			testFunction:    PhantomRead,
		},
		{
			//// 序列化 事务 隔离
			//// //// 事务1 将读取所有 admin 邮箱的用户的 cash
			/// 事务2 将 jack 移动到 admin 邮箱组
			//// 事务1 再次读取 全部 admin 邮箱 用户的 cash
			//// 事务1 操作 对全部所选 admin 邮箱用户 的cash + 15
			////  提交事务 1
			////  提交事务 2   ##  在 序列化隔离级 事务2将提交失败，报错 由于 事务之间的读/写依赖 而无法序列化访问
			////  Failed to commit tx2: ERROR: could not serialize access due to read/write dependencies among transactions (SQLSTATE 40001)
			name: "Serialization anomaly",
			// "READ-UNCOMMITTED", 隔离级别 当 其他事务修改了主事务的同一条数据时，导致主事务 lock超时
			isolationLevels: []string{"SERIALIZABLE"}, // "READ-COMMITTED", "REPEATABLE-READ",
			testFunction:    SerializationAnomaly,
		},
	}
	for _, phenomena := range phenomenas {
		Logg.Printf("%s\n", phenomena.name)
		for i, isolationLevel := range phenomena.isolationLevels {
			Logg.Printf("\nIsolation level - %s\n", isolationLevel)
			/// 每次重置 表
			rbool := InitDb(db)
			Logg.Println("init db:", rbool)
			/// 查询 db数据，
			isolationLevels := strings.Replace(isolationLevel, " ", "", -1)
			pointName := fmt.Sprintf("%v%v", isolationLevels, i)
			// Logg.Println("dbs:", db, db2, "iso level:", isolationLevels, "pointName:", pointName)
			phenomena.testFunction(db, db2, isolationLevel, pointName)
			//// 查询主 事务最后结果
			printTable(db)
		}
		Logg.Printf("\n---\n\n")
	}
}
func printTable(conn *sql.DB) {
	Logg.Printf("from conn Final table state:\n")
	rows, _ := conn.Query("SELECT participantid, name, email, cash FROM participant ORDER BY participantid")
	for rows.Next() {
		var name, email []byte
		var id int
		var cash float64
		rows.Scan(&id, &name, &email, &cash)
		Logg.Printf("%2d | %10s | %50s | %+v\n", id, name, email, cash)
	}
}
func main() {
	// users, _ := UsersByArtist("Tom")
	// Logg.Printf("get users:%+v\n", users)

	IsolationDo()
}

// 查询多行
// UsersByArtist queries for Users that have the specified artist name.
func UsersByArtist(name string) ([]User, error) {
	// An Users slice to hold data from returned rows.
	// 声明您定义Users的类型的切片。User这将保存来自返回行的数据。
	// 结构字段名称和类型对应于数据库列名称和类型。
	var Users []User
	// 用于DB.Query执行SELECT语句以查询具有指定艺术家姓名的专辑。
	// Query的第一个参数是 SQL语句。
	// 在参数之后，您可以传递零个或多个任何类型的参数。
	//这些为您提供了在 SQL 语句中指定参数值的位置。
	//通过将 SQL 语句与参数值分开（而不是将它们连接起来，比如说，fmt.Sprintf），
	//您可以让 database/sql包将值与 SQL 文本分开发送，从而消除任何 SQL 注入风险。
	rows, err := db.Query("SELECT name,cash FROM participant WHERE name = ?", name)
	if err != nil {
		return nil, fmt.Errorf("participant %s: %v", name, err)
	}
	// 延迟关闭rows，以便在函数退出时释放它持有的任何资源。
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	// 循环遍历返回的行， Rows.Scan用于将每行的列值分配给User结构字段。
	// Scan获取指向 Go 值的指针列表，列值将被写入其中。alb在这里，
	//您将指针传递给使用运算符创建的变量中的字段 &。Scan通过指针写入以更新结构字段。
	for rows.Next() {
		var alb User
		//在循环内部，检查将列值扫描到结构字段中是否存在错误。
		if err := rows.Scan(&alb.Name, &alb.Cash); err != nil {
			return nil, fmt.Errorf("participant user: %s: %v", name, err)
		}
		//在循环内部，将新的附加alb到Users切片。
		Users = append(Users, alb)
	}
	if err := rows.Err(); err != nil {
		//在循环之后，使用 . 检查整个查询中的错误 rows.Err。
		//请注意，如果查询本身失败，则在此处检查错误是找出结果不完整的唯一方法。
		return nil, fmt.Errorf("participant user: %s: %v", name, err)
	}
	Logg.Println("get user from participant", Users)
	return Users, nil
}

// 查询单行
// UserByID queries for the User with the specified ID.
func UserByID(id int64) (User, error) {
	// An User to hold data from the returned row.
	var alb User
	// 用于DB.QueryRow 执行SELECT语句查询指定ID的数据。
	// 它返回一个sql.Row. 为了简化调用代码（您的代码！），QueryRow不返回错误。
	// 相反，它安排从以后返回任何查询错误（例如sql.ErrNoRows）Rows.Scan
	row := db.QueryRow("SELECT * FROM User WHERE id = ?", id)
	// 用于Row.Scan将列值复制到结构字段中。
	if err := row.Scan(&alb.Name, &alb.Cash); err != nil {
		// 检查来自 Scan的错误。
		if err == sql.ErrNoRows {
			// 特殊错误sql.ErrNoRows表示查询未返回任何行。通常，该错误值得用更具体的文本替换，例如此处的“没有这样的专辑”。
			return alb, fmt.Errorf("UsersById %d: no such User", id)
		}
		return alb, fmt.Errorf("UsersById %d: %v", id, err)
	}
	return alb, nil
}

//添加数据
// addUser adds the specified User to the database,
// returning the User ID of the new entry
func addUser(alb User) (int64, error) {
	// 用于DB.Exec执行INSERT语句。Like Query，Exec接受一条 SQL 语句，后跟 SQL 语句的参数值。
	result, err := db.Exec("INSERT INTO User (title, artist, price) VALUES (?, ?)", alb.Name, alb.Cash)
	if err != nil {
		// 检查尝试中的错误INSERT。
		return 0, fmt.Errorf("addUser: %v", err)
	}
	//使用 检索插入的数据库行的 ID Result.LastInsertId。
	id, err := result.LastInsertId()
	if err != nil {
		//检查尝试检索 ID 的错误。
		return 0, fmt.Errorf("addUser: %v", err)
	}
	return id, nil
}
