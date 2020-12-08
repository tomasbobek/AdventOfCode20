package main

import (
    "bufio"
    "fmt"
    "os"
    "regexp"
    "strconv"
)

const InstructionRegex = `^(acc|jmp|nop)\s(\+|\-)(\d+)$`

func main() {
    workingDir, err := os.Getwd()
    if err != nil {
        panic(fmt.Sprintf("Could not establish working directory, error: %v", err))
    }

    var program program
    program.reset()
    err = program.loadInstructions(workingDir + "\\input")
    if err != nil {
        panic(fmt.Sprintf("Could not load input file, error: %v\n", err))
    }

    // PART 1 ----->

    program.run()
    fmt.Printf("PART 1: Accumulator value before going into loop is: %d\n", program.accumulator)

    // PART 2 ----->

    // Switch "jmp" and "nop" operations until the program finishes with "success" exit code.
    instLoop: for _, inst := range program.instructions {
        program.reset()
        switch inst.code {
        case Jump:
            inst.code = NoOperation
            if program.run() == Success {
                break instLoop
            } else {
                inst.code = Jump
            }
        case NoOperation:
            inst.code = Jump
            if program.run() == Success {
                break instLoop
            } else {
                inst.code = NoOperation
            }
        }
    }
    fmt.Printf("PART 2: Accumulator value of fixed program is: %d\n", program.accumulator)
}

type instructionCode string

const (
    Accumulate instructionCode = "acc"
    Jump = "jmp"
    NoOperation = "nop"
)

type exitCode int

const (
    Success exitCode = iota
    InfiniteLoop
)

type instruction struct {
    code      instructionCode
    argument  int
    execCount int
}

type program struct {
    accumulator  int
    execIndex    int
    instructions []*instruction
}

// Sets the whole program to initial (zeroed) state.
func (p *program) reset() {
    p.execIndex = 0
    p.accumulator = 0
    for _, inst := range p.instructions {
        inst.execCount = 0
    }
}

// Loads every individual instruction from input file into a program structure.
// Invalid input data (which don't pass through regular expression) are ignored.
func (p *program) loadInstructions(filePath string) error {
    file, err := os.Open(filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        instructionRegex := regexp.MustCompile(InstructionRegex)
        instructionData := instructionRegex.FindStringSubmatch(scanner.Text())
        if len(instructionData) > 0 {
            if instruction, ok := parseInstruction(instructionData); ok {
                p.instructions = append(p.instructions, &instruction)
            }
        }
    }

    return nil
}

func (p *program) run() exitCode {
    for {
        // Terminate with "success" if the execution index points after last instruction.
        if p.execIndex >= len(p.instructions) {
            return Success
        }

        currentInst := p.instructions[p.execIndex]
        currentInst.execCount++

        // Terminate with "infinite loop" is any instruction is to be executed twice.
        if currentInst.execCount > 1 {
            return InfiniteLoop
        }

        switch currentInst.code {
        case Accumulate:
            p.accumulator += currentInst.argument
            p.execIndex++
        case Jump:
            p.execIndex += currentInst.argument
        case NoOperation:
            p.execIndex++
        }
    }
}

// Fills instruction structure with data-parts provided on input.
// Flag marking successful parsing is returned as well.
func parseInstruction(data []string) (instruction, bool) {
    argIntValue, err := strconv.Atoi(data[3])
    if err != nil {
        return instruction{}, false
    }

    if data[2] == "-" {
        argIntValue = -argIntValue
    }

    return instruction{
        code:      instructionCode(data[1]),
        argument:  argIntValue,
        execCount: 0,
    }, true
}
