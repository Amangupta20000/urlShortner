package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Database struct {
		DBName         string `yaml:"dbName"`
		CollectionName string `yaml:"collectionName"`
	} `yaml:"database"`
}

// MOST IMPORTANT
var collection *mongo.Collection
var config Config

func loadConfig() error {
	// Use os.ReadFile instead of ioutil.ReadFile
	yamlFile, err := os.ReadFile("config/config.yaml")
	if err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	// Unmarshal the YAML file into the config struct
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return fmt.Errorf("error unmarshalling config file: %v", err)
	}
	return nil
}

// Connect with MongoDB
func init() {

	// Get MongoDB URL from environment variables
	connectionString := os.Getenv("MONGODB_URL")
	if connectionString == "" {
		log.Fatal("MONGODB_URL not set in environment variables")
	}

	// Load the YAML config
	// err = loadConfig()
	err := loadConfig()
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	// Client options
	clientOptions := options.Client().ApplyURI(connectionString)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("MongoDB connection success")

	collection = client.Database(config.Database.DBName).Collection(config.Database.CollectionName)

	// If collection instance is ready
	fmt.Println("Collection instance/reference is ready")
}

func ShortenURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	json.NewEncoder(w).Encode("All Okay")
}
