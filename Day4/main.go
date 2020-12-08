package main

import (
    "bufio"
    "fmt"
    "os"
    "reflect"
    "regexp"
    "strconv"
    "strings"
)

func main() {
    workingDir, err := os.Getwd()
    if err != nil {
        panic(fmt.Sprintf("Could not establish working directory, error: %v", err))
    }

    passports, err := loadPassports(workingDir + "\\input")
    if err != nil {
        panic(fmt.Sprintf("Could not load input file, error: %v\n", err))
    }

    // PART 1 ----->

    numberOfValid := 0
    for _, passport := range passports {
        if passport.isDocumentValid() {
            numberOfValid++
        }
    }
    fmt.Printf("PART 1: Found %d valid passports\n", numberOfValid)

    // PART 2 ----->

    numberOfValid = 0
    for _, passport := range passports {
        if passport.areDocumentRecordsValid() {
            numberOfValid++
        }
    }
    fmt.Printf("PART 2: Found %d valid passports\n", numberOfValid)
}

type passport struct {
    Byr string
    Iyr string
    Eyr string
    Hgt string
    Hcl string
    Ecl string
    Pid string
    Cid string
}

func (p passport) isDocumentValid() bool {
    return p.Byr != "" && p.Iyr != "" && p.Eyr != "" && p.Hgt != "" && p.Hcl != "" && p.Ecl != "" && p.Pid != ""
}

func (p passport) areDocumentRecordsValid() bool {
    return validateAsNumberInRange(p.Byr, 1920, 2002) &&
        validateAsNumberInRange(p.Iyr, 2010, 2020) &&
        validateAsNumberInRange(p.Eyr, 2020, 2030) &&
        (validateAsNumberWithUnitInRange(p.Hgt, "cm", 150, 193) || validateAsNumberWithUnitInRange(p.Hgt, "in", 59, 76)) &&
        validateAsHexColor(p.Hcl) &&
        validateAsColorCode(p.Ecl) &&
        validateAsNumberOfDigits(p.Pid, 9)
}

// Checks whether the given value is numerical and within the given range.
func validateAsNumberInRange(number string, min int, max int) bool {
    num, err := strconv.Atoi(number)
    if err == nil && num >= min && num <= max {
        return true
    }
    return false
}

// Checks whether the given value has suffix of given unit and whether the number before this suffix belongs to the
// given range. Example inputs are: "180cm", "68in", "10m".
func validateAsNumberWithUnitInRange(length string, unit string, min int, max int) bool {
    if strings.HasSuffix(length, unit) {
        number, err := strconv.Atoi(length[:len(length)-len(unit)])
        if err == nil && number >= min && number <= max {
            return true
        }
    }

    return false
}

// Checks whether the given value meets the criteria of standard hexadecimal color code (including leading "#").
func validateAsHexColor(color string) bool {
    hexColorRegex := regexp.MustCompile(`^#[0-9a-f]{6}$`)
    return hexColorRegex.MatchString(color)
}

// Checks whether the given value matches one of the predefined color codes.
func validateAsColorCode(color string) bool {
    validCodes := []string{"amb", "blu", "brn", "gry", "grn", "hzl", "oth"}
    for _, code := range validCodes {
        if color == code {
            return true
        }
    }

    return false
}

// Checks whether given value is numerical and contains exactly the given number of digits.
func validateAsNumberOfDigits(digits string, numberOfDigits int) bool {
    numberOfDigitsRegex := regexp.MustCompile(fmt.Sprintf(`^[0-9]{%d}$`, numberOfDigits))
    return numberOfDigitsRegex.MatchString(digits)
}

// Loads file rows collection of Passport structures.
// Invalid records within passport data feed are ignored.
func loadPassports(filePath string) ([]passport, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var passports []passport

    scanner := bufio.NewScanner(file)
    var currentPassport passport
    for scanner.Scan() {
        // Passport record can span over multiple lines, but wholly empty line indicates new passport record.
        if scanner.Text() == "" {
            passports = append(passports, currentPassport)
            currentPassport = passport{}
        } else {
            parsePassportData(scanner.Text(), &currentPassport)
        }
    }
    passports = append(passports, currentPassport)

    return passports, scanner.Err()
}

// All passport properties are treated as string values. This parser uses reflection to map the property values from
// data stream to the respective fields of Passport structure.
func parsePassportData(data string, passport *passport) {
    for _, datum := range strings.Split(data, " ") {
        tuple := strings.Split(datum, ":")
        if len(tuple) != 2 {
            continue
        }

        r := reflect.ValueOf(passport)
        s := r.Elem()
        f := s.FieldByName(strings.Title(tuple[0]))
        if f.CanSet() {
            if f.Kind() == reflect.String {
                f.SetString(tuple[1])
            }
        }
    }
}
