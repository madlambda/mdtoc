<!-- mdtocstart -->
# Table of Contents

- [mdtoc](#mdtoc)
    - [Why ?](#why-)
    - [Install](#install)
    - [Usage](#usage)
<!-- mdtocend -->

# mdtoc

A Table of contents generator for markdown.

## Why ?

Feeling like having good TOCs on my markdowns and
wanted to hack something simple for it.

It is not a comprehensive markdown parser, it does
not have a lot of cool features, it probably do not
even covers all markdown syntax use cases.

Right now it only works for [atx](https://daringfireball.net/projects/markdown/syntax#header)
syntax headers.

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

The result will be a markdown with the TOC on its beginning.
The TOC is generated based on the parsed headers.

You can also pass a file:

```sh
mdtoc somemarkdownfile.md > newfile.md
```

Or make the change in place:

```sh
mdtoc -w somemarkdownfile.md
```
