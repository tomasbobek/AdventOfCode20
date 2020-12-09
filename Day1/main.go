package main

import (
    "bufio"
    "fmt"
    "math"
    "os"
    "sort"
    "strconv"
)

func main() {
    workingDir, err := os.Getwd()
    if err != nil {
        panic(fmt.Sprintf("Could not establish working directory, error: %v", err))
    }

    numbers, err := loadInput(workingDir + "\\input")
    if err != nil {
        panic(fmt.Sprintf("Could not load input file, error: %v\n", err))
    }

    // PART 1 ----->

    if first, second, found := getPairWhichTotalsTo(numbers, 2020); found {
        fmt.Printf("PART 1: Result is: %d\n", first*second)
    } else {
        fmt.Println("PART 1: Desired result could not be completed!")
    }

    // PART 2 ----->

    if first, second, third, found := getTripletWhichTotalsTo(numbers, 2020); found {
        fmt.Printf("PART 2: Result is: %d\n", first*second*third)
    } else {
        fmt.Println("PART 2: Desired result could not be completed!")
    }
}

// Loads file rows into slice of integers.
// Non-numeric rows are logged and skipped.
func loadInput(filePath string) ([]int, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var records []int
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        record, err := strconv.Atoi(scanner.Text())
        if err == nil {
            records = append(records, record)
        } else {
            fmt.Printf("skipping non-numeric input line %q\n", scanner.Text())
        }
    }
    return records, scanner.Err()
}

// Looks up two numbers in slice that add together given total number.
// Returns the numbers and flag indicating whether the lookup was successful.
func getPairWhichTotalsTo(numbers []int, total int) (int, int, bool) {
    // First sort the slice
    sort.Ints(numbers)

    // Combine smallest numbers with greatest until the total match or the result is smaller than expected total
    // (following results could only get smaller so it's meaningless to continue).
    for _, first := range numbers {
        for topIndex := len(numbers)-1; topIndex >= 0; topIndex-- {
            second := numbers[topIndex]

            if first + second == total {
                return first, second, true
            }
            if first + second < total {
                break
            }
        }
    }

    return 0, 0, false
}

// Looks up three numbers in slice that add together given total number.
// Returns the numbers and flag indicating whether the lookup was successful.
func getTripletWhichTotalsTo(numbers []int, total int) (int, int, int, bool) {
    // First sort the slice
    sort.Ints(numbers)

    // Combine smallest numbers with greatest and "average" until the total match expectation.
    // Algorithm is trying to combine smallest and greatest numbers together with third number which is picked from the
    // middle of slice and then is increasing or decreasing based on first comparison with desired total.
    for _, first := range numbers {
        for topIndex := len(numbers)-1; topIndex >= 0; topIndex-- {
            second := numbers[topIndex]

            middleIndex := int(math.Floor(float64(len(numbers)/2)))
            third := numbers[middleIndex]
            runningTotal := first + second + third
            offsetDirection := getOffsetDirection(runningTotal, total)
            for {
                third = numbers[middleIndex]
                runningTotal = first + second + third
                if runningTotal == total {
                    return first, second, third, true
                }

                middleIndex += offsetDirection

                // If we search for the third number in increasing direction and get the greater running total than
                // we're looking for, we don't need to continue. Same with decreasing direction and too low running total.
                if (offsetDirection == 1 && runningTotal > total) || (offsetDirection == -1 && runningTotal < total) ||
                    middleIndex < 0 || middleIndex >= len(numbers) {
                    break
                }
            }
        }
    }

    return 0, 0, 0, false
}

func getOffsetDirection(a, b int) int {
    if a >= b {
        return -1
    } else {
        return 1
    }
}
