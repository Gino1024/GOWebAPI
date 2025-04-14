# 🛠️ golang-migrate 快速參考手冊

## 📦 安裝 migrate CLI

```bash
brew install golang-migrate
# 或使用 curl
curl -L https://github.com/golang-migrate/migrate/releases/latest/download/migrate.darwin-amd64.tar.gz | tar xvz
```

---

## 📁 建立 Migration 檔案

```bash
migrate create -ext sql -dir db/migration -seq create_users_table
```

### 🔧 檔名版本模式（三選一）：

| 選項         | 範例檔名                                  | 說明                      |
|--------------|--------------------------------------------|---------------------------|
| `-seq`       | `000001_create_users_table.up.sql`         | 預設遞增編號，易讀、開發方便 |
| `-timestamp` | `20250414114312_create_users_table.up.sql` | 用 Unix timestamp，更適合團隊與 CI/CD |
| `-version`   | `42_create_users_table.up.sql`             | 自定義版本號              |

---

## ⬆️ 執行 Up Migration

```bash
migrate -path db/migration -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" up
```

## ⬇️ 執行 Down Migration

```bash
migrate -path db/migration -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" down
```

## 🔁 回退一步

```bash
migrate ... down 1
```

## 🔄 重新執行最後一次 migration（回退 + 執行 up）

```bash
migrate ... redo
```

## 🔍 查看目前版本

```bash
migrate ... version
```

## 🧹 解決 dirty 狀態（強制設定版本）

```bash
migrate ... force 0
```

---

## 📌 Migration 檔案撰寫順序建議

### `init_schema.up.sql` 建議順序：
1. 建立資料表
2. 建立 index
3. 建立 foreign key

### `init_schema.down.sql` 建議順序（反順序）：
1. 移除 foreign key
2. 移除 index
3. 刪除資料表

---

## 💡 小技巧：

- `DROP IF EXISTS` + `CREATE IF NOT EXISTS` 可以避免 migration 出錯
- 不要在正式環境亂用 `force`，會跳過執行記錄（務必備份 DB）