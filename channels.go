package main

ch := make(chan T)

// value v is going into the channel ch
ch <- v

// v is receiving the values from the channel ch
v = <- ch