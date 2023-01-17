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

func FormData(body []byte, delimiter []byte) map[string]interface{} {

	result := make(map[string]interface{})

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

		name_map := make(map[string]interface{})

		//Images have two funky \n, one at the beginning and one at the end
		name_map["Data"] = value[2 : len(value)-2]
		name_map["FileName"] = filename
		name_map["Content-Type"] = ct[1]

		file_data := make(map[string]interface{})

		file_data[name] = name_map

		result["Files"] = file_data
	}

	return result
}
