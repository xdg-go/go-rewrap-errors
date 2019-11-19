# go-rewrap-errors

[![GoDoc](https://godoc.org/github.com/xdg-go/go-rewrap-errors?status.svg)](https://godoc.org/github.com/xdg-go/go-rewrap-errors) [![Build Status](https://travis-ci.org/xdg-go/go-rewrap-errors.svg?branch=master)](https://travis-ci.org/xdg-go/go-rewrap-errors) [![codecov](https://codecov.io/gh/xdg-go/go-rewrap-errors/branch/master/graph/badge.svg)](https://codecov.io/gh/xdg-go/go-rewrap-errors) [![Go Report Card](https://goreportcard.com/badge/github.com/xdg-go/go-rewrap-errors)](https://goreportcard.com/report/github.com/xdg-go/go-rewrap-errors) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Rewrite Go source files to replace pkg/errors with Go 1.13 error wrapping.

[Still under development.]

This program reads a Go source file and rewraps errors:

```
errors.Wrap(err, "text")        -> fmt.Errorf("text: %w", err)
errors.Wrapf(err, "text %s", s) -> fmt.Errorf("text %s: %w", s, err)
```

If the string argument to `Wrap` or `Wrapf` is `fmt.Sprintf`, it will be
unwrapped:

```
errors.Wrap(err, fmt.Sprintf("text %s", s)) -> fmt.Errorf("text %s: %w", s, err)
```

Non-literal error/format strings will be concatenated with `: %w`:

```
const errFmt = "text %s"
errors.Wrapf(err, errFmt, s) -> fmt.Errorf(errFmt+": %w", s, err)
```

Output defaults to stdout or the original file can be overwritten with the
`-w` option.

It does not modify package imports: run the resulting source through
`goimports` for that.

Currently, only `Wrap` and `Wrapf` are supported.

## Installation

```
go get github.com/xdg-go/go-rewrap-errors
```

### Usage

```
go-rewrap-errors source.go > new-source.go
```

# Copyright and License

Copyright 2019 by David A. Golden. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License").
You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
