# generate-struct
generate struct from mysql table for golang, a small userful tool


how to use

cd generate-struct

1、change const variables for db config
vi generate.go
<pre><code>
const (
	DB_TYPE = "mysql"
	DB_HOST = "127.0.0.1"
	DB_PORT = "3306"
	DB_USER = "root"
	DB_PASS = "root"
	DB_NAME = "dbname"
)
</code></pre>

2、type this
//table_name 数据库表名
<pre><code>
go run generate.go  table_name  
</code></pre>


