# enumarr

Generates an array of your enums!

## Quick start
Add go:generate directive in your file:
```go
//go:generate go run github.com/powehihihi/enumarr@latest -type YourEnumType
type YourEnumTyp int
```
And run:
```
go generate ./...
```

## Example
From this:
```go
//git/status.go
package git

//go:generate go run github.com/powehihihi/enumarr@latest -type Status
type Status int8

const (
	Unmodified Status = iota
	Modified
	Added
	Deleted
	Renamed
)
```
you can generate this:
```go
//git/status_array.go
package git

var _StatusArray = [...]Status{
	Unmodified,
	Modified,
	Added,
	Deleted,
	Renamed,
}

func StatusAll() []Status {
	return _StatusArray[:]
}
```

