package main

import (
    "bufio"
    "fmt"
    "os"
    "regexp"
    "strconv"
)

const (
    BagColorRegex = `^([a-z]+\s[a-z]+)\sbags\scontain.*$`
    InnerBagsRegex = `((?P<count>\d)\s(?P<color>[a-z]+\s[a-z]+)\sbag[s]?)`
)

func main() {
    workingDir, err := os.Getwd()
    if err != nil {
        panic(fmt.Sprintf("Could not establish working directory, error: %v", err))
    }

    baggageRules, err := loadBaggageRules(workingDir + "\\input")
    if err != nil {
        panic(fmt.Sprintf("Could not load input file, error: %v\n", err))
    }

    // PART 1 ----->

    allowedBags := make(map[string]bool)
    findBagsThatCanContain(baggageRules, "shiny gold", allowedBags)
    fmt.Printf("PART 1: Shiny Gold bag can be contained in %d bag types\n", len(allowedBags))

    // PART 2 ----->

    numberOfBagsInside := 0
    findNumberOfBagsContainedIn(baggageRules, "shiny gold", 1, &numberOfBagsInside)
    fmt.Printf("PART 2: Shiny Gold bag contains %d other bags\n", numberOfBagsInside)
}

// Recursively searches through available baggage rules and find color of those bags that can (even indirectly) contain
// a bag of color given on input. Results are stored in referenced map passed as a 'result' parameter.
func findBagsThatCanContain(rules []baggageRule, color string, result map[string]bool) {
    for _, rule := range rules {
        if rule.canContainBag(color) {
            result[rule.color] = true
            findBagsThatCanContain(rules, rule.color, result)
        }
    }
}

// Recursively browses through baggage rules in order to find a total number of bags that are contained in a bag of color
// given on input. Result is stored in referenced variable passed as 'total' parameter. There is also 'multiplier' param
// which signals how many bags of given color are contained in given bag and therefore the number by which we need to
// multiply subsequent results (example: blue bag contains 3 red bags and red bag contains 2 green bags -> blue bag therefore
// contains total of 3 red bags and 6 green bags (2 per each red)).
func findNumberOfBagsContainedIn(rules []baggageRule, color string, multiplier int, total *int) {
    if rule, ok := findRuleByColor(rules, color); ok {
        for color, count := range rule.canContain {
            *total += multiplier * count
            findNumberOfBagsContainedIn(rules, color, multiplier*count, total)
        }
    }
}

func findRuleByColor(rules []baggageRule, color string) (baggageRule, bool) {
    for _, rule := range rules {
        if rule.color == color {
            return rule, true
        }
    }

    return baggageRule{}, false
}

type baggageRule struct {
    color      string
    canContain map[string]int
}

func (br baggageRule) canContainBag(color string) bool {
    _, ok := br.canContain[color]
    return ok
}

// Loads baggage rules from input file. Each file line represents individual rule.
func loadBaggageRules(filePath string) ([]baggageRule, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var baggageRules []baggageRule

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        if parsedRule, ok := parseBaggageRule(scanner.Text()); ok {
            baggageRules = append(baggageRules, parsedRule)
        }
    }

    return baggageRules, scanner.Err()
}

// Each baggage rule is given in predefined format which is parsed using regular expressions.
// Parser fist looks for base bag color and if it was found, it parses the rules about bags it contains (those are
// composed of bag color and count).
func parseBaggageRule(data string) (baggageRule, bool) {
    var rule baggageRule

    baseColorRegex := regexp.MustCompile(BagColorRegex)
    baseColorParts := baseColorRegex.FindStringSubmatch(data)
    if len(baseColorParts) < 2 {
        return baggageRule{}, false
    }
    rule.color = baseColorParts[1]

    rulesRegex := regexp.MustCompile(InnerBagsRegex)
    rule.canContain = make(map[string]int)
    rulesData := rulesRegex.FindAllStringSubmatch(data, -1)
    for _, ruleData := range rulesData {
        bagCount, err := strconv.Atoi(ruleData[2])
        if err == nil {
            rule.canContain[ruleData[3]] = bagCount
        }
    }

    return rule, true
}
