package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/apsdehal/go-logger"
	"github.com/stretchr/testify/assert"
	"github.com/traherom/memstream"
)

func TestUnitNoBlank(t *testing.T) {

	_, assert := SetUp(t)
	assert.Equal("Du lieber Gott", noBlank(" Du lieber Gott"), "they should be equal")
	assert.Equal("a", noBlank(" a"), "they should be equal")
	assert.Equal("", noBlank("                "), "they should be equal")
	assert.Equal("", noBlank(" "), "they should be equal")
	assert.Equal("\n         ", noBlank("   \n         "), "they should be equal")
}

func TestUnitReplaceLineBreaks(t *testing.T) {

	_, assert := SetUp(t)
	assert.Equal("Ananas schmeckt ;              gut; ", replaceLineBreaks("Ananas schmeckt \n\n             gut\n\n"), "they should be equal")
	assert.Equal("; Zukini schmekt mit gar  nicht", replaceLineBreaks("\nZukini schmekt mit gar  nicht"), "they should be equal")
	assert.Equal(" ", replaceLineBreaks(" "), "they should be equal")
	assert.Equal("; ", replaceLineBreaks("\n\n\n"), "they should be equal")
	assert.Equal("   ;          ", replaceLineBreaks("   \n         "), "they should be equal")
}

func TestIntegrationFindNewCommitsAndPrint(t *testing.T) {
	mylog, assert := SetUp(t)
	var writer1 bytes.Buffer
	//var writer2 bytes.Buffer
	arrexpectedString1 := []string{"Results in  .", "Tests", "githelperlog-integrationtest1", "githelperlog-integrationtest2", "githelperlog-integrationtest3\nlog is \n'2020-10-12 14:49 , b965de2 | no change now' \n-----------------------------------------------------------------------------------------------------------------"}
	expectedString1 := strings.Join(arrexpectedString1, string(os.PathSeparator))
	arrexpectedString2 := []string{"Results in  .", "Tests", "githelperlo"}
	expectedString2 := strings.Join(arrexpectedString2, string(os.PathSeparator))
	arrexpectedString3 := []string{"g-integrationtest1", "githelperlog-integrationtest2\nlog is \n'2020-10-1"}
	expectedString3 := strings.Join(arrexpectedString3, string(os.PathSeparator))
	arrexpectedString4 := []string{"2 14:46 , 9663a36 | no change' \n-----------------------------------------------------------------------------------------------------------------\n"}
	expectedString4 := strings.Join(arrexpectedString4, string(os.PathSeparator))
	expectedString := expectedString2 + expectedString3 + expectedString4 + expectedString1
	writer1.Write([]byte(expectedString))
	directoryArray := []string{".", "Tests", "githelperlog-integrationtest1"}
	directory := strings.Join(directoryArray, string(os.PathSeparator))
	firstNameArray := []string{".", "Tests", "githelperlog-integrationtest1", "List1"}
	firstName := strings.Join(firstNameArray, string(os.PathSeparator))
	secondNameArray := []string{".", "Tests", "githelperlog-integrationtest1", "List2"}
	secondName := strings.Join(secondNameArray, string(os.PathSeparator))
	repositoryInfosFirst1, err1 := readGithelperListJSONToMap4(firstName, mylog)
	repositoryInfosFirst := repositoryInfoConvert(repositoryInfosFirst1)
	repositoryInfosSecond1, err2 := readGithelperListJSONToMap4(secondName, mylog)
	assert.Equal(nil, err1, "they should be equal")
	assert.Equal(nil, err2, "they should be equal")
	repositoryInfosSecond := repositoryInfoConvert(repositoryInfosSecond1)
	var writer = memstream.New()
	findNewCommitsAndPrint(repositoryInfosFirst, repositoryInfosSecond, directory, writer)
	myString := string(writer.Bytes())
	mylog.DebugF("myoutput is %s\n", myString)
	assert.Equal(expectedString, myString, "they should be equal")

}
func TestUnitBarOneBlank2(t *testing.T) {
	_, assert := SetUp(t)

	assert.Equal("Ananas schmeckt gut", oneBlank2("Ananas schmeckt              gut"), "they should be equal")
	assert.Equal("Zukini schmekt mit gar nicht", oneBlank2("Zukini schmekt mit gar  nicht"), "they should be equal")
	assert.Equal(" ", oneBlank2(" "), "they should be equal")
	assert.Equal("\n\n\n", oneBlank2("\n\n\n"), "they should be equal")
	assert.Equal(" \n ", oneBlank2(" \n "), "they should be equal")
	assert.Equal(" \n ", oneBlank2("  \n "), "they should be equal")
	assert.Equal(" \n \n ", oneBlank2("                 \n         \n  "), "they should be equal")
	assert.Equal(" ", oneBlank2("                               "), "they should be equal")
}

func TestUnitRepositoryInfoConvert(t *testing.T) {
	_, assert := SetUp(t)
	inputRep := make(map[string]RepositoryInfo4)
	outputRep := make(map[string]RepositoryInfo)
	var rep4 RepositoryInfo4
	var rep RepositoryInfo
	rep4.Name = "name"
	rep4.Folder = "folder"
	rep4.ResolvedCommit = "commit"
	rep4.FinalPath = "finalpath"
	rep.Name = "name"
	rep.Folder = "folder"
	rep.ResolvedCommit = "commit"
	rep.FinalPath = "finalpath"
	inputRep["1"] = rep4
	outputRep["1"] = rep

	assert.Equal(outputRep, repositoryInfoConvert(inputRep), "they should be equal")
}

func SetUp(t *testing.T) (*logger.Logger, *assert.Assertions) {
	assert := assert.New(t)
	var myerr error
	mylog, myerr := logger.New("test", 1, os.Stdout)
	if myerr != nil {
		panic(myerr) // Check for error
	}
	return mylog, assert
}
func TestUnitReadGithelperListJSONMap4(t *testing.T) {
	mylog, assert := SetUp(t)
	repositoryInfoMap := make(map[string]RepositoryInfo4)
	myarray1 := []string{"rel4", "external_repos", "fipMedia", "buildhelper"}
	finalPath1 := strings.Join(myarray1, string(os.PathSeparator))
	expectedRep := repositoryInfoMap["rel4/external_repos/fipMedia/buildhelper"]
	expectedRep.ResolvedCommit = "a62c79d1bda3c815df3cf9468da8b8b054630547"
	outputRrad, err := readGithelperListJSONToMap4("./Read/githelperlist_v4.json", mylog)
	assert.Equal(nil, err, "they should be equal")
	output := outputRrad[finalPath1]

	expectedRep2 := repositoryInfoMap["rel4/external_repos/fipMedia/libuv"]
	expectedRep2.ResolvedCommit = "d3e5239698763dd1651c598a8a3ce27434886ab6"
	myarray2 := []string{"rel4", "external_repos", "fipMedia", "libuv"}
	finalPath2 := strings.Join(myarray2, string(os.PathSeparator))
	output2 := outputRrad[finalPath2]

	assert.Equal(expectedRep.ResolvedCommit, output.ResolvedCommit, "they should be equal")
	assert.Equal(expectedRep2.ResolvedCommit, output2.ResolvedCommit, "they should be equal")

}

func TestUnitReadGithelperListJSONMap(t *testing.T) {

	mylog, assert := SetUp(t)
	repositoryInfoMap := make(map[string]RepositoryInfo)
	myarray3 := []string{"", "rel4", "external_repos", "fipMedia", "buildhelper"}
	finalPath3 := strings.Join(myarray3, string(os.PathSeparator))
	expectedRep := repositoryInfoMap[finalPath3]
	expectedRep.ResolvedCommit = "a62c79d1bda3c815df3cf9468da8b8b054630547"
	outputRrad, err := readGithelperListJSONToMap("./Read/githelperlist_v3.json", mylog)
	assert.Equal(nil, err, "they should be equal")
	output := outputRrad[finalPath3]
	k := output.ResolvedCommit
	assert.Equal(expectedRep.ResolvedCommit, k, "they should be equal")
}
