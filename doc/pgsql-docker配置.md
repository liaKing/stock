# PostgreSQL Docker 本地运行配置

## 1. 启动容器

若本地没有镜像，可先拉取：

```bash
docker pull postgres:16
```

在项目根目录执行（已使用与 `doc/schema_pgsql.sql` 一致的库名）：

```bash
docker run -d \
  --name stock-pgsql \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=123456 \
  -e POSTGRES_DB=stock_db \
  -p 5432:5432 \
  postgres:16
```

## 2. 连接配置

| 配置项   | 值          |
|----------|-------------|
| Host     | localhost   |
| Port     | 5432        |
| User     | postgres    |
| Password | 123456      |
| Database | stock_db    |

**连接字符串（DSN）：**

```
postgres://postgres:123456@localhost:5432/stock_db?sslmode=disable
```

**Go 常用格式：**

```
host=localhost port=5432 user=postgres password=123456 dbname=stock_db sslmode=disable
```

## 3. 常用命令

```bash
# 启动已有容器
docker start stock-pgsql

# 停止
docker stop stock-pgsql

# 删除容器（需先 stop）
docker stop stock-pgsql && docker rm stock-pgsql

# 进入 psql
docker exec -it stock-pgsql psql -U postgres -d stock_db
```

## 4. 初始化表结构

容器启动后，在项目根目录执行建表：

```bash
docker exec -i stock-pgsql psql -U postgres -d stock_db < doc/schema_pgsql.sql
```

或在已进入的 psql 中：

```sql
\i /path/to/doc/schema_pgsql.sql
```

（容器内需先将文件挂载或复制进去，推荐用上面的 `docker exec -i ... < file` 方式。）

## 5. 配置文件示例（供应用使用）

若应用通过配置文件连接 PostgreSQL，可参考（如 `conf/data.json` 增加 pgsql 段）：

```json
{
  "pgsql": {
    "host": "localhost",
    "port": 5432,
    "user": "postgres",
    "password": "123456",
    "dbname": "stock_db",
    "sslmode": "disable"
  }
},
```

---

*生成时间：2025-02-11*
