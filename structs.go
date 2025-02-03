package main

import (
  "strconv"
  "math"
  "regexp"
  "errors"
  "strings"
)

type itemRequest struct {
	ShortDescription string `json:"shortDescription"`
	Price string `json:"price"`
}

type item struct {
	ShortDescription string `json:"shortDescription"`
	Price float64 `json:"price"`
}

type receiptRequest struct {
	Retailer string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items []itemRequest `json:"items"`
	Total string `json:"total"`
}

type receipt struct {
	Points int `json:"points"`
	Retailer string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items []item `json:"items"`
	Total float64 `json:"total"`
}

type postResponse struct {
	ID string `json:"id"`
}

type getResponse struct {
	Points int `json:"points"`
}

// From yml file
var moneyMatcher, _ = regexp.Compile("^\\d+\\.\\d{2}$")
var stringMatcher, _ = regexp.Compile("^[\\w\\s\\-&]+$")

// Since googling how to handle dates / times in golang was getting irksome only the example formats
// specified in the yml will be accepted (i.e. YYYY-MM-DD for PurchaseDate and HH:MM for PurchaseTime)
var dateMatcher, _ = regexp.Compile("^\\d{4}-(0[1-9]|1[1,2])-(0[1-9]|[1,2][0-9]|3[0,1])$")
var timeMatcher, _ = regexp.Compile("^([0,1][0-9]|2[0-3]):([0-5][0-9])$")

// Validate everything other than the itemRequest objects
func (request *receiptRequest) ValidateNonItems() (error) {
  // I decided that all whitespace strings are not valid. A string made entirely of dashes is though. Does not have
  // to contain an alphanumeric character, simply not be entirely whitespace and adhere to the regex defined in the yml
  if(!stringMatcher.MatchString(request.Retailer) || len(strings.TrimSpace(request.Retailer)) == 0){
    return errors.New("Retailer is not valid.")
  }

  if(!moneyMatcher.MatchString(request.Total)){
    return errors.New("Total is not valid.")
  }

  if(!dateMatcher.MatchString(request.PurchaseDate)){
    return errors.New("Date is not valid. Must be in format YYYY-MM-DD.")
  }

  // Check edge case of being in February with a day that is too large
  date := strings.Split(request.PurchaseDate, "-")
  if(date[1] == "02"){

    // Here the compiler didn't like "if(date[2][0] == "3")" or string(date[2])[0] and 
    // I don't have an environment set up to debug. So we hard code even more =)
    if(date[2] == "30" || date[2] == "31"){
      return errors.New("February does not have 30+ days.")
    }

    if(date[2] == "29"){
      year, err := strconv.ParseFloat(date[0], 32)

      if(err != nil){
        return err
      }

      if(math.Mod(year, 4) != 0){
        return errors.New("February only has 29 days during leap years.")
      }
    }
  }

  if(!timeMatcher.MatchString(request.PurchaseTime)){
    return errors.New("Time is not valid. Must be in format HH:MM.")
  }

  return nil
}

func (request *receiptRequest) ValidateAndConvertItems() ([]item, error){
  if(len(request.Items) < 1){
    return nil, errors.New("Must have at least one item.")
  }

  convertedItems := []item{}

  for _, myItem := range request.Items{

    // I decided that all whitespace strings are not valid. A string made entirely of dashes is though. Does not have
    // to contain an alphanumeric character, simply not be entirely whitespace and adhere to the regex defined in the yml
    if (!stringMatcher.MatchString(myItem.ShortDescription) || len(strings.TrimSpace(myItem.ShortDescription)) == 0){
      return nil, errors.New("Description: " + myItem.ShortDescription + " is not valid.")
    }

    if (!moneyMatcher.MatchString(myItem.Price)){
      return nil, errors.New("Price: " + myItem.Price + " is not valid.")
    }

    price, err := strconv.ParseFloat(myItem.Price, 64)

    // Shouldn't happen but why not
    if(err != nil){
      return nil, err
    }

    convertedItems = append(convertedItems, item{ShortDescription:myItem.ShortDescription, Price:price})
  }

  return convertedItems, nil
}

func (request *receiptRequest) convertReceiptRequestToReceipt() (*receipt, error){
  err := request.ValidateNonItems()

  if(err != nil){
    return nil, err
  }

  convertedItems, err := request.ValidateAndConvertItems() 

  if(err != nil){
    return nil, err
  }

  total, err := strconv.ParseFloat(request.Total, 64)

  if(err != nil){
    return nil, err
  }

  // I'm sure there must be a better way in Go. Looking forward to learning what that is
  newReceipt := receipt{
    Points: 0, 
    Retailer: request.Retailer, 
    PurchaseDate: request.PurchaseDate, 
    PurchaseTime: request.PurchaseTime,
    Items: convertedItems,
    Total: total}

  return &newReceipt, nil
}



