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

type Results struct {
	TestStart           string
	Executable          string
	Config              Config
	CpuCycles           int
	WrdWriteTotal       int
	WrdWriteMin         int
	WrdWriteMax         int
	WrdWriteAvg         int
	WrdWriteDataPoints  []int
	WrdReadTotal        int
	WrdReadMin          int
	WrdReadMax          int
	WrdReadAvg          int
	WrdReadDataPoints   []int
	WrdDeleteTotal      int
	WrdDeleteMin        int
	WrdDeleteMax        int
	WrdDeleteAvg        int
	WrdDeleteDataPoints []int
	RoReadTotal         int
	RoReadMin           int
	RoReadMax           int
	RoReadAvg           int
	RoReadDataPoints    []int
}

var config Config
var results Results
var silent bool = false
var logLinePrefix string = ""
var testFolder string = ""
var configFile string = ""

var (
	//go:embed VERSION
	version string
	//go:embed report_template.html
	reportTemplate string
)

func main() {

	results.TestStart = time.Now().Format(time.RFC3339)
	results.Executable = os.Args[0]

	mode := "help"
	config = getDefaultConfig()

	dir, err := os.Getwd()
	errorCheck(err)
	testFolder = dir + "/nshptt"

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
		}
	}

	configFile = testFolder + "/config.json"
	resultsCsvPath := testFolder + "/results/" + (time.Now().Format("2006-01-02T15-04-05")) + ".csv"
	resultsHtmlPath := testFolder + "/results/" + (time.Now().Format("2006-01-02T15-04-05")) + ".html"

	resultsCsv := "type;time_ms;value;timestamp\n"

	if mode == "help" {
		fmt.Println("===Nullix Server Hardware Performance Test Tool v" + version + " (NSHPTT)===")
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
		fmt.Printf("  ATTENTION: This will use %s of disk space", getFilesizeDesc(config.DiskTestRoFileSize*config.DiskTestRoFiles))
		fmt.Println("")
		fmt.Println("  --run")
		fmt.Println("  Runs the tests")
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
		logPrint("Using testfolder: " + testFolder)
		readConfig()
		if config.CpuTestRuntime != 0 {
			logLinePrefix = "CPU Test: "
			logPrintf("Started with total expected runtime %dms\n", config.CpuTestRuntime)
			start := time.Now()
			for {
				// do some heavy calc
				math.Sqrt(rand.Float64()/rand.Float64()*rand.Float64() + rand.Float64() - rand.Float64())
				results.CpuCycles++
				if config.CpuTestRuntime > 0 && config.CpuTestRuntime <= int(time.Since(start).Milliseconds()) {
					break
				}
			}
			elapsed := int(time.Since(start).Milliseconds())
			logPrint("Ended")
			logPrintf("%d cycles done\n", results.CpuCycles)
			logPrintf("%dms total runtime\n", int(time.Since(start).Milliseconds()))
			resultsCsv += fmt.Sprintf("CPU Test Cycles;%d;%d;%s\n", elapsed, results.CpuCycles, (time.Now().Format(time.RFC3339)))
		}

		if config.DiskTestWrdFiles > 0 {
			logLinePrefix = "Disk WRD Test: "
			logPrintf("Started with %d files, each %s\n", config.DiskTestWrdFiles, getFilesizeDesc(config.DiskTestWrdFileSize))
			logPrintf("Total Filesize %s\n", getFilesizeDesc(config.DiskTestWrdFiles*config.DiskTestWrdFileSize))
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
				elapsed := int(time.Since(start).Milliseconds())
				resultsCsv += fmt.Sprintf("Disk WRD Test #%d, Write;%d;;%s\n", totalCount, elapsed, (time.Now().Format(time.RFC3339)))
				results.WrdWriteDataPoints = append(results.WrdWriteDataPoints, elapsed)
				results.WrdWriteTotal += elapsed
				if elapsed > results.WrdWriteMax {
					results.WrdWriteMax = elapsed
				}
				if elapsed < results.WrdWriteMin || results.WrdWriteMin == 0 {
					results.WrdWriteMin = elapsed
				}
			}
			results.WrdWriteAvg = results.WrdWriteTotal / config.DiskTestWrdFiles
			logPrintf("Writes done in %dms, min: %dms, max: %dms, average: %dms\n", results.WrdWriteTotal, results.WrdWriteMin, results.WrdWriteMax, results.WrdWriteAvg)
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
				elapsed := int(time.Since(start).Milliseconds())
				resultsCsv += fmt.Sprintf("Disk WRD Test #%d, Read;%d;;%s\n", totalCount, elapsed, (time.Now().Format(time.RFC3339)))
				results.WrdReadDataPoints = append(results.WrdReadDataPoints, elapsed)
				results.WrdReadTotal += elapsed
				if elapsed > results.WrdReadMax {
					results.WrdReadMax = elapsed
				}
				if elapsed < results.WrdReadMin || results.WrdReadMin == 0 {
					results.WrdReadMin = elapsed
				}
			}
			results.WrdReadAvg = results.WrdReadTotal / config.DiskTestWrdFiles
			logPrintf("Read done in %dms, min: %dms, max: %dms, average: %dms\n", results.WrdReadTotal, results.WrdReadMin, results.WrdReadMax, results.WrdReadAvg)
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
				elapsed := int(time.Since(start).Milliseconds())
				resultsCsv += fmt.Sprintf("Disk WRD Test #%d, Delete;%d;;%s\n", totalCount, elapsed, (time.Now().Format(time.RFC3339)))
				results.WrdDeleteDataPoints = append(results.WrdDeleteDataPoints, elapsed)
				errorCheck(err)
				results.WrdDeleteTotal += elapsed
				if elapsed > results.WrdDeleteMax {
					results.WrdDeleteMax = elapsed
				}
				if elapsed < results.WrdDeleteMin || results.WrdDeleteMin == 0 {
					results.WrdDeleteMin = elapsed
				}
			}
			results.WrdDeleteAvg = results.WrdDeleteTotal / config.DiskTestWrdFiles
			logPrintf("Delete done in %dms, min: %dms, max: %dms, average: %dms\n", results.WrdDeleteTotal, results.WrdDeleteMin, results.WrdDeleteMax, results.WrdDeleteAvg)

			logPrint("Ended")
			logPrintf("%dms total runtime\n", results.WrdWriteTotal+results.WrdReadTotal+results.WrdDeleteTotal)
		}
		if config.DiskTestRoFiles > 0 {
			logLinePrefix = "Disk Read-Only Test: "
			logPrintf("Started with %d files, each %s\n", config.DiskTestRoFiles, getFilesizeDesc(config.DiskTestRoFileSize))
			logPrintf("Total Filesize %s\n", getFilesizeDesc(config.DiskTestRoFiles*config.DiskTestRoFileSize))
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
				elapsed := int(time.Since(start).Milliseconds())
				resultsCsv += fmt.Sprintf("Disk Read-Only #%d;%d;;%s\n", totalCount, elapsed, (time.Now().Format(time.RFC3339)))
				results.RoReadDataPoints = append(results.RoReadDataPoints, elapsed)
				results.RoReadTotal += elapsed
				if elapsed > results.RoReadMax {
					results.RoReadMax = elapsed
				}
				if elapsed < results.RoReadMin || results.RoReadMin == 0 {
					results.RoReadMin = elapsed
				}
			}
			results.RoReadAvg = results.RoReadTotal / config.DiskTestRoFiles
			logPrintf("Read done in %dms, min: %dms, max: %dms, average: %dms\n", results.RoReadTotal, results.RoReadMin, results.RoReadMax, results.RoReadAvg)
			logPrint("Ended")
			logPrintf("%dms total runtime\n", results.RoReadTotal)
		}
		os.WriteFile(resultsCsvPath, []byte(resultsCsv), 0777)
		resultsJsonStr, err := json.Marshal(results)
		errorCheck(err)
		reportTemplate = strings.Replace(reportTemplate, "// REPORT DATA", "window.reportData = "+string(resultsJsonStr), -1)
		os.WriteFile(resultsHtmlPath, []byte(reportTemplate), 0777)
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

	results.Config = config
}

func getFilesizeDesc(size int) string {
	str := strconv.Itoa(size) + " bytes ("
	if size > 1024*1024 {
		str += strconv.FormatFloat(float64(float64(size)/1024.0/1024.0), 'f', 2, 64) + " MB"
	} else if size > 1024 {
		str += strconv.FormatFloat(float64(float64(size)/1024.0), 'f', 2, 64) + " KB"
	}
	str += ")"
	return str
}
