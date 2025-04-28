package main

import (
	"fmt"
	"github.com/RamelHenderson/ggs/api/gamer"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"log"
	"sync"
	"time"

	"gorm.io/gorm"
)

var GormDBList []*gorm.DB

// Initializes the application
func init() {
	// Load environment variables from .env file
	godotenv.Load()

	// Command to start MySQL and Postgres containers
	// docker run --name ggs-mysql -p 3306:3306 -e MYSQL_DATABASE=ggs -e MYSQL_ROOT_PASSWORD=password -d mysql:latest && docker run --name ggs-postgres -p 5432:5432 -e POSTGRES_PASSWORD=password -e POSTGRES_USER=postgres -e POSTGRES_DB=ggs -d postgres:latest

	// Create a MySQL connection
	{
		// docker run --name ggs-mysql -p 3306:3306 -e MYSQL_DATABASE=ggs MYSQL_ROOT_PASSWORD=password -d mysql:latest
		mysqlDsn := "root:password@tcp(localhost:3306)/ggs?charset=utf8mb4&parseTime=True&loc=Local"
		MySqlDB, err := gorm.Open(mysql.Open(mysqlDsn), &gorm.Config{})
		if err != nil {
			panic("failed to connect to mysql database")
		}
		// Add MySQL DB to the list of GormDBList
		GormDBList = append(GormDBList, MySqlDB)
		e, err := GormDBList[0].DB()
		if err != nil || e.Ping() != nil {
			panic("failed to connect database")
		}
		log.Printf("[%s] Connected to MySQL database: %s", MySqlDB.Config.Name(), mysqlDsn)
	}

	// Create a Postgres connection
	{
		// docker run --name ggs-postgres -p 5432:5432 -e POSTGRES_PASSWORD=password -e POSTGRES_USER=postgres -e POSTGRES_DB=ggs -d postgres:latest
		postgresDSN := "host=localhost user=postgres password=password dbname=ggs port=5432 sslmode=disable TimeZone=UTC"
		PostgresDB, err := gorm.Open(postgres.Open(postgresDSN), &gorm.Config{})
		if err != nil {
			panic("failed to connect postgres database")
		}
		log.Printf("[%s] Connected to Postgres database: %s", PostgresDB.Config.Name(), postgresDSN)
		// Add Postgres DB to the list of GormDBList
		GormDBList = append(GormDBList, PostgresDB)
	}

	// Create a SQLite connection
	{
		sqlLiteDSN := "file::memory:?cache=shared"
		SQLiteDB, err := gorm.Open(sqlite.Open(sqlLiteDSN), &gorm.Config{})
		if err != nil {
			panic("failed to connect to sqlite database")
		}
		log.Printf("[%s] Connected to SQLite database: %s", SQLiteDB.Config.Name(), sqlLiteDSN)
		// Add SQLite DB to the list of GormDBList
		GormDBList = append(GormDBList, SQLiteDB)
	}

	// Create a function that will migrate and seed the database
	migrateAndSeed := func(gormDB *gorm.DB, seedCount int) {
		configName := gormDB.Config.Name()
		log.Printf("[%s] Starting migration and seeding", configName)

		log.Printf("[%s] Dropping tables...", configName)
		err := gormDB.Migrator().DropTable(&gamer.Gamer{})
		if err != nil {
			log.Fatalf("[%s]Error dropping table: %v", configName, err)
		}

		// Migrate the schema for MySQL
		log.Printf("[%s] Migrating tables...", configName)
		if err := gormDB.AutoMigrate(&gamer.Gamer{}); err != nil {
			log.Fatalf("[%s]Error migrating table: %v", configName, err)
		}

		// Seed the database with some data
		log.Printf("[%s] Seeding tables...", configName)
		for i := 0; i < seedCount; i++ {
			user := &gamer.Gamer{
				Name:     fmt.Sprintf("%s_user_%d", configName, i),
				Email:    fmt.Sprintf("%s_useremail_%d@email.com", configName, i),
				Password: fmt.Sprintf("%s_password-%d", configName, i),
			}
			if err := gormDB.Create(&user).Error; err != nil {
				log.Fatalf("[%s]Error seeding table: %v", configName, err)
			}
		}
	}

	// Create a wait group to wait for all migrations to finish
	wg := sync.WaitGroup{}
	// Add the number of databases to the wait group
	wg.Add(len(GormDBList))

	// Iterate over the list of GormDBList and run the migration and seeding in parallel
	for _, db := range GormDBList {
		go func() {
			timer := time.Now()
			defer wg.Done()
			// Migrate and seed the database for MySQL
			migrateAndSeed(db, 1000)
			log.Printf("[%s] ✅ Migration and seeding took <%s>", db.Config.Name(), time.Since(timer))
		}()
	}
	// Wait for all migrations to finish
	wg.Wait()

	log.Printf("🚀Initialization complete. All databases are migrated and seeded.")
}

// Main function
func main() {}
