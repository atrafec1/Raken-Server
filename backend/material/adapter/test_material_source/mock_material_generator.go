package test_material_source

import (
	"math/rand"
	"prg_tools/external/rakenapi"
	"prg_tools/material/domain"
	"time"

	"github.com/bxcodec/faker/v4"
)

var materialNames = []string{
	"Item-1-Hydroseed", "Item-2-Hydromulch", "Item-4-Fiberrolls", "Item-5-Sprinkler", "Item-6-Pipe", "Item-7-Gravel",
	"Item-8-Fiberrolls", "Item-9-Fiberrolls", "Item-10-Fiberrolls", "Item-11-Fiberrolls", "Item-12-Fiberrolls", "Item-13-Fiberrolls",
}

var materialUnits = []string{
	"acres", "pounds", "ls", "tons", "cubic yards", "gallons", "bags", "rolls", "feet", "meters", "square feet", "square meters",
}

func randomDate(fromDate, toDate string) string {
	const layout = "2006-01-02"

	start, err := time.Parse(layout, fromDate)
	if err != nil {
		return ""
	}

	end, err := time.Parse(layout, toDate)
	if err != nil {
		return ""
	}

	if end.Before(start) {
		start, end = end, start // swap if reversed
	}

	// random duration between start and end
	diff := end.Sub(start)
	randomOffset := time.Duration(rand.Int63n(int64(diff)))

	randomTime := start.Add(randomOffset)
	return randomTime.Format(layout)
}

func randomInteger(from, to int) int {
	randdomIntResp, _ := faker.RandomInt(from, to, 1)
	return randdomIntResp[0]
}

func randomFloat(from, to int) float64 {
	randInt := randomInteger(from, to)
	randFloat := rand.Float64() * float64(randInt)
	return randFloat
}
func randomProjectName() string {
	projectNames := []string{"Project Alpha", "Project Beta", "Project Gamma", "Project Delta", "Project Epsilon"}
	return projectNames[randomInteger(0, len(projectNames)-1)]
}

func randomMaterial() rakenapi.Material {
	return rakenapi.Material{
		UUID: faker.UUIDDigit(),
		Name: materialNames[randomInteger(0, len(materialNames)-1)],
		Unit: rakenapi.MaterialUnit{
			Name: materialUnits[randomInteger(0, len(materialUnits)-1)],
		},
	}
}

func generateMockCollections(n int, fromDate, toDate string) []rakenapi.MaterialLogResponse {
	var collections []rakenapi.MaterialLogResponse
	for i := 0; i < n; i++ {
		collections = append(collections, generateMock(100, fromDate, toDate))
	}
	return collections
}

func generateMock(n int, fromDate, toDate string) rakenapi.MaterialLogResponse {
	var logs []rakenapi.MaterialLog
	for i := 0; i < n; i++ {
		logs = append(logs, rakenapi.MaterialLog{
			Date:     randomDate(fromDate, toDate),
			Material: randomMaterial(),
			Quantity: randomFloat(0, 1000),
		})
	}
	return rakenapi.MaterialLogResponse{Collection: logs}
}

func generateMockJob() domain.Job {
	jobNames := []string{"Project Alpha", "Project Beta", "Project Gamma", "Project Delta", "Project Epsilon"}
	numbers := []string{"222222", "333333", "444444", "555555", "666666"}
	return domain.Job{
		Name:   jobNames[randomInteger(0, len(jobNames)-1)],
		Number: numbers[randomInteger(0, len(numbers)-1)],
	}
}
