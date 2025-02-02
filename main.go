package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type item struct {
	ShortDescription string `json:"shortDescription"`
	Price string `json:"price"`
}

type receipt struct {
  Points int `json:"points"`
  Retailer string `json:"retailer"`
  PurchaseDate string `json:"purchaseDate"`
  PurchaseTime string `json:"purchaseTime"`
  Items []item `json:"items"`
  Total string `json:"total"`
}

type postResponse struct {
	ID string `json:"id"`
}

type getResponse struct {
	Points int `json:"points"`
}

var receipts = make(map[string]receipt)

func calculatePoints(newReceipt receipt){
	
}

func main() {
	router := gin.Default()
	router.GET("/receipts/:id/points", getReceiptPointsById)
	router.POST("/receipts/process", postReceipt)

	router.Run("localhost:8080")
}


func getReceiptPointsById(context *gin.Context){
	id := context.Param("id")

	val, ok := receipts[id]

	if(ok){
		context.IndentedJSON(http.StatusOK, val)
	} else {
		context.IndentedJSON(http.StatusNotFound, "No receipt found for that ID.")
	}

}

//GET curl http://localhost:8080/receipts/f5a0773e-3199-4b0b-977c-99adfdf20a92/points
//POST curl

func postReceipt(context *gin.Context){
	var newReceipt receipt

	if err := context.BindJSON(&newReceipt); err != nil {
		return
	}

	calculatePoints(newReceipt)

	var id = uuid.New().String()

	receipts[id] = newReceipt

	context.IndentedJSON(http.StatusCreated, postResponse{ID:id})
}

