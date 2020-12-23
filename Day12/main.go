package main

import (
    "bufio"
    "fmt"
    "math"
    "os"
    "strconv"
)

func main() {

    workingDir, err := os.Getwd()
    if err != nil {
        panic(fmt.Sprintf("Could not establish working directory, error: %v", err))
    }

    instructions, err := loadInstructions(workingDir + "\\input")
    if err != nil {
        panic(fmt.Sprintf("Could not load input file, error: %v\n", err))
    }

    // PART 1 ----->

    ferry := ferry{
        program:   instructions,
        position:  position{0,0},
        direction: East,
    }
    ferry.setSail()
    fmt.Printf("PART 1: Manhattan distance of the ferry is: %d\n", int(math.Abs(float64(ferry.position.ns)) + math.Abs(float64(ferry.position.ew))))

    // PART 2 ----->

    // Second calculation decides movement of 2 objects so it's not implemented as method of the ferry.
    ferry.position = position{0, 0}
    waypoint := position{1, 10}
    for _, instruction := range instructions {
        switch instruction.action {
        case 'N':
            waypoint.ns += instruction.value
        case 'S':
            waypoint.ns -= instruction.value
        case 'E':
            waypoint.ew += instruction.value
        case 'W':
            waypoint.ew -= instruction.value
        case 'L':
            waypoint.rotate(instruction.value)
        case 'R':
            waypoint.rotate(-instruction.value)
        case 'F':
            ferry.moveToWaypoint(&waypoint, instruction.value)
        }
    }
    fmt.Printf("PART 2: Manhattan distance of the ferry is: %d\n", int(math.Abs(float64(ferry.position.ns)) + math.Abs(float64(ferry.position.ew))))

}

type direction int
const (
    North direction = iota
    East
    South
    West
)

type position struct {
    ns int
    ew int
}

// Rotates point around base [0,0] by given amount of angle degrees.
// Rotation is allowed clockwise (negative input) or counter-clockwise (positive input).
func (p *position) rotate(degrees int) {
    radians := float64(degrees) * (math.Pi / 180)
    ew := int(math.Round(float64(p.ew) * math.Cos(radians) - float64(p.ns) * math.Sin(radians)))
    ns := int(math.Round(float64(p.ew) * math.Sin(radians) + float64(p.ns) * math.Cos(radians)))
    p.ew = ew
    p.ns = ns
}

type ferry struct {
    program   []instruction
    position  position
    direction direction
}

// Execute whole instruction flow of onboard computer.
func (f *ferry) setSail() {
    for _, instruction := range f.program {
        f.executeInstruction(instruction)
    }
}

// Determines and executes action based on current instruction.
func (f *ferry) executeInstruction(i instruction) {
    switch i.action {
    case 'N':
        f.position.ns += i.value
    case 'S':
        f.position.ns -= i.value
    case 'E':
        f.position.ew += i.value
    case 'W':
        f.position.ew -= i.value
    case 'L':
        f.rotate(-i.value)
    case 'R':
        f.rotate(i.value)
    case 'F':
        f.moveForward(i.value)
    }
}

// Moves the ferry forward by given distance in direction the ferry is currently facing (N, E, S, W).
func (f *ferry) moveForward(distance int) {
    switch f.direction {
    case North:
        f.position.ns += distance
    case South:
        f.position.ns -= distance
    case East:
        f.position.ew += distance
    case West:
        f.position.ew -= distance
    }
}

// Rotates the facing direction of ferry by certain amount of degrees (which has to be multiple of 90).
// Rotation is allowed clockwise (positive input) or counter-clockwise (negative input).
func (f *ferry) rotate(degrees int) {
    turns := degrees / 90
    rotationVector := int(f.direction) + turns
    if rotationVector > 3 {
        f.direction = direction(rotationVector % 4)
    } else if rotationVector < 0 {
        f.direction = direction(4 + rotationVector)
    } else {
        f.direction = direction(rotationVector)
    }
}

// Moves the ferry toward relatively positioned waypoint by the certain amount steps (given by multiplier param).
// If the waypoint position is NS=3 and EW=-5 and multiplier is 4, final movement is by 12 units North and 20 units West.
func (f *ferry) moveToWaypoint(w *position, multiplier int) {
    f.position.ns += w.ns * multiplier
    f.position.ew += w.ew * multiplier
}

type instruction struct {
    action uint8
    value  int
}

// Loads file rows into slice of instructions.
// Invalid rows are logged and skipped.
func loadInstructions(filePath string) ([]instruction, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var instructions []instruction
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        if inst, ok := parseInstruction(scanner.Text()); ok {
            instructions = append(instructions, inst)
        } else {
            fmt.Printf("skipping invalid input line %q\n", scanner.Text())
        }
    }
    return instructions, scanner.Err()
}

func parseInstruction(data string) (instruction, bool) {
    allowedValues := map[uint8]bool{
        'N': true, 'S': true, 'E': true, 'W': true, 'L': true, 'R': true, 'F': true,
    }

    var result instruction
    if _, ok := allowedValues[data[0]]; ok {
        result.action = data[0]

        value, err := strconv.Atoi(data[1:])
        if err == nil {
            result.value = value
            return result, true
        }
    }

    return result, false
}
