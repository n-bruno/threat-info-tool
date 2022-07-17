# Threat Info Tool

## About
This is a simple data enrichment API for IP Addresses. It interfaces with AbuseIpDb's API, geo-location information, malicious activity history, and more.

## Usage

### Starting the server
```
./threatInfoTool -c /configFilePath.json -v
```
Use the `-h` flag for further information.

### Requests to the server
```
curl --request GET \
  --url 'http://localhost:8080/ipInfo?ipAddress=$IP_ADDRESS' \
  --header 'Key: INSERT_API_KEY_HERE'
```

**Note:** The API key specified in the header is not an API key for a third party tool, it's specific to this tool. This tool is designed to be publicly accessible to the world (for a brief period of time), as a result, having some sort of restrictive access is best.

## Config file
Included inside the repository is a config file. A default config file is provided, albeit lacks information.

The `apiKeys` property is blank. This is where the apiKeys are specified which would then be used in request's headers. Here's what the format looks like:
```json
"apiKeys":
{
    "user1": "KEY1",
    "user2": "KEY2",
},
```
Enable the program by setting the `enable` property at the top to `true`.

You can enable SSL by setting `enable` property under `ssl` to `true`. You'll then have to specify the certificate and private key file paths.

The tool is also reliant on [AbuseIPDB's API](https://www.abuseipdb.com/api). Register for free to receive an API key.

## Shortcomings
In order for this project to not get too out of hand given a short development cycle, the tool's API key management is primitive and is managed in the config file. Here's an example:

```json
"apiKeys":
{
    "user1": "KEY1",
    "user2": "KEY2",
},
```

If a username and API key is to be retired, the pair must be removed from the config file. The program must also be restarted for the change to take effect. Ideally, the keys would be stored in a database.

## Testing 
Inside the `"validate"` folder is a test file to test the ip address validation. I realize there's an `"net/http/httptest"` package to test endpoints, but alas.

## Deployment
I am deploying this program to Google Cloud. The process was straight forward and required me to follow this guide:

https://cloud.google.com/docs/terraform/get-started-with-terraform

I deployed with Terraform with the Cloud Shell. Once the VM is created, I sshed into the VM and deployed the application.