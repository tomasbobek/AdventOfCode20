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

    numbers, err := loadInput(workingDir + "\\input")
    if err != nil {
        panic(fmt.Sprintf("Could not load input file, error: %v\n", err))
    }

    // PART 1 ----->

    if firstError, ok := findFirstInvalidNumber(numbers, 25); ok {
        fmt.Printf("PART 1: First number that doesn't pass the check is: %d\n", firstError)

        // PART 2 ----->

        if sequence, ok := findContiguousListThatAddTo(numbers, firstError); ok {
            sort.Ints(sequence)
            fmt.Printf("PART 2: Sum of interval start and end is: %d\n", sequence[0]+sequence[len(sequence)-1])
        }
    }
}

// Invalid number is one that cannot be represented by sum of any of 2 numbers from X (preamble param) numbers that
// precede it. First number in the list that's invalid is returned. There's also flag whether such number was even found.
func findFirstInvalidNumber(numbers []int, preamble int) (int, bool) {
    for i := preamble; i < len(numbers); i++ {
        if !isSumOfTwoNumbersInSlice(numbers[i-preamble:i], numbers[i]) {
            return numbers[i], true
        }
    }
    return 0, false
}

// Checks whether given number (sum param) can be represented as sum of any 2 numbers (which are not identical) from
// given slice (numbers param).
func isSumOfTwoNumbersInSlice(numbers []int, sum int) bool {
    for _, num1 := range numbers {
        for _, num2 := range numbers {
            if num1 != num2 && num1+num2 == sum {
                return true
            }
        }
    }
    return false
}

// Looks for contiguous sub-slice in a slice (numbers param) where the sum of all numbers in this sub-slice equals to
// given number (sum param). There's also flag whether such sequence was even found.
func findContiguousListThatAddTo(numbers []int, sum int) ([]int, bool) {
    for i, startNum := range numbers {
        total := startNum
        for j, num := range numbers[i+1:] {
            total += num

            if total == sum {
                return numbers[i:i+j+1], true
            }
            if total > sum {
                break
            }
        }
    }
    return []int{}, false
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
