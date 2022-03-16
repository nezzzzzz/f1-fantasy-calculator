package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

type Data struct {
	Constructors []Constructor `json:"constructors"`
	Drivers      []Driver      `json:"drivers"`
}

type Driver struct {
	Name   string `json:"name"`
	Points int    `json:"points"`
	Price  int    `json:"price"`
}

type Constructor struct {
	Name   string `json:"name"`
	Points int    `json:"points"`
	Price  int    `json:"price"`
}

type Setup struct {
	Constructor Constructor
	Drivers     []Driver
	Points      int
	Price       int
}

type Requirements struct {
	MaximumPrice  int
	MinimumPrice  int
	MinimumPoints int
}

func setupIsGood(setup *Setup, requirements *Requirements) bool {
	for driver := 0; driver < len(setup.Drivers); driver++ {
		setup.Points += setup.Drivers[driver].Points
		setup.Price += setup.Drivers[driver].Price
	}

	return setup.Points >= requirements.MinimumPoints && setup.Price >= requirements.MinimumPrice && setup.Price <= requirements.MaximumPrice
}

func createDriverSetups(budgetLeft int, drivers *[]Driver, setupTemplate Setup, indexOfNextUsableDriver int,
	requirements *Requirements, setups *[]Setup) *[]Setup {

	for driverCandidateIndex := indexOfNextUsableDriver; driverCandidateIndex < len(*drivers); driverCandidateIndex++ {
		setupInProgress := setupTemplate

		if budgetLeft-(*drivers)[driverCandidateIndex].Price >= 0 {
			setupInProgress.Drivers = append(setupInProgress.Drivers, (*drivers)[driverCandidateIndex])

			remainingBudget := budgetLeft - (*drivers)[driverCandidateIndex].Price

			if remainingBudget > 0 {
				createDriverSetups(remainingBudget, drivers, setupInProgress, driverCandidateIndex+1, requirements, setups)
			}

			// Must have exactly 5 drivers.
			if len(setupInProgress.Drivers) != 5 {
				continue
			} else {
				if setupIsGood(&setupInProgress, requirements) {
					*setups = append(*setups, setupInProgress)
				}
			}
		}
	}

	return setups
}

func createSetups(data *Data, requirements *Requirements) []Setup {
	var setups []Setup

	for constructorIndex := 0; constructorIndex < len(data.Constructors); constructorIndex++ {
		constructor := data.Constructors[constructorIndex]
		setupTemplate := Setup{
			Constructor: constructor, Drivers: []Driver{}, Points: constructor.Points, Price: constructor.Price}

		createDriverSetups(requirements.MaximumPrice-constructor.Price, &data.Drivers, setupTemplate, 0, requirements, &setups)
	}

	return setups
}

func printSetups(setups []Setup) {

	fmt.Print("\nDone! " + strconv.Itoa(len(setups)) + " setups created.\n")

	for i := 0; i < len(setups); i++ {
		fmt.Println("\n----- Setup " + strconv.Itoa(i+1) + " -----\n" +
			"\n - Constructor: " + setups[i].Constructor.Name +
			"\n - Drivers:")

		for j := 0; j < len(setups[i].Drivers); j++ {
			fmt.Println("   - " + setups[i].Drivers[j].Name)
		}

		fmt.Println("\n - Points:  " + strconv.Itoa(setups[i].Points) +
			" - Price: " + strconv.Itoa(setups[i].Price))
	}
}

func main() {

	fmt.Print("Welcome to F1 Fantasy Calculator!" + "\n")

	filename := "data.json"
	jsonFile, err := os.Open(filename)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully opened " + filename + ".")

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var data Data
	json.Unmarshal(byteValue, &data)

	var requirements Requirements

	fmt.Println("What is your maximum budget?")
	fmt.Scanln(&requirements.MaximumPrice)
	fmt.Println("What is your minimum allowed price?")
	fmt.Scanln(&requirements.MinimumPrice)
	fmt.Print("Ok, price will be between " + strconv.Itoa(requirements.MinimumPrice) + " - " + strconv.Itoa(requirements.MaximumPrice) + " coins of money.\n")

	fmt.Println("What about your minimum allowed points? ")
	fmt.Scanln(&requirements.MinimumPoints)
	fmt.Print("And points will be at least " + strconv.Itoa(requirements.MinimumPoints) + ".\n")

	fmt.Println("Starting to create setups.")

	printSetups(createSetups(&data, &requirements))
}
