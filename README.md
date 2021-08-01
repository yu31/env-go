# loader

loader if a golang library for managing configuration data from environment variables or K/V storage

## Features
* User-define struct tag name
* User-define prefix
* Set default value in tag label
* Struct nesting
* User-define Setter to deserialize values
* User-define Getter to get value by specified tag key

## Supported Struct Field Types

envconfig supports these struct field types:

  * string
  * int, int8, int16, int32, int64
  * uint, uint8, uint16, uint32, uint64
  * bool
  * float32, float64
  * slice of any supported type
  * array of any supported type
  * map (keys and values of any supported type)
  * [encoding.TextUnmarshaler](https://golang.org/pkg/encoding/#TextUnmarshaler)
  * [encoding.BinaryUnmarshaler](https://golang.org/pkg/encoding/#BinaryUnmarshaler)
  * [time.Duration](https://golang.org/pkg/time/#Duration)

Embedded structs using these fields are also supported.

## Installation

```bash
go get -u github.com/DataWorkbench/loader
```
Used in go modules
```bash
go get -insecure github.com/DataWorkbench/loader
```

## Usage

#### Load config from environment variables

Set some environment variables:

```Bash
export MYAPP_ADDRESS="127.0.0.1"
export MYAPP_PORT=8080
export MYAPP_TIMEOUT="30s"
export MYAPP_USERS="rob ken robert"
export MYAPP_COLORCODES="red:1 green:2 blue:3"
export MYAPP_EMBEDDED_NUMBER=1024
export MYAPP_NAME1="name1"
export MYAPP_NAME2="name1"
```

The field will be ignore which are not set tag or tag key is "-".

```Go
package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/DataWorkbench/loader"
)

type Embedded struct {
	Number int64 `loader:"NUMBER"`
}

type Config struct {
	Address    string         `loader:"ADDRESS"`
	Port       int            `loader:"PORT"`
	Timeout    time.Duration  `loader:"TIMEOUT"`
	Users      []string       `loader:"USERS"`
	Rate       float32        `loader:"ROTE"`
	ColorCodes map[string]int `loader:"COLORCODES"`
	Embedded   *Embedded      `loader:"Embedded"`
	Name1      string         `loader:"-"`
	Name2      string
}

func main() {
	var c Config
	l := loader.New(loader.WithPrefix("MYAPP"))
	if err := l.Load(&c); err != nil {
		fmt.Println(err)
	}

	b, err := json.MarshalIndent(&c, "", "\t")
	if err != nil {
		return
	}

	fmt.Println(string(b))

	/* output:
	{
		"Address": "127.0.0.1",
		"Port": 8080,
		"Timeout": 30000000000,
		"Users": [
			"rob",
			"ken",
			"robert"
		],
		"Rate": 0,
		"ColorCodes": {
			"blue": 3,
			"green": 2,
			"red": 1
		},
		"Embedded": {
			"Number": 1024
		}
	}
	*/
}
```

####  User-define the tag name
Set some environment variables:

```Bash
export TEST_NAME1="name1"
export TEST_NAME2="name2"
```

```go
package main

import (
	"fmt"

	"github.com/DataWorkbench/loader"
)

type CustomTag struct {
	Name1 string `env:"NAME1"`
	Name2 string `env:"NAME2"`
}

func main() {
	var c CustomTag
	l := loader.New(loader.WithPrefix("TEST"), loader.WithTagName("env"))
	if err := l.Load(&c); err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v\n", c)

	/* output:
	{Name1:name1 Name2:name2}
	*/
}
```

#### Default label and User-define Setter 

```bash
export EMB="127.0.0.1:9090"
export LISTS="Joe/man Lisa/woman"
```

```go
package main

import (
	"fmt"
	"strings"

	"github.com/DataWorkbench/loader"
)

type List struct{
	Name string
	Sex  string
}

// String for test output.
func (l *List) String() string {
	return fmt.Sprintf("{Name: %s, Sex: %s}", l.Name, l.Sex)
}

func (l *List) Set(value string) error {
	x := strings.Split(value, "/")
	l.Name = x[0]
	l.Sex = x[1]
	return nil
}

type Embedded struct {
	IP  string
	Port string
}

func (e *Embedded) Set(value string) error {
	x := strings.Split(value, ":")
	e.IP = x[0]
	e.Port = x[1]
	return nil
}

type Config struct {
	Retry int `env:"RETRY,default=10"`
	Message string `env:"MESSAGE,default=Hello world"`
	Embedded Embedded `env:"EMB"`
	Lists    []*List  `env:"LISTS"`
}

func main() {
	//os.Setenv("EMB", "127.0.0.1:9090")
	//os.Setenv("LISTS", "Joe/man Lisa/woman")

	var c Config
	l := loader.New(loader.WithTagName("env"))
	if err := l.Load(&c); err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v\n", c)

	/* output:
	{Retry:10 Message:Hello world Embedded:{IP:127.0.0.1 Port:9090} Lists:[{Name: Joe, Sex: man} {Name: Lisa, Sex: woman}]}
	*/
}
```


