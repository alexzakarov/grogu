# Metadata Cache Demonstration

This demonstration program validates the metadata caching functionality within the Grogu library.

## Execution

From the project root directory:

```bash
go run examples/cache_demo/main.go
```

## Test Scenarios

The program executes the following validation tests:

### Test 1: Initial Call - Cache Miss

- First conversion of User struct
- Metadata generation and cache storage
- Expected duration: ~28µs

### Test 2: Subsequent Call - Cache Hit

- Second conversion of User struct
- Metadata retrieved from cache (no reflection overhead)
- Expected duration: ~700ns
- Performance improvement: ~39x faster

### Test 3: Different Type - Cache Miss

- First conversion of Product struct
- Metadata generation for new type

### Test 4: Type Cache Hit

- Second conversion of Product struct
- Performance improvement: ~7x faster

### Test 5: Pointer Validation

- Pointers share cache with value types
- Single metadata entry per type

### Test 6: High-Volume Performance

- 10,000 iterations
- Throughput: ~1.7M operations/second
- Average latency: ~563ns per operation

## Expected Results

The metadata cache demonstrates:

- Cache miss on first call per struct type
- Cache hit provides significant performance improvement
- Single metadata build per struct type
- Thread-safe operation via sync.Map

## Custom Struct Testing

Modify `main.go` to test custom struct types:

```go
type CustomStruct struct {
    ID   int64  `json:"id"`
    Name string `json:"name"`
}

data := CustomStruct{ID: 1, Name: "Test"}
fields, values, err := utils.Convert(data)
```

## Performance Characteristics

1. Initial call per type: full reflection cost
2. Subsequent calls: cache lookup only
3. Performance improvement: 10-40x range
4. Concurrent access: thread-safe
