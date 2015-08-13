# goplist [![Build Status](https://travis-ci.org/zach-klippenstein/goplist.svg)](https://travis-ci.org/zach-klippenstein/goplist) [![GoDoc](https://godoc.org/github.com/zach-klippenstein/goplist?status.svg)](https://godoc.org/github.com/zach-klippenstein/goplist)

A library for reading and writing OSX plist files.

```go
import "os"
import "github.com/zach-klippenstein/goplist"

// Keys sorted alphabetically.
goplist.Marshal(map[string]interface{}{
	"aString": "foobar",
	"anInt":   42,
}, os.Stdout)

// Keys sorted by the order they're added.
dict := goplist.NewEmptyDict()
dict.Set("aString", "foobar")
dict.Set("anInt", 42)
goplist.Marshal(dict, os.Stdout)
```

## TODO

* reading
* []byte (<data>)
* time.Time (<date>)