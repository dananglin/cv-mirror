package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

const (
	aspellLang       string = "en_GB"
	defaultImageName string = "cv-builder"
	defaultImageTag  string = "latest"
)

// ImageName generates the docker image name from the environment
// variables IMAGE_NAME and IMAGE_TAG. Default vaules will be used
// if these variables are not set or empty.
func ImageName() string {
	name := os.Getenv("IMAGE_NAME")
	if len(name) == 0 {
		name = defaultImageName
	}

	tag := os.Getenv("IMAGE_TAG")
	if len(tag) == 0 {
		tag = defaultImageTag
	}

	return name + ":" + tag
}

// CreateCVTex generates the CV docuemnt as a Tex file.
func CreateCVTex(cvJsonDataFile, cvTemplateDir, cvOutputDir, cvOutputFileName string) error {
	fmt.Printf("INFO: Attempting to read data from %s...\n", cvJsonDataFile)

	data, err := ioutil.ReadFile(cvJsonDataFile)
	if err != nil {
		return fmt.Errorf("ERROR: Unable to read data from file. %s\n", err.Error())
	}
	fmt.Println("INFO: Successfully read data.")

	fmt.Println("INFO: Attempting to unmarshal JSON data...")

	var cv Cv
	if err = json.Unmarshal(data, &cv); err != nil {
		return fmt.Errorf("ERROR: Unable to unmarshal JSON data. %s\n", err.Error())
	}
	fmt.Println("INFO: JSON unmarshalling was successful.")

	// if CV_CONTACT_PHONE is set then add it to the CV
	phone := os.Getenv("CV_CONTACT_PHONE")
	if len(phone) > 0 {
		cv.Contact.Phone = phone
	}

	cvOutputPath := cvOutputDir + "/" + cvOutputFileName
	fmt.Printf("INFO: Attempting to create output file %s...\n", cvOutputPath)

	if err = os.MkdirAll(cvOutputDir, 0750); err != nil {
		return fmt.Errorf("ERROR: Unable to create output directory %s. %s\n", cvOutputDir, err.Error())
	}

	output, err := os.Create(cvOutputPath)
	if err != nil {
		return fmt.Errorf("ERROR: Unable to create output file %s. %s\n", cvOutputPath, err.Error())
	}
	fmt.Printf("INFO: Successfully created output file %s.\n", cvOutputPath)
	defer output.Close()

	fmt.Println("INFO: Attempting template execution...")
	fmap := template.FuncMap{
		"notLastElement":   notLastElement,
		"join":             join,
		"durationToString": durationToString,
	}

	t := template.Must(template.New("cv.tmpl.tex").Funcs(fmap).Delims("<<", ">>").ParseGlob(cvTemplateDir + "*.tmpl.tex"))

	if err = t.Execute(output, cv); err != nil {
		return fmt.Errorf("ERROR: Unable to execute the CV template. %s\n", err.Error())
	}
	fmt.Println("INFO: Template execution successful.")

	return nil
}

func SpellCheck(dataFile, aspellPersonalWordlist string) error {
	fmt.Printf("INFO: Reading data from %s...\n", dataFile)

	data, err := ioutil.ReadFile(dataFile)
	if err != nil {
		return fmt.Errorf("unable to read data from %s, %s", dataFile, err.Error())
	}

	// declare the aspell command and its standard input pipe.
	cmd := exec.Command("aspell", "-d", aspellLang, "-p", aspellPersonalWordlist, "list")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	// write the CV data to standard input for piping.
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, string(data))
	}()

	// run aspell and get the list of mispelt words (if any).
	// (the output is a string.)
	fmt.Println("Running aspell...")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	list := strings.Split(string(out), "\n")
	if list[len(list)-1] == "" {
		list = list[:len(list)-1]
	}

	if len(list) > 0 {
		var b strings.Builder
		errMsg := fmt.Sprintf("the following spelling errors were found in %s:", dataFile)
		b.WriteString(errMsg)

		for _, v := range list {
			s := "\n- " + v
			b.WriteString(s)
		}

		return fmt.Errorf(b.String())
	} else {
		fmt.Println("No spelling errors were found.")
	}

	return nil
}
