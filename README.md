<h1 align="left">RunGo</h1>
<br />

RunGo is a free cross-platform Go playground, that allows users to experiment,
prototype and get instant feedback. It provides support for running Go versions
from 1.16+, and is built on top of [Fyne](https://fyne.io), a cross-platform GUI
toolkit made with Go and inspired by Material Design.

[![Go Version](https://img.shields.io/github/go-mod/go-version/itsksrof/run-go)](https://github.com/itsksrof/run-go/blob/master/go.mod)
[![License](https://img.shields.io/github/license/itsksrof/run-go)](https://github.com/itsksrof/run-go/blob/master/LICENSE)

## Motivation
The idea behind RunGo basically emerged from me being lazy. I wanted to try out
prototypes, ideas and isolate blocks of code during the development of several
projects, each of them running different Go versions, and I found it to be a hassle.
I either had to do it through the browser, using [The Go Playground](https://go.dev/play)
or [Better Go Playground](https://goplay.tools) and if none of both supported the
version I was running, create a standalone project in my editor and go from there.

For sometime I told myself that it wasn't too bad, but the thought of making a
playground never left my mind. One day at the office, a good colleague showed me
RunJS which as the name suggests is a playground for running Javascript code, and
more importantly he showed me, how he was able to try the same things I tried,
but in a much more efficient way. After that the decision was clear as day and I
set out to learn how to make desktop applications, because, I had and (still have)
no idea on how I could do it :)

## Building RunGo
To build RunGo you will firstly need to install Fyne, you can do so by heading to
their [Getting Started](https://developer.fyne.io/started/) page, from there follow
the steps laid out in the guide, and then comeback.

After installing Fyne, clone the repository with the following command:
```sh
git clone https://github.com/itsksrof/run-go
```

And inside the directory run `make build`. This will automatically generate
a binary of the program for your operating system and architecture.

Beware that at the moment the playground only supports the following operating
system and architecture combinations:

| Operating System | Processor Architecture |
| ---------------- | ---------------------- |
| Darwin | ARM64, AMD64 |
| Linux | AMD64 |
| Windows | AMD64 |

More processor architectures will be supported in the near future.

## Next Steps
Currently RunGo is in a pre-alpha state, the core functionallity is working, as you
are able to run Go code from 1.16 and beyond, create new tabs save snippets and open
snippets. But there are much more things that I wish to add and polish, below is a 
list of the features and improvements that I want to implement in RunGo before moving 
to the alpha stage.

- [ ] Proper code editor with line numbers, indentation and syntax highlighting
- [ ] Autocomplete engine for the code editor
- [ ] Minor improvements
    - [ ] Add caching to the various requests performed in the application
    - [ ] Automatically change the Go version when a snippet is opened and has a different Go version
    - [ ] Automatically create a new tab when opening a snippet in a tab that already has content

## Contributing
All contributions are extremely appreciated, if you find an issue that is interesting
to you, do not hesitate and say something so I know that you are hacking on that. Also
do not be afraid to ask for help or feel like you are missing some information, I will
be delighted to help :)


### Git commit message guidelines
The most important part is that each commit messages should have a title/subject in imperative
mood starting with a capital letter and no trailing period: `generator: Fix articles list being sorted in ascending order`
**NOT** `list sorted right.`. If you still unsure about how to write a good commit message 
this [blog article](https://cbea.ms/git-commit/) is a good resource for learning that.

Most title/subjects should have a lower-cased prefix with a colon and one whitespace. The prefix can be:
- The name of the package where (most of) the changes are made (e.g. `parse: Add RawTitle to metadata struct`)
- If the commit touches several packages with a common functional topic, use that as a prefix (e.g. `errors: Resolve correct line numbers`)
- If the commit touches several packages without a common functional topic, prefix with `all:` (e.g. `all: Reformat go code`)
- If this is a documentation update, prefix with (e.g. `docs:`)
- If nothing of the above applies, just leave the prefix out

Also, if your commit references some or more Github issues, always end your commit message body with *See #1234* or *Fixes #1234*.
Replace *1234* with the Github issue ID.

An example:
```text
generator: Fix articles list being sorted in ascending order

Added a function that returns a sorted "map" by passing it an unsorted one, 
it creates an "array" with the keys of the first "map", sorts the given "array"
using sort.StableSort and iterates over those keys to create new entries in the
new "map" and sets its value by accessing the first "map".

Fixes #1234
```

## Credits
RunGo makes use of a variety of open-source projects including:
- [github.com/fyne-io/fyne](github.com/fyne-io/fyne)
- [github.com/PuerkitoBio/goquery](github.com/PuerkitoBio/goquery)
- [github.com/golang/mod](github.com/golang/mod)
