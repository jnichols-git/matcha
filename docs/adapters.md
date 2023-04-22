# Adapting Router

`matcha` includes the `Adapter` interface to help define functionality that receives HTTP requests through methods other than direct HTTP/S, such as with serverless computing or through a message queue. External adapters aren't required to use this, but it may help.

## Implementing the Adapter Interface

`Adapter` has two type parameters, type `In` and type `Out`, that denote the input and output data being used to emulate an HTTP request. Out should be a pointer type. It has one function, `Adapt`, which takes in a value of type `In` and returns a few values:

- `http.ResponseWriter, *http.Request`: These values should be passed to a `Router` via `ServeHTTP`. The ResponseWriter should write response data to the `Out` value.
- `Out`: This value should be returned in place of normally writing an HTTP response.
- `error`: Adapters may optionally return errors if there could be missing data/incorrect formatting.
