package REST

import (
	//"SWIFT/structs"

	"strings"

	"net/http"
	"github.com/gin-gonic/gin"
)

// User struct represents a user
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Sample data (simulating a database)
var users = []User{
	{ID: 1, Name: "Alice", Email: "alice@example.com"},
	{ID: 2, Name: "Bob", Email: "bob@example.com"},
}

func RunTheServer() *gin.Engine{
	server := gin.Default()

	// specifying the routes for the server requests
	server.GET("/v1/swift-codes/:swift-code", getBySWIFTcode)
	server.GET("/v1/swift-codes/country/:ISO2", getAllFromCountry)
	server.POST("/v1/swift-codes", addEntry)
	server.DELETE("/v1/swift-codes/:swift-code", deleteEntry)

	return server;
}

func getBySWIFTcode(c *gin.Context) {
	SWIFTcode := c.Param("swift-code")

	if strings.Contains(SWIFTcode, "XXX") { // SWIFT belongs to the headquarter

		c.JSON(http.StatusBadRequest, gin.H{"message" : "The headquarter with the requested SWIFT code doesn't figure in the database"})



	} else { // SWIFT belongs to the branch
		

		c.JSON(http.StatusBadRequest, gin.H{"message" : "The branch with the requested SWIFT code doesn't figure in the database"})
	}
}

func getAllFromCountry(c *gin.Context) {

}

func addEntry(c *gin.Context) {

}

func deleteEntry(c *gin.Context) {

}

// // Get all users
// func getUsers(c *gin.Context) {
// 	c.JSON(http.StatusOK, users)
// }

// // Get a user by ID
// func getUserByID(c *gin.Context) {
// 	id := c.Param("id")

// 	for _, user := range users {
// 		if id == string(rune(user.ID)) {
// 			c.JSON(http.StatusOK, user)
// 			return
// 		}
// 	}
// 	c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
// }

// // Create a new user
// func createUser(c *gin.Context) {
// 	var newUser User

// 	if err := c.ShouldBindJSON(&newUser); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Simulate ID assignment
// 	newUser.ID = len(users) + 1
// 	users = append(users, newUser)

// 	c.JSON(http.StatusCreated, newUser)
// }

// // Update a user by ID
// func updateUser(c *gin.Context) {
// 	id := c.Param("id")
// 	var updatedUser User

// 	if err := c.ShouldBindJSON(&updatedUser); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	for i, user := range users {
// 		if id == string(rune(user.ID)) {
// 			users[i].Name = updatedUser.Name
// 			users[i].Email = updatedUser.Email
// 			c.JSON(http.StatusOK, users[i])
// 			return
// 		}
// 	}

// 	c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
// }

// // Delete a user by ID
// func deleteUser(c *gin.Context) {
// 	id := c.Param("id")

// 	for i, user := range users {
// 		if id == string(rune(user.ID)) {
// 			users = append(users[:i], users[i+1:]...) // Remove user from slice
// 			c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
// 			return
// 		}
// 	}

// 	c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
// }