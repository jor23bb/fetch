package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"math"
	"strings"
	"regexp"
	//"time"
)

var iLikeToCheat = false

var receipts = make(map[string]receipt)

func (newReceipt *receipt) calculatePoints(){
	totalPoints := 0

	// Is this a common problem these days?
	if(iLikeToCheat && newReceipt.Total > 10){
		totalPoints += 5
	}

	re := regexp.MustCompile(`[a-zA-Z0-9]`)
	totalPoints += len(re.FindAllString(newReceipt.Retailer, -1))

	if(math.Mod(newReceipt.Total, .25) == 0){
		totalPoints += 25
	}

	// Even though a round dollar means that it will always get the 25 points for being a multiple of .25
	// this is separated out so that if the point value for being a multiple .25 ever changed then the code would
	// only need to be updated in one place
	if(math.Mod(newReceipt.Total, 1.0) == 0){
		totalPoints += 50
	}

	totalPoints += 5 * (len(newReceipt.Items) / 2)

	for _, item := range newReceipt.Items{
		totalPoints += calculateItemDescriptionPoints(item)
	}

	//if()

	newReceipt.Points = totalPoints
}

func calculateItemDescriptionPoints(currItem item) int{
	if(math.Mod(float64(len(strings.TrimSpace(currItem.ShortDescription))), 3) == 0){
		return int(math.Ceil(currItem.Price * .2))
	}
	return 0
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
		context.IndentedJSON(http.StatusOK, getResponse{Points:val.Points})
	} else {
		context.IndentedJSON(http.StatusNotFound, "No receipt found for that ID.")
	}
}

//GET curl http://localhost:8080/receipts/98f6f6e3-62c8-4722-8547-2cf13135ed1d/points


func postReceipt(context *gin.Context){
	var request receiptRequest

	if err := context.BindJSON(&request); err != nil {
		context.IndentedJSON(http.StatusBadRequest, "The receipt is invalid.")
		return
	}

	newReceipt, err := request.convertReceiptRequestToReceipt()

	if(err != nil){
		context.IndentedJSON(http.StatusBadRequest, "The receipt is invalid." /*err.Error()*/ )
		return
	}

	newReceipt.calculatePoints()

	var id = uuid.New().String()

	receipts[id] = *newReceipt

	context.IndentedJSON(http.StatusCreated, postResponse{ID:id})
}

