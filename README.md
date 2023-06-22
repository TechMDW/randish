# Randish

`Randish` is a Go package that provides pseudo-random number generators which are unique and thread-safe.

Do not use this for cryptography! If you need a package for that purpose, please use the `crypto/rand` package instead.

## Features

- `Rand()`: Generates a new pseudo-random number generator that is initialized with a unique seed value.
- `RandS()`: Singleton and thread-safe variant of the `Rand()` function. Returns a globally unique instance of a pseudo-random number generator.
- `RandSA()`: Initializes an array of pseudo-random number generators and returns one from the array using pseudo-random selection.
- `Seed()`: Function that generates a seed for pseudo-random number generation using various system and context-specific elements.

## Getting Started

### Prerequisites

- Go 1.7 or later

### Installation

1. Download `randish` using the go get command:
   ```
   go get github.com/TechMDW/randish
   ```
2. Import `randish` in your code:
   ```go
   import "github.com/TechMDW/randish"
   ```
3. Use the functions as needed.

## Usage

Here is an example of how to use `randish` in your Go code:

```go
package main

import (
    "github.com/TechMDW/randish"
)

func main() {
    // Generate a new random number generator
    //
    // Returns the *rand.Rand from "math/rand" package
    r := randish.Rand()

    // Generate a random number
    num := r.Int()

    println(num)
}
```

## TODO

- [x] Write README
- [ ] Version control

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contact

Maintained by @TechMDW - contact@techmdw.com - Feel free to contact us if you have any questions.
