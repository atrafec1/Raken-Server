package main

import (
	"fmt"
	"prg_tools/material"
)

func main() {
	progSvc := material.NewTestProgressEstimateService()

	materialCollections, err := progSvc.GetMaterialLogs("2025-11-09", "2025-11-15")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Material Collections: %+v\n", materialCollections)

}
