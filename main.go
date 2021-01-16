package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type SubscriptionDetails struct {
	Description string
	FolderName string
	PublisherID string
	Tags string
	Thumbnail string
	Title string
	Type string
}

type Subscription struct {
	Path string
	Details SubscriptionDetails
}

const bo3WorkshopContentRelativePath = "\\steamapps\\workshop\\content\\311210"

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("ERROR: expecting steam path as input.\n")
		fmt.Printf("Usage:\n\t%s <Steam path here>\n", os.Args[0])
		fmt.Printf("Example:\n\t%s \"%s\"\n", os.Args[0], "C:\\Program Files (x86)\\Steam")
		return
	}

	steamPath := filepath.Clean(os.Args[1])
	path := filepath.Join(steamPath + bo3WorkshopContentRelativePath)
	fmt.Printf("Steam Path:\n\t%s\n", steamPath)
	fmt.Printf("BO3 Workshop Subscriptions Path:\n\t%s\n", path)

	subscriptionPaths, err := getWorkshopSubscriptions(path)
	check(err)
	subscriptions, err := getAllSubscriptionDetails(subscriptionPaths)
	check(err)

	if len(subscriptions) == 0 {
		fmt.Printf("\nNo subscriptions found\n")
		return
	}

	fmt.Printf("\nSubscriptions:\n")
	for i, sub := range subscriptions {
		fmt.Printf("\t%d: %s\n", i, sub.Details.Title)
	}

	fmt.Print("\nDelete which? ")
	subscriptionsToRemove, err := ReadInputNumbers()
	check(err)
	if len(subscriptionsToRemove) == 0 {
		fmt.Println("None selected")
		return
	}

	fmt.Printf("\nSelected Subscriptions:\n")
	for _, i := range subscriptionsToRemove {
		if i >= len(subscriptions) {
			panic(fmt.Sprintf("input out of bounds: %d", i))
		}
		sub := subscriptions[i]
		fmt.Printf("\t%d: %s (%s)\n", i, sub.Details.Title, sub.Path)
	}

	fmt.Printf("\nDelete (y/n)? ")
	yes, err := ReadInputYesNo()
	check(err)
	if !yes {
		fmt.Println("\ncancelled")
		return
	}

	for _, i := range subscriptionsToRemove {
		sub := subscriptions[i]
		if strings.Contains(sub.Path, steamPath) && strings.Contains(sub.Path, bo3WorkshopContentRelativePath) {
			fmt.Printf("\t...Deleting %s\n", sub.Details.Title)
			err = os.RemoveAll(sub.Path)
			check(err)
		} else {
			fmt.Printf("\tNot deleting %s because checks failed.\n", sub.Path)
		}
	}
}

func getWorkshopSubscriptions(workshopDirectory string) ([]string, error) {
	var workshopJSON []string
	err := filepath.Walk(workshopDirectory, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, "workshop.json") {
			workshopJSON = append(workshopJSON, path)
		}

		return nil
	})

	return workshopJSON, err
}

func getSubscriptionDetails(workshopJsonPath string) (details SubscriptionDetails, err error) {
	data, err := ioutil.ReadFile(workshopJsonPath)
	if err != nil {
		return details, err
	}

	err = json.Unmarshal(data, &details)
	return details, err
}

func getAllSubscriptionDetails(subscriptionPaths []string) ([]Subscription, error) {
	var subscriptions []Subscription

	for _, filePath := range subscriptionPaths {
		details, err := getSubscriptionDetails(filePath)
		if err != nil {
			return nil, err
		}

		subscriptions = append(subscriptions, Subscription{
			Path: filepath.Dir(filePath),
			Details: details,
		})
	}

	return subscriptions, nil
}
