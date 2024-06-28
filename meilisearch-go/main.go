package main

import (
	// "context"
	"fmt"
	"log"
    "net/http"
    "io/ioutil"
    "strings"
    "encoding/json"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
	// "github.com/meilisearch/meilisearch-go"
)

const (
    // host     = "localhost"  // or the Docker service name if running in another container
    // port     = 5421         // default PostgreSQL port
    // user     = "myuser"     // as defined in docker-compose.yml
    // password = "mypassword" // as defined in docker-compose.yml
    // dbname   = "mydatabase" // as defined in docker-compose.yml
    meiliSearchURL = "http://localhost:7700"  // Replace with your MeiliSearch server's URL
    apiKey         = "masterKey"    
    indexName      = "movies"  
)

func main() {
	// Connect to PostgreSQL database
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"localhost", 5421, "myuser", "mypassword", "mydatabase")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

    // for _, movie := range movies {
    //     db.Create(&movie)
    // }

	// Migrate the Movie schema
	db.AutoMigrate(&Movie{})
	fmt.Println("Database migration completed!")

    UpdateField(db)

	// // Index movies into MeiliSearch
	// err = indexMovies(movies)
	// if err != nil {
	// 	log.Fatalf("Error indexing movies: %v", err)
	// }

	// // Perform a search in MeiliSearch
	// searchQuery := "Drama"
	// searchResponse, err := searchMovies(searchQuery)
	// if err != nil {
	// 	log.Fatalf("Error searching movies: %v", err)
	// }

	// fmt.Printf("Search Response: %+v\n", searchResponse)
}

func UpdateField(db *gorm.DB) {
    for _, movie := range movies {
        var existingMovie Movie
        result := db.First(&existingMovie, movie.ID)
        if result.Error == nil {
            fmt.Printf("Movie '%s' already exists in the database\n", movie.Title)
        } else if result.Error == gorm.ErrRecordNotFound {
            if err := db.Create(&movie).Error; err != nil {
                log.Fatalf("Error creating movie record: %v", err)
            }
            fmt.Printf("Movie '%s' created successfully\n", movie.Title)
        } else {
            log.Fatalf("Error checking movie record: %v", result.Error)
        }
    }

    fmt.Println("Data insertion completed!")
}

// Function to index movies into MeiliSearch
func indexMovies(movies []Movie) error {
	url := fmt.Sprintf("%s/indexes/%s/documents", meiliSearchURL, indexName)

	// Prepare documents to send
	var documents []map[string]interface{}
	for _, movie := range movies {
		doc := map[string]interface{}{
			"id":     movie.ID,
			"title":  movie.Title,
			"genres": movie.Genres,
		}
		documents = append(documents, doc)
	}

	// Convert documents to JSON
	payload, err := json.Marshal(documents)
	if err != nil {
		return fmt.Errorf("error marshaling documents: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", url, strings.NewReader(string(payload)))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Meili-API-Key", apiKey)

	// Send HTTP request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	fmt.Println("Movies indexed successfully")
	return nil
}

// Function to search movies in MeiliSearch
func searchMovies(query string) ([]byte, error) {
	url := fmt.Sprintf("%s/indexes/%s/search?q=%s", meiliSearchURL, indexName, query)

	// Create HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("X-Meili-API-Key", apiKey)

	// Send HTTP request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	return body, nil
}