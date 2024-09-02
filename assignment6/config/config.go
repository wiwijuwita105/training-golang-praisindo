package config

const (
	//DBWalletDSN   = "postgresql://postgres:postgres@postgres-db-assignment6-wallet:5436/db_wallets?sslmode=disable"
	DBWalletDSN   = "postgresql://postgres:postgres@localhost:5436/db_wallets?sslmode=disable"
	DBUserDSN     = "postgresql://postgres:postgres@postgres-db-assignment6-user:5435/db_users?sslmode=disable"
	RedisHost     = "redis"
	RedisiPort    = "6378"
	RedisPassword = ""
	RedisDb       = 0
)
