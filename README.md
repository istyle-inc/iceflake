# iceflake

[![CircleCI](https://circleci.com/gh/istyle-inc/iceflake/tree/master.svg?style=svg)](https://circleci.com/gh/istyle-inc/iceflake/tree/master) [![Coverage Status](https://coveralls.io/repos/github/istyle-inc/iceflake/badge.svg?branch=master&service=github)](https://coveralls.io/github/istyle-inc/iceflake?branch=master&service=github)

## What is This

iceflake is a Unique ID generator using 'snowflake' algorithm.
connect using unix domain socket(we might implement tcp connection mode optionally, also, someday),
data are transferred using [Protocol Buffer](https://developers.google.com/protocol-buffers/).

## Usage

### IceFlake Server
You can download from release page or `go get -u github.com/istyle-inc/iceflake`

```
$ iceflake -w 1 -s YOUR_SOCKET_FILE_PATH
```

"YOUR_SOCKET_FILE_PATH" here, is to be an absolute path.

### IceFlake Client
Also you can use through this package.
So, execute

```
go get -u github.com/istyle-inc/iceflake
```

and here's a client code example
```
client := iceflake.NewClient("unix", "YOUR_SOCKET_FILE_PATH")
flake, err := client.Get()
if err != nil {
    logger.Error("Error: Failed connect socket or get data: ", err)
}
fmt.Println(flake.GetId())
```

You can access "Id" on result struct, but more safety to access through GetId() func, which return zero value when receiver is nil.


## License

[MIT](https://github.com/istyle-inc/iceflake/blob/master/LICENSE)

## Author

[istyle inc.](http://www.istyle.co.jp/)