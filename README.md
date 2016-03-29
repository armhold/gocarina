# go-carina
Golang version of OCARINA


![logo](https://github.com/armhold/gocarina/blob/master/gocarina-logo.png "gocarina Logo")

This is a Go port of the [Ruby project](https://github.com/armhold/ocarina) I did a few years back.

Gocarina builds a feed-forward neural network to do simple Optical Character Recognition. It's trained on [Letterpress](http://www.atebits.com/letterpress) game boards, which makes it a handy way to automate a Letterpress cheat.

## Usage

`$ go run cmd/train/train.go -max 1000 -save ocr.save`

You now have a trained neural network in `ocr.save`.

You can ask it decipher Letterpress boards like this:

`$ go run cmd/recognize/recognize.go  -network ocr.save  -board board-images/board5.png`
```
  C Y M T I
  P Z Y L Y
  D W O H S
  D W H A S
  O Z X G K
```


## How it works

