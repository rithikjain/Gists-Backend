package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/rithikjain/GistsBackend/api/handler"
	"github.com/rithikjain/GistsBackend/pkg/gists"
	"github.com/rithikjain/GistsBackend/pkg/user"
	"log"
	"net/http"
	"os"
)

func dbConnect(host, port, user, dbname, password, sslmode string) (*gorm.DB, error) {
	// In the case of heroku
	if os.Getenv("DATABASE_URL") != "" {
		return gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	}
	db, err := gorm.Open(
		"postgres",
		fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", host, port, user, dbname, password, sslmode),
	)
	return db, err
}

func GetPort() string {
	var port = os.Getenv("PORT")
	if port == "" {
		fmt.Println("INFO: No PORT environment variable detected, defaulting to 3000")
		return "localhost:3000"
	}
	return ":" + port
}

func main() {
	if os.Getenv("onServer") != "True" {
		// Loading the .env file
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	// Setting up DB
	db, err := dbConnect(
		os.Getenv("dbHost"),
		os.Getenv("dbPort"),
		os.Getenv("dbUser"),
		os.Getenv("dbName"),
		os.Getenv("dbPass"),
		os.Getenv("sslmode"),
	)
	if err != nil {
		log.Fatalf("Error connecting to the database: %s", err.Error())
	}

	db.AutoMigrate(&user.User{})

	defer db.Close()
	fmt.Println("Connected to DB...")
	db.LogMode(true)

	// Setting up the router
	r := http.NewServeMux()

	gistsSvc := gists.NewService(db)
	handler.MakeGistsHandler(r, gistsSvc)

	userRepo := user.NewRepo(db)
	userSvc := user.NewService(userRepo)
	handler.MakeUserHandler(r, userSvc)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Hello There"))
		return
	})
	fmt.Println("Serving...")
	log.Fatal(http.ListenAndServe(GetPort(), r))
}
