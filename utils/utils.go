package utils

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GenRandomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	return randSeq(n)
}

func NewPulumiStringArray(strArray []string) pulumi.StringArray {

	if len(strArray) == 0 {
		return nil
	}

	var pStrArray pulumi.StringArray

	for i := 0; i <= len(strArray)-1; i++ {
		pStrArray = append(pStrArray, pulumi.String(strArray[i]))
	}

	return pStrArray
}

func HttpGet(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return string(body)
}

func ReadFile(fileName string) (text string) {
	// Read File Content
	byteArray, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Provided file name was:%v\n. Error: %v", fileName, err.Error())
	}

	text = string(byteArray)

	return text
}

type BoundaryConfigFile struct {
	Controller_PublicClusterAddress string
	Worker_PublicAddress            string
	Worker_Name                     string
	Worker_Controllers              string
}

func RenderTemplate(fileName string, config BoundaryConfigFile) string {

	// Read config file
	text := ReadFile(fileName)

	// Create Template from text file
	tmpl, err := template.New("test").Parse(text)
	if err != nil {
		log.Fatal(err)
	}

	// Render template with provided config
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, config)
	if err != nil {
		log.Fatal(err)
	}

	return buf.String()
}

func ExitOnError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func GetConfString(ctx *pulumi.Context, key string) string {

	var value string
	envVarName := strings.Replace(strings.ToUpper(key), ".", "_", -1) // Convert "." to underscore in env var

	// Read config file config file
	key = strings.ToLower(key)
	conf := config.New(ctx, "")
	value = conf.Get(key)

	// Check env variable as a fallback for pulumi config
	if value == "" {
		value = os.Getenv(envVarName)
		if value == "" {
			fmt.Println("Could not find config key '" + key + "' in config file or environment variable '" + envVarName + "' \n")
			os.Exit(1)
		}
	}
	return value
}

func Find(a []string, x string) int {
	for i, n := range a {
		if x == n {
			return i
		}
	}

	return -1
}
