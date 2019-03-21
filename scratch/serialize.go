package main

import (
	"encoding/binary"
	"fmt"
	"os"

	database "github.com/ear7h/2019-03-database"
)

func main() {
	row := database.Country{
		Name: "me town",
		Population: 1,
		// ignore rest
	}

	//copy(row.Name[:], []byte("me town"))

	if binary.Size(row) == -1 {
		fmt.Fprintln(os.Stderr, "variable size")

		os.Stdout.Write(append([]byte(row.Name), 0))
		//os.Stdout.Write(append(row.Name[:], 0))
		binary.Write(os.Stdout, binary.LittleEndian, row.Population)
		binary.Write(os.Stdout, binary.LittleEndian, row.North)
		binary.Write(os.Stdout, binary.LittleEndian, row.South)
		binary.Write(os.Stdout, binary.LittleEndian, row.East)
		binary.Write(os.Stdout, binary.LittleEndian, row.West)
	} else {
		fmt.Fprintln(os.Stderr, "static size: ", binary.Size(row))

		binary.Write(os.Stdout, binary.LittleEndian, row)
	}
}
