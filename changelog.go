package main

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
)

type packageInfo struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	Architecture string `json:"architecture"`
}

func main() {
	newPackages := getPackageVersions("base.lock.json")
	oldPackages := getPackageVersions("old.lock.json")
	changes := ""
	for key := range newPackages {
		archChanges := getDifferences(oldPackages[key], newPackages[key])
		if len(archChanges) > 0 {
			changes += fmt.Sprintf(" - %s\n%s", key, archChanges)
		}
	}
	if len(changes) > 0 {
		fmt.Printf("Updating to %s with the following changes: \n%s", os.Getenv("VERSTRING"), changes)
	} else {
		fmt.Printf("Updating to %s with the misc non package changes\n", os.Getenv("VERSTRING"))
	}
}

func getPackageVersions(filename string) map[string]map[string]string {
	var fileBytes []byte
	var err error
	if fileBytes, err = os.ReadFile(filename); err != nil {
		panic("Unable to read file: " + err.Error())
	}
	versions := struct {
		Contents struct {
			Packages []packageInfo `json:"packages"`
		} `json:"contents"`
	}{}
	if err = json.Unmarshal(fileBytes, &versions); err != nil {
		panic("Unable to parse file: " + filename + ": " + err.Error())
	}
	packages := map[string]map[string]string{}
	for key := range versions.Contents.Packages {
		if packages[versions.Contents.Packages[key].Architecture] == nil {
			packages[versions.Contents.Packages[key].Architecture] = map[string]string{}
		}
		packages[versions.Contents.Packages[key].Architecture][versions.Contents.Packages[key].Name] = versions.Contents.Packages[key].Version
	}
	return packages
}

func getDifferences(old, new map[string]string) string {
	maps.DeleteFunc(new, func(newPackage string, newVersion string) bool {
		return old[newPackage] == newVersion
	})
	change := ""
	if len(new) > 0 {
		for key, value := range new {
			change += fmt.Sprintf("    - %s: %s => %s\n", key, old[key], value)
		}
	}
	return change
}
