# Hightouch

## To run the script:
 > go run main.go -config [config-file-path] -data [data-file-path]
 
example:
 > go run main.go -config ./configuration.json -data ./data.json

### Findings
With URI as https://track.customer.io/api/v1sdasdasdasdasd/customers 
& Site_id and Api_key invalid response is 200 OK
A not valid path should return 404.

### Product Question
Yes it is supported in the customer.io
the documentation [here](https://customer.io/docs/api/#operation/identify) state that having
 > set _update:true
Customer.io will not create a new profile, even if the identifier in the path isn't found.