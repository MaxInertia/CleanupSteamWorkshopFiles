package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type SubscriptionDetails struct {
	Description string
	FolderName  string
	PublisherID string
	Tags        string
	Thumbnail   string
	Title       string
	Type        string
	Size        int64
}

type Subscription struct {
	Path    string
	Details SubscriptionDetails
}

const bo3WorkshopContentRelativePath = "\\steamapps\\workshop\\content\\311210"

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

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
	var totalSize int64
	for i, sub := range subscriptions {
		fmt.Printf("\t%d: %s (%s) [%s]\n", i, sub.Details.Title, sub.Details.Type, formatSize(sub.Details.Size))
		totalSize += sub.Details.Size
	}
	fmt.Printf("\nSpace used: %s\n", formatSize(totalSize))

	fmt.Print("\nDelete which? ")
	subscriptionsToRemove, err := ReadInputNumbers()
	check(err)
	if len(subscriptionsToRemove) == 0 {
		fmt.Println("None selected")
		return
	}

	fmt.Printf("\nSelected Subscriptions:\n")
	var spaceToReclaim int64
	for _, i := range subscriptionsToRemove {
		if i >= len(subscriptions) {
			panic(fmt.Sprintf("input out of bounds: %d", i))
		}
		sub := subscriptions[i]
		spaceToReclaim += sub.Details.Size
		fmt.Printf("\t%d: %s (%s)\n", i, sub.Details.Title, sub.Path)
	}
	fmt.Printf("\nThis will reclaim %s\n", formatSize(spaceToReclaim))

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

			// First walk the directory from deepest to shallowest
			var paths []string
			err = filepath.Walk(sub.Path, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.Mode()&os.ModeSymlink != 0 {
					fmt.Printf("\tSkipping symlink: %s\n", path)
					return nil // Just skip this file, not the whole directory
				}
				paths = append(paths, path)
				return nil
			})
			check(err)

			// Reverse the paths to delete from deepest to shallowest
			for i := len(paths) - 1; i >= 0; i-- {
				path := paths[i]
				err = os.Remove(path)
				if err != nil {
					fmt.Printf("\tCould not remove %s: %v\n", path, err)
					continue
				}
			}
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
	data, err := os.ReadFile(workshopJsonPath)
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

		dirSize, err := getDirSize(filepath.Dir(filePath))
		if err != nil {
			return nil, err
		}
		details.Size = dirSize

		subscriptions = append(subscriptions, Subscription{
			Path:    filepath.Dir(filePath),
			Details: details,
		})

		sort.Slice(subscriptions, func(i, j int) bool {
			return subscriptions[i].Details.Size > subscriptions[j].Details.Size
		})
	}

	return subscriptions, nil
}

func getDirSize(dirPath string) (int64, error) {
	var size int64
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return filepath.SkipDir
		}
		size += info.Size()
		return nil
	})
	return size, err
}
