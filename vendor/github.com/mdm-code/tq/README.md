<h1 align="center">
  <div >
    <img
      src="https://raw.githubusercontent.com/mdm-code/mdm-code.github.io/main/tq_logo.png"
      alt="logo"
      style="object-fit: contain"
      width="30%"
    />
  </div>
</h1>

<h4 align="center">Query TOML configuration files with the Tq terminal utility</h4>

<div align="center">
<p>
    <a href="https://github.com/mdm-code/tq/actions?query=workflow%3ACI">
        <img alt="Build status" src="https://github.com/mdm-code/tq/workflows/CI/badge.svg">
    </a>
    <a href="https://app.codecov.io/gh/mdm-code/tq">
        <img alt="Code coverage" src="https://codecov.io/gh/mdm-code/tq/branch/main/graphs/badge.svg?branch=main">
    </a>
    <a href="https://opensource.org/licenses/MIT" rel="nofollow">
        <img alt="MIT license" src="https://img.shields.io/github/license/mdm-code/tq">
    </a>
    <a href="https://goreportcard.com/report/github.com/mdm-code/tq">
        <img alt="Go report card" src="https://goreportcard.com/badge/github.com/mdm-code/tq">
    </a>
    <a href="https://pkg.go.dev/github.com/mdm-code/tq">
        <img alt="Go package docs" src="https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white">
    </a>
</p>
</div>

The `tq` program lets you query TOML configuration files with a sequence of
intuitive filters. It works as a regular Unix filter program reading input data
from the standard input and producing results to the standard output. Consult the
[package documentation](https://pkg.go.dev/github.com/mdm-code/tq) or check the
[Usage](#usage) section to see how you can use `tq`.


## Installation

Install the program and use `tq` on the command-line to filter TOML files on
the terminal.

```sh
go install github.com/mdm-code/tq/cmd/tq@latest
```

Here is how you can get the whole Go package downloaded to fiddle with, but it
does not expose any public interfaces in code.

```sh
go get github.com/mdm-code/tq
```


## Usage

Enter `tq -h` to get usage information and the list of options that can be used
with the command. Here is table with the supported filter expressions and some
examples to get you going on how to use `tq` in your workflow.


### Supported filters

| <a href="#supported-filters"><img width="1000" height="0"></a><p>Filter</p> | <a href="#supported-filters"><img width="1000" height="0"></a><p>Expression</p> |
| :-------------------------------------------------------------------------: | :-----------------------------------------------------------------------------: |
| <kbd><b>identity</b></kbd>                                                  | <kbd><b>.</b></kbd>                                                             |
| <kbd><b>key</b></kbd>                                                       | <kbd><b>["string"]</b></kbd>                                                    |
| <kbd><b>index</b></kbd>                                                     | <kbd><b>[0]</b></kbd>                                                           |
| <kbd><b>iterator</b></kbd>                                                  | <kbd><b>[]</b></kbd>                                                            |
| <kbd><b>span</b></kbd>                                                      | <kbd><b>[:]</b></kbd>                                                           |


### Retrieve IPs from a table of server tables

In the example below, the TOML input file is (1) queried with the key
`["servers"]`, then (2) the retrieved table is converted to an iterator of
objects with `[]`, and then (3) the IP address is recovered from each of the
objects with the key `["ip"]`.

```sh
tq -q '["servers"][]["ip"]' <<EOF
[servers]

[servers.prod]
ip = "10.0.0.1"
role = "backend"

[servers.staging]
ip = "10.0.0.2"
role = "backend"
EOF

Output:

'10.0.0.1'
'10.0.0.2'
```


### Retrieve selected ports from a list of databases

This example queries the TOML input for the for the all ports aside from the
first one assigned to the first database record on the list.

```sh
tq -q '.["databases"][0]["ports"][1:][]' <<EOF
databases = [ {enabled = true, ports = [ 5432, 5433, 5434 ]} ]
EOF

Output:

5433
5434
```


### Run inside of a container

If you don't feel like installing `tq` with `go install`, you can test `tq` out
running inside of a container with this command:

```sh
docker run -i ghcr.io/mdm-code/tq:latest tq -q "['dependencies']['ignore']" <<EOF
[dependencies]
anyhow = "1.0.75"
bstr = "1.7.0"
grep = { version = "0.3.1", path = "crates/grep" }
ignore = { version = "0.4.22", path = "crates/ignore" }
lexopt = "0.3.0"
log = "0.4.5"
serde_json = "1.0.23"
termcolor = "1.1.0"
textwrap = { version = "0.16.0", default-features = false }
EOF
```


## Development

Go through the [Makefile](Makefile) to get an idea of the formatting, testing
and linting that can be used locally for development purposes.


## License

Copyright (c) 2024 Michał Adamczyk.

This project is licensed under the [MIT license](https://opensource.org/licenses/MIT).
See [LICENSE](LICENSE) for more details.
