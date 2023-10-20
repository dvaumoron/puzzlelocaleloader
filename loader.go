/*
 *
 * Copyright 2023 puzzlelocaleloader authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package puzzlelocaleloader

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

var errNoLocale = errors.New("no locales declared")

// TODO use PasswordStrengthService.GetRules
func Load(localesPath string, allLang []string) (map[string]map[string]string, error) {
	if len(allLang) == 0 {
		return nil, errNoLocale
	}

	messages := map[string]map[string]string{}
	for _, lang := range allLang {
		messagesLang := map[string]string{}
		messages[lang] = messagesLang

		var pathBuilder strings.Builder
		pathBuilder.WriteString(localesPath)
		pathBuilder.WriteString("/messages_")
		pathBuilder.WriteString(lang)
		pathBuilder.WriteString(".properties")
		path := pathBuilder.String()

		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if len(line) != 0 && line[0] != '#' {
				if equal := strings.Index(line, "="); equal > 0 {
					if key := strings.TrimSpace(line[:equal]); key != "" {
						if value := strings.TrimSpace(line[equal+1:]); value != "" {
							messagesLang[key] = value
						}
					}
				}
			}
		}
		if err = scanner.Err(); err != nil {
			return nil, err
		}
	}

	defaultLang := allLang[0]
	messagesDefaultLang := messages[defaultLang]
	for _, lang := range allLang {
		if lang == defaultLang {
			continue
		}
		messagesLang := messages[lang]
		for key, value := range messagesLang {
			if value == "" {
				messagesLang[key] = messagesDefaultLang[key]
			}
		}
	}
	return messages, nil
}
