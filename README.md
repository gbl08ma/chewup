# chewup
The world's most barebones static site generator, sloppily developed in record time.

Written in Go, it's mainly a shortcut into using Go's template engine to produce HTML for a static website.

## Installation

`go get -u github.com/gbl08ma/chewup`

(Feel free to fork to add your favorite dependency management solution, etc. Don't complain: did you miss the "record time" part?)

## Usage

`chewup -in input -out output`

If the input folder is not specified, it will be the working directory.
If the output folder is not specified (but the input folder is), it will be the working directory.
If neither is specified, the input folder will be the working directory, and the output folder will be `generated` inside the input folder.

If you are unsure, by using the flag `-test` Chewup will execute and show the input and output directories it will use, without writing anything.

Chewup expects to find files with extension `.html` and `.template` inside the input directory.
It will look for those files recursively.
All matching files will be evaluated by Go's `html/template` ([Godoc](https://golang.org/pkg/html/template/)) and they will be [associated](https://golang.org/pkg/text/template/#hdr-Associated_templates) with each other, which means you can invoke other files from each template (which is the whole point of this program, really).

Each `.html` file will be executed and have a corresponding file in the output directory, with the same relative path and name.
`.template` files are reserved for inclusion by other files.