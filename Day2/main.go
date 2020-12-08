package main

import (
    "bufio"
    "fmt"
    "os"
    "regexp"
    "strconv"
)

const passwordRecordRegex = `(\d*)-(\d*)\s([a-z]):\s(.*)`

func main() {
    workingDir, err := os.Getwd()
    if err != nil {
        panic(fmt.Sprintf("Could not establish working directory, error: %v", err))
    }

    passwordRecords, err := loadInput(workingDir + "\\input")
    if err != nil {
        panic(fmt.Sprintf("Could not load input file, error: %v\n", err))
    }

    // PART 1 ----->

    counter := 0
    for _, record := range passwordRecords {
        if record.isValidForFirstMethod() {
            counter++
        }
    }
    fmt.Printf("PART 1: Result is %d\n", counter)

    // PART 1 ----->

    counter = 0
    for _, record := range passwordRecords {
        if record.isValidForSecondMethod() {
            counter++
        }
    }
    fmt.Printf("PART 2: Result is %d\n", counter)
}

type passwordRecord struct {
    policy   passwordPolicy
    password string
}

// First validation method counts the number of letter occurrences and compares it with specified minimum and maximum
// (letter, minimum and maximum are defined in password policy).
func (pr passwordRecord) isValidForFirstMethod() bool {
    counter := 0
    for _, letter := range pr.password {
        if string(letter) == pr.policy.letter {
            counter++
        }
    }

    return counter >= pr.policy.min && counter <= pr.policy.max
}

// Second validation method looks fox exactly one letter occurrence of given letter on positions (indexes) given by
// min and max variable from password policy. Occurrences on both positions are not allowed.
// Important: Positions are indexed from 1, not 0.
func (pr passwordRecord) isValidForSecondMethod() bool {
    // Prevent index out of bounds error
    if len(pr.password) < pr.policy.max {
        return false
    }

    firstPositionMatch := string(pr.password[pr.policy.min-1]) == pr.policy.letter
    secondPositionMatch := string(pr.password[pr.policy.max-1]) == pr.policy.letter

    return (firstPositionMatch && !secondPositionMatch) || (!firstPositionMatch && secondPositionMatch)
}

type passwordPolicy struct {
    min    int
    max    int
    letter string
}

// Loads file rows into slice password records.
// Rows that do not match password record format (enforced by regex) are skipped.
func loadInput(filePath string) ([]passwordRecord, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var records []passwordRecord
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        recordRegex := regexp.MustCompile(passwordRecordRegex)
        recordParts := recordRegex.FindStringSubmatch(scanner.Text())
        if len(recordParts) != 5 {
            fmt.Printf("skipping invalid input line %q\n", scanner.Text())
            continue
        }

        if parsedRecord, ok := parsePasswordRecord(recordParts[1:]); ok {
            records = append(records, parsedRecord)
        }

    }
    return records, scanner.Err()
}

// Converts the slice of password record components into final structure.
// Besides the structure, it also reports flag for conversion success.
func parsePasswordRecord(parts []string) (passwordRecord, bool) {
    var record passwordRecord
    var policy passwordPolicy

    if min, err := strconv.Atoi(parts[0]); err == nil {
        policy.min = min
    } else {
        return passwordRecord{}, false
    }

    if max, err := strconv.Atoi(parts[1]); err == nil {
        policy.max = max
    } else {
        return passwordRecord{}, false
    }

    policy.letter = parts[2]
    record.password = parts[3]
    record.policy = policy

    return record, true
}
