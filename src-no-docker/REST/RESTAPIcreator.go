package REST

import (
	"SWIFT/src-no-docker/databaseControl"
	"SWIFT/src-no-docker/structs"
	"fmt"

	"strings"
	"errors"
	"encoding/json"
	"github.com/go-playground/validator/v10"

	"net/http"

	"github.com/gin-gonic/gin"
)

func RunTheServer() *gin.Engine{
	gin.SetMode(gin.ReleaseMode)
	server := gin.Default()

	// specifying the routes for the server requests
	server.GET("/v1/swift-codes", getAll)
	server.GET("/v1/swift-codes/:swift-code", getBySWIFTcode)
	server.GET("/v1/swift-codes/country/:ISO2", getAllFromCountry)
	server.POST("/v1/swift-codes", addEntry)
	server.DELETE("/v1/swift-codes/:swift-code", deleteEntry)

	return server;
}

func getAll(c *gin.Context) {
	entries, valid, err := databaseControl.GetAll("")
	if !valid || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message" : fmt.Sprintf("Error: %s", err.Error())})
		return
	}
	c.JSON(http.StatusOK, entries)
}

func getBySWIFTcode(c *gin.Context) {
	SWIFTcode := c.Param("swift-code")

	if strings.Contains(SWIFTcode, "XXX") { // SWIFT belongs to the headquarter
		hq, valid, err := databaseControl.GetHeadquarter("", SWIFTcode)
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
		br, valid, err := databaseControl.GetBranch("", SWIFTcode)
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
	country, valid, err := databaseControl.GetCountry("", iso2)
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
	err := c.ShouldBindJSON(&newEntry)
	if err != nil {
        // Special handling for different error types
        var jsonErr *json.UnmarshalTypeError
        var syntaxErr *json.SyntaxError
        var bindErr validator.ValidationErrors
        
        switch {
        case errors.As(err, &jsonErr):
            c.JSON(http.StatusBadRequest, gin.H{
                "message": fmt.Sprintf("Error: invalid JSON type in %s (expected %s)", 
				jsonErr.Field, jsonErr.Type.String()),
            })
        case errors.As(err, &syntaxErr):
            c.JSON(http.StatusBadRequest, gin.H{
                "message": fmt.Sprintf("Error: malformed JSON at %sth character", syntaxErr.Offset),
            })
        case errors.As(err, &bindErr):
			message := "Error: validation failed"
            for _, fieldErr := range bindErr {
                message = message + fmt.Sprintf(" - field '%s' doesn't satisfy the  '%s' condition",
				fieldErr.Field(), fieldErr.Tag())
            }
		
            c.JSON(http.StatusBadRequest, gin.H{
                "message": message,
            })
        default:
            c.JSON(http.StatusBadRequest, gin.H{
                "message": fmt.Sprintf("Error: Invalid request, %s error occured", err.Error()),
            })
        }
        return
    }

	var added bool
	var entry_type string

	if *newEntry.IsHeadquarter {
		added, err = databaseControl.AddHeadquarter(databaseControl.Dsn, newEntry)
		entry_type = "headquarter"
	} else {
		added, err = databaseControl.AddBranch(databaseControl.Dsn, newEntry)
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