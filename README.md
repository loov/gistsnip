# gistsnip

Create multiple gists from code based on annotations.

As an example you may have `main.go`:

```go
func main() {
    //gistsnip:start:loop
    for i := 0; i < 100; i++ {
        fmt.Println(i)
    }
    //gistsnip:end:loop

    //gistsnip:start:example
    fmt.Println("hello world")
    //gistsnip:end:example
}
```

`gistsnip` will create two gists:

```go
for i := 0; i < 100; i++ {
    fmt.Println(i)
}
```

and

```go
fmt.Println("hello world")
```

For multiple see https://github.com/egonelbre/db-demo/

## Usage

1. Create a Github Personal Access token at https://github.com/settings/tokens.
2. Either set it as `GISTSNIP_TOKEN` environment variable or pass it in as a command-line argument.
3. Run `gistsnip` in the target folder.

To update the gists, simply rerun `gistsnip`. `gistsnip` will cache the last gist state and update things as needed.