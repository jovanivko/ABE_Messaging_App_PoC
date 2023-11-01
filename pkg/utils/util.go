package utils

import (
	"abeProofOfConcept/internal/logger"
	"bufio"
	"bytes"
	"encoding/gob"
	"os"
	"regexp"
	"strings"
)

func EncodeToBytes(obj interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(obj); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DecodeToObject(data []byte, obj interface{}) error {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	return decoder.Decode(obj)
}

func ParseMessageFile(filePath string) ([]string, []string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		logger.Logger.Errorln("Error opening the file:", err)
		return nil, nil, err
	}
	defer file.Close()

	var conditions []string
	var contents []string

	scanner := bufio.NewScanner(file)
	var currentContent string
	var currentCondition string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "/**") && strings.HasSuffix(line, "**/") {
			// Condition starts and ends on a single line
			currentCondition = strings.TrimPrefix(line, "/**")
			currentCondition = strings.TrimSuffix(currentCondition, "**/")
			currentCondition = strings.TrimSpace(currentCondition)
			conditions = append(conditions, currentCondition)
			if currentContent != "" {
				contents = append(contents, currentContent)
				currentContent = ""
			}
		} else if !strings.HasPrefix(line, "/**") && !strings.HasSuffix(line, "**/") {
			// It's content
			currentContent += line + "\n"
		}
	}

	// After the loop, append the remaining content
	if currentContent != "" {
		contents = append(contents, currentContent)
	}

	if err := scanner.Err(); err != nil {
		logger.Logger.Errorln("Error reading the file:", err)
		return nil, nil, err
	}

	//for _, condition := range conditions {
	//	if !isValidConditionFormat(condition) {
	//		logger.Logger.Errorln(
	//			"Invalid condition format",
	//			"condition", condition,
	//		)
	//		return nil, nil, errors.New("Invalid condition format: " + condition)
	//	}
	//}

	return conditions, contents, nil
}

func isValidConditionFormat(condition string) bool {
	// Regular expression to match the condition format
	regex := `^(?:\w+(@\w+)?(?::\w+)?(?:\s*(?:AND|OR)\s*(?:\(\w+(@\w+)?(?::\w+)?(?:\s*(?:AND|OR)\s*\w+(@\w+)?(?::\w+)?)?\))?)+|\(\w+(@\w+)?(?::\w+)?(?:\s*(?:AND|OR)\s*\w+(@\w+)?(?::\w+)?)?\)|\w+(@\w+)?(?::\w+)?)$
`
	return regexp.MustCompile(regex).MatchString(condition)
}

func CreateAttributesFromEmails(sender string, recipients []string) []string {
	attributes := []string{"email:" + sender}
	for _, rec := range recipients {
		attributes = append(attributes, "email:"+rec)
	}
	return attributes
}

func CreatePolicyFromEmails(sender string, recipients []string) string {
	attributes := CreateAttributesFromEmails(sender, recipients)
	return strings.Join(attributes, " OR ")
}
