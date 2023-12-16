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
package git
//go:generate go run github.com/powehihihi/enumarr@latest -type GitStatus
type GitStatus int8

const (
	Unmodified GitStatus = iota
	Modified
	Added
	Deleted
	Renamed
)
```
you can generate this:
```go
package git

var _GitStatusArray = [...]GitStatus{
	Unmodified,
	Modified,
	Added,
	Deleted,
	Renamed,
}

func GitStatusAll() []GitStatus {
	return _GitStatusArray[:]
}
```

