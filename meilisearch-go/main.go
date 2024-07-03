package main

import (
	// "context"
	"fmt"
	"log"
    "net/http"
    "io/ioutil"
    "strings"
    "encoding/json"
	"time"

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
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "localhost", 5421, "myuser", "mypassword", "mydatabase")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	// for _, movie := range movies {
	// 	    db.Create(&movie)
	// 	}

	// if err := deleteAllMoviesFromMeiliSearch(); err != nil {
	// 	log.Fatalf("failed to delete all movies from MeiliSearch: %v", err)
	// }

	// Migrate the Movie schema
	db.AutoMigrate(&Movie{})
	fmt.Println("Database migration completed!")

	// Fetch movies from PostgreSQL
	var movies []Movie
	if err := db.Find(&movies).Error; err != nil {
		log.Fatalf("Error fetching movies from database: %v", err)
	}
	fmt.Printf("Fetched %d movies from PostgreSQL\n", len(movies))

	for _, movie := range movies {
		fmt.Printf("Movie ID: %d, Title: %s, Genres: %v\n", movie.ID, movie.Title, movie.Genres)
	}

	// Create index if it doesn't exist
	indexExists, err := checkIndexExists()
	if err != nil {
		log.Fatalf("Error checking if index exists: %v", err)
	}

	if !indexExists {
		taskUID, err := createIndex()
		if err != nil {
			log.Fatalf("Error creating index: %v", err)
		}

		// Wait for the index creation task to complete
		err = waitForTask(taskUID)
		if err != nil {
			log.Fatalf("Error waiting for index creation task: %v", err)
		}
	} else {
		fmt.Println("Index already exists, skipping index creation")
	}

	// movieID := uint(1) // The ID of the movie to delete
	// if err := deleteMovie(db, movieID); err != nil {
	// 	log.Fatalf("Error deleting movie: %v", err)
	// }

	// Index movies into MeiliSearch
	taskUID, err := indexMovies(movies)
	if err != nil {
		log.Fatalf("Error indexing movies: %v", err)
	}

	// Wait for the indexing task to complete
	err = waitForTask(taskUID)
	if err != nil {
		log.Fatalf("Error waiting for indexing task: %v", err)
	}

	// Perform a search in MeiliSearch
	searchQuery := "Drama"
	searchResponse, err := searchMovies(searchQuery)
	if err != nil {
		log.Fatalf("Error searching movies: %v", err)
	}

	fmt.Printf("Search Response: %+v\n", string(searchResponse))
}

func checkIndexExists() (bool, error) {
	url := fmt.Sprintf("%s/indexes/%s", meiliSearchURL, indexName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	} else if resp.StatusCode == http.StatusNotFound {
		return false, nil
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		return false, fmt.Errorf("unexpected response status: %s, body: %s", resp.Status, string(body))
	}
}

func createIndex() (int64, error) {
	url := fmt.Sprintf("%s/indexes", meiliSearchURL)
	payload := fmt.Sprintf(`{"uid": "%s"}`, indexName)

	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		return 0, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusAccepted {
		body, _ := ioutil.ReadAll(resp.Body)
		return 0, fmt.Errorf("unexpected response status: %s, body: %s", resp.Status, string(body))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading response body: %v", err)
	}

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, fmt.Errorf("error unmarshaling response: %v", err)
	}

	taskUID, ok := response["taskUid"].(float64)
	if !ok {
		return 0, fmt.Errorf("task UID not found in response")
	}

	fmt.Println("Index creation task enqueued successfully")
	return int64(taskUID), nil
}

func deleteAllMoviesFromMeiliSearch() error {
	url := fmt.Sprintf("%s/indexes/%s/documents", meiliSearchURL, indexName)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("unexpected response status: %s, body: %s", resp.Status, string(body))
	}

	fmt.Println("All movies deleted from MeiliSearch successfully")
	return nil
}


func indexMovies(movies []Movie) (int64, error) {
	url := fmt.Sprintf("%s/indexes/%s/documents", meiliSearchURL, indexName)

	// Prepare documents to send
	var documents []map[string]interface{}
	for _, movie := range movies {
		doc := map[string]interface{}{
			"id":     movie.ID,
			"title":  movie.Title,
			"genres": []string(movie.Genres),
		}
		documents = append(documents, doc)
	}

	// Convert documents to JSON
	payload, err := json.Marshal(documents)
	if err != nil {
		return 0, fmt.Errorf("error marshaling documents: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", url, strings.NewReader(string(payload)))
	if err != nil {
		return 0, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	// Send HTTP request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		body, _ := ioutil.ReadAll(resp.Body)
		return 0, fmt.Errorf("unexpected response status: %s, body: %s", resp.Status, string(body))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading response body: %v", err)
	}

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, fmt.Errorf("error unmarshaling response: %v", err)
	}

	taskUID, ok := response["taskUid"].(float64)
	if !ok {
		return 0, fmt.Errorf("task UID not found in response")
	}

	fmt.Println("Indexing task enqueued successfully")
	return int64(taskUID), nil
}

func waitForTask(taskUID int64) error {
	url := fmt.Sprintf("%s/tasks/%d", meiliSearchURL, taskUID)

	for {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return fmt.Errorf("error creating request: %v", err)
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("error sending request: %v", err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error reading response body: %v", err)
		}

		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		if err != nil {
			return fmt.Errorf("error unmarshaling response: %v", err)
		}

		status, ok := response["status"].(string)
		if !ok {
			return fmt.Errorf("status not found in response")
		}

		if status == "succeeded" {
			fmt.Printf("Task %d completed successfully\n", taskUID)
			return nil
		} else if status == "failed" {
			taskError, _ := json.Marshal(response)
			return fmt.Errorf("task %d failed, details: %s", taskUID, string(taskError))
		}

		time.Sleep(1 * time.Second) // Wait before checking again
	}
}

func searchMovies(query string) ([]byte, error) {
	url := fmt.Sprintf("%s/indexes/%s/search?q=%s", meiliSearchURL, indexName, query)

	// Create HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

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