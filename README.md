# GOCARINA - simple Optical Character Recognition in Go

![logo](https://github.com/armhold/gocarina/blob/master/gocarina-logo.png "gocarina Logo")

Gocarina uses a neural network to do simple Optical Character Recognition (OCR).
It's trained on [Letterpress(™)](http://www.atebits.com/letterpress) game boards.

This is a Go port of the [Ruby project](https://github.com/armhold/ocarina) I did a few years back.


## Usage

First, install the software:

`$ go install github.com/armhold/gocarina/...`

Next, we need to create and train a network. Be sure to first connect to the source directory
(`train` expects the game boards to appear in `board-images/`):

`$ cd $GOPATH/github.com/armhold/gocarina`
`$ train`

You now have a trained neural network in `ocr.save`.

Now you can ask it decipher game boards like this:

`$ recognize -board board-images/board3.png`
```
 L H F L M
 R V P U K
 V O E E X
 I N R I T
 V N S I Q
```

You can also ask it to give you a list of words that can be formed with the board:

`$ recognize -board board-images/board3.png -words`
```
 L H F L M
 R V P U K
 V O E E X
 I N R I T
 V N S I Q


overmultiplies
relinquishment
feuilletonism
fluorimetries
interinvolves
pluviometries
reptiliferous
[etc...]
```



## How it works

We start with three "known" game boards. We split them up into individual tiles, one per letter.
This covers the entire alphabet, and gives us our training set. We feed the tiles into the network one at a time,
repeatedly, until the network is trained.


## Representation & Encoding

### Input

The tiles are quantized to black & white, bounding boxed, and finally scaled down to a small rectangular bitmap.
These bits are then fed directly into the inputs of the network.


### Output

We use a bit string to represent a given letter. 8 bits allows us to represent 256 different characters, which is
more than sufficient to cover the 26 characters used in Letterpress (we could certainly get away with using only
5 bits, but I wanted to hold the door open for potentially doing more than just Letterpress).

For convenience, we use the natural ASCII/Unicode mapping where 'A' = 65, aka 01000001. So our network has 8
outputs, corresponding to the 8 bits of our letters.


### What's with the name?

Original project: **Oc**a**r**ina, i.e. OCR. Go + Ocarina => Gocarina.


###  Credits

The file `words-en.txt` is in the Public Domain, licensed under CC0 thanks to https://github.com/atebits/Words.

