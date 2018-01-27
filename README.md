# `runes`: A gentle introduction to TDD in Go

In this repo you will find the step-by-step development of the `runes` example: a command-line utility to find Unicode characters by name.

Each step in the development is documented explaining the Go language features used in the code.

Only a very basic knowdledge of Go is required to follow.

## Our goal

By the end of this tutorial, we'll have a command-line utility which works like this:

```
$ runes face eyes
U+1F601	😁	GRINNING FACE WITH SMILING EYES
U+1F604	😄	SMILING FACE WITH OPEN MOUTH AND SMILING EYES
U+1F606	😆	SMILING FACE WITH OPEN MOUTH AND TIGHTLY-CLOSED EYES
U+1F60A	😊	SMILING FACE WITH SMILING EYES
U+1F60D	😍	SMILING FACE WITH HEART-SHAPED EYES
U+1F619	😙	KISSING FACE WITH SMILING EYES
U+1F61A	😚	KISSING FACE WITH CLOSED EYES
U+1F61D	😝	FACE WITH STUCK-OUT TONGUE AND TIGHTLY-CLOSED EYES
U+1F638	😸	GRINNING CAT FACE WITH SMILING EYES
U+1F63B	😻	SMILING CAT FACE WITH HEART-SHAPED EYES
U+1F63D	😽	KISSING CAT FACE WITH CLOSED EYES
U+1F644	🙄	FACE WITH ROLLING EYES
```

You give `runes` one or more words as arguments, and it displays a list of Unicode characters whose names contains all the words you provided.

Learn more in the [project page (in Portuguese for now)](https://ThoughtWorksInc.github.io/sinais/).


## Credits

This tutorial is based in the `charfinder` example from chapter 18 of [Fluent Python](https://www.amazon.com/_/dp/1491946008), by Luciano Ramalho. The Go version named `runefinder`, was started in the [Garoa Gophers](https://garoa.net.br/wiki/Garoa_Gophers), study group by Afonso Coutinho (@afonso), Alexandre Souza (@alexandre), Andrews Medina (@andrewsmedina), João "JC" Martins (@jcmartins), Luciano Ramalho (@ramalho), Marcio Ribeiro (@mmr), and Michael Howard.
