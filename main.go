package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	// io/ioutil устарел и я не вижу смысла его использовать, в официальной документации предлагается os
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("incorrect args")
		return
	}

	inpFilePath := os.Args[1]
	outFilePath := os.Args[2]

	inpFile, outFile, err := initFiles(inpFilePath, outFilePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer inpFile.Close()
	defer outFile.Close()

	err = calc(inpFile, outFile)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func initFiles(inpFilePath, outFilePath string) (inpFile, outFile *os.File, err error) {
	_, err = os.Stat(inpFilePath)
	if os.IsNotExist(err) {
		return nil, nil, fmt.Errorf("inp file does not exist")
	}

	inpFile, err = os.Open(inpFilePath)
	if err != nil {
		return nil, nil, fmt.Errorf("open inp file error: %w", err)
	}

	outFile, err = os.OpenFile(outFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, nil, fmt.Errorf("create out file error: %w", err)
	}

	return inpFile, outFile, nil
}

func calc(inpFile, outFile *os.File) (err error) {
	var (
		calcRe            *regexp.Regexp
		num1, num2        int
		operation, result string
	)

	writer := bufio.NewWriterSize(outFile, 32)
	calcRe, err = regexp.Compile(`^(\w+)([\+\-\*\/])(\w+)=\?$`)

	scanner := bufio.NewScanner(inpFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		groups := calcRe.FindStringSubmatch(scanner.Text())
		if len(groups) != 4 {
			continue
		}

		num1, err = strconv.Atoi(groups[1])
		if err != nil {
			continue
		}

		operation = groups[2]

		num2, err = strconv.Atoi(groups[3])
		if err != nil {
			continue
		}

		switch operation {
		case "+":
			result = fmt.Sprintf("%d%s%d=%d\n", num1, operation, num2, num1+num2)
		case "-":
			result = fmt.Sprintf("%d%s%d=%d\n", num1, operation, num2, num1-num2)
		case "*":
			result = fmt.Sprintf("%d%s%d=%d\n", num1, operation, num2, num1*num2)
		case "/":
			result = fmt.Sprintf("%d%s%d=%d\n", num1, operation, num2, num1/num2)
		default:
			continue
		}

		_, err = writer.WriteString(result)
		if err != nil {
			return fmt.Errorf("write to out file error: %w", err)
		}
		//_, err = writer.WriteRune('\n')
		//if err != nil {
		//	return fmt.Errorf("write to out file error: %w", err)
		//}
	}

	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("write to out file error: %w", err)
	}

	return nil
}
