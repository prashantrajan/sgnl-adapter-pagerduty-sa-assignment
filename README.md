# SGNL Adapter for PagerDuty

The SGNL Adapter for PagerDuty is a gRPC-based service that integrates with PagerDuty's API to fetch and transform data into a format suitable for SGNL's ingestion service. This adapter is designed to handle pagination, validate requests, and provide seamless integration with with SGNL's adapter framework.


## Features

- **Entity Support**: Currently supports querying the `teams` entity from PagerDuty.
- **Pagination**: Implements classic pagination.
- **Validation**: Ensures proper configuration and validates incoming requests.
- **Extensibility**: Easily extendable to support additional PagerDuty entities.


## Prerequisites

- **Golang**: Version 1.21 or higher.
- **Docker**: For containerized builds and deployments (optional).
- **Postman**: For testing gRPC requests (optional).
- **Basic Knowledge**: Familiarity with gRPC, protocol buffers, and the PagerDuty API.


## Build

### Configure Authentication Tokens

Create `authTokens.json` file:

```sh
cp authTokens.json.example authTokens.json
```

Set the `AUTH_TOKENS_PATH` environment variable to the path of this file:

```sh
export AUTH_TOKENS_PATH=/path/to/authTokens.json
```

### Building a Binary

To build the adapter as a binary, use the following commands:

```sh
go build -o ./bin/sgnl-adapter-pagerduty ./cmd/adapter`
```

### Building a Docker Image

```sh
docker build -t sgnl-adapter-pagerduty:latest .
```


## Run

```sh
export AUTH_TOKENS_PATH=/path/to/authTokens.json

# Run main.go
go run cmd/adapter/main.go

# OR if you have a previously built binary, you can run
./bin/sgnl-adapter-pagerduty

# OR run as a Docker container
docker run -p 8080:8080 --rm -it -e AUTH_TOKENS_PATH=/local/path/to/authTokens.json sgnl-adapter-pagerduty:latest

```


## Testing

### Using Postman

Fork this public Postman workspace that contains collections for localhost and a live test server: {FIXME: POSTMAN PUBLIC WORKSPACE URL}

1. Define the [`GetPage` Protobuf](https://github.com/SGNL-ai/adapter-framework/blob/f2cafb0d963b54c350350967906ce59776d720a1/api/adapter/v1/adapter.proto) schema.

![Define the `GetPage` Protobuf definition](/docs/assets/postman_proto_definition.png)

2. Add the Pager Duty API token as a Workspace variable.

![Workspace Variables](/docs/assets/postman_workspace_variables.png)

3. Add the shared auth token (between SGNL and adapter) using Metadata. This token should match the token value stored in `authTokens.json`:

![Collection Metadata ](/docs/assets/postman_collection_metadata.png)

4. Select the GetPage Protobuf definition in the Service Definition.

![Collection Service Definition](/docs/assets/postman_collection_service_definition.png)

5. Add a Before Invoke Script to handle Base64 encoding.

![Collection Before Invoke Script](/docs/assets/postman_collection_before_invoke_script.png)

6. Select one of the Test Case (TC) messages to Invoke to see a live response.
![Create a new gRPC request](/docs/assets/postman_new_grpc_request.png)

### Unit Tests

Coming soon...


## References

### SGNL Adapter Template

The code in this repo was adapter from the [SGNL Adapter Template](https://github.com/SGNL-ai/adapter-template).

### PagerDuty API

API references:
- [Overview](https://developer.pagerduty.com/docs/rest-api-overview)
- [Authentication](https://developer.pagerduty.com/docs/authentication)
- [Pagination](https://developer.pagerduty.com/docs/pagination)
- [List Teams Endpoint](https://developer.pagerduty.com/api-reference/0138639504311-list-teams)


## Future Improvements

- Add unit tests for all components.
- Implement SSL/TLS for local development.
- Extend support for additional PagerDuty entities.
