package example

import (
	"github.com/gin-gonic/gin"
	"github.com/kundank78/rate-limiter/pkg"
	"log"
	"net/http"
)

type Album struct {
	Id         string  `json:"id"`
	AlbumName  string  `json:"albumName"`
	ArtistName string  `json:"artistName"`
	Price      float64 `json:"price"`
}

var albums = []Album{
	{Id: "1", AlbumName: "Blue Train", ArtistName: "John Coltrane", Price: 56.99},
	{Id: "2", AlbumName: "Jeru", ArtistName: "Gerry Mulligan", Price: 17.99},
	{Id: "3", AlbumName: "Sarah Vaughan and Clifford Brown", ArtistName: "Sarah Vaughan", Price: 39.99},
}

func AlbumServer() {
	router := gin.Default()

	limiter := pkg.Init(pkg.SlidingWindowDistributedAlgo, 10, 60)

	router.GET("/albums", pkg.RateLimiter(limiter, getAlbums))
	router.GET("/albums/:id", pkg.RateLimiter(limiter, getAlbumById))
	router.POST("/album", pkg.RateLimiter(limiter, postAlbum))

	err := router.Run("localhost:8000")
	if err != nil {
		log.Fatalln("Error starting the http server", err)
	}
}

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func getAlbumById(c *gin.Context) {
	id := c.Param("id")

	for _, album := range albums {
		if album.Id == id {
			c.IndentedJSON(http.StatusOK, album)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func postAlbum(c *gin.Context) {
	var album Album
	if err := c.BindJSON(album); err != nil {
		c.IndentedJSON(http.StatusBadRequest, nil)
	}
	albums = append(albums, album)
	c.IndentedJSON(http.StatusCreated, album)
}
