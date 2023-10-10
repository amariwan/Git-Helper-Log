package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/apsdehal/go-logger"
	. "github.com/go-git/go-git/_examples"
)

// an array of users
type RepositoryInfos struct {
	Infos []RepositoryInfo `json:"users"`
}
type detailsOfList struct {
	commit    string
	Autor     string
	Date      string
	Leer1     string
	commitlog string
	Leer2     string
	Leer3     string
}

// RepositoryInfo struct which contains a name
// a type and a list of social links
type RepositoryInfo struct {
	Name                string              `json:"name"`
	Folder              string              `json:"folder"`
	Repository          string              `json:"repository"`
	Ignoreprefix        string              `json:"ignorePrefix"`
	Versions            []map[string]string `json:"versions"`
	ParentFolder        string              `json:"parentFolder"`
	FinalPath           string              `json:"finalPath"`
	ParentVersion       string              `json:"parentVersion"`
	ParentRepository    string              `json:"parentRepository"`
	AvailableVersions   []string            `json:"availableVersions"`
	ValidSemverVersions []string            `json:"validSemverVersions"`
	ResolvedVersion     string              `json:"resolvedVersion"`
	ResolvedCommit      string              `json:"resolvedCommit"`
}

type RepositoryInfo4 struct {
	Name             string              `json:"name"`
	Folder           string              `json:"folder"`
	Repository       string              `json:"repository"`
	Ignoreprefix     string              `json:"ignorePrefix"`
	Versions         []map[string]string `json:"versions"`
	FinalPath        string              `json:"finalPath"`
	MatchingVersion  string              `json:"matchingVersion"`
	ResolvedCommit   string              `json:"resolvedCommit"`
	ParentVersion    string              `json:"parentVersion"`
	ParentRepository string              `json:"parentRepository"`
	Dependencies     []RepositoryInfo4   `json:"dependencies"`
}

func repositoryInfoConvert(inputrepositoryMap map[string]RepositoryInfo4) map[string]RepositoryInfo {
	outputRepositoryInfoMap := make(map[string]RepositoryInfo)

	for finalPath, repositoryInfo := range inputrepositoryMap {
		k := finalPath
		var rep RepositoryInfo
		rep.Name = repositoryInfo.Name
		rep.Folder = repositoryInfo.Folder
		rep.ResolvedCommit = repositoryInfo.ResolvedCommit
		rep.FinalPath = repositoryInfo.FinalPath
		outputRepositoryInfoMap[k] = rep

	}
	return outputRepositoryInfoMap
}

func readGithelperListJSONToArray4(name string, mylog *logger.Logger) ([]RepositoryInfo4, error) {
	jsonFile, err := os.Open(name)

	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	// read our opened jsonFile as a byte array.
	byteValue, err2 := ioutil.ReadAll(jsonFile)
	err = err2
	if err != nil {
		return nil, err
	}
	//fmt.Printf(" ist %s\n", byteValue)
	// we initialize our Users array
	var repositoryInfos []RepositoryInfo4
	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above

	err = json.Unmarshal([]byte(byteValue), &repositoryInfos)
	//fmt.Print("that is ma Reposotory %s\n", repositoryInfos)
	// if we os.Open returns an error then handle it
	if err != nil {
		return nil, err
	}

	mylog.Debug("Successfully Opened users.json")
	//fmt.Printf("%s\n", repositoryInfos)
	return repositoryInfos, err
}

func loopRpository(repositoryInfos []RepositoryInfo4, mylog *logger.Logger) map[string]RepositoryInfo4 {
	repositoryInfoMap := make(map[string]RepositoryInfo4)
	map2 := make(map[string]RepositoryInfo4)
	fmt.Printf("in loopRepository\n")

	for i := 0; i < len(repositoryInfos); i++ {
		fmt.Printf("in loopRepository loop %d\n", i)
		ss := repositoryInfos[i].FinalPath

		arrayss := strings.Split(ss, "\\")

		finalPath := strings.Join(arrayss, string(os.PathSeparator))
		fmt.Printf("finalpath working Path in loopRpository is %s\n", finalPath)

		repositoryInfoMap[finalPath] = repositoryInfos[i]
		myrep := repositoryInfos[i]
		fmt.Printf("commit is %s\n", myrep.ResolvedCommit)
		if len(repositoryInfos[i].Dependencies) != 0 {
			map2 = loopRpository(repositoryInfos[i].Dependencies, mylog)
		}
	}

	return merge4(map2, repositoryInfoMap)
}

func readGithelperListJSONToMap4(name string, mylog *logger.Logger) (map[string]RepositoryInfo4, error) {

	repositoryInfos, err := readGithelperListJSONToArray4(name, mylog)

	repositoryInfoMap := make(map[string]RepositoryInfo4)
	map2 := make(map[string]RepositoryInfo4)
	for i := 0; i < len(repositoryInfos); i++ {

		ss := repositoryInfos[i].FinalPath

		arrayss := strings.Split(ss, "\\")

		finalPath := strings.Join(arrayss, string(os.PathSeparator))
		fmt.Printf("finalpath working Path in readGithelperListJSONToMap4 is %s\n", finalPath)

		repositoryInfoMap[finalPath] = repositoryInfos[i]

		fmt.Printf("len is %v\n", len(repositoryInfos[i].Dependencies))
		myrep := repositoryInfos[i]
		fmt.Printf("commit is %s\n", myrep.ResolvedCommit)
		if len(repositoryInfos[i].Dependencies) != 0 {
			map2 = loopRpository(repositoryInfos[i].Dependencies, mylog)
		}

	}

	return merge4(map2, repositoryInfoMap), err
}
func merge4(m1, m2 map[string]RepositoryInfo4) map[string]RepositoryInfo4 {
	res := make(map[string]RepositoryInfo4)
	for finalPath, _ := range m1 {
		res[finalPath] = m1[finalPath]
	}
	for finalPath, _ := range m2 {
		res[finalPath] = m2[finalPath]
	}
	return res
}

func replaceLineBreaks(message string) string {
	re1 := regexp.MustCompile("\n$")
	re2 := regexp.MustCompile("\n+")
	return re2.ReplaceAllString(re1.ReplaceAllString(message, ""), "; ")
	//replacer := strings.NewReplacer("\n", "; ")
	//return replacer.Replace(message)
}

func oneBlank2(message string) string {
	re := regexp.MustCompile("  +")
	return re.ReplaceAllString(message, " ")
}

func noBlank(word string) string {
	t := strings.Split(word, "")
	for i := 0; i < len(t); i++ {
		if t[i] == " " {
			t[i] = ""
		} else {
			break
		}
	}
	return strings.Join(t, "")
}

func smallestPath(myRep map[string]RepositoryInfo) string {
	var array [200]string
	i := 0
	for finalPath := range myRep {
		array[i] = myRep[finalPath].ParentFolder
		i++
	}
	smallest := array[0]
	for i := 0; i < len(myRep); i++ {
		if len(smallest) > len(array[i]) {
			smallest = array[i]
		}
	}
	return smallest
}

func smallestPathArray(myRep []RepositoryInfo, err error) (string, error) {

	smallest := myRep[0].ParentFolder
	for i := 0; i < len(myRep); i++ {
		if len(smallest) > len(myRep[i].ParentFolder) {
			smallest = myRep[i].ParentFolder
		}
	}

	fmt.Printf("smallestPathArray %s\n", smallest)
	return smallest, err
}

func pathFunc(function string, mylog *logger.Logger) string {
	path, err1 := exec.LookPath(filepath.FromSlash(function))
	if err1 != nil {
		log.Fatal("installing fortune is in your future: ", err1)
	}
	mylog.DebugF("fortune is available at %s\n", path)
	return path
}

func readGithelperListJSONToArray(name string, mylog *logger.Logger) ([]RepositoryInfo, error) {
	fmt.Printf("name is %s\n", name)
	jsonFile, err := os.Open(name)
	if err != nil {
		fmt.Printf("err1 is %s", err)
		return nil, err
	}
	defer jsonFile.Close()
	// read our opened jsonFile as a byte array.
	byteValue, err2 := ioutil.ReadAll(jsonFile)
	err = err2
	if err != nil {
		fmt.Printf("err2 is %s", err)
		return nil, err
	}
	// we initialize our Users array
	var repositoryInfos []RepositoryInfo
	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above

	err = json.Unmarshal([]byte(byteValue), &repositoryInfos)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Printf("err3 is %s", err)
		return nil, err
	}

	mylog.Debug("Successfully Opened users.json in vERSION 3")
	return repositoryInfos, err
}

func readGithelperListJSONToMap(name string, mylog *logger.Logger) (map[string]RepositoryInfo, error) {
	arraySmallest, err := smallestPathArray(readGithelperListJSONToArray(name, mylog))
	fmt.Printf("arraySmallest %s", arraySmallest)
	if arraySmallest == "" || err != nil {
		mylog.DebugF("err is %s", err)
		return nil, err

	}
	repositoryInfos, _ := readGithelperListJSONToArray(name, mylog)
	LengthSmallPath := len(strings.Split(arraySmallest, "/"))

	mylog.DebugF("arraySmallest is %s\n", arraySmallest)
	repositoryInfoMap := make(map[string]RepositoryInfo)

	for i := 0; i < len(repositoryInfos); i++ {
		finalPath := repositoryInfos[i].FinalPath
		mylog.DebugF("finalPath is %s\n", finalPath)
		arrayfinalPath := strings.Split(finalPath, "/")
		k := arrayfinalPath[LengthSmallPath:]
		k1 := strings.Join(k, " ")
		k1 = noBlank(k1)
		k1 = " " + k1
		arrayfinalPath = strings.Split(k1, " ")
		k1 = strings.Join(arrayfinalPath, string(os.PathSeparator))
		repositoryInfoMap[k1] = repositoryInfos[i]
	}
	return repositoryInfoMap, err
}

func executeCommand(command string, arguments []string, dir string) string {
	path, err1 := exec.LookPath(filepath.FromSlash(command))
	if err1 != nil {
		log.Fatal("installing fortune is in your future: ", err1)
	}
	cmd := exec.Command(path, arguments...)
	cmd.Dir = dir
	var stdoutBuffer bytes.Buffer
	cmd.Stdout = &stdoutBuffer
	var stderrBuffer bytes.Buffer
	cmd.Stderr = &stderrBuffer
	err2 := cmd.Run()
	if err2 != nil {
		fmt.Print(stderrBuffer.String())
		log.Fatal(err2)
	}
	//fmt.Println(stdoutBuffer.String())
	return stdoutBuffer.String()
}

func getNewFile(firstName, directory string) string {
	tmpfile, err := ioutil.TempFile("", "secondFile")
	if err != nil {
		log.Fatal(err)
	}

	arguments := []string{"list", "-o", tmpfile.Name()}
	fmt.Printf("SecondFolder is %s\n", tmpfile.Name())
	secondFile := executeCommand("githelper", arguments, directory)
	if _, err := tmpfile.Write([]byte(secondFile)); err != nil {
		log.Fatal(err)
	}

	return tmpfile.Name()

}

func parseCommandLineArguments(currentWorkingDirectory string) (string, string, string, bool) {
	var directory, firstName, secondName string
	flag.StringVar(&directory, "d", currentWorkingDirectory, "location of root repository")
	flag.StringVar(&firstName, "f", "", "first githelper list JSON")
	flag.StringVar(&secondName, "s", "", "second githelper list JSON")
	flag.Parse()
	fmt.Printf("secondName %s\n", secondName)
	deleteSecondFile := false
	if secondName == "" {
		secondName = getNewFile(firstName, directory)
		deleteSecondFile = true
	}
	fmt.Printf("secondNameNew---- %s\n", secondName)
	return directory, firstName, secondName, deleteSecondFile
}

func findNewCommitsAndPrint(repositoryInfosFirst map[string]RepositoryInfo, repositoryInfosSecond map[string]RepositoryInfo, directory string, writer2 io.Writer) /*bytes.Buffer*/ {
	myoutput := ""
	//var writer2 bytes.Buffer
	i := 0
	for finalPath, repositoryInfo := range repositoryInfosFirst {

		secondResolvedCommitHash := repositoryInfo.ResolvedCommit
		firstResolvedCommitHash := repositoryInfosSecond[finalPath].ResolvedCommit
		myArray := []string{directory, finalPath}
		u := strings.Join(myArray, string(os.PathSeparator))
		arguments := []string{"log", "--pretty=format:'%cI, %h | %s' ", secondResolvedCommitHash + ".." + firstResolvedCommitHash}
		log := executeCommand("git", arguments, u)
		log = desin(log)
		d := fmt.Sprintf("Results in  %s\n", u)
		e := fmt.Sprintf("log is \n%s\n", log)
		f := fmt.Sprintf("-----------------------------------------------------------------------------------------------------------------")
		if i == 0 {
			myoutput = d + e + f + myoutput
		} else {
			myoutput = d + e + f + "\n" + myoutput
		}
		i = i + 1
	}
	writer2.Write([]byte(myoutput))
	//return writer2
}

func desin(logInput string) string {
	arrayInput := strings.Split(logInput, "\n")

	for i := 0; i < len(arrayInput); i++ {
		if len(arrayInput[i]) > 0 {
			slys1 := arrayInput[i][0:11]
			slys2 := arrayInput[i][12:17]
			slys3 := arrayInput[i][26:]
			words := []string{slys1, slys2, slys3}
			arrayInput[i] = strings.Join(words, " ")
		}
	}
	return strings.Join(arrayInput, "\n")

}

func removeSecondJsonIfRequired(fileName string, doDelete bool) {
	if doDelete {
		os.Remove(fileName)
	}
}

func main() {
	var myerr error
	mylog, myerr := logger.New("test", 1, os.Stdout)
	if myerr != nil {
		panic(myerr) // Check for error
	}
	mylog.SetLogLevel(logger.DebugLevel)

	currentWorkingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	mylog.DebugF("Current Wroking Direcoty: %s\n", currentWorkingDirectory)

	directory, firstName, secondName, deleteSecondFile := parseCommandLineArguments(currentWorkingDirectory)
	defer removeSecondJsonIfRequired(secondName, deleteSecondFile)
	repositoryInfosSecond, err1 := readGithelperListJSONToMap(secondName, mylog)
	if repositoryInfosSecond == nil {
		repositoryInfosSecond1, _ := readGithelperListJSONToMap4(secondName, mylog)
		_, err1 = readGithelperListJSONToMap4(secondName, mylog)
		repositoryInfosSecond = repositoryInfoConvert(repositoryInfosSecond1)
	}
	if err1 != nil {
		log.Fatal(err1)
	}
	repositoryInfosFirst, err2 := readGithelperListJSONToMap(firstName, mylog)
	if repositoryInfosFirst == nil {
		repositoryInfosFirst1, _ := readGithelperListJSONToMap4(firstName, mylog)
		_, err2 = readGithelperListJSONToMap4(firstName, mylog)
		repositoryInfosFirst = repositoryInfoConvert(repositoryInfosFirst1)

	}
	if err2 != nil {
		log.Fatal(err2)
	}
	findNewCommitsAndPrint(repositoryInfosFirst, repositoryInfosSecond, directory, os.Stdout)

	// defer the closing of our jsonFile so that we can parse it later on
	//mypATh := smallestPath(repositoryInfosFirst)
	//arraySmallest := strings.Split(mypATh, "/")
	//mylog.DebugF("arraySmallest is %s\n", arraySmallest)

	// Gets the HEAD history from HEAD, just like this command:
	Info("git log")
	CheckIfError(err)

}
