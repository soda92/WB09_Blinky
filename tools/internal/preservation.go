package internal

import (
	"fmt"
	"os"
	"regexp"
)

// UserCodeMap maps the tag name (e.g. "PV") to the content string
type UserCodeMap map[string]string

// Regex to find blocks: /* USER CODE BEGIN Tag */ ... /* USER CODE END Tag */
// Go regexp doesn't support backreferences (\1), so we match both tags and verify equality in code.
var blockRegex = regexp.MustCompile(`/\* USER CODE BEGIN ([a-zA-Z0-9_]+) \*/([\s\S]*?)/\* USER CODE END ([a-zA-Z0-9_]+) \*/`)

// ExtractUserCode reads a file and returns a map of all user code blocks.
func ExtractUserCode(path string) (UserCodeMap, error) {
	codes := make(UserCodeMap)
	
	if !FileExists(path) {
		return codes, nil // Return empty map if file doesn't exist
	}

	contentBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	content := string(contentBytes)

	matches := blockRegex.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) == 4 {
			startTag := match[1]
			code := match[2]
			endTag := match[3]
			
			if startTag == endTag {
				codes[startTag] = code
			}
		}
	}
	return codes, nil
}

// InjectUserCode takes file content and a map of user codes, and replaces the empty blocks with user content.
func InjectUserCode(fileContent string, codes UserCodeMap) string {
	// We use ReplaceAllStringFunc to find all blocks in the *new* content
	// and replace them if we have matching user code.
	return blockRegex.ReplaceAllStringFunc(fileContent, func(match string) string {
		// match is the whole string "/* USER ... */ ... /* USER ... */"
		subMatches := blockRegex.FindStringSubmatch(match)
		if len(subMatches) < 4 {
			return match // Should not happen
		}
		startTag := subMatches[1]
		endTag := subMatches[3]
		
		if startTag != endTag {
			return match // Mismatched tags, ignore
		}
		
		if userCode, ok := codes[startTag]; ok {
			// Found saved code for this tag!
			// Construct the block again with saved code.
			return fmt.Sprintf("/* USER CODE BEGIN %s */%s/* USER CODE END %s */", startTag, userCode, startTag)
		}
		
		// If no user code saved, keep the template's content
		return match
	})
}

// RestoreUserCodeInFile reads the current file (which is the Fresh Template),
// injects the saved user code, and writes it back.
func RestoreUserCodeInFile(path string, codes UserCodeMap) error {
	if len(codes) == 0 {
		return nil
	}

	contentBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	
	newContent := InjectUserCode(string(contentBytes), codes)
	
	return os.WriteFile(path, []byte(newContent), 0644)
}
