package main

import (
	"bufio"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"ownbucket/orgdata"
	"strings"
	"sync"
	"time"
)

const (
	Reset   = "\033[0m"
	Green   = "\033[32m"
	Gray    = "\033[90m"
	Red     = "\033[31m"
	Blue    = "\033[34m"
	Cyan    = "\033[36m"
	Yellow  = "\033[33m"
	Dim     = "\033[2m"
	Magenta = "\033[35m"
	Reset2  = "\033[49m"
)

func banner() {
	version := "1.0"
	banner := fmt.Sprintf(`%s
   ___                 ____             _        _
  / _ \__      ___ __ | __ ) _   _  ___| | _____| |_
 | | | \ \ /\ / / '_ \|  _ \| | | |/ __| |/ / _ \ __|
 | |_| |\ V  V /| | | | |_) | |_| | (__|   <  __/ |_
  \___/  \_/\_/ |_| |_|____/ \____|\___|_|\_\___|\__| v%s
%s
   `, Red, version, Reset)

	fmt.Println(banner)
	fmt.Printf("%s\t OwnBucket : Bucket Enumeration Tool %s [BY : @mayank_pandey01]%s\n\n\n", Magenta, Reset, Reset)
}

func createWordlist(target string) []string {
	env := []string{"dev", "development", "stage", "s3", "staging", "prod", "production", "test", "frontend", "backend", "temp"}
	var bucketPrefix []string
	suff := []string{"-", "."}
	pref := []string{"-", "."}
	file, err := os.Open("bucket_prefixes.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		bucketPrefix = append(bucketPrefix, line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	var final []string

	var tmp1 []string
	for _, i := range bucketPrefix {
		tmp1 = append(tmp1, fmt.Sprintf("%s-%s", target, i))
	}

	var tmp2 []string
	for _, i := range env {
		tmp2 = append(tmp2, fmt.Sprintf("%s-%s", target, i))
	}

	var tmp3 []string
	for _, i := range bucketPrefix {
		for _, j := range env {
			for _, k := range pref {
				for _, l := range suff {
					tmp3 = append(tmp3, fmt.Sprintf("%s%s%s%s%s", target, k, i, l, j))
				}
			}
		}
	}
	final = append(final, target)
	final = append(final, tmp1...)
	final = append(final, tmp2...)
	final = append(final, tmp3...)
	return final
}

type ListBucketResult struct {
	XMLName  xml.Name  `xml:"ListBucketResult"`
	Contents []Content `xml:"Contents"`
}

type Content struct {
	Key string `xml:"Key"`
}

func countKeys(response *http.Response) (int, error) {
	var result ListBucketResult
	decoder := xml.NewDecoder(response.Body)

	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return 0, err
		}

		switch t := token.(type) {
		case xml.StartElement:
			if t.Name.Local == "Contents" {
				var content Content
				if err := decoder.DecodeElement(&content, &t); err != nil {
					return 0, err
				}
				result.Contents = append(result.Contents, content)
			}
		}
	}

	return len(result.Contents), nil
}
func checkAWS(bucket string) {
	url := fmt.Sprintf("https://%s.s3.amazonaws.com", bucket)
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	lineCount, err := countKeys(resp)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	if resp.StatusCode == 200 {
		fmt.Printf("%s%s : %d  [%d] \n%s", Green, bucket, resp.StatusCode, lineCount, Reset)
	}
	if resp.StatusCode == 403 {
		fmt.Printf("%s%s : %d\n%s", Red, bucket, resp.StatusCode, Reset)
	}
	if resp.StatusCode == 400 {
		fmt.Printf("%s%s : %d\n%s", Yellow, bucket, resp.StatusCode, Reset)
	}
}

func checkGCP(bucket string) {

	url := fmt.Sprintf("http://storage.googleapis.com/%s", bucket)
	resp, err := http.Get(url)
	if err != nil {

		return
	}
	defer resp.Body.Close()

	lineCount, err := countKeys(resp)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	if resp.StatusCode == 200 {
		fmt.Printf("%s%s : %d  [%d] \n%s", Green, bucket, resp.StatusCode, lineCount, Reset)
	}
	if resp.StatusCode == 403 {
		fmt.Printf("%s%s : %d\n%s", Red, bucket, resp.StatusCode, Reset)
	}
}

func checkAzure(bucket string) {

	url := fmt.Sprintf("https://%s.blob.core.windows.net/", bucket)
	resp, err := http.Get(url)
	if err != nil {

		return
	}
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)

	lineCount := 0

	for scanner.Scan() {
		lineCount++
	}
	if resp.StatusCode == 400 {
		fmt.Printf("%s%s : %d  : Valid Azure Storage Blob \n%s", Green, bucket, resp.StatusCode, Reset)
	}
}

func initSem(bucket string, wg *sync.WaitGroup, semaphore chan struct{}, scan_type string) {
	defer wg.Done()

	semaphore <- struct{}{}
	defer func() {

		<-semaphore
	}()
	if scan_type == "aws" {
		checkAWS(bucket)
	}
	if scan_type == "gcp" {
		checkGCP(bucket)
	}
	if scan_type == "azure" {
		checkAzure(bucket)
	}

}
func main() {
	var maxConcurrent = 1000
	startTime := time.Now()
	banner()
	buckets := make([]string, 0)
	var scanGCP, scanAzure, scanAll, enum_org bool
	var target string
	flag.StringVar(&target, "t", "", "Target to search Buckets for")
	flag.BoolVar(&scanGCP, "gcp", false, "Scan for GCP")
	flag.BoolVar(&scanAzure, "azure", false, "Scan for Azure")
	flag.BoolVar(&scanAll, "all", false, "Scan for all services")
	flag.BoolVar(&enum_org, "enumerate", false, "Search for organization names, then look for buckets.")
	flag.Parse()

	if target == "" {
		fmt.Println("Target is required")
		flag.Usage()
		os.Exit(1)
	}
	scanAWS := true
	if scanGCP || scanAzure || scanAll {
		scanAWS = false
	}
	if (scanGCP || scanAzure) && scanAll {
		fmt.Println("Too many arguments provided, use either --all, --azure, or --gcp")
		os.Exit(1)
	}
	if scanAll {
		scanAWS = true
		scanGCP = true
		scanAzure = true
	}
	if enum_org {
		fmt.Printf("[-] %sEnumerate flag Enabled\n%s", Yellow, Reset)
		fmt.Printf("[-] %sEnabling the Enumerate flag may increase false positives and prolong scan durations.\n\n%s", Yellow, Reset)

	}
	if enum_org {
		maxConcurrent = 4000
		orgs, err := orgdata.GetOrgs(target)
		fmt.Printf("[+] %sFound %d organizations%s\n", Gray, len(orgs), Reset)

		if err != nil {
			fmt.Printf("[-] %sAn Error Occured while Fetching Organizations\n%s : %s", Yellow, Reset, err)
			os.Exit(0)
			return
		}

		for _, domain := range orgs {
			fmt.Println(domain)

			tmp_buckets := createWordlist(domain)
			buckets = append(buckets, tmp_buckets...)
		}

		fmt.Printf("[+] %sCreated Wordlist of %d buckets%s\n", Gray, len(buckets), Reset)
	} else {
		buckets = createWordlist(target)

		fmt.Printf("[+] %sCreated Wordlist of %d buckets%s\n", Gray, len(buckets), Reset)
	}

	if scanAWS {
		fmt.Printf("\n[+] Checking AWS S3 Buckets for '%s'\n\n", target)

		var wg sync.WaitGroup

		semaphore := make(chan struct{}, maxConcurrent)

		for _, bucket := range buckets {
			wg.Add(1)
			go initSem(bucket, &wg, semaphore, "aws")
		}

		wg.Wait()
	}

	if scanGCP {
		fmt.Printf("\n[+] Checking GCP Buckets for '%s'\n\n", target)

		var wg sync.WaitGroup

		semaphore := make(chan struct{}, maxConcurrent)

		for _, bucket := range buckets {
			wg.Add(1)
			go initSem(bucket, &wg, semaphore, "gcp")
		}

		wg.Wait()
	}
	if scanAzure {
		fmt.Printf("\n[+] Checking Azure Storage for '%s'\n\n", target)

		var wg sync.WaitGroup

		semaphore := make(chan struct{}, maxConcurrent)

		for _, bucket := range buckets {
			wg.Add(1)
			go initSem(bucket, &wg, semaphore, "azure")
		}

		wg.Wait()
	}

	elapsedTime := time.Since(startTime)
	fmt.Printf("\n\nTotal time taken: %s\n", elapsedTime)
}
