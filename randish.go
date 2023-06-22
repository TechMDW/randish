package randish

import (
	randC "crypto/rand"
	"encoding/binary"
	"hash/fnv"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const arrSize = 10

var blacklist = []string{"once.go"}

var randishRand *rand.Rand
var onceS sync.Once
var onceSA sync.Once

// Variables for RandSA
var randArr [arrSize]*rand.Rand
var lastIdx int32

// Rand{}
// RandS{ingleton}
// RandSA{rray}

// Rand returns a new pseudo-random number generator, *rand.Rand, which is
// initialized with a unique seed value obtained from the Seed function. The
// Seed function uses a mix of various system and context-specific elements
// to generate a unique seed value. This guarantees that the sequence of
// pseudo-random numbers generated by the returned *rand.Rand object is distinct.
func Rand() *rand.Rand {
	seed := Seed()
	return rand.New(rand.NewSource(seed))
}

// RandS is a singleton and thread-safe variant of the Rand function. It returns a
// globally unique instance of a *rand.Rand, a pseudo-random number generator. This
// generator is initialized with a unique seed value obtained from the Seed function,
// which utilizes various system and context-specific elements. The singleton nature
// of this function ensures that the same *rand.Rand instance is returned across
// multiple calls, across all threads. The thread-safety guarantees that the singleton
// instance is correctly initialized even in a multi-threaded context.
func RandS() *rand.Rand {
	onceS.Do(func() {
		seed := Seed()
		randishRand = rand.New(rand.NewSource(seed))
	})
	return randishRand
}

// RandSA initializes an array of *rand.Rand instances, if not already done,
// and returns a *rand.Rand chosen from the array using pseudo-random selection.
func RandSA() *rand.Rand {
	onceSA.Do(func() {
		for i := 0; i < arrSize; i++ {
			seed := Seed()
			randArr[i] = rand.New(rand.NewSource(seed))
		}

		lastIdx = 0
	})

	newIdx := randArr[lastIdx].Intn(arrSize)
	atomic.StoreInt32(&lastIdx, int32(newIdx))

	return randArr[newIdx]
}

// Don't use this function directly. Use Rand or RandS instead.
func RandTest() (randish *rand.Rand, seed int64) {
	seed = Seed()
	randish = rand.New(rand.NewSource(seed))
	return
}

// Seed is a function that generates a seed for pseudo-random number generation.
// It uses various unique and variable elements from the system as a basis for the seed,
// which is then used to generate unique random numbers across multiple calls.
//
// The seed is initially set to the current Unix timestamp in nanoseconds.
// A new random source is then created based on this timestamp, which is used for further calculations.
//
// The function gathers system-specific data, including the current hostname, process ID, user ID,
// group ID, working directory, and the number of logical CPUs on the system.
//
// It also gathers some context-specific data by calling the runtime.Caller function.
// This returns information about the function that is invoking Seed, including its file name, line number,
// and a counter that represents the current instruction address.
//
// The function then reads a random number from the crypto/rand package to further enhance the randomness of the seed.
//
// After that, the function iterates through each environment variable, converting each one to an integer hash.
// The generated hash is used in a random mathematical operation (addition, subtraction, or multiplication)
// to modify the seed value.
//
// Each piece of system-specific and context-specific data, as well as the randomly generated number from crypto/rand,
// is also used in the same way, to further modify the seed value.
//
// The resulting seed is a complex, unique value that can be used to seed a pseudo-random number generator
// (such as the one in the math/rand package) to ensure that the generated sequence of numbers
// is as random and unique as possible across different contexts and systems.
func Seed() int64 {
	t := time.Now()
	randT := rand.New(rand.NewSource(t.UnixNano()))
	hostname, _ := os.Hostname()
	pid := os.Getpid()
	uid := os.Getuid()
	gid := os.Getgid()
	ps := os.Getpagesize()
	wd, _ := os.Getwd()
	cpus := runtime.NumCPU()
	cnt, file, line, _ := getExternalCaller(&blacklist)

	// Get random number from crypto/rand.
	var n int32
	binary.Read(randC.Reader, binary.BigEndian, &n)

	var seed int64 = t.UnixNano()

	// Get random number from crypto/rand
	randomMathOperation(randT, &seed, int64(n))
	// Page size
	randomMathOperation(randT, &seed, int64(ps))
	// cpu count
	randomMathOperation(randT, &seed, int64(cpus))
	// line number
	randomMathOperation(randT, &seed, int64(line))
	// file name + counter
	randomMathOperation(randT, &seed, int64(stringToHash(file))+int64(cnt))
	// Get the current hostname
	randomMathOperation(randT, &seed, int64(stringToHash(hostname)))
	// Get the current process id
	randomMathOperation(randT, &seed, int64(pid))
	// Get the current user id
	randomMathOperation(randT, &seed, int64(uid))
	// Get the current group id
	randomMathOperation(randT, &seed, int64(gid))
	// Get the current working directory
	randomMathOperation(randT, &seed, int64(stringToHash(wd)))
	// Time to execute Seed
	randomMathOperation(randT, &seed, int64(time.Since(t).Nanoseconds()))

	return seed
}

// TODO: Maybe don't allow overflow?
func randomMathOperation(rnd *rand.Rand, seed *int64, val int64) {
	switch rnd.Intn(3) {
	case 0:
		*seed += val
	case 1:
		*seed -= val
	case 2:
		if *seed == 0 || val == 0 {
			return
		}
		*seed *= val
	}
}

func stringToHash(s string) uint32 {
	h := fnv.New32()
	h.Write([]byte(s))
	return h.Sum32()
}

// getExternalCaller returns the PC, file, line, and ok for the first caller
// outside of the current package and not in the given ignore list.
func getExternalCaller(blacklist *[]string) (uintptr, string, int, bool) {
	for skip := 1; ; skip++ {
		pc, file, line, ok := runtime.Caller(skip)
		if !ok {
			return 0, "", 0, false
		}

		// Check if the file belongs to the current package or the ignore list.
		// This is a simple check and may not always work...
		// TODO: Find a better way to check if the file belongs to the current package.
		if !strings.Contains(file, "randish") && !inBlacklist(file, blacklist) {
			return pc, file, line, ok
		}
	}
}

// inBlacklist checks if the given file is in the blacklist.
func inBlacklist(file string, blacklist *[]string) bool {
	for _, ignoreFile := range *blacklist {
		if strings.Contains(file, ignoreFile) {
			return true
		}
	}
	return false
}
