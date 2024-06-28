package main

import (
	// "gorm.io/gorm"
	// "time"
	// "github.com/golang-jwt/jwt/v4"
	// "github.com/valyala/fasthttp"
	"github.com/lib/pq"
)

type Book struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

// var books = []Book{
// 	{ID: 1, Title: "In Search of Lost Time", Author: "Marcel"},
// }

type Movie struct {
	ID     uint     `json:"id" gorm:"primaryKey"`
	Title  string   `json:"title"`
	Genres pq.StringArray `json:"genres" gorm:"type:text[]"`
}

var movies = []Movie{
    {ID: 1, Title: "Carol", Genres: pq.StringArray{"Romance", "Drama"}},
    {ID: 2, Title: "Wonder Woman", Genres: pq.StringArray{"Action", "Adventure"}},
    {ID: 3, Title: "Life of Pi", Genres: pq.StringArray{"Adventure", "Drama"}},
    {ID: 4, Title: "Mad Max: Fury Road", Genres: pq.StringArray{"Adventure", "Science Fiction"}},
    {ID: 5, Title: "Moana", Genres: pq.StringArray{"Fantasy", "Action"}},
    {ID: 6, Title: "Philadelphia", Genres: pq.StringArray{"Drama"}},
    {ID: 7, Title: "The Shawshank Redemption", Genres: pq.StringArray{"Drama"}},
    {ID: 8, Title: "The Dark Knight", Genres: pq.StringArray{"Action", "Crime", "Drama"}},
    {ID: 9, Title: "Forrest Gump", Genres: pq.StringArray{"Drama", "Romance"}},
    {ID: 10, Title: "The Matrix", Genres: pq.StringArray{"Action", "Sci-Fi"}},
    {ID: 11, Title: "Inception", Genres: pq.StringArray{"Action", "Adventure", "Sci-Fi"}},
    {ID: 12, Title: "Pulp Fiction", Genres: pq.StringArray{"Crime", "Drama"}},
    {ID: 13, Title: "The Lord of the Rings: The Fellowship of the Ring", Genres: pq.StringArray{"Action", "Adventure", "Drama"}},
    {ID: 14, Title: "Titanic", Genres: pq.StringArray{"Drama", "Romance"}},
    {ID: 15, Title: "Avatar", Genres: pq.StringArray{"Action", "Adventure", "Fantasy"}},
    {ID: 16, Title: "The Godfather", Genres: pq.StringArray{"Crime", "Drama"}},
    {ID: 17, Title: "Gladiator", Genres: pq.StringArray{"Action", "Drama"}},
    {ID: 18, Title: "The Avengers", Genres: pq.StringArray{"Action", "Adventure", "Sci-Fi"}},
    {ID: 19, Title: "Jurassic Park", Genres: pq.StringArray{"Action", "Adventure", "Sci-Fi"}},
    {ID: 20, Title: "The Silence of the Lambs", Genres: pq.StringArray{"Crime", "Drama", "Thriller"}},
    {ID: 21, Title: "Interstellar", Genres: pq.StringArray{"Adventure", "Drama", "Sci-Fi"}},
    {ID: 22, Title: "The Lion King", Genres: pq.StringArray{"Animation", "Adventure", "Drama"}},
    {ID: 23, Title: "Back to the Future", Genres: pq.StringArray{"Adventure", "Comedy", "Sci-Fi"}},
    {ID: 24, Title: "The Terminator", Genres: pq.StringArray{"Action", "Sci-Fi"}},
    {ID: 25, Title: "Die Hard", Genres: pq.StringArray{"Action", "Thriller"}},
}

// type Client struct {
// 	// config     ClientConfig
// 	httpClient *fasthttp.Client
// }

// type Index struct {
// 	UID        string    `json:"uid"`
// 	CreatedAt  time.Time `json:"createdAt"`
// 	UpdatedAt  time.Time `json:"updatedAt"`
// 	PrimaryKey string    `json:"primaryKey,omitempty"`
// 	client     *Client
// }

// type Key struct {
// 	Name        string    `json:"name"`
// 	Description string    `json:"description"`
// 	Key         string    `json:"key,omitempty"`
// 	UID         string    `json:"uid,omitempty"`
// 	Actions     []string  `json:"actions,omitempty"`
// 	Indexes     []string  `json:"indexes,omitempty"`
// 	CreatedAt   time.Time `json:"createdAt,omitempty"`
// 	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
// 	ExpiresAt   time.Time `json:"expiresAt"`
// }

type SearchRequest struct {
	Offset                  int64
	Limit                   int64
	AttributesToRetrieve    []string
	AttributesToSearchOn    []string
	AttributesToCrop        []string
	CropLength              int64
	CropMarker              string
	AttributesToHighlight   []string
	HighlightPreTag         string
	HighlightPostTag        string
	MatchingStrategy        string
	Filter                  interface{}
	ShowMatchesPosition     bool
	ShowRankingScore        bool
	ShowRankingScoreDetails bool
	Facets                  []string
	PlaceholderSearch       bool
	Sort                    []string
	Vector                  []float32
	HitsPerPage             int64
	Page                    int64
	IndexUID                string
	Query                   string
	Hybrid                  *SearchRequestHybrid
}

type SearchRequestHybrid struct {
	SemanticRatio float64
	Embedder      string
}

type MultiSearchRequest struct {
	Queries []SearchRequest `json:"queries"`
}

// type MeiliSearchBook struct {
//     ID          int64  `json:"id"`
//     Name        string `json:"name"`
//     Author      string `json:"author"`
//     Publisher   string `json:"publisher"`
//     Description string `json:"description"`
//     Writer      string `json:"writer"`
// }
