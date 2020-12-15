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

    // PART 1 ----->

    seatingMap, err := loadSeatingMap(workingDir + "\\input")
    if err != nil {
        panic(fmt.Sprintf("Could not load input file, error: %v\n", err))
    }

    for {
        changed := seatingMap.runSimulation()
        if !changed {
            break
        }
    }
    fmt.Printf("PART 1: Number of taken seats is: %d\n", seatingMap.countTakenSeats())

    // PART 2 ----->

    seatingMap, err = loadSeatingMap(workingDir + "\\input")
    if err != nil {
        panic(fmt.Sprintf("Could not load input file, error: %v\n", err))
    }

    for {
        changed := seatingMap.runSimulationWithVectors()
        if !changed {
            break
        }
    }
    fmt.Printf("PART 2: Number of taken seats is: %d\n", seatingMap.countTakenSeats())
}

const (
    TakenSeat = '#'
    FreeSeat  = 'L'
)

type coordinate struct {
    x int
    y int
}

type object struct {
    currentState int32
    nextState    int32
}

func (o object) isSeat() bool {
    return o.currentState == TakenSeat || o.currentState == FreeSeat
}

func (o object) isTakenSeat() bool {
    return o.currentState == TakenSeat
}

type seatingMap struct {
    objects map[coordinate]*object
}

// Simulates one step of seating distribution which is always based on the current state of seating plan.
// This simulation take in following rules:
// - seat gets occupied when there is no adjacent seat taken
// - seat is freed when there are 4 or more adjacent seats taken
func (sm *seatingMap) runSimulation() bool {
    hasStateChanged := false
    for objPosition, obj := range sm.objects {
        switch {
        case !obj.isSeat():
            obj.nextState = obj.currentState
            break
        case sm.countAdjacentSeatsTaken(objPosition) == 0:
            obj.nextState = TakenSeat
            break
        case sm.countAdjacentSeatsTaken(objPosition) >= 4:
            obj.nextState = FreeSeat
            break
        default:
            obj.nextState = obj.currentState
        }
        if obj.currentState != obj.nextState {
            hasStateChanged = true
        }
    }
    sm.syncState()
    return hasStateChanged
}

// Simulates one step of seating distribution which is always based on the current state of seating plan.
// This simulation take in following rules:
// - seat gets occupied when there is no taken seat in sight
// - seat is freed when there are 5 or more taken seats in sight
// - seat in sight means one that's located in base perpendicular axis or their 45-degree rotation
func (sm *seatingMap) runSimulationWithVectors() bool {
    hasStateChanged := false
    for objPosition, obj := range sm.objects {
        switch {
        case !obj.isSeat():
            obj.nextState = obj.currentState
            break
        case sm.countVisibleSeatsTaken(objPosition) == 0:
            obj.nextState = TakenSeat
            break
        case sm.countVisibleSeatsTaken(objPosition) >= 5:
            obj.nextState = FreeSeat
            break
        default:
            obj.nextState = obj.currentState
        }
        if obj.currentState != obj.nextState {
            hasStateChanged = true
        }
    }
    sm.syncState()
    return hasStateChanged
}

// Saves the new state of the seating configuration as current state for the next step.
func (sm *seatingMap) syncState() {
    for _, obj := range sm.objects {
        obj.currentState = obj.nextState
    }
}

// Returns the number of taken seats in the whole seating plan.
func (sm seatingMap) countTakenSeats() int {
    takenCount := 0
    for _, obj := range sm.objects {
        if obj.isTakenSeat() {
            takenCount++
        }
    }

    return takenCount
}

// Counts the number of taken seats that are directly in touch with given position (there are up to 8 neighbouring
// points for each individual point).
func (sm seatingMap) countAdjacentSeatsTaken(c coordinate) int {
    // Listing of all adjacent points
    toCheck := []coordinate{
        {x: c.x-1, y: c.y-1},
        {x: c.x, y: c.y-1},
        {x: c.x+1, y: c.y-1},
        {x: c.x-1, y: c.y},
        {x: c.x+1, y: c.y},
        {x: c.x-1, y: c.y+1},
        {x: c.x, y: c.y+1},
        {x: c.x+1, y: c.y+1},
    }

    takenSeatCount := 0
    for _, position := range toCheck {
        // It doesn't matter if there is nothing on given coordinates (over the plan edge)
        if obj, ok := sm.objects[position]; ok && obj.isTakenSeat() {
            takenSeatCount++
        }
    }

    return takenSeatCount
}

// Counts the number of taken seats that are visible from given position. Visibility is possible in horizontal and
// vertical axis and their 45-degree rotation. Any object except the floor block the further view.
func (sm seatingMap) countVisibleSeatsTaken(c coordinate) int {
    // Listing of all 45-degree rotations of base vector <0,1>
    vectors := [][2]int{{0,1},{1,1},{1,0},{1,-1},{0,-1},{-1,-1},{-1,0},{-1,1}}

    takenSeatCount := 0
    for _, vector := range vectors {
        if obj, ok := sm.getVisibleObject(c, vector); ok && obj.isTakenSeat() {
            takenSeatCount++
        }
    }

    return takenSeatCount
}

// Retrieves the first visible object in sight. Visibility is possible in horizontal and vertical axis and their
// 45-degree rotation.
func (sm seatingMap) getVisibleObject(c coordinate, vector [2]int) (*object, bool) {
    currentPosition := coordinate{x: c.x + vector[0], y: c.y + vector[1]}
    for {
        if obj, ok := sm.objects[currentPosition]; ok {
            if obj.isSeat() {
                return obj, true
            }
        } else {
            return nil, false
        }

        currentPosition = coordinate{x: currentPosition.x + vector[0], y: currentPosition.y + vector[1]}
    }
}

// Loads file rows into slice of integers.
// Non-numeric rows are logged and skipped.
func loadSeatingMap(filePath string) (seatingMap, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return seatingMap{}, err
    }
    defer file.Close()

    var seatingMap seatingMap
    seatingMap.objects = make(map[coordinate]*object)
    yPos := 0
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        for xPos, letter := range scanner.Text() {
            seatingMap.objects[coordinate{x: xPos, y: yPos}] = &object{currentState: letter}
        }
        yPos++
    }

    return seatingMap, scanner.Err()
}
