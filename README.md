# Hightouch

## To run the script:
 > go run main.go -config [config-file-path] -data [data-file-path]
 
example:
 > go run main.go -config ./configuration.json -data ./data.json

### Findings
With URI as https://track.customer.io/api/v1sdasdasdasdasd/customers 
& invalid Site_id/Api_key, response is 200 OK.
A not valid path should return 404.

### Product Question
Yes it is supported in the customer.io
The documentation [here](https://customer.io/docs/api/#operation/identify) states that having
 > set _update:true

Customer.io will not create a new profile, even if the identifier in the path isn't found. Hence allowing only updates.
We would need to add the above line for each of the users in the data.config file.
On the scale of 1-10 this is not hacky since this is how customer.io has specifications defined so I would say 1.
Edge cases I would be concerned about would be what if a user gives "false" and expects customer.io to start inserting new users.