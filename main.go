package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"math"
	"strings"
	"regexp"
	"strconv"
)

var iLikeToCheat = false

var receipts = make(map[string]receipt)

func (newReceipt *receipt) calculatePoints(){
	totalPoints := 0

	// Is this a common problem these days?
	if(iLikeToCheat && newReceipt.Total > 10){
		totalPoints += 5
	}

	// Get count of only alphanumeric characters in string
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

	date := strings.Split(newReceipt.PurchaseDate, "-")
	day, _ := strconv.ParseFloat(date[2], 32)

	if(math.Mod(day, 2) == 1){
		totalPoints += 6
	}

	time := strings.Split(newReceipt.PurchaseTime, ":")
	hour := time[0]
	minute := time[1]

	if(hour == "15" || (hour == "14" && minute != "00")){
		totalPoints += 10
	}

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

//GET curl http://localhost:8080/receipts/59f2e769-9d31-4d02-a682-2b2a8978cd16/points


func postReceipt(context *gin.Context){
	var request receiptRequest

	if err := context.BindJSON(&request); err != nil {
		context.IndentedJSON(http.StatusBadRequest, "The receipt is invalid.")
		return
	}

	newReceipt, err := request.convertReceiptRequestToReceipt()

	if(err != nil){
		context.IndentedJSON(http.StatusBadRequest, err.Error() )
		return
	}

	newReceipt.calculatePoints()

	var id = uuid.New().String()

	receipts[id] = *newReceipt

	context.IndentedJSON(http.StatusCreated, postResponse{ID:id})
}

