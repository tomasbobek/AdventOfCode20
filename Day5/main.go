package main

import (
    "bufio"
    "fmt"
    "os"
    "regexp"
    "strconv"
)

const (
    BoardingPassFormat = `^([F,B]{7})([L,R]{3})$`
)

func main() {
    workingDir, err := os.Getwd()
    if err != nil {
        panic(fmt.Sprintf("Could not establish working directory, error: %v", err))
    }

    passes, err := loadBoardingPasses(workingDir + "\\input")
    if err != nil {
        panic(fmt.Sprintf("Could not load input file, error: %v\n", err))
    }

    // SHARED PRE-CALCULATION FOR BOTH PARTS

    highestSeatId := 0
    allTakenSeatIds := make(map[int]bool)
    for _, pass := range passes {
        seatId := pass.getSeatId()
        allTakenSeatIds[seatId] = true
        if seatId > highestSeatId {
            highestSeatId = seatId
        }
    }

    // PART 1 ----->

    fmt.Printf("PART 1: Highest Seat ID is %d\n", highestSeatId)

    // PART 2 ----->

    mySeatId := 0
    for seatId := 0; seatId < highestSeatId; seatId++ {
        if _, ok := allTakenSeatIds[seatId]; ok {
            continue
        }

        _, prevSeatTaken := allTakenSeatIds[seatId - 1]
        _, nextSeatTaken := allTakenSeatIds[seatId + 1]
        if prevSeatTaken && nextSeatTaken {
            mySeatId = seatId
        }
    }
    fmt.Printf("PART 2: My Seat ID is %d\n", mySeatId)
}

type boardingPass struct {
    rowCode    string
    columnCode string
}

func (bp boardingPass) getSeatId() int {
    return bp.getRowPosition() * 8 + bp.getColumnPosition()
}

func (bp boardingPass) getRowPosition() int {
    rowPosition, err := turnBitStringToNumber(bp.rowCode, 'F', 'B')

    if err != nil {
        return 0
    } else {
        return rowPosition
    }
}

func (bp boardingPass) getColumnPosition() int {
    rowPosition, err := turnBitStringToNumber(bp.columnCode, 'L', 'R')

    if err != nil {
        fmt.Println(err)
        return 0
    } else {
        return rowPosition
    }
}

// Translates string of repeating pair of letters to a binary number which is in turn converted to decimal number.
// Letters for one and zero bits are given on input.
// Example: 'X' is one bit, 'Y' is zero bit, "XXYX" translates to "1101" which converts to 7.
func turnBitStringToNumber(value string, zeroBitLetter, oneBitLetter int32) (int, error) {
    binaryString := ""
    for _, letter := range value {
        switch letter {
        case zeroBitLetter:
            binaryString += "0"
        case oneBitLetter:
            binaryString += "1"
        }
    }

    rowIndex, err := strconv.ParseInt(binaryString, 2, 64)
    if err != nil {
        return 0, err
    }

    return int(rowIndex), nil
}

// Loads file rows into slice password records.
// Rows that do not match password record format (enforced by regex) are skipped.
func loadBoardingPasses(filePath string) ([]boardingPass, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var passes []boardingPass
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        recordRegex := regexp.MustCompile(BoardingPassFormat)
        if !recordRegex.MatchString(scanner.Text()) {
            fmt.Printf("skipping invalid input line %q\n", scanner.Text())
            continue
        }

        passParts := recordRegex.FindStringSubmatch(scanner.Text())
        passes = append(passes, boardingPass{
            rowCode:    passParts[1],
            columnCode: passParts[2],
        })
    }

    return passes, scanner.Err()
}
