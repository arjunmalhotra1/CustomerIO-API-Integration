# CustomerIO Api Integration.
The goal of the
script is to sync user data from local files to the Customer.io API.
The descriptions of `data.json` and `configuration.json` can be found below.

## Requirements
User data is upserted into Customer.io.
If a user already exists, then their properties are updated (as configured
by mappings).
If a user doesn't exist, then the user is inserted into Customer.io with an initial set of
properties (as configured by mappings).

After running the script, one can lookup users in Customer.io and
see that they all exist and have attributes set as defined from the input files. The
specifics of how the syncing occurs is documented below.

## Inputs
The script will take in two inputs.
Configuration file - contains configuration about how to sync the data.
Data file - contains the user data to be synced.
Both files should be in JSON format and passed in via command line parameters.

## Functionalities
### Retrying
The code implements some basic retrying. Handles error cases for
each of the following two scenarios.
Retryable errors
Non-retryable errors

If the Customer.io API tells whether or not an error is retryable, it uses that.

### Idempotency
After calling this script once with a set of input files, you should be able to call it
again any number of times without any functional changes to the data in
Customer.io.

## To run the script:
 > go run main.go -config [config-file-path] -data [data-file-path]
 
example:
 > go run main.go -config ./configuration.json -data ./data.json


## Configuration File - configuration.json
The configuration file will define how the user data will be synced into
Customer.io. 

There are three top-level keys that configuration files includes.
### parallelism - controls the number of API requests to make in parallel
The parallelism is set to 25, meaning that
we should be making 25 API requests concurrently.
### userId - chooses the key from the data file to use as the Customer.io user ID
### mappings - configure how fields from user data map to attributes in Customer.io

Additionally, we can add key(s) for Customer.io API authentication into this file.
These must be read in from the configuration file.

Example of Configuration file can be found below:
The key id should be used as the user ID sent to Customer.io. The key
computed_ltv from the source data should end up as user attribute ltv and the key
name from the source data should end up as the user attribute name .

```
{
"parallelism": 25,
"userId": "id",
"mappings": [
{
"from": "computed_ltv",
"to": "ltv"
},
{
"from": "name",
"to": "name"
}
]
}
```

## Data File - data.json
Data file
The data file contains a JSON array of user data. Each object in the array
represents an individual user. The data making up a user can be completely
different from one data file to another.