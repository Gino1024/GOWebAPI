postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=@a123456 -d postgres:latest
createdb:
	docker exec -it postgres createdb --username=postgres --owner=postgres simple_bank
dropdb:
	docker exec -it postgres dropdb --username=postgres simple_bank
migrateup:
	migrate -path db/migration -database "postgresql://postgres:@a123456@localhost:5432/simple_bank?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://postgres:@a123456@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

#強制執行命令內指令, 預設會先檢查是否有對應檔案在目錄
.PHONY: postgres createdb dropdb migrateup migratedown sqlc
