[![Build Status](https://travis-ci.org/jesusrmoreno/gordita.svg?branch=master)](https://travis-ci.org/jesusrmoreno/gordita)
# gordita
Simple Ledis Backed Queue for Golang

## Why?
I wanted a simple Queue without all of the bells and whistles that I could embed into my own applications. I was already using LedisDB so I thought why not?

## How to Use
```[Go]
  import (
    "github.com/jesusrmoreno/gordita"
  )
  
  func main() {
    // Initialize ledis db...
    q := gordita.New(db)
    q.Push([]byte("First Item")
    q.Push([]byte("Next Item")
    size := q.Size()
    fmt.Println(size) // 2
    val := q.Pop()
    fmt.Println(string(val)) // "First Item"
    val2 := q.Pop()
    fmt.Println(string(val2)) // "Next Item"
  }
  
```
