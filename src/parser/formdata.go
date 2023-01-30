package parser

import (
	"bytes"
	"regexp"
)

const (
	contentType_rex_s = "^Content-Type: (.+?)$"
	hasFile_rex_s     = `^Content-Disposition: form-data; name="(.+?)"; filename="(.+?)"$`
	noFile_rex_s      = `^Content-Disposition: form-data; name="(.+?)"$`
)

func FormData(body []byte, delimiter []byte) any {

	result := make(map[string]any)

	contentType_rex := regexp.MustCompile(contentType_rex_s)
	hasFile_rex := regexp.MustCompile(hasFile_rex_s)
	noFile_rex := regexp.MustCompile(noFile_rex_s)

	formdata_segments := bytes.Split(body, delimiter)

	for index, sub := range formdata_segments {

		if len(sub) == 0 || index == len(formdata_segments)-1 {
			continue
		}

		formdata := bytes.SplitN(sub, []byte("\r\n"), 4)

		disposition := string(formdata[1])

		contentType := string(formdata[2])

		//Removing the funky \n at the end (I have literally no idea why its there)
		value := formdata[3]

		ct := contentType_rex.FindStringSubmatch(contentType)

		//if there is no file
		if len(ct) == 0 {

			name := noFile_rex.FindStringSubmatch(disposition)[1]

			result[name] = string(value[0 : len(value)-2])

			continue
		}

		names := hasFile_rex.FindStringSubmatch(disposition)

		if len(names) < 2 {
			continue
		}

		name := names[1]
		filename := names[2]

		file_content := make(map[string]any)
		file_bytes := value[2 : len(value)-2]

		//Images have two funky \n, one at the beginning and one at the end
		file_content["Data"] = file_bytes
		file_content["Size"] = len(file_bytes)
		file_content["Filename"] = filename
		file_content["Content-Type"] = ct[1]

		if result["Files"] == nil {
			result["Files"] = make(map[string]any)
		}

		result["Files"].(map[string]any)[name] = file_content
	}

	return result
}
