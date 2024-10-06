package controller

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Amangupta20000/urlShortner/model"
	"go.mongodb.org/mongo-driver/bson"
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

// helper functions
func generateShortURL(original_url string) string {
	hasher := md5.New()
	hasher.Write([]byte(original_url))
	data := hasher.Sum(nil)

	hash := hex.EncodeToString(data)
	fmt.Println("hasher : ", hash[:8])
	return hash[:8]
}

func createURL(original_url string) (model.URL, error) {
	shortURL := generateShortURL(original_url)
	id := shortURL // use short url for id for simplicity

	response := model.URL{
		ID:           id,
		OriginalURL:  original_url,
		ShortURL:     shortURL,
		CreationDate: time.Now(),
	}

	inserted, err := collection.InsertOne(context.Background(), response)
	if err != nil {
		return model.URL{}, fmt.Errorf("failed to insert URL: %v", err)
	}

	fmt.Println("Inserted 1 URL with id: ", inserted.InsertedID)

	// Return the inserted object
	return response, nil
}

func getURL(id string) (model.URL, error) {
	// Define the filter to search by the custom 'id' field (which is a string)
	filter := bson.M{"id": id}

	var urlObj model.URL

	// Find the URL in the collection
	err := collection.FindOne(context.Background(), filter).Decode(&urlObj)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			// No document was found, handle it accordingly
			return model.URL{}, errors.New("URL not Found")
		}
		// Other errors
		return model.URL{}, fmt.Errorf("failed to fetch URL: %v", err)
	}

	return urlObj, nil
}
