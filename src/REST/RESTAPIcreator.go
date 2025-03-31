package REST

import (
	"SWIFT/src/databaseControl"
	"SWIFT/src/structs"
	
	"fmt"
	"strings"
	"errors"
	"encoding/json"
	"github.com/go-playground/validator/v10"

	"net/http"

	"github.com/gin-gonic/gin"
)

// local wrappers for swagger 
type BranchResponse structs.ReqBranch
type HeadquarterResponse structs.ReqHeadquarter
type CountryResponse structs.ReqCountry
type AllResponse structs.ReqAll
type MessageResponse structs.ReqErr

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

// getAll godoc
// @Summary      Get all SWIFT entities (headquarters and branches)
// @Description  Retrieves all headquarters (without their subsequent branches) and all branches from the SWIFT database
// @Tags         swift-codes headquarter branch
// @Accept       json
// @Produce      json
// @Success      200  {object}  AllResponse  "List of SWIFT entries"
// @Failure      400  {object}  MessageResponse  "Invalid request"
// @Failure      404  {object}  MessageResponse  "No entries found"
// @Failure      500  {object}  MessageResponse  "Internal server error"
// @Router       /swift-codes [get]
func getAll(c *gin.Context) {
    entries, valid, err := databaseControl.GetAll("")
    if !valid || err != nil {
        c.JSON(http.StatusBadRequest, structs.ReqErr{
            Message: fmt.Sprintf("Error: %s", err.Error()),
        })
        return
    }
    
    c.JSON(http.StatusOK, entries)
}

// getBySWIFTcode godoc
// @Summary Get SWIFT entity
// @Description Returns the information about the headquarter (with subsequent branches) or the branch
// @Tags swift-codes
// @Accept json
// @Produce json
// @Param swift-code path string true "SWIFT code" example("MRWORLDWXXX")
// @Success 200 {object} BranchResponse "SWIFT branch entity" 
// @Success 200 {object} HeadquarterResponse "SWIFT headquarter entity" 
// @Example {"address":"Sesame Street 8","bankName":"WorldWide Bank at Cracow","countryISO":"PB","countryName":"Pitbulland","isHeadquarter":false,"swiftCode":"MRWORLDWIDE"}
// @Example {"address":"Sesame Street 8","bankName":"WorldWide Bank","countryISO":"PB","countryName":"Pitbulland","isHeadquarter":true,"swiftCode":"MRWORLDWXXX","branches":[{"address":"Sesame Street 8","bankName":"WorldWide Bank at Cracow","countryISO":"PB","countryName":"Pitbulland","isHeadquarter":false,"swiftCode":"MRWORLDWIDE"}]}
// @Failure 404  {object}  MessageResponse  "No country found"
// @Failure 500  {object}  MessageResponse  "Internal server error"
// @Router /swift-codes/{swift-code} [get]
func getBySWIFTcode(c *gin.Context) {
	SWIFTcode := c.Param("swift-code")

	if strings.Contains(SWIFTcode, "XXX") { // SWIFT belongs to the headquarter
		hq, valid, err := databaseControl.GetHeadquarter("", SWIFTcode)
		if !valid || err != nil {
			if strings.Contains(err.Error(), "no headquarter found") {
				c.JSON(http.StatusNotFound, structs.ReqErr{
					Message : fmt.Sprintf("Error: %s", err.Error()),
				})
				return
			}
			c.JSON(http.StatusBadRequest, structs.ReqErr{
				Message : fmt.Sprintf("Error: %s", err.Error()),
			})
			return
		}
		c.JSON(http.StatusOK, hq)
		return

	} else { // SWIFT belongs to the branch
		br, valid, err := databaseControl.GetBranch("", SWIFTcode)
		if !valid || err != nil {
			if strings.Contains(err.Error(), "no branch found") {
				c.JSON(http.StatusNotFound, structs.ReqErr{
					Message : fmt.Sprintf("Error: %s", err.Error()),
				})
				return
			}
			c.JSON(http.StatusBadRequest, structs.ReqErr{
				Message : fmt.Sprintf("Error: %s", err.Error()),
			})
			return
		}
		c.JSON(http.StatusOK, br)
		return
	}
}

// getAllFromCountry godoc
// @Summary      Get data about the country and its subsequent entries
// @Description  Retrieves a country's name, ISO2 and its subsequent headquarters (without the subsequent branches) and branches from the SWIFT database
// @Tags         country 
// @Accept       json
// @Produce      json
// @Param ISO2 path string true "ISO2" example("PL")
// @Success      200  {object}  CountryResponse  "Country exists"
// @Failure      400  {object}  MessageResponse  "Invalid request"
// @Failure      404  {object}  MessageResponse  "No country found"
// @Failure      500  {object}  MessageResponse  "Internal server error"
// @Router       /swift-codes/country/{ISO2} [get]
func getAllFromCountry(c *gin.Context) {
	iso2 := c.Param("ISO2")
	country, valid, err := databaseControl.GetCountry("", iso2)
	if !valid || err != nil {
		if strings.Contains(err.Error(), "no country found") {
			c.JSON(http.StatusNotFound, structs.ReqErr{
				Message : fmt.Sprintf("Error: %s", err.Error()),
			})
			return
		}
		c.JSON(http.StatusBadRequest, structs.ReqErr{
			Message : fmt.Sprintf("Error: %s", err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, country)
}

// addEntry godoc
// @Summary      Add a new entry
// @Description  Add new headquarter or branch to the SWIFT database
// @Tags         headquarter branch
// @Accept       json
// @Produce      json
// @Param        request body BranchResponse true "Entry data"
// @Success      200  {object}  MessageResponse  "Successful addition"
// @Failure      400  {object}  MessageResponse  "Invalid request"
// @Failure      500  {object}  MessageResponse  "Internal server error"
// @Router       /swift-codes [post]
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
            c.JSON(http.StatusBadRequest, structs.ReqErr{
                Message: fmt.Sprintf("Error: invalid JSON type in %s (expected %s)", 
				jsonErr.Field, jsonErr.Type.String()),
            })
        case errors.As(err, &syntaxErr):
            c.JSON(http.StatusBadRequest, structs.ReqErr{
                Message: fmt.Sprintf("Error: malformed JSON at %sth character", fmt.Sprint(syntaxErr.Offset)),
            })
        case errors.As(err, &bindErr):
			message := "Error: validation failed"
            for _, fieldErr := range bindErr {
                message = message + fmt.Sprintf(" - field '%s' doesn't satisfy the  '%s' condition",
				fieldErr.Field(), fieldErr.Tag())
            }
		
            c.JSON(http.StatusBadRequest, structs.ReqErr{
                Message: message,
            })
        default:
            c.JSON(http.StatusBadRequest, structs.ReqErr{
                Message: fmt.Sprintf("Error: Invalid request, %s error occured", err.Error()),
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
		c.JSON(http.StatusBadRequest, structs.ReqErr{
			Message: fmt.Sprintf("Error: %s", err.Error()),
		})
		return
	}
	c.JSON(http.StatusCreated, structs.ReqErr{
		Message: fmt.Sprintf("New %s has been added succesfully", entry_type),
	})
}

// deleteEntry godoc
// @Summary      Delete a SWIFT determined entry
// @Description  Delete a headquarter (without the subsequent branches) or branch specified by the SWIFT code
// @Tags         headquarter branch 
// @Accept       json
// @Produce      json
// @Param swift-code path string true "SWIFT code" example("MRWORLDWXXX")
// @Success      200  {object}  MessageResponse  "Entry successfully deleted"
// @Failure      400  {object}  MessageResponse  "Invalid request"
// @Failure      404  {object}  MessageResponse  "No entry found"
// @Failure      500  {object}  MessageResponse  "Internal server error"
// @Router       /swift-codes/{swift-code} [delete]
func deleteEntry(c *gin.Context) {
	SWIFTcode := c.Param("swift-code")
	deleted, err := databaseControl.DeleteEntry("", SWIFTcode)
	if !deleted || err != nil {
		if strings.Contains(err.Error(), "no entry found") {
			c.JSON(http.StatusNotFound, structs.ReqErr{
				Message: fmt.Sprintf("Error: %s", err.Error()),
			})
			return
		}
		c.JSON(http.StatusBadRequest, structs.ReqErr{
			Message: fmt.Sprintf("Error: %s", err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, structs.ReqErr{
		Message : "Entry has been deleted succesfully",
	})
}