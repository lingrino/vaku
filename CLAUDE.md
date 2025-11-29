# Claude Code Guidelines for Vaku

## Project Overview

Vaku is a CLI and Go API for managing HashiCorp Vault secrets. It provides operations for reading, writing, copying, moving, and deleting secrets at both path and folder levels, with support for KV v1 and v2 secret engines.

## Project Structure

```
api/           # Core API library
cmd/           # CLI implementation using Cobra
docs/cli/      # Auto-generated CLI documentation
www/           # Website
.github/       # GitHub Actions workflows and configuration
```

## Development Commands

```bash
# Run all tests
go test ./...

# Run linter
golangci-lint run ./...

# Regenerate CLI docs (REQUIRED after adding/changing CLI flags)
go run main.go docs docs/cli

# Build
go build ./...
```

## CI Requirements

Before pushing, ensure:
1. `go test ./...` passes
2. `golangci-lint run ./...` passes (max line length: 120 chars)
3. `go run main.go docs docs/cli` has been run if CLI flags changed
4. `gofmt -w .` has been applied

## Code Patterns

### Adding a New API Method

1. Create `api/<operation>.go` with:
   ```go
   var (
       ErrOperationName = errors.New("operation name")
   )

   func (c *Client) OperationName(...) error {
       // implementation
   }
   ```

2. Add method to `ClientInterface` in `api/client.go`

3. Create `api/<operation>_test.go` following existing test patterns

4. Add mock implementation in `cmd/main_test.go`:
   ```go
   func (c *testVakuClient) OperationName(...) error {
       return nil
   }
   ```

### Error Handling

Use the `newWrapErr` function for consistent error wrapping:
```go
return newWrapErr(path, ErrOperationName, err)
```

Error chains in tests must match exactly:
```go
compareErrors(t, err, []error{ErrOuterOp, ErrInnerOp, ErrVaultRead})
```

### Adding CLI Commands/Flags

1. Define constants for new flags:
   ```go
   const (
       flagXxxName    = "xxx"
       flagXxxUse     = "description of flag"
       flagXxxDefault = false
   )
   ```

2. Add flag in command constructor:
   ```go
   cmd.Flags().Bool(flagXxxName, flagXxxDefault, flagXxxUse)
   ```

3. Handle flag in run function with error checking:
   ```go
   xxx, err := cmd.Flags().GetBool(flagXxxName)
   if err != nil {
       return err
   }
   ```

4. **IMPORTANT**: Regenerate docs after adding flags:
   ```bash
   go run main.go docs docs/cli
   ```

### Folder Operations with Workers

For concurrent folder operations, use the worker pattern:
```go
eg, ctx := errgroup.WithContext(ctx)

pathC, errC := c.FolderListChan(ctx, src)
eg.Go(func() error {
    return <-errC
})

for i := 0; i < c.workers; i++ {
    eg.Go(func() error {
        return c.workerFunc(...)
    })
}

return eg.Wait()
```

### KV v2 Only Operations

For operations that only work on KV v2, validate upfront:
```go
if err := c.validateKV2Mount(src, c); err != nil {
    return newWrapErr(src, ErrOperationName, err)
}
if err := c.validateKV2Mount(dst, c.dc); err != nil {
    return newWrapErr(dst, ErrOperationName, err)
}
```

## Testing Patterns

### API Tests

```go
func TestOperationName(t *testing.T) {
    t.Parallel()

    tests := []struct {
        giveSrc    string
        giveDst    string
        wantErr    []error
        wantNilDst bool
    }{
        // test cases
    }

    for _, tt := range tests {
        t.Run(testName(tt.giveSrc, tt.giveDst), func(t *testing.T) {
            t.Parallel()
            for _, prefixPair := range seededPrefixProduct(t) {
                // Skip kv1 for v2-only operations
                if strings.HasPrefix(prefixPair[0], "kv1") {
                    continue
                }

                t.Run(testName(prefixPair[0], prefixPair[1]), func(t *testing.T) {
                    t.Parallel()
                    // test implementation
                })
            }
        })
    }
}
```

### CLI Tests

```go
func TestCommandName(t *testing.T) {
    t.Parallel()

    tests := []struct {
        name     string
        giveArgs []string
        wantOut  string
        wantErr  string
    }{
        {
            name:     "basic",
            giveArgs: []string{"arg1", "arg2"},
            wantOut:  "",
            wantErr:  "",
        },
        {
            name:     "with flag",
            giveArgs: []string{"--flag", "arg1", "arg2"},
            wantOut:  "",
            wantErr:  "",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            args := append([]string{"path", "command"}, tt.giveArgs...)
            cli, outW, errW := newTestCLIWithAPI(t, args)
            // assertions
        })
    }
}
```

## Common Pitfalls

1. **Forgetting to regenerate docs** - CI will fail if CLI flags are added without running `go run main.go docs docs/cli`

2. **Line length** - Keep lines under 120 characters; break long error slices across multiple lines

3. **Error chain mismatch** - Test error chains must exactly match the wrapping order

4. **Missing interface update** - New public API methods must be added to `ClientInterface`

5. **Missing mock** - New API methods need mock implementations in `cmd/main_test.go`
