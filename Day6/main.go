package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    workingDir, err := os.Getwd()
    if err != nil {
        panic(fmt.Sprintf("Could not establish working directory, error: %v", err))
    }

    answerGroups, err := loadAnswers(workingDir + "\\input")
    if err != nil {
        panic(fmt.Sprintf("Could not load input file, error: %v\n", err))
    }

    // PART 1 ----->

    answerSum := 0
    for _, group := range answerGroups {
        answerSum += len(group.uniqueAnswers)
    }
    fmt.Printf("PART 1: Found %d unique answers\n", answerSum)

    // PART 2 ----->

    commonAnswerSum := 0
    for _, group := range answerGroups {
        commonAnswerSum += group.getCommonAnswerCount()
    }
    fmt.Printf("PART 2: Found %d unique answers common to all respondents\n", commonAnswerSum)
}

type answerGroup struct {
    uniqueAnswers map[int32]bool
    personAnswers []map[int32]bool
}

// Calculates count of answers that are answered by all people in given group (meaning that the letter is present on
// each row of answer group).
func (ag answerGroup) getCommonAnswerCount() int {
    commonAnswerCount := 0

    baseLoop: for baseAnswer := range ag.personAnswers[0] {
        for _, answers := range ag.personAnswers {
            if _, ok := answers[baseAnswer]; !ok {
                continue baseLoop
            }
        }
        commonAnswerCount++
    }

    return commonAnswerCount
}

// Loads questionnaire (input) answers and serializes them into groups (separated by blank lines). Answers themselves
// are also categorized per person which is represented by a new line. Answer is identified any represented by single letter.
func loadAnswers(filePath string) ([]answerGroup, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var answerGroups []answerGroup

    scanner := bufio.NewScanner(file)
    var currentAnswerGroup answerGroup
    for scanner.Scan() {
        // Answers can span over multiple lines, but wholly empty line indicates new group.
        if scanner.Text() == "" {
            answerGroups = append(answerGroups, currentAnswerGroup)
            currentAnswerGroup = answerGroup{}
        } else {
            parseAnswers(scanner.Text(), &currentAnswerGroup)
        }
    }
    answerGroups = append(answerGroups, currentAnswerGroup)

    return answerGroups, scanner.Err()
}

// We need to parse answer data per person, but we can also get set of unique answers (for the whole data input) right
// away (which is useful for getting first output).
func parseAnswers(data string, answerGroup *answerGroup) {
    answers := make(map[int32]bool)

    for _, letter := range data {
        if answerGroup.uniqueAnswers == nil {
            answerGroup.uniqueAnswers = make(map[int32]bool)
        }

        answerGroup.uniqueAnswers[letter] = true
        answers[letter] = true
    }

    answerGroup.personAnswers = append(answerGroup.personAnswers, answers)
}
