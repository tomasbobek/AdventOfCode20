package main

import (
    "bufio"
    "fmt"
    "os"
    "sort"
    "strconv"
)

func main() {

    workingDir, err := os.Getwd()
    if err != nil {
        panic(fmt.Sprintf("Could not establish working directory, error: %v", err))
    }

    ratings, err := loadRatings(workingDir + "\\input")
    if err != nil {
        panic(fmt.Sprintf("Could not load input file, error: %v\n", err))
    }
    sort.Ints(ratings)

    // PART 1 ----->

    // Add the rating for the end device which is 3 jolts higher than the most powerful adapter.
    ratings = append(ratings, ratings[len(ratings)-1]+3)
    if oneHops, threeHops, ok := countIncreasesInAdapterSequence(ratings); ok {
        fmt.Printf("PART 1: Number of 1-hops multiplied by 3-hops is: %d\n", oneHops*threeHops)
    }

    // PART 1 ----->

    // Prepend ratings with initial value of outlet (which is 0).
    ratings = append([]int{0}, ratings...)
    sequenceCount := countAdapterSequences(ratings)
    fmt.Printf("PART 2: Number of possible adapter sequences is: %d\n", sequenceCount)

}

// Counts number of increases by 1 and by 3 in sorted integer slice. Example: in slice of [1,3,4,7,10] the number of
// 1-hops is 2 (from 0 - initial value - to 1, from 3 to 4) and 3-hops is 2 as well (from 4 to 7, from 7 to 10).
// Function also returns "success" flag, which turns false once there's any hop larger than 3.
func countIncreasesInAdapterSequence(ratings []int) (int, int, bool) {
    lastRating, oneHops, threeHops := 0, 0, 0
    for _, rating := range ratings {
        switch {
        case rating == lastRating+1:
            oneHops++
            break
        case rating == lastRating+3:
            threeHops++
            break
        case rating > lastRating+3:
            return 0, 0, false
        }

        lastRating = rating
    }

    return oneHops, threeHops, true
}

// Go through sorted list of rating in backwards mode and for every rating it counts the number of rating sequences
// that can follow. Count of following subsequences is stored for each rating and the number grows toward reaching
// the first rating (value 0). Number of all sequences is returned.
func countAdapterSequences(ratings []int) int {
    options := make(map[int]int)
    options[ratings[len(ratings) - 1]] = 1

    for i := range ratings {
        index := len(ratings) - i - 1
        currentVal := ratings[index]
        limit := index + 4

        // Do not point to non-existent position
        if limit > len(ratings) {
            limit = len(ratings)
        }

        // Sum the numbers of following ratings sub-sequences which is the number of sub-sequences for current rating
        for _, nextRating := range ratings[index+1:limit] {
            if nextRating - currentVal <= 3 {
                if nextOpts, ok := options[nextRating]; ok {
                    options[currentVal] += nextOpts
                }
            }
        }
    }

    return options[0]
}

// Loads file rows into slice of integers.
// Non-numeric rows are logged and skipped.
func loadRatings(filePath string) ([]int, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var ratings []int
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        record, err := strconv.Atoi(scanner.Text())
        if err == nil {
            ratings = append(ratings, record)
        } else {
            fmt.Printf("skipping non-numeric input line %q\n", scanner.Text())
        }
    }
    return ratings, scanner.Err()
}
