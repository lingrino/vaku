package vaku

// MaxConcurrency is the maximum number of threads/workers to use when calling Folder-based
// functions that execute concurrently. The default value is 10, but a stable and well-tuned
// Vault server should be able to handle up to 100 without issues. Use with caution and tune
// specifically to your environment and storage backend.
var MaxConcurrency = 10
