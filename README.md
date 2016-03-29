# gocarina - Golang version of OCARINA

![logo](https://github.com/armhold/gocarina/blob/master/gocarina-logo.png "gocarina Logo")

This is a Go port of the [Ruby project](https://github.com/armhold/ocarina) I did a few years back.

Gocarina uses a neural network to do simple Optical Character Recognition (OCR).
It's trained on [Letterpress](http://www.atebits.com/letterpress) game boards.

## Usage

First, we need to train the network:

`$ go run cmd/train/train.go -max 1000 -save ocr.save`

You now have a trained neural network in `ocr.save`.

Now you can ask it decipher Letterpress boards like this:

`$ go run cmd/recognize/recognize.go  -network ocr.save  -board board-images/board5.png`
```
  C Y M T I
  P Z Y L Y
  D W O H S
  D W H A S
  O Z X G K
```


## How it works

We start with three "known" Letterpress boards. We split them up into individual tiles, one per letter.
This covers the entire alphabet, and gives us our training set. We feed the tiles into the network one at a time.


## Representation & Encoding

### Input

Tiles are fed into the network as a series of bits. Tiles are quantized to black & white, bounding boxed, and finally
scaled down to a small rectangular bitmap. These bits are then fed directly into the inputs of the network.


### Output

We use a bit string to represent a given letter. 8 bits allows us to represent 256 different characters, which is
more than sufficient to cover the 26 characters used in Letterpress (we could certainly get away with using only
5 bits, but I wanted to hold the door open for potentially doing more than just Letterpress).

For convenience, we use the natural ASCII/Unicode mapping where 'A' = 65, aka 01000001. So our network has 8
outputs, corresponding to the 8 bits of our letters.


