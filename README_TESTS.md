# Unit Tests for Go Hello World Server

This project includes comprehensive unit tests for the main HTTP server functionality.

## Test Files

- `main_test.go` - Contains all unit tests for the application

## Running Tests

### Using Okteto CLI

1. **Run all unit tests:**
   ```bash
   okteto test unit
   ```

2. **Run tests with coverage:**
   ```bash
   okteto test coverage
   ```
   This will generate:
   - `coverage.out` - Coverage data file
   - `coverage.html` - HTML coverage report

### Using Go directly (if Go is installed locally)

1. **Run all tests:**
   ```bash
   go test -v ./...
   ```

2. **Run tests with coverage:**
   ```bash
   go test -v -coverprofile=coverage.out ./...
   go tool cover -html=coverage.out -o coverage.html
   ```

3. **Run benchmarks:**
   ```bash
   go test -bench=. -benchmem
   ```

## Test Coverage

The tests cover:

1. **HTTP Handler Tests (`TestHelloServer`):**
   - GET, POST, PUT, DELETE, HEAD requests
   - Requests with query parameters
   - Requests to different paths
   - All requests should return "Hello world!"

2. **Additional Handler Tests:**
   - `TestHelloServerWithHeaders` - Tests with custom headers
   - `TestHelloServerWithBody` - Tests with request body

3. **Integration Tests:**
   - `TestMainFunction` - Tests the server setup and handler registration

4. **Performance Tests:**
   - `BenchmarkHelloServer` - Benchmarks the handler performance

## Current Coverage

The current test coverage is approximately 20%. This is because:
- The `main()` function contains `http.ListenAndServe()` which blocks and cannot be directly tested
- The `helloServer` handler function is fully covered

## Test Configuration

The test configuration is defined in `okteto.yml`:

```yaml
test:
  unit:
    image: okteto/golang:1
    caches:
      - /go
      - /root/.cache
    commands:
      - go test -v ./...
  coverage:
    image: okteto/golang:1
    artifacts:
      - coverage.out
      - coverage.html
    caches:
      - /go
      - /root/.cache
    commands:
      - go test -v -coverprofile=coverage.out ./...
      - go tool cover -html=coverage.out -o coverage.html
```

## Future Improvements

To improve test coverage, consider:
1. Refactoring the `main()` function to extract the server setup logic into a separate testable function
2. Adding integration tests that start the server on a test port
3. Adding tests for error scenarios (though the current simple handler doesn't have error cases)