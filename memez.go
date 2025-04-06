package main

import (
    "encoding/json"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "math/rand"
    "net/http"
    "os"
    "strconv"
    "time"
)

type Meme struct {
    ID    int    `json:"id"`
    Title string `json:"title"`
    URL   string `json:"url"`
}

var memes []Meme

func main() {
    if err := loadMemesFromFile("memes.json"); err != nil {
        panic("Failed to load memes: " + err.Error())
    }

    r := gin.Default()
    r.Use(cors.Default()) // allows access from websites

    r.GET("/memes", getAllMemes)
    r.GET("/memes/random", getRandomMeme)
    r.GET("/memes/:id", getMemeByID)
    r.POST("/memes", addMeme)

    r.Run(":8080") // runs on http://localhost:8080
}

func getAllMemes(c *gin.Context) {
    c.JSON(http.StatusOK, memes)
}

func getMemeByID(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    for _, meme := range memes {
        if meme.ID == id {
            c.JSON(http.StatusOK, meme)
            return
        }
    }

    c.JSON(http.StatusNotFound, gin.H{"error": "Meme not found"})
}

func getRandomMeme(c *gin.Context) {
    rand.Seed(time.Now().UnixNano())
    if len(memes) == 0 {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "No memes available"})
        return
    }
    meme := memes[rand.Intn(len(memes))]
    c.JSON(http.StatusOK, meme)
}

func addMeme(c *gin.Context) {
    var newMeme Meme
    if err := c.ShouldBindJSON(&newMeme); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
        return
    }

    newMeme.ID = getNextID()
    memes = append(memes, newMeme)

    if err := saveMemesToFile("memes.json"); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save meme"})
        return
    }

    c.JSON(http.StatusOK, newMeme)
}

func loadMemesFromFile(filename string) error {
    data, err := os.ReadFile(filename)
    if err != nil {
        return err
    }
    return json.Unmarshal(data, &memes)
}

func saveMemesToFile(filename string) error {
    data, err := json.MarshalIndent(memes, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(filename, data, 0644)
}

func getNextID() int {
    max := 0
    for _, meme := range memes {
        if meme.ID > max {
            max = meme.ID
        }
    }
    return max + 1
}
