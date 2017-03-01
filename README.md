# mdtoc

A Table of contents generator for markdown.

## Why ?

Feeling like having good TOCs on my markdowns and
wanted to hack something simple for it.

It is not a comprehensive markdown parser, it does
not have a lot of cool features, it probably do not
even covers all markdown syntax use cases.

## Install

If you have Go, go-get it:

```sh
# Make sure GOPATH/bin is in your PATH
go get github.com/katcipis/mdtoc/cmd/mdtoc
```
If not, install Go :-).

## Usage

Input is read from stdin, results on stdout, Just run:

```sh
cat somemarkdownfile.md | mdtoc > newfile.md
```

That is it.
