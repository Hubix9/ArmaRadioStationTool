package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"github.com/g0dsCookie/gopbo/pbo"
)

var addonName string
var addonDir string
var authorName string
var addonDescription string
var addonPath string
var workDir string
var cfgPatches string
var cfgRadioStations string
var modCpp string

type radioStation struct {
	name string
	url  string
}

var radioStations []radioStation

func getUserInput() {
	//Getting addon name
	fmt.Print("Please enter name for your addon: ")
	fmt.Scanln(&addonName)
	//Getting user's nickname
	fmt.Print("Enter your nickname: ")
	fmt.Scanln(&authorName)
	//Getting addon description
	fmt.Print("Enter addon's description: ")
	fmt.Scanln(&addonDescription)
}

//Clearing unwanted characters from user input
func clearUserInput() {
	clearRegex := regexp.MustCompile(`-|&|@|\[|\]|\#|\%|\*|\^|\!|\'|\"|\.|,| |\(|\)`)
	addonName = clearRegex.ReplaceAllString(addonName, "_")
	authorName = clearRegex.ReplaceAllString(authorName, "_")
	addonDescription = clearRegex.ReplaceAllString(addonDescription, "_")
}

func createFolderStruct() {
	if _, err := os.Stat(addonPath); !os.IsNotExist(err) {
		os.RemoveAll(addonPath)
	}
	os.MkdirAll(path.Join(addonPath, "Addons", addonName), os.ModePerm)
}

func readStations() {
	files, err := ioutil.ReadDir(path.Join(workDir, "stations"))
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		fileExt := filepath.Ext(file.Name())
		fileName := file.Name()[0 : len(file.Name())-len(fileExt)]
		fileData, err := ioutil.ReadFile(filepath.Join(workDir, "stations", file.Name()))
		if err != nil {
			log.Fatal(err)
		}
		radioStations = append(radioStations, radioStation{fileName, string(fileData)})
	}
}

func generateCfgPatches() {
	cfgPatches = fmt.Sprintf(`class CfgPatches
	{
		class %s
		{
			units[]={};
			weapons[]={};
			requiredAddons[]={"A3_Characters_F_BLUFOR"};
		};
	};`, addonName)
}

func genereateCfgRadioStations() {
	clearRegex := regexp.MustCompile(`-|&|@|\[|\]|\#|\%|\*|\^|\!|\'|\"|\.|,| |\(|\)`)
	cfgRadioStations += "class CfgRadioStations { \n"
	for _, station := range radioStations {
		cfgRadioStations += fmt.Sprintf(`	class %s {
			name = "%s";
			url = "%s"
		};
		`, clearRegex.ReplaceAllString(station.name, "_"), station.name, station.url)
	}
	cfgRadioStations += "\n };"
}

func generateModCpp() {
	modCpp = fmt.Sprintf(`name = %s;
	tooltip = %s;
	overview = %s;
	author = %s;`, addonName, addonDescription, addonDescription, authorName)
}

func writeModCpp() {
	byteData := []byte(modCpp)
	ioutil.WriteFile(path.Join(addonPath, "mod.cpp"), byteData, os.ModePerm)
}

func writeConfigCpp() {
	byteData := []byte(cfgPatches + "\n" + cfgRadioStations)
	ioutil.WriteFile(path.Join(addonPath, "Addons", addonName, "config.cpp"), byteData, os.ModePerm)
}

func createPbo() {
	pbo.Pack(filepath.Join(addonPath, "Addons", addonName), filepath.Join(addonPath, "Addons", addonName+".pbo"), true)
}

func cleanTempFiles() {
	err := os.RemoveAll(filepath.Join(addonPath, "Addons", addonName))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	workDir, _ = os.Getwd()
	getUserInput()
	addonDir = "@" + addonName
	addonPath = path.Join(workDir, addonDir)
	createFolderStruct()
	readStations()
	generateCfgPatches()
	genereateCfgRadioStations()
	generateModCpp()
	writeModCpp()
	writeConfigCpp()
	createPbo()
	cleanTempFiles()
}
