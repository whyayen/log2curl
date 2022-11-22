# log2curl

log2curl is a tool that transforms logs on AWS Cloud Watch to cURL command easily with query-id(AWS given).
Something needs to be done before running this script:
1. Retrieve query-id with [start-query](https://docs.aws.amazon.com/cli/latest/reference/logs/start-query.html) via AWS CLI or API
2. Define the patterns to render cURL command

## How to retrieve query-id
There are several way to retrieve query-id (e.g, CLI, API, SDK). Here is an example to retrieve query-id via aws-cli.

```bash
aws logs start-query --log-group-names "/aws/apigateway/welcome" "/aws/lambda/Test01" \
--start-time 1598936400000 --end-time 1611464400000 --query-string "fields @timestamp, @message"
```

## Installation

```bash
wget https://github.com/whyayen/log2curl/raw/main/log2curl

# change permission
chmod u+x ./log2curl
```


## Getting Started
![log2curl](https://raw.githubusercontent.com/whyayen/log2curl/main/example.png)

Assume we have some logs on Cloud Watch. We get query results like below format when we search on Log Insights page on AWS Console or AWS CLI:
```json
{
    "method": "GET",
    "headers.scheme": "https",
    "host": "api.example.com",
    "path": "/users",
    "headers.Authorization": "Bearer abcdef12345",
    "parameters.user": "Lil Wayne",
    "parameters.page": "1",
    "parameters.per_page": "20"
}
```

And we have a query-id from executing `aws logs start-query`: 
```
4ffc3f36-2979-4558-88e8-dbe256d05d20
```

We can run log2curl to transform logs to cURL commands now:
```bash
$ log2curl --scheme "headers.scheme" --method "method" --host "host" --path "path" --headers-prefix "headers" --parameters-prefix "parameters" cloud-watch -q "4ffc3f36-2979-4558-88e8-dbe256d05d20"
```

Example response:
```
Start to get query results...
Finished. Save file in: log2curl.1669150193.txt
```

log2curl will generate a txt file in your current directory:
```bash
$ ls | grep 'log2curl'
log2curl.1669150193.txt
```

You can change the output by yourself:

```bash
$ log2curl -o ~/test.txt cloud-watch -q "4ffc3f36-2979-4558-88e8-dbe256d05d20"
```

### Generate config
Set patterns each time is annoying. We can generate a config file, and save our settings to config.

```bash
$ log2curl generate --config
```

Generate successfully
```
Generate default config to $HOME/.log2curl.json successfully
```

```bash
$ cat ~/.log2curl.json

{
  "custom": {
    "host": ""
  },
  "key": {
    "headers_prefix": "headers",
    "host": "host",
    "method": "method",
    "parameters_prefix": "parameters",
    "path": "path",
    "scheme": "scheme"
  },
  "whitelist_headers": [
    "Content-Type",
    "Authorization"
  ]
}

# Just modify configuration if you would like to customize it.
$ vim ~/.log2curl.json
```

### Whitelist Headers
We could decide which fields could be used in cURL if we have multiple headers field in the log.

For example, if we have the log:
```json
{
    "method": "GET",
    "headers.scheme": "https",
    "host": "api.example.com",
    "path": "/users",
    "headers.Authorization": "Bearer abcdef12345",
    "headers.User-Agent": "something...",
    "headers.Version": "HTTP/1.1",
    "headers.Content-Type": "application/json",
    "parameters.user": "Lil Wayne",
    "parameters.page": "1",
    "parameters.per_page": "20"
}
```

But our cURL request just need Authorization & Content-Type, we can set whitelist_headers in `~/.log2curl.json`:

```json
{
  ...,
  "whitelist_headers": [
    "Authorization",
    "Content-Type"
  ]
}
```

or

```bash
log2curl --whitelist-headers "Authorization,Content-Type" cloud-watch -q "4ffc3f36-2979-4558-88e8-dbe256d05d20"
```

### Custom Host
Sometimes we want to reproduce a **Production** log on **Staging**, we can replace host with `custom.host`.

```json
{
    "custom": {
        "host": "staging.example.com"
    },
    ...,
}
```

or

```bash
log2curl --custom-host "staging.example.com" cloud-watch -q "4ffc3f36-2979-4558-88e8-dbe256d05d20"
```

### Failed
Sometimes, log2curl can't parse field or transform to cURL. The result will be printed `Failed` in the file.

```
curl -X GET https://example.com/v1/users/page=1&per_page=50 \ 
 -H 'Authorization: Bearer JzCfIfzMGo'

Failed

Failed

curl -X GET https://example.com/v1/items \ 
 -H 'Authorization: Bearer Czx2341xa'
```

## Known Issues
- Just support AWS default profile now (It is unaviable to choose specific profile to connect AWS).
- All fields become string after transforming to cURL.
