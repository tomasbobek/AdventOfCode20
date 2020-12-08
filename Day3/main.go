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

    slopeMap, err := loadMap(workingDir + "\\input")
    if err != nil {
        panic(fmt.Sprintf("Could not load input file, error: %v\n", err))
    }

    // PART 1 ----->

    encounteredTrees := slopeMap.countEncounteredTreesForMovement(3, 1)
    fmt.Printf("PART 1: Encountered %d trees during the ride\n", encounteredTrees)

    // PART 2 ----->

    encounteredTreesRide1 := slopeMap.countEncounteredTreesForMovement(1, 1)
    encounteredTreesRide2 := slopeMap.countEncounteredTreesForMovement(3, 1)
    encounteredTreesRide3 := slopeMap.countEncounteredTreesForMovement(5, 1)
    encounteredTreesRide4 := slopeMap.countEncounteredTreesForMovement(7, 1)
    encounteredTreesRide5 := slopeMap.countEncounteredTreesForMovement(1, 2)
    fmt.Printf("PART 2: Result is %d", encounteredTreesRide1 * encounteredTreesRide2 * encounteredTreesRide3 * encounteredTreesRide4 * encounteredTreesRide5)
}

type coordinate struct {
    x int
    y int
}

type slopeMap struct {
    width  int
    height int
    trees  map[coordinate]bool
}

func (sm slopeMap) countEncounteredTreesForMovement(xMov, yMov int) int {
    xPos, yPos, encounteredTrees := 0, 0, 0

    for {
        if _, ok := sm.trees[coordinate{
            x: xPos,
            y: yPos,
        }]; ok {
            encounteredTrees++
        }

        // If we reach the right boundary of the map, we should continue from the (left) beginning as the map repeats
        // infinitely to the right direction.
        if xPos+xMov >= sm.width {
            xPos = xPos+xMov-sm.width
        } else {
            xPos += xMov
        }
        yPos += yMov

        // While the position is equal to the map height, we have reached the slope end.
        if yPos >= sm.height {
            break
        }
    }

    return encounteredTrees
}

// Loads file rows into Slope Map structure where "#" indicates tree and "." empty space.
// Invalid characters are treated as empty space.
func loadMap(filePath string) (slopeMap, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return slopeMap{}, err
    }
    defer file.Close()

    var slopeMap slopeMap
    slopeMap.trees = make(map[coordinate]bool)
    scanner := bufio.NewScanner(file)
    yCoord, mapWidth := 0, 0
    for scanner.Scan() {
        mapWidth = len(scanner.Text())
        for xCoord, char := range scanner.Text() {
            if char == '#' {
                slopeMap.trees[coordinate{
                    x: xCoord,
                    y: yCoord,
                }] = true
            }
        }
        yCoord++
    }

    slopeMap.height = yCoord
    slopeMap.width = mapWidth

    return slopeMap, scanner.Err()
}