acme playground
===============

acme playground enables 'playground' style programming in the [Acme Editor](https://en.wikipedia.org/wiki/Acme_(text_editor)). It works by passing the contents of an Acme window into a program via stdin and passing the output of that program into a different Acme window every time a keyboard event is processed.

## Install


```bash
go get github.com/sewh/acme_playground/cmd/playground
go install github.com/sewh/acme_playground/cmd/playground

# Ensure $(go env GOPATH) is in your path
```


## Usage

From an Acme window:

```bash
playground [program]
```

For example, in a Python file:

```bash
playground python3
```