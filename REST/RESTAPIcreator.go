package REST

import (
	"SWIFT/databaseControl"
	"SWIFT/structs"
	"fmt"

	"strings"

	"net/http"

	"github.com/gin-gonic/gin"
)

func RunTheServer() *gin.Engine{
	gin.SetMode(gin.ReleaseMode)
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
		hq, valid, err := databaseControl.GetHeadquarter(SWIFTcode)
		if !valid || err != nil {
			if strings.Contains(err.Error(), "no headquarter found") {
				c.JSON(http.StatusNotFound, gin.H{"message" : fmt.Sprintf("Error: %s", err.Error())})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"message" : fmt.Sprintf("Error: %s", err.Error())})
			return
		}
		c.JSON(http.StatusOK, hq)
		return

	} else { // SWIFT belongs to the branch
		br, valid, err := databaseControl.GetBranch(SWIFTcode)
		if !valid || err != nil {
			if strings.Contains(err.Error(), "no branch found") {
				c.JSON(http.StatusNotFound, gin.H{"message" : fmt.Sprintf("Error: %s", err.Error())})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"message" : fmt.Sprintf("Error: %s", err.Error())})
			return
		}
		c.JSON(http.StatusOK, br)
		return
	}
}

func getAllFromCountry(c *gin.Context) {
	iso2 := c.Param("ISO2")
	country, valid, err := databaseControl.GetCountry(iso2)
	if !valid || err != nil {
		if strings.Contains(err.Error(), "no country found") {
			c.JSON(http.StatusNotFound, gin.H{"message" : fmt.Sprintf("Error: %s", err.Error())})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"message" : fmt.Sprintf("Error: %s", err.Error())})
		return
	}
	c.JSON(http.StatusOK, country)
}

func addEntry(c *gin.Context) {
	var newEntry structs.ReqBranch

	if err := c.ShouldBindJSON(&newEntry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Error: %s", err.Error())})
		return
	}

	var added bool
	var err error
	var entry_type string

	if newEntry.IsHeadquarter {
		added, err = databaseControl.AddHeadquarter(newEntry)
		entry_type = "headquarter"
	} else {
		added, err = databaseControl.AddBranch(newEntry)
		entry_type = "branch"
	}
	
	if !added || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Error: %s", err.Error())})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": fmt.Sprintf("New %s has been added succesfully", entry_type)})
}

func deleteEntry(c *gin.Context) {
	SWIFTcode := c.Param("swift-code")
	deleted, err := databaseControl.DeleteEntry(SWIFTcode)
	if !deleted || err != nil {
		if strings.Contains(err.Error(), "no entry found") {
			c.JSON(http.StatusNotFound, gin.H{"message" : fmt.Sprintf("Error: %s", err.Error())})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"message" : fmt.Sprintf("Error: %s", err.Error())})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message" : "Entry has been deleted succesfully"})
}