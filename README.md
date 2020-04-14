# sql orm 数据库的连贯操作
* 这个包的目的是方便使用数据库操作而不是使用原生sql来进行增删改查
* 目前支持标准库
```
sqlite github.com/mattn/go-sqlite3
```

### 定义开始

```
//创建一个模型的结构体并且需要根据表的字段定义对应的属性
//属性标签可以通过`model:"-"` 来忽略全局字段的查询的操作,对删改查操作无效
type Users struct {
	Id int8 `model:"id"`
	Name string `model:"name"`
	Age int8 `model:"-"`
}

//初始化操作 connect.Start 需要提供对应数据库的连接方式配置,模型本身
//例如sqlite只需要db文件的路径就行了
func (this *Users) Init(connect *sqlite.Connect) *sqlite.Connect {
	connect = connect.Start(
		sqlite.Config{Path:os.Getenv("GOPATH") + "/databases/SQLite3/base.db"},
		this,
	)
	return connect
}

//定义操作的表名称，如果为定义的话默认采用结构体名称，或是连贯操作的Table 
//优先级 连贯操作 Table > func TableName() > 结构体名称 
func (this *Users) TableName() string {
	return "users"
}
```

### 连贯where操作(生成连贯操作所需要的where条件)
```
//连贯操作开始之前需要先介绍下 databases.Where 的使用
//为了方便解耦把复杂的where条件单独定义一个结构体

//初始化 where 需要一个WhereOption 对象
var where databases.Where

//一个用and 关联操作的where条件 会形成一个 age = 18 AND age = 17 的关联操作 

where.add(option = databases.WhereOption{
      Option: databases.DATABASE_WHERE_HANDLE_AND,
      Operation: []databases.WhereOperation{
         {
             Key: "age",
             Handle: "=",
             Value: "18",
         },
         {
             Key: "age",
             Handle: "=",
             Value: "17",
         },
     },
  })

//一个用Or 关联操作的where条件 会形成一个 age = 18 OR age = 17 的关联操作 

where.add(databases.WhereOption{
      Option: databases.DATABASE_WHERE_HANDLE_OR,
      Operation: []databases.WhereOperation{
          {
              Key: "age",
              Handle: "=",
              Value: "18",
          },
          {
              Key: "age",
              Handle: "=",
              Value: "17",
          },
     },
  })

//如果存在复杂的关联操作，where也可以使用连贯操作 `age` = 17 OR (`age` = 18 AND `age` = 17)
//通过多层嵌套的关系来形成一个(`age` = 18 AND `age` = 17 ) 的括号条件

where.add(databases.WhereOption{
      Option: databases.DATABASE_WHERE_HANDLE_OR,
      Operation: []databases.WhereOperation{
          {
              Key: "age",
              Handle: "=",
              Value: "17",
          },
     },
  })

where.add(databases.WhereOption{
      Option: databases.DATABASE_WHERE_HANDLE_OR,
      Operation: []databases.WhereOperation{
          {
              WhereOption: databases.WhereOption{
                     Option: databases.DATABASE_WHERE_HANDLE_AND,
                     Operation: []databases.WhereOperation{
                        {
                            Key: "age",
                            Handle: "=",
                            Value: "18",
                        },
                        {
                            Key: "age",
                            Handle: "=",
                            Value: "17",
                        },
                    },
                 }
          },
     },
  })
```

### 连贯操作
```
//初始化自定义结构体
var table = new(model.Users)

//初始化where条件
var where databases.Where

//sql的连贯操作 sqlite.Connect
model := table.Init(new(sqlite.Connect))

//定义一个简单的where条件
where.add(databases.WhereOption{
      Option: databases.DATABASE_WHERE_HANDLE_AND,
      Operation: []databases.WhereOperation{
          {
              Key: "age",
              Handle: "=",
              Value: "18",
          },
     },
  })

//这样就算定义来一个简单的连贯操作,其中的内容会转换成如下的sql执行
//SELECT `id`,`name` FROM users WHERE `age` = 18 ORDER BY `id` DESC 
model.Where(where).order("id",databases.DATABASE_ORDER_HANDLE_DESC).Select().All()

//model 的select 提供来暂时的缓存如果想要再次使用数据可以通过
model.All() //获取上一次的sql条件

//在你不需要使用数据库操作的时候，需要关闭资源
//通过Select(),Insert(),Delete(),Update() 会把关闭的资源重新申请所以不用担心关闭的影响
model.Close()
```

### 事务操作

```
//事务的开始 model是连贯操作中的model
var tx = new(sqlite.TXExec)

//可以通过某个model.Config 来获取所需要的db连接配置,这样就不用重新定义了
tx.Start(model.Config)
//加入需要事务操作的model可以通过连贯操作来操作多个表的事务
tx.Add(model).Add(model2)

//提交事务操作
tx.Commit()

//回滚事务
tx.Rollback()
```

