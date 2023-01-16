package parser

import (
	"bytes"
	"regexp"
)

//^Content-Disposition: form-data; name="(.+?)"; filename="(.+?)"$
//^Content-Disposition: form-data; name="(.+?)"$

func FormData(body []byte, delimiter []byte) map[string]interface{} {

	result := make(map[string]interface{})

	contentType_rex := regexp.MustCompile("^Content-Type: (.+?)$")

	hasFile_rex := regexp.MustCompile(`^Content-Disposition: form-data; name="(.+?)"; filename="(.+?)"$`)
	noFile_rex := regexp.MustCompile(`^Content-Disposition: form-data; name="(.+?)"$`)

	formdata_segments := bytes.Split(body, delimiter)

	for index, sub := range formdata_segments {

		if len(sub) == 0 || index == len(formdata_segments)-1 {
			continue
		}

		formdata := bytes.SplitN(sub, []byte("\r\n"), 4)

		disposition := string(formdata[1])

		contentType := string(formdata[2])

		value := formdata[3]

		ct := contentType_rex.FindStringSubmatch(contentType)

		//if there is no file
		if len(ct) == 0 {

			name := noFile_rex.FindStringSubmatch(disposition)[1]

			//Removing the funky \n at the end (I have literally no idea why its there)
			result[name] = string(value[0 : len(value)-2])

			continue
		}

		names := hasFile_rex.FindStringSubmatch(disposition)

		name := names[1]
		filename := names[2]

		file_data := make(map[string]map[string]interface{})
		file_data[name] = make(map[string]interface{})

		file_data[name]["FileName"] = value
		file_data[name]["Data"] = []byte{13, 10}
		file_data[name]["FileName"] = filename
		file_data[name]["Content-Type"] = ct[1]

		result["Files"] = file_data
	}

	return result
}
