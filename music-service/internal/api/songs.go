package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"music-service/internal/db"
	"music-service/internal/logging"
	"music-service/internal/model"
	"net/http"
	"strconv"
	"strings"
)

func Get(c *gin.Context) {
	logging.Logger.Debug("Getting a list of songs")

	var songs []model.Song
	query := db.DB

	if group := c.Query("group"); group != "" {
		logging.Logger.Debugf("Applying Group filtering: %s", group)
		query = query.Where("group = ?", group)
	}

	if song := c.Query("song"); song != "" {
		logging.Logger.Debugf("Applying filtering by song name: %s", song)
		query = query.Where("song = ?", song)
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	logging.Logger.Debugf("Using Pagination: page=%d, limit=%d", page, limit)

	query = query.Offset(offset).Limit(limit).Find(&songs)
	if query.Error != nil {
		logging.Logger.Error("Error getting the list of songs: ", query.Error)
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": "Error receiving songs"})
		return
	}

	logging.Logger.Debugf("Song list received successfully")
	c.JSON(http.StatusOK, songs)
}

func GetText(c *gin.Context) {
	var song model.Song
	id := c.Param("id")
	if err := db.DB.First(&song, id).Error; err != nil {
		c.JSON(http.StatusNotFound,
			gin.H{"error": "Song not found"})
		return
	}

	paragraphs := splitText(song.Text)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit
	if offset > len(paragraphs) {
		c.JSON(http.StatusOK, gin.H{"data": []string{}})
		return
	}

	end := offset + limit
	if end > len(paragraphs) {
		end = len(paragraphs)
	}

	c.JSON(http.StatusOK,
		gin.H{"data": paragraphs[offset:end]})
}

func splitText(text string) []string {
	return strings.Split(text, "\n\n")
}

func Add(c *gin.Context) {
	logging.Logger.Debug("Adding a song")

	var newSong model.Song
	if err := c.ShouldBind(&newSong); err != nil {
		logging.Logger.Warn("Invalid data received for adding a song")
		c.JSON(http.StatusBadRequest,
			gin.H{"error": "Invalid data format"})
		return
	}

	logging.Logger.Debugf("Adding a new song: group=%s, name=%s",
		newSong.Group, newSong.SongName)

	externalAPI := fmt.Sprintf("http://localhost:8080/info?group=%s&song=%s", newSong.Group, newSong.SongName)
	resp, err := http.Get(externalAPI)
	if err != nil || resp.StatusCode != http.StatusOK {
		logging.Logger.Error("Error accessing the external API")
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": "Error accessing the external API"})
		return
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var songDetail struct {
		ReleaseDate string `json:"releaseDate"`
		Text        string `json:"text"`
		Link        string `json:"link"`
	}

	if err := json.Unmarshal(body, songDetail); err != nil {
		logging.Logger.Error("Error when processing data from the API")
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": "Error when processing data from the API"})
		return
	}

	newSong.ReleaseDate = songDetail.ReleaseDate
	newSong.Text = songDetail.Text
	newSong.Link = songDetail.Link

	if err := db.DB.Create(&newSong).Error; err != nil {
		logging.Logger.Error("Error adding a song")
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": "Error when adding a song to the database"})
		return
	}

	logging.Logger.Info("Song added successfully")
	c.JSON(http.StatusCreated, newSong)
}

func Update(c *gin.Context) {
	var song model.Song
	id := c.Param("id")

	if err := db.DB.First(&song, id).Error; err != nil {
		c.JSON(http.StatusNotFound,
			gin.H{"error": "Song not found"})
		return
	}

	var updatedSong model.Song
	if err := c.ShouldBindJSON(&updatedSong); err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{"error": "Invalid data format"})
		return
	}

	song.Group = updatedSong.Group
	song.SongName = updatedSong.SongName
	song.ReleaseDate = updatedSong.ReleaseDate
	song.Text = updatedSong.Text
	song.Link = updatedSong.Link

	if err := db.DB.Save(&song).Error; err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": "Error updating a song"})
		return
	}

	c.JSON(http.StatusOK, song)
}

func Delete(c *gin.Context) {
	id := c.Param("id")

	if err := db.DB.Delete(&model.Song{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": "Error deleting a song"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Song deleted successfully"})
}

func Routes(router *gin.Engine) {
	router.GET("/songs", Get)
	router.GET("/songs/:id/text", GetText)
	router.POST("/songs", Add)
	router.PUT("/songs/:id", Update)
	router.DELETE("/songs/:id", Delete)
}
