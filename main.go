package main

import (
	"net/http"
	// "messages.go"
	"github.com/gin-gonic/gin"
)

type dog struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Breed string `json:"breed"`
	Owner string `json:"owner"`
}

var dogs = []dog{
	{ID: "10000", Name: "San", Breed: "Jack Russle", Owner: "Ori"},
	{ID: "10001", Name: "Star", Breed: "Border Collie", Owner: "Itai"},
	{ID: "10002", Name: "Murray", Breed: "Mixed", Owner: "Uria"},
	{ID: "10003", Name: "Chuchu", Breed: "Mixed", Owner: "Shira"},
}

const DogNotFound = "The dog with the given ID does not exist in the database."
const ConflictingIDs = "Conflicting IDs - please make sure that the ID in your URL is correct, and remove the ID field from your request body."

func (d *dog) Update(other dog) {
	if other.Name != "" {
		d.Name = other.Name
	}
	if other.Breed != "" {
		d.Breed = other.Breed
	}
	if other.Owner != "" {
		d.Owner = other.Owner
	}
}

func main() {
	router := gin.Default()
	router.GET("/dogs", getDogs)
	router.GET("/dogs/:id", getDogByID)
	router.POST("/dogs", createDog)
	router.PATCH("/dogs/:id", updateDog)
	router.DELETE("/dogs/:id", deleteDog)

	router.Run("localhost:8080")
}

// getDogs returns a JSON with a list of all dogs in the database.
func getDogs(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, dogs)
}

// getDogByID returns a JSON with the details of the dog with the given ID.
func getDogByID(c *gin.Context) {
	id := c.Param("id")

	for _, dog := range dogs {
		if dog.ID == id {
			c.IndentedJSON(http.StatusOK, dog)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": DogNotFound})
}

// createDog adds a new dog to the database.
func createDog(c *gin.Context) {
	var newDog dog

	// check for errors when binding JSON
	if err := c.BindJSON(&newDog); err != nil {
		return
	}

	// check that ID and name are provided
	if newDog.ID == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "ID must be provided when creating a new dog."})
		return
	} else if newDog.Name == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Name must be provided when creating a new dog."})
		return
	}

	dogs = append(dogs, newDog)
	c.IndentedJSON(http.StatusCreated, newDog)
}

// updateDog updates the entity of the dog with the given ID.
// empty fields will not be accepted.
func updateDog(c *gin.Context) {
	var updatedDog dog
	id := c.Param("id")

	if err := c.BindJSON(&updatedDog); err != nil {
		return
	}
	if updatedDog.ID != "" && updatedDog.ID != id {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": ConflictingIDs})
		return
	}

	for i, dog := range dogs {
		if dog.ID == id {
			dog.Update(updatedDog)
			c.IndentedJSON(http.StatusOK, dog)
			dogs[i] = dog
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": DogNotFound})
}

// deleteDog removes the dog with the given ID from the database.
func deleteDog(c *gin.Context) {
	id := c.Param("id")

	for i, dog := range dogs {
		if dog.ID == id {
			dogs = append(dogs[:i], dogs[i+1:]...)
			c.IndentedJSON(http.StatusOK, "")
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": DogNotFound})

}
