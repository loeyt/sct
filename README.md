# ssh-config-template

A tool I use to manage my `~/.ssh/config`

## CAUTION

Use this at your own risk. This _WILL_ overwrite `~/.ssh/config` no questions
asked.

## Installation

Get [Go], and `go get luit.eu/ssh-config-template`.

## Use

`ssh-config-template` looks at all folders in `~/.ssh/config-template/`. In
those folders it expects to find a `template` file, and a `hosts.json` file.
The `template` file must be a valid [text/template]. The hosts.json must be a
valid JSON Object, and each value in the Object will be passed to the template
as `{{.Name}}` and `{{.Value}}`. Folders are processed in alphabetical order,
and for each folder the hosts in `hosts.json` will be processed in alphabetical
order.

More README coming soon, maybe.

[Go]: https://golang.org/
[text/template]: https://golang.org/pkg/text/template/

