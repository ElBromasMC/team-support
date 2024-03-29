package main

import (
	"alc/model/cart"
	"encoding/gob"
)

func init() {
	gob.Register([]cart.ItemRequest{})
}
