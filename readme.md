# Fetch Receipt Processor Challenge
## Build Docker Image
```
docker build -t receipt-processor .
```
## Run Docker Image
```
docker run -p 8080:8080 -it --rm --name receipt-processor-container receipt-processor
```
## Send POST Request to /receipts/process
```
curl -H "Content-Type: application/json" -X POST --data-binary @examplePayloads/example1.json "http://localhost:8080/receipts/process"
```
```
curl -H "Content-Type: application/json" -X POST --data-binary @examplePayloads/example2.json "http://localhost:8080/receipts/process"
```