package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"upload-files/utils"

	"net/http"

	tus "github.com/eventials/go-tus"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var files []string

	root := os.Getenv("ROOT_DIR")
	err = filepath.Walk(root, visit(&files))
	if err != nil {
		panic(err)
	}

	ssoToken, err := utils.GetSsoToken()
	if err != nil {
		panic(err)
	}
	sub := os.Getenv("MEERA_CLIENT_SUB")

	fileToken := utils.GenFileToken(sub)

	urls := upload(files, ssoToken, fileToken)
	dumpListFiles(urls)
}

func visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		*files = append(*files, path)
		return nil
	}
}

func upload(paths []string, ssoToken, fileToken string) []string {
	var res []string

	for _, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			panic(err)
		}

		//Http Headers
		headers := make(http.Header)
		headers.Add("Authorization", fmt.Sprintf("Bearer %s", ssoToken))
		headers.Add("X-Meera-Storage-Token", fmt.Sprintf("Bearer %s", fileToken))

		config := &tus.Config{
			ChunkSize:           5 * 1024 * 1024, // Cloudflare Stream requires a minimum chunk size of 5MB.
			Resume:              false,
			OverridePatchMethod: false,
			Store:               nil,
			Header:              headers,
			HttpClient:          nil,
		}

		// create the tus client.
		client, _ := tus.NewClient("https://api.dev.meeraspace.com/meerastorage/files/", config)

		// create an upload from a file.
		upload, err := tus.NewUploadFromFile(f)
		if err != nil {
			fmt.Println("fail to Create NewUploadFromFile :" + err.Error())

		}

		// create the uploader.
		uploader, err := client.CreateUpload(upload)
		if err != nil {
			fmt.Println("fail to CreateUpload :" + err.Error())
		}

		// start the uploading process.
		err = uploader.Upload()
		if err != nil {
			fmt.Println("fail to upload file :" + err.Error())

		}
		url := uploader.Url()
		fmt.Println(url)
		res = append(res, url)
	}
	return res
}

func dumpListFiles(urls []string) {
	f, err := os.Create("data.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for _, url := range urls {
		_, err2 := f.WriteString(fmt.Sprintf("%s\n", url))

		if err2 != nil {
			log.Fatal(err2)
		}
	}

	fmt.Println("done")
}
