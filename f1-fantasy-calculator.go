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
	Budget int
	Points int
	Price  int
}

func setupIsGood(setup *Setup, requirements *Requirements) bool {
	for driverIndex := 0; driverIndex < len(setup.Drivers); driverIndex++ {
		setup.Points += setup.Drivers[driverIndex].Points
		setup.Price += setup.Drivers[driverIndex].Price
	}

	if setup.Price > requirements.Budget || setup.Price < requirements.Price || setup.Points < requirements.Points {
		return false
	} else {
		return true
	}
}

func finalizeSetups(budget int, drivers *[]Driver, setupTemplate Setup, indexOfNextNewDriver int,
	requirements *Requirements, setups *[]Setup) *[]Setup {

	for driverCandidateIndex := indexOfNextNewDriver; driverCandidateIndex < len(*drivers); driverCandidateIndex++ {
		setupInProgress := setupTemplate

		if budget-(*drivers)[driverCandidateIndex].Points >= 0 {
			setupInProgress.Drivers = append(setupInProgress.Drivers, (*drivers)[driverCandidateIndex])

			remainingBudget := budget - (*drivers)[driverCandidateIndex].Points

			if remainingBudget > 0 {
				finalizeSetups(remainingBudget, drivers, setupInProgress, driverCandidateIndex+1, requirements, setups)
			}

			// Must have exactly 5 drivers
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

		finalizeSetups(requirements.Budget, &data.Drivers, setupTemplate, 0, requirements, &setups)
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

	fmt.Println("What is your budget? (if you own 10 million money, answer 1000000)")
	fmt.Scanln(&requirements.Budget)
	fmt.Print("Ok, " + strconv.Itoa(requirements.Budget) + " coins of money.\n")

	fmt.Println("What's your minimum allowed price? ")
	fmt.Scanln(&requirements.Price)
	fmt.Print("Ok, price will be over " + strconv.Itoa(requirements.Price) + ".\n")

	fmt.Println("What about your minimum allowed points? ")
	fmt.Scanln(&requirements.Points)
	fmt.Print("And points will be over " + strconv.Itoa(requirements.Points) + ".\n")

	fmt.Println("Starting to create setups.")

	printSetups(createSetups(&data, &requirements))
}
