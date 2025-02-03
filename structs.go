package main

import (
  "strconv"
  //"math"
 // "time"
  "regexp"
  "errors"
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
  	PurchaseTime float32 `json:"purchaseTime"`
  	Items []item `json:"items"`
  	Total float64 `json:"total"`
}

type postResponse struct {
	ID string `json:"id"`
}

type getResponse struct {
	Points int `json:"points"`
}

func (request *receiptRequest) convertReceiptRequestToReceipt() (*receipt, error){
  moneyMatcher, _ := regexp.Compile("^\\d+\\.\\d{2}$")
  stringMatcher, _ := regexp.Compile("^[\\w\\s\\-&]+$")

  if(!stringMatcher.MatchString(request.Retailer)){
    return nil, errors.New("Retailer is not valid.")
  }

  if(!moneyMatcher.MatchString(request.Total)){
    return nil, errors.New("Total is not valid.")
  }

  if(len(request.Items) < 1){
    return nil, errors.New("Must have at least one item.")
  }

  convertedItems := []item{}

  for _, myItem := range request.Items{

    if (!stringMatcher.MatchString(myItem.ShortDescription)){
      return nil, errors.New("Description: " + myItem.ShortDescription + " is not valid.")
    }

    if (!moneyMatcher.MatchString(myItem.Price)){
      return nil, errors.New("Price: " + myItem.Price + " is not valid.")
    }

    price, err := strconv.ParseFloat(myItem.Price, 32)

    // Shouldn't happen but why not
    if(err != nil){
      return nil, err
    }

    convertedItems = append(convertedItems, item{ShortDescription:myItem.ShortDescription, Price:price})
  }

  /*input_date_layout := ""
  output_date_layout := ""

  input_time_layout := ""
  output_time_layout := ""
*/

  total, err := strconv.ParseFloat(request.Total, 64)

  if(err != nil){
    return nil, err
  }

  newReceipt := receipt{Points: 0}
  newReceipt.Total = total

  return &newReceipt, nil
}



