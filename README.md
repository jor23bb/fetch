Fetch: Receipt Processor challenge

How to run (assuming Go installed):

Clone repo

Navigate to directory w/ main.go

(Possibly have to call "go get ." here. But I don't think so)

Call "go run ."

Make request(s)


Example GET request:

curl http://localhost:8080/receipts/8d63b0e3-2aaf-4af8-83bf-bdcbf14c0e31/points


Example POST request:

curl http://localhost:8080/receipts/process --include --header "Content-Type: application/json" --request "POST" --data '{"retailer": "Target","purchaseDate": "2022-01-01","purchaseTime": "13:01","items": [{"shortDescription": "Emils Cheese Pizza", "price": "12.25"}], "total": "35.43"}'
