package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	host     = "192.168.30.129"
	port     = 5432
	username = "postgre"
	password = "post.2021"
	dbname   = "pgstate"
)

var (
	ctx  context.Context
	Logg = log.New(os.Stderr, "INFO -", 18)
)

type User struct {
	Name string
	Cash float64
}

//// 脏读
func dirtyRead(conn1, conn2 *pgxpool.Pool, isolationLevel string, pointName string) {
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
	*/
	ptx, perr := conn1.Begin(ctx)
	if perr != nil {
		panic(perr)
	}
	///启动事务  测试两次isolationLevel {"READ UNCOMMITTED", "READ COMMITTED"}
	ptx.Exec(ctx, "SET TRANSACTION ISOLATION LEVEL "+isolationLevel)

	//// 执行交易 从Tom 扣除12.99
	_, perr = ptx.Exec(ctx, "UPDATE participant SET cash=cash-12.99 WHERE participantid=1;")
	if perr != nil {
		Logg.Printf("Failed to update Tom cash in tx: %v\n", perr)
	}
	//// 执行交易 给Jack 增加12.99
	_, perr2 := ptx.Exec(ctx, "UPDATE participant SET cash=cash+12.99 WHERE participantid=2;")
	if perr2 != nil {
		Logg.Printf("Failed to update Tom cash in tx: %v\n", perr2)
	}
	//// 执行交易 把物品从 Jack交还给Tom
	_, perr3 := ptx.Exec(ctx, "UPDATE trunk SET participantid = 1 WHERE name = 'ComputerABC' AND participantid=2;")
	if perr != nil {
		Logg.Printf("Failed to trunk SET  in tx: %v\n", perr3)
	}

	/////保存事务执行点
	sql_savepoint := fmt.Sprintf("SAVEPOINT  %v", pointName)
	ptx.Exec(ctx, sql_savepoint)
	/// 事务提交前 当前事务连接 和 其他连接 查询对比
	var balance float64
	row := ptx.QueryRow(ctx, "SELECT cash FROM participant WHERE name='Tom'")
	row.Scan(&balance)
	Logg.Printf("提交事务前:连接内部值Tom cash from main transaction after update: %f\n", balance) /// 事务连接内部值 tom cash 1087.01

	var othercash float64
	row1 := conn2.QueryRow(ctx, "SELECT cash FROM participant WHERE name='Tom'")
	row1.Scan(&othercash)
	/// 其他连接的 脏读
	Logg.Printf("提交事务前:其他连接Tom cash from conn2 : %f\n", othercash) /// 其他连接tom cash 1100
	Logg.Printf("same eq?:%+v\n", balance+12.99 == othercash)

	if err := ptx.Commit(ctx); err != nil {
		Logg.Printf("Failed to commit: %v\n", err)
	}

	///ComputerABC
	var participantOwn string
	row2 := conn2.QueryRow(ctx, "SELECT *  FROM trunk WHERE name = 'ComputerABC'")
	row2.Scan(&participantOwn)
	Logg.Printf("提交事务后:Tom trunk from conn2 : %s\n", participantOwn)
}

///// 重复读 当一个事务重新读取前面读取过的数据时，发现该数据已经被另一个已提交事务修改了
//// 两次读取之间，数据被其他事务修改。 看起来数据就不对
func NonrepeatableRead(conn1, conn2 *pgxpool.Pool, isolationLevel string, pointName string) {
	//// 测试两个级别 都使 非主事务连接读取的数值为旧的。 {"READ COMMITTED", "REPEATABLE READ"}
	ptx, err := conn1.Begin(ctx)
	if err != nil {
		panic(err)
	}
	ptx.Exec(ctx, "SET TRANSACTION ISOLATION LEVEL "+isolationLevel)
	////事务查 tom 的cash值
	row := ptx.QueryRow(ctx, "SELECT cash FROM participant WHERE name='Tom'")
	var balance float64
	row.Scan(&balance)
	Logg.Printf("Tom cash at the beginning of transaction: %f\n", balance)

	///其他连接 更新tom 的cash值
	Logg.Printf("Updating Tom cash to 1101 from connection 2\n")
	_, err = conn2.Exec(ctx, "UPDATE participant SET cash = 1101 WHERE name='Tom'")
	if err != nil {
		Logg.Printf("Failed to update Tom cash from conn2  %e", err)
	}
	Logg.Println("conn2 read:")
	printTable(conn2)
	///事务 修改 Tom的cash，增加10，， 此时该连接应该cash = 1120
	Logg.Printf("ptx translations add Tom cash 1100 to 1110 from connection 1\n")
	_, err = ptx.Exec(ctx, "UPDATE participant SET cash = $1 WHERE name='Tom'", balance+10)
	if err != nil {
		Logg.Printf("Failed to update Bob balance in tx conn1: %v\n", err)
	}
	Logg.Println("Conn1 Read:")
	rows, _ := ptx.Query(ctx, "SELECT participantid, name, email, cash FROM participant ORDER BY participantid")
	for rows.Next() {
		var name, email []byte
		var id int
		var cash float64
		rows.Scan(&id, &name, &email, &cash)
		Logg.Printf("conn1: %2d | %10s | %50s | %+v\n", id, name, email, cash)
	}

	/////保存事务执行点
	sql_savepoint := fmt.Sprintf("SAVEPOINT  %v", pointName)
	ptx.Exec(ctx, sql_savepoint)

	if err := ptx.Commit(ctx); err != nil {
		///提交事务失败
		Logg.Printf("Failed to commit: %v\n", err)
	}
}

////幻读， 当事务对 admin邮箱用户进行操作时，其他连接修改了 一个用户的邮箱为 admin
/// 这样 事务连接 的admin 邮箱用户对象更多了，当处于 读提交 隔离级别，这可以完成事务
//// 当处于 重复读隔离级别，事务将 只更改 在事务启动的时间点的 admin邮箱用户，不更改其他连接在事务中新增的admin用户
func PhantomRead(conn1, conn2 *pgxpool.Pool, isolationLevel string, pointName string) {
	//// 测试两个隔离级别 {"READ COMMITTED", "REPEATABLE READ"}
	ptx, err := conn1.Begin(ctx)
	if err != nil {
		panic(err)
	}
	ptx.Exec(ctx, "SET TRANSACTION ISOLATION LEVEL "+isolationLevel)

	var users []User
	var user User
	rows, _ := ptx.Query(ctx, "SELECT name, cash FROM participant WHERE  email = 'Admin@example.com';")
	for rows.Next() { /// 遍历返回的所有行
		var user User
		rows.Scan(&user.Name, &user.Cash)
		users = append(users, user)
	}
	fmt.Printf("Users at the beginning of transaction:\n%v\n", users)

	//// 其他连接 更改了 Jack的 邮箱，导致相同邮箱的用户变多
	fmt.Printf("Cuncurrent transaction moves Jack email to same tom\n")
	conn2.Exec(ctx, "UPDATE participant SET email = 'Admin@example.com' WHERE name='Jack'")

	/// 当前事务连接 admin邮箱查询
	users = []User{}
	rows, _ = ptx.Query(ctx, "SELECT name, cash FROM participant WHERE email = 'Admin@example.com'")
	for rows.Next() {
		rows.Scan(&user.Name, &user.Cash)
		users = append(users, user)
	}
	fmt.Printf("Users same email after cuncurrent transaction:\n%v\n", users)

	fmt.Printf("Update selected users cash by +15\n")
	for _, user := range users {
		_, err = ptx.Exec(ctx, "UPDATE participant SET cash = $1 WHERE name=$2", user.Cash+15, user.Name)
		if err != nil {
			fmt.Printf("Failed to update in tx: %v\n", err)
		}
	}
	/////保存事务执行点
	sql_savepoint := fmt.Sprintf("SAVEPOINT  %v", pointName)
	ptx.Exec(ctx, sql_savepoint)

	/// 提交事务
	if err := ptx.Commit(ctx); err != nil {
		fmt.Printf("Failed to commit: %v\n", err)
	}
}

//// 可序列化 隔离级别
//// 如果有多个 事务连接正在执行操作，而我们希望 最终表的 数据的状态取决于 事务执行和提交的顺序，
//// 这时候的 隔离级别就是 序列化隔离级 Serialization anomaly, 这将阻止 对同一批数据的 序列化操作
func SerializationAnomaly(conn1, conn2 *pgxpool.Pool, isolationLevel string, pointName string) {
	//// 事务1 将读取所有 admin 邮箱的用户的 cash
	/// 事务2 将 jack 移动到 admin 邮箱组
	//// 事务1 再次读取 全部 admin 邮箱 用户的 cash
	//// 事务1 操作 对全部所选 admin 邮箱用户 的cash + 15
	////  提交事务 1
	////  提交事务 2   ##  在 序列化隔离级 事务2将提交失败，报错 由于 事务之间的读/写依赖 而无法序列化访问

	/// 启动事务1 设置隔离级别
	ptx, err := conn1.Begin(ctx)
	if err != nil {
		panic(err)
	}
	ptx.Exec(ctx, "SET TRANSACTION ISOLATION LEVEL "+isolationLevel)

	///启动事务 2 设置隔离级别与 事务 1相同
	tx2, err := conn2.Begin(ctx)
	if err != nil {
		panic(err)
	}
	tx2.Exec(ctx, "SET TRANSACTION ISOLATION LEVEL "+isolationLevel)

	/// 事务1 选择 admin 邮箱的所有用户
	var sum int
	row := ptx.QueryRow(ctx, "SELECT SUM(cash) FROM participant WHERE email = 'Admin@example.com'")
	row.Scan(&sum)

	/// 事务 2 修改 jack 到 admin 邮箱组
	tx2.Exec(ctx, "UPDATE participant SET email = 'Admin@example.com' WHERE name='Jack'")
	if err != nil {
		fmt.Printf("Error in tx2: %v\n", err)
	}

	/// 事务1 再次选择 admin 邮箱组
	rows, _ := ptx.Query(ctx, "SELECT name, cash FROM participant WHERE email = 'Admin@example.com'")
	type User struct {
		Name    string
		Balance int
	}
	var users []User
	for rows.Next() {
		var user User
		rows.Scan(&user.Name, &user.Balance)
		users = append(users, user)
	}

	/// 事务 1更新 全部 admin 邮箱组用户的 cash
	for _, user := range users {
		_, err = ptx.Exec(ctx, "UPDATE participant SET cash = $1 WHERE name=$2", user.Balance+sum, user.Name)
		if err != nil {
			fmt.Printf("Failed to update in tx: %v\n", err)
		}
	}

	/////保存事务执行点
	sql_savepoint := fmt.Sprintf("SAVEPOINT  %v", pointName)
	ptx.Exec(ctx, sql_savepoint)

	sql_savepoint2 := fmt.Sprintf("SAVEPOINT  %v", pointName)
	tx2.Exec(ctx, sql_savepoint2)

	/// 提交事务 1
	if err := ptx.Commit(ctx); err != nil {
		fmt.Printf("Failed to commit tx: %v\n", err)
	}

	/// 提交事务 2  在序列化 隔离级别 将被阻止提交
	if err := tx2.Commit(ctx); err != nil {
		fmt.Printf("Failed to commit tx2: %v\n", err)
	}
}

func InitDb(conn *pgxpool.Pool) {
	/// 初始化  DROP ... CASCADE 强制删除依赖，将导致其他表的依赖 列 混乱
	sql := `
	
	DROP TABLE IF EXISTS trunk CASCADE;	
	DROP TABLE IF EXISTS participant CASCADE;
	CREATE TABLE  participant (
		participantid SERIAL,
		name TEXT NOT NULL,
		email TEXT NOT NULL,
		cash NUMERIC(15,2) NOT NULL,   

		PRIMARY KEY (participantid)
	);
	CREATE TABLE  trunk (
		trunkid SERIAL,
		participantid INTEGER NOT NULL REFERENCES participant(participantid),
		name TEXT NOT NULL,
		price NUMERIC(5,2) NOT NULL,
		description TEXT NOT NULL,
		PRIMARY KEY  (trunkid)
	);

	INSERT INTO participant (name, email, cash) VALUES ('Tom', 'Admin@example.com', '1100.00');
 
	INSERT INTO participant (name, email, cash) VALUES ('Jack', 'User@example.com', '1150.00');
	
	INSERT INTO trunk (participantid, name,price, description) VALUES (1,'Linux CD', '1.00', 'Complete OS on a CD'); 
	INSERT INTO trunk (participantid, name,price, description) VALUES (2,'ComputerABC', '12.90', 'a book about OS computer!');
	INSERT INTO trunk (participantid, name,price, description) VALUES (2,'Magazines', '6.90', 'Stack of Computer Magezines computer!');
`
	_, err := conn.Exec(context.Background(), sql)
	if err != nil {
		Logg.Println("failed to panic", err)
		panic(err)
	}
}

func IsolationDo() {
	////
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbname)
	ctx = context.Background()

	/*
		dbpool.QueryRow(ctx, "select 'Hello, world!'").Scan(&greeting)
	*/
	/// 创建第一个pgsql连接
	conn1, err := pgxpool.Connect(ctx, connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn1.Close()

	/// 创建第2个pgsql连接 用于在第一个连接执行时 执行对比操作，验证 事务隔离级别
	conn2, err := pgxpool.Connect(ctx, connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn2.Close()

	type ReadCash struct {
		name            string
		isolationLevels []string
		testFunction    func(conn1, conn2 *pgxpool.Pool, isolationLevel string, pointName string)
	}

	phenomenas := []ReadCash{
		{
			/// 脏读 测试 读未提交 和 读已提交对比，脏读结果一致
			//// 其他连接 总是读取到 事务未提交 的 cash 值，所以是脏读
			name:            "Dirty read",
			isolationLevels: []string{"READ UNCOMMITTED", "READ COMMITTED"},
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
			isolationLevels: []string{"READ COMMITTED", "REPEATABLE READ"},
			testFunction:    NonrepeatableRead,
		},
		{
			/////幻读
			//  幻读 类似于 不可重复读，在事务中选择 一组 行。 如果外部更改了一行，则一组行都更改，这就是幻读的场景
			//// 可重复读 Repeatable read 隔离级别，可以防止这种场景。 此隔离级别将保存从事务开始的数据，而 1 和 3的读取将返回同一行集，与并发更改隔离
			//// 可重复读 隔离级，不包括在 事务操作中，新加入的 同一email地址的 Jack，
			//// 而在 读提交 隔离级 将包括 在事务中新增相同email的Jack
			name:            "Phantom read",
			isolationLevels: []string{"READ COMMITTED", "REPEATABLE READ"},
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
			name:            "Serialization anomaly",
			isolationLevels: []string{"REPEATABLE READ", "SERIALIZABLE"},
			testFunction:    SerializationAnomaly,
		},
	}

	for _, phenomena := range phenomenas {
		Logg.Printf("%s\n", phenomena.name)
		for i, isolationLevel := range phenomena.isolationLevels {
			Logg.Printf("\nIsolation level - %s\n", isolationLevel)
			/// 每次重置 表
			InitDb(conn1)
			/// 查询 db数据，
			isolationLevels := strings.Replace(isolationLevel, " ", "", -1)
			pointName := fmt.Sprintf("%v%v", isolationLevels, i)
			phenomena.testFunction(conn1, conn2, isolationLevel, pointName)
			printTable(conn1)
		}
		Logg.Printf("\n---\n\n")
	}
}

func TransactionDo(conn *pgxpool.Pool, pointName string) {
	//// 执行交易 事务
	sql := `
	START TRANSACTION;
	UPDATE participant SET cash=cash-12.99 WHERE participantid=1;
	UPDATE participant SET cash=cash+12.99 WHERE participantid=2;
	UPDATE trunk SET participantid = 1 WHERE name = 'ComputerABC' AND participantid=2;
	SAVEPOINT %s;
	SELECT * FROM participant;
	COMMIT;
	`
	fullSql := fmt.Sprintf(sql, pointName)
	Logg.Println("fullSql:", fullSql)
	_, err := conn.Exec(context.Background(), fullSql)
	if err != nil {
		panic(err)
	}
}

func printTable(conn *pgxpool.Pool) {
	Logg.Printf("from conn Final table state:\n")
	rows, _ := conn.Query(ctx, "SELECT participantid, name, email, cash FROM participant ORDER BY participantid")
	for rows.Next() {
		var name, email []byte
		var id int
		var cash float64
		rows.Scan(&id, &name, &email, &cash)
		Logg.Printf("%2d | %10s | %50s | %+v\n", id, name, email, cash)
	}
}

func DefaultConnSinger() {
	///// 官方教程 需要导入
	//// 使用环境变量种指定的数据库 URL DATABASE_URL pgx支持标准PostgreSQL环境变量
	/// 例如PGHOST PGDATABASE 使用与上述测试时相同的连接设置 psql，如果您的psql不需要任何参数，则不为pgx指定任何参数
	//// pgx 使用与psql默认连接值类似的逻辑
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbname)
	ctx = context.Background()
	// conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL")) // 从OS环境变量获得 数据库连接地址
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		Logg.Println(os.Stderr, "Unable to connecto to database:%+v\n", err)
		fmt.Fprintf(os.Stderr, "Unable to connecto to database:%+v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	var greeting string
	//// 从pgsql 查询 hello world，回显一个hello world
	err = conn.QueryRow(context.Background(), "select 'Hello, world'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed:%+v\n", err)
		Logg.Println(os.Stderr, "QueryRow failed:%+v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "from pgsql:%+v\n", greeting)
	Logg.Println(greeting)
}

func PoolConnSinger() {
	//// 使用连接池 *pgx.Conn返回的 表示 pgx.Connect()单个连接，并且不是并发(线程)安全的
	//// 这完全适用于上面的简单命令行示例，但是对于 许多用途，例如web应用程序服务器，需要并发性，则使用连接池
	//// 请将import github.com/jackc/pgx/v4 替换为 github.com/jackc/pgx/v4/pgxpool,
	//// 连接数据库操作 使用 pgxpool.Connect() 而不是 pgx.Connect()
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbname)
	ctx = context.Background()
	dbpool, err := pgxpool.Connect(ctx, connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connecto to database:%+v\n", err)
		Logg.Println(os.Stderr, "Unable to connecto to database:%+v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	var greeting string
	err = dbpool.QueryRow(ctx, "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed:%+v\n", err)
		Logg.Println(os.Stderr, "QueryRow failed:%+v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "QueryRow fmt:%+v\n", greeting)
	Logg.Printf("QueryRow from db:%+v\n", greeting)
}
func main() {

	// DefaultConnSinger()
	// PoolConnSinger()
	IsolationDo()
}
