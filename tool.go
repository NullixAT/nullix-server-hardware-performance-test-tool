package main

import (
	crand "crypto/rand"
	_ "embed"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	CpuTestRuntime      int `json:"cpuTestRuntime"`
	DiskTestWrdFiles    int `json:"diskTestWrdFiles"`
	DiskTestWrdFileSize int `json:"diskTestWrdFileSize"`
	DiskTestRoFiles     int `json:"diskTestRoFiles"`
	DiskTestRoFileSize  int `json:"diskTestRoFileSize"`
}

type TestSetupData struct {
	TestStart       string
	Executable      string
	ToolVersionInfo string
	Config          Config
}

type TestRunCpu struct {
	DataPoints []int
}

type TestRunWrd struct {
	WriteTotal       int
	WriteMin         int
	WriteMax         int
	WriteAvg         int
	WriteDataPoints  []int
	ReadTotal        int
	ReadMin          int
	ReadMax          int
	ReadAvg          int
	ReadDataPoints   []int
	DeleteTotal      int
	DeleteMin        int
	DeleteMax        int
	DeleteAvg        int
	DeleteDataPoints []int
}

type TestRunRo struct {
	Total      int
	Min        int
	Max        int
	Avg        int
	DataPoints []int
}

var config Config
var testSetupData TestSetupData
var silent bool = false
var logLinePrefix string = ""
var testFolder string = ""
var configFile string = ""

var (
	version    string = "dev"
	commitHash string = ""
	buildTime  string = ""
	//go:embed report_template.html
	reportTemplate string
)

func main() {

	toolVersionInfo := "Version: " + version + ", BuildTime: " + buildTime + ", CommitHash: " + commitHash

	testSetupData.TestStart = time.Now().Format(time.RFC3339)
	testSetupData.Executable = os.Args[0]
	testSetupData.ToolVersionInfo = toolVersionInfo

	mode := "help"
	config = getDefaultConfig()

	dir, err := os.Getwd()
	errorCheck(err)
	testFolder = dir + "/nshptt"

	testRuns := 1

	for _, value := range os.Args {
		if strings.HasPrefix(value, "--test-folder") {
			testFolder = value[14:]
			continue
		}
		if value == "--silent" {
			silent = true
			continue
		}
		if value == "--create-config" {
			mode = "createConfig"
		} else if value == "--create-test-files" {
			mode = "createTestFiles"
		} else if value == "--run" {
			mode = "run"
		} else if strings.HasPrefix(value, "--run=") {
			mode = "run"
			testRuns, err = strconv.Atoi(value[6:])
			errorCheck(err)
		}
	}

	configFile = testFolder + "/config.json"
	resultsCsvPath := testFolder + "/results/" + (time.Now().Format("2006-01-02T15-04-05")) + ".csv"
	resultsHtmlPath := testFolder + "/results/" + (time.Now().Format("2006-01-02T15-04-05")) + ".html"

	if mode == "help" {
		fmt.Println("===Nullix Server Hardware Performance Test Tool (NSHPTT)===")
		fmt.Println(toolVersionInfo)
		fmt.Println("https://github.com/NullixAT/nullix-server-hardware-performance-test-tool")
		fmt.Println("Notice: Make sure you have enough disk space for all test files")
		fmt.Println("Notice: Make sure you have enough RAM to load one testfile completely into RAM")
		fmt.Println("Command Line Parameters:")
		fmt.Println("  --test-folder=XXX")
		fmt.Println("  Will use the xxx test folder path instead of the default.")
		fmt.Println("  You can use it to create multiple test configs.")
		fmt.Println("  Also you can use it to run tests on a network drive folder,")
		fmt.Println("  to test network performance as well.")
		fmt.Println("  Default: " + testFolder)
		fmt.Println("")
		fmt.Println("  --silent")
		fmt.Println("  Disables console output while --run")
		fmt.Println("")
		fmt.Println("  --create-config")
		fmt.Println("  This will create the testfolder and a config.json with all default test values.")
		fmt.Println("  You can modify it to your needs.")
		fmt.Println("  Goto our github repo wiki for all details.")
		fmt.Println("  Testfolder: " + testFolder)
		fmt.Println("")
		fmt.Println("  --create-test-files")
		fmt.Println("  This will create all test files that are required for some disk tests.")
		fmt.Printf("  ATTENTION: This will use %s of disk space", getFilesizeDesc(config.DiskTestRoFileSize*config.DiskTestRoFiles, ""))
		fmt.Println("")
		fmt.Println("  --run")
		fmt.Println("  Runs the tests once")
		fmt.Println("")
		fmt.Println("  --run=xx")
		fmt.Println("  Runs the tests xx times and accumalte all results into one result file")
	} else if mode == "createConfig" {
		if !fileExists(testFolder) {
			os.Mkdir(testFolder, 0777)
			fmt.Println("Testfolder " + testFolder + " created")
		}
		if !fileExists(testFolder + "/results") {
			os.Mkdir(testFolder+"/results", 0777)
			fmt.Println("Testfolder " + testFolder + "/results created")
		}
		if !fileExists(configFile) {
			str, err := json.MarshalIndent(getDefaultConfig(), "", "  ")
			errorCheck(err)
			errorCheck(os.WriteFile(configFile, str, 0777))
			fmt.Println("Config File " + configFile + " created")
		}
	} else if mode == "createTestFiles" {
		readConfig()
		logLinePrefix = "Create temporary test files: "
		logPrintf("Started with %d files, each %d bytes\n", config.DiskTestRoFiles, config.DiskTestRoFileSize)
		// write all files
		totalCount := 0
		for {
			totalCount++
			if totalCount > config.DiskTestRoFiles {
				break
			}
			path := testFolder + "/ro_testfile_" + strconv.Itoa(totalCount)
			if fileExists(path) {
				os.Remove(path)
			}
			data := make([]byte, config.DiskTestRoFileSize)
			crand.Read(data)
			os.WriteFile(path, data, 0777)
			logPrint(path + " created")
		}
		logPrint("Ended")
	} else if mode == "run" {

		readConfig()
		logPrint("Using testfolder: " + testFolder)

		resultsCsvFile, err := os.Create(resultsCsvPath)
		errorCheck(err)
		resultsCsvFile.WriteString("type;time_microseconds;value;timestamp\n")

		resultsHtmlFile, err := os.Create(resultsHtmlPath)
		errorCheck(err)
		resultsHtmlFile.WriteString(reportTemplate)
		writeStructToHtmlJsTag(*resultsHtmlFile, "testSetupData", testSetupData)

		for i := 0; i < testRuns; i++ {
			logLinePrefix = ""
			if testRuns > 1 {
				logPrintf("Testrun %d of %d\n", i+1, testRuns)
			}
			if config.CpuTestRuntime != 0 {
				var testResults TestRunCpu
				logLinePrefix = "CPU Test: "
				logPrintf("Started with total expected runtime %dms\n", config.CpuTestRuntime)
				start := time.Now()
				cycles := 0
				for {
					// do some heavy calc
					math.Sqrt(rand.Float64()/rand.Float64()*rand.Float64() + rand.Float64() - rand.Float64())
					cycles++
					if config.CpuTestRuntime > 0 && config.CpuTestRuntime <= int(time.Since(start).Milliseconds()) {
						break
					}
				}
				testResults.DataPoints = append(testResults.DataPoints, cycles)
				elapsed := int(time.Since(start).Milliseconds())
				logPrint("Ended")
				logPrintf("%d cycles done\n", cycles)
				logPrintf("%dms total runtime\n", int(time.Since(start).Milliseconds()))
				resultsCsvFile.WriteString(fmt.Sprintf("CPU Test Cycles;%d;%d;%s\n", elapsed, cycles, (time.Now().Format(time.RFC3339Nano))))
				writeStructToHtmlJsTag(*resultsHtmlFile, "testRunCpu", testResults)
			}

			if config.DiskTestWrdFiles > 0 {
				var testResults TestRunWrd
				logLinePrefix = "Disk WRD Test: "
				logPrintf("Started with %d files, each %s\n", config.DiskTestWrdFiles, getFilesizeDesc(config.DiskTestWrdFileSize, ""))
				totalFileSize := config.DiskTestWrdFiles * config.DiskTestWrdFileSize
				logPrintf("Total Filesize %s\n", getFilesizeDesc(totalFileSize, ""))
				// write all files
				totalCount := 0
				for {
					totalCount++
					if totalCount > config.DiskTestWrdFiles {
						break
					}
					path := testFolder + "/wrd_testfile_" + strconv.Itoa(totalCount)
					if fileExists(path) {
						os.Remove(path)
					}
					data := make([]byte, config.DiskTestWrdFileSize)
					crand.Read(data)
					start := time.Now()
					os.WriteFile(path, data, 0777)
					elapsed := int(time.Since(start).Microseconds())
					resultsCsvFile.WriteString(fmt.Sprintf("Disk WRD Test #%d, Write;%d;;%s\n", totalCount, elapsed, (time.Now().Format(time.RFC3339Nano))))
					testResults.WriteDataPoints = append(testResults.WriteDataPoints, elapsed)
					testResults.WriteTotal += elapsed
					if elapsed > testResults.WriteMax {
						testResults.WriteMax = elapsed
					}
					if elapsed < testResults.WriteMin || testResults.WriteMin == 0 {
						testResults.WriteMin = elapsed
					}
				}
				testResults.WriteAvg = testResults.WriteTotal / config.DiskTestWrdFiles
				logPrintf("Writes done in %dms, min: %dms, max: %dms, average: %dms, %s\n", testResults.WriteTotal/1000, testResults.WriteMin/1000, testResults.WriteMax/1000, testResults.WriteAvg/1000, getAvgTransferSpeed(testResults.WriteTotal, totalFileSize))
				// read all files
				totalCount = 0
				for {
					totalCount++
					if totalCount > config.DiskTestWrdFiles {
						break
					}
					path := testFolder + "/wrd_testfile_" + strconv.Itoa(totalCount)
					if !fileExists(path) {
						panic("Testfile " + path + " removed during test, aborted")
					}
					start := time.Now()
					bytes, err := os.ReadFile(path)
					errorCheck(err)
					if len(bytes) != config.DiskTestWrdFileSize {
						panic("Testfile has " + strconv.Itoa(len(bytes)) + "bytes but expected " + strconv.Itoa(config.DiskTestWrdFileSize))
					}
					elapsed := int(time.Since(start).Microseconds())
					resultsCsvFile.WriteString(fmt.Sprintf("Disk WRD Test #%d, Read;%d;;%s\n", totalCount, elapsed, (time.Now().Format(time.RFC3339Nano))))
					testResults.ReadDataPoints = append(testResults.ReadDataPoints, elapsed)
					testResults.ReadTotal += elapsed
					if elapsed > testResults.ReadMax {
						testResults.ReadMax = elapsed
					}
					if elapsed < testResults.ReadMin || testResults.ReadMin == 0 {
						testResults.ReadMin = elapsed
					}
				}
				testResults.ReadAvg = testResults.ReadTotal / config.DiskTestWrdFiles
				logPrintf("Read done in %dms, min: %dms, max: %dms, average: %dms, %s\n", testResults.ReadTotal/1000, testResults.ReadMin/1000, testResults.ReadMax/1000, testResults.ReadAvg/1000, getAvgTransferSpeed(testResults.ReadTotal, totalFileSize))
				// delete all files
				totalCount = 0
				for {
					totalCount++
					if totalCount > config.DiskTestWrdFiles {
						break
					}
					path := testFolder + "/wrd_testfile_" + strconv.Itoa(totalCount)
					if !fileExists(path) {
						panic("Testfile " + path + " removed during test, aborted")
					}
					start := time.Now()
					err := os.Remove(path)
					elapsed := int(time.Since(start).Microseconds())
					resultsCsvFile.WriteString(fmt.Sprintf("Disk WRD Test #%d, Delete;%d;;%s\n", totalCount, elapsed, (time.Now().Format(time.RFC3339Nano))))
					testResults.DeleteDataPoints = append(testResults.DeleteDataPoints, elapsed)
					errorCheck(err)
					testResults.DeleteTotal += elapsed
					if elapsed > testResults.DeleteMax {
						testResults.DeleteMax = elapsed
					}
					if elapsed < testResults.DeleteMin || testResults.DeleteMin == 0 {
						testResults.DeleteMin = elapsed
					}
				}
				testResults.DeleteAvg = testResults.DeleteTotal / config.DiskTestWrdFiles
				logPrintf("Delete done in %dms, min: %dms, max: %dms, average: %dms, %s\n", testResults.DeleteTotal/1000, testResults.DeleteMin/1000, testResults.DeleteMax/1000, testResults.DeleteAvg/1000, getAvgTransferSpeed(testResults.DeleteTotal, totalFileSize))

				logPrint("Ended")
				logPrintf("%dms total runtime\n", (testResults.WriteTotal+testResults.ReadTotal+testResults.DeleteTotal)/1000)
				writeStructToHtmlJsTag(*resultsHtmlFile, "testRunWrd", testResults)
			}
			if config.DiskTestRoFiles > 0 {
				var testResults TestRunRo
				logLinePrefix = "Disk Read-Only Test: "
				totalFileSize := config.DiskTestRoFiles * config.DiskTestRoFileSize
				logPrintf("Started with %d files, each %s\n", config.DiskTestRoFiles, getFilesizeDesc(config.DiskTestRoFileSize, ""))
				logPrintf("Total Filesize %s\n", getFilesizeDesc(totalFileSize, ""))
				totalCount := 0
				for {
					totalCount++
					if totalCount > config.DiskTestRoFiles {
						break
					}
					path := testFolder + "/ro_testfile_" + strconv.Itoa(totalCount)
					if !fileExists(path) {
						panic("Testfile " + path + " does not exist, create with --create-test-files, aborted")
					}
					start := time.Now()
					bytes, err := os.ReadFile(path)
					errorCheck(err)
					if len(bytes) != config.DiskTestRoFileSize {
						panic("Testfile has " + strconv.Itoa(len(bytes)) + "bytes but expected " + strconv.Itoa(config.DiskTestWrdFileSize))
					}
					elapsed := int(time.Since(start).Microseconds())
					resultsCsvFile.WriteString(fmt.Sprintf("Disk Read-Only #%d;%d;;%s\n", totalCount, elapsed, (time.Now().Format(time.RFC3339))))
					testResults.DataPoints = append(testResults.DataPoints, elapsed)
					testResults.Total += elapsed
					if elapsed > testResults.Max {
						testResults.Max = elapsed
					}
					if elapsed < testResults.Min || testResults.Min == 0 {
						testResults.Min = elapsed
					}
				}
				testResults.Avg = testResults.Total / config.DiskTestRoFiles
				logPrintf("Read done in %dms, min: %dms, max: %dms, average: %dms, %s\n", testResults.Total/1000, testResults.Min/1000, testResults.Max/1000, testResults.Avg/1000, getAvgTransferSpeed(testResults.Total, totalFileSize))
				logPrint("Ended")
				logPrintf("%dms total runtime\n", testResults.Total/1000)
				writeStructToHtmlJsTag(*resultsHtmlFile, "testRunRo", testResults)
			}
		}
		resultsCsvFile.Close()
	}
}

func logPrint(msg string) {
	if silent {
		return
	}
	fmt.Println(logLinePrefix + msg)
}

func logPrintf(msg string, a ...any) {
	if silent {
		return
	}
	fmt.Printf(logLinePrefix+msg, a...)
}

func errorCheck(e error) {
	if e != nil {
		panic(e)
	}
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func getDefaultConfig() Config {
	type defaults Config
	opts := defaults{
		CpuTestRuntime:      5000,
		DiskTestWrdFiles:    100,
		DiskTestWrdFileSize: 1024 * 1024 * 20, // 20Mb
		DiskTestRoFiles:     20,
		DiskTestRoFileSize:  1024 * 1024 * 400, // 400Mb
	}
	return Config(opts)
}

func readConfig() {
	if !fileExists(configFile) {
		panic("Missing " + configFile)
	}
	logPrint("Read config file " + configFile)

	// read our opened jsonFile as a byte array.
	byteValue, err := os.ReadFile(configFile)
	errorCheck(err)

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &config)

	testSetupData.Config = config
}

func getFilesizeDesc(size int, unitAffix string) string {
	str := strconv.Itoa(size) + " bytes" + unitAffix + " ("
	if size > 1024*1024 {
		str += strconv.FormatFloat(float64(float64(size)/1024.0/1024.0), 'f', 2, 64) + " MB"
	} else if size > 1024 {
		str += strconv.FormatFloat(float64(float64(size)/1024.0), 'f', 2, 64) + " KB"
	}
	str += unitAffix + ")"
	return str
}

func getAvgTransferSpeed(us int, bytes int) string {
	if us <= 0 {
		return ""
	}
	return getFilesizeDesc((bytes / us * 1000 * 1000), "/s")
}

func writeStructToHtmlJsTag(file os.File, dataType string, data interface{}) {
	json, err := json.Marshal(data)
	errorCheck(err)
	file.WriteString("<script>addData('" + dataType + "',")
	file.Write(json)
	file.WriteString(")</script>\n")
}
