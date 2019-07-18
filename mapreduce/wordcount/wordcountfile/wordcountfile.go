package wordcountfile

import (
	"SDCC-Project-WorkerNode/mapreduce/wordcount/cloud/amazon"
	"SDCC-Project-WorkerNode/utility"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
)

type WordCountFile struct {
	Name       string
	ParentPath string
	FullPath   string

	OutputSplitTaskPath  string
	OutputMapTaskPath    string
	OutputReduceTaskPath string
	OutputPath           string
}

func New(filename string) (*WordCountFile, error) {

	output := new(WordCountFile)

	(*output).Name = filename
	(*output).ParentPath = path.Join("~/WordCountCache", (*output).Name)
	(*output).FullPath = path.Join((*output).ParentPath, (*output).Name)
	(*output).OutputSplitTaskPath = path.Join((*output).ParentPath, "OutputSplitTaskPath")
	(*output).OutputMapTaskPath = path.Join((*output).ParentPath, "OutputMapTaskPath")
	(*output).OutputReduceTaskPath = path.Join((*output).ParentPath, "OutputReduceTaskPath")
	(*output).OutputPath = path.Join((*output).ParentPath, "Output")

	if err := utility.GenerateDirectory((*output).ParentPath); err != nil {
		return nil, err
	}

	if err := utility.GenerateDirectory((*output).OutputSplitTaskPath); err != nil {
		return nil, err
	}

	if err := utility.GenerateDirectory((*output).OutputMapTaskPath); err != nil {
		return nil, err
	}

	if err := utility.GenerateDirectory((*output).OutputReduceTaskPath); err != nil {
		return nil, err
	}

	if err := utility.GenerateDirectory((*output).OutputPath); err != nil {
		return nil, err
	}

	return output, nil
}

func (obj *WordCountFile) DownloadFromCloud() {

	s3Client := amazon.New()
	(*s3Client).Download((*obj).Name, (*obj).FullPath)
}

func (obj *WordCountFile) Split() error {

	var inputFile *os.File
	var fileInfo os.FileInfo
	var averageChunkSize int64
	var err error

	parts := 6 // zookeeper.GetLocalClusterWorkerPopulation() / system.DefaultArbitraryFaultToleranceLevel

	if inputFile, err = os.Open((*obj).FullPath); err != nil {
		return err
	}
	defer func() {
		utility.CheckError(inputFile.Close())
	}()
	if fileInfo, err = inputFile.Stat(); err != nil {
		return err
	}

	averageChunkSize = fileInfo.Size() / int64(parts)
	readByte := make([]byte, 1)

	for index, currentStartByte, currentEndByte := 0, int64(0), int64(averageChunkSize); ; {

		if _, err = inputFile.Seek(currentEndByte, 0); err != nil {
			return err
		}
		if _, err = inputFile.Read(readByte); err != nil {
			return err
		}

		currentChar := string(readByte[0])

		if strings.Compare(currentChar, " ") == 0 {

			readData := make([]byte, currentEndByte-currentStartByte)

			if _, err = inputFile.Seek(currentStartByte, 0); err != nil {
				return err
			}
			if _, err = inputFile.Read(readData); err != nil {
				return err
			}
			if err := ioutil.WriteFile(path.Join((*obj).OutputSplitTaskPath, strconv.Itoa(index)), readData, 0644); err != nil {
				return err
			}
			index++

			if (currentEndByte + averageChunkSize) < fileInfo.Size() {

				currentStartByte = currentEndByte
				currentEndByte += averageChunkSize

			} else {

				readData := make([]byte, fileInfo.Size()-currentEndByte)

				if _, err = inputFile.Seek(currentEndByte, 0); err != nil {
					return err
				}
				if _, err = inputFile.Read(readData); err != nil {
					return err
				}
				if err := ioutil.WriteFile(path.Join((*obj).OutputSplitTaskPath, strconv.Itoa(index)), readData, 0644); err != nil {
					return err
				}
				break
			}

		} else {
			currentEndByte++
		}

	}
	return nil
}
