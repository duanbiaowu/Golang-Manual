package gorm

import (
	"log"
	"os"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func TestSqlite(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Println(err)
		t.Skip()
	}

	// Migrate the schema
	err = db.AutoMigrate(&Product{})
	if err != nil {
		log.Println(err)
		t.Skip()
	}

	// Create
	db.Create(&Product{Code: "D42", Price: 100})

	// Read
	var product Product
	db.First(&product, 1)                 // find product with integer primary key
	db.First(&product, "code = ?", "D42") // find product with code D42

	// Update - update product's price to 200
	db.Model(&product).Update("Price", 200)
	// Update - update multiple fields
	db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // non-zero fields
	db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - delete product
	db.Delete(&product, 1)

	err = os.Remove("test.db")
	if err != nil {
		log.Println("Remove file", err)
	}
}

func TestMySql(t *testing.T) {
	dsn := "root:123456@tcp(:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println(err)
		t.Skip()
	}

	// Migrate the schema
	//err = db.AutoMigrate(&Product{})
	err = db.Set("gorm:table_options", "CHARSET=utf8mb4").AutoMigrate(&Product{})
	if err != nil {
		log.Println(err)
		t.Skip()
	}

	// Create
	db.Create(&Product{Code: "D42", Price: 100})

	// Read
	var product Product
	db.First(&product, 1)                 // find product with integer primary key
	db.First(&product, "code = ?", "D42") // find product with code D42

	// Update - update product's price to 200
	db.Model(&product).Update("Price", 200)
	// Update - update multiple fields
	db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // non-zero fields
	db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - delete product
	db.Delete(&product, 1)

	// Truncate
	db.Create(&Product{Code: "D42", Price: 100})
	db.Create(&Product{Code: "D43", Price: 200})
	db.Exec("TRUNCATE TABLE products")

	// Drop Table
	err = db.Migrator().DropTable(&Product{})
	if err != nil {
		log.Println(err)
	}
}
