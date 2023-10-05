package orgdata

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	parser "github.com/Cgboal/DomainParser"
)

type Certificate struct {
	CommonName string `json:"common_name"`
}

func GetOrgs(org string) ([]string, error) {
	final := make([]string, 0)
	url := "https://crt.sh/?o=" + org + "&output=json"
	commonNames, err := fetchAndFilterCommonNames(url, org)
	if err != nil {
		return nil, err
	}

	parser := parser.NewDomainParser()
	for _, commonName := range commonNames {
		domain := parser.GetDomain(commonName)
		final = deDupe(final, domain, org)

	}
	return final, nil
}

func fetchAndFilterCommonNames(url string, org string) ([]string, error) {
	resp, err := http.Get(url)
	if resp.StatusCode == 503 {
		fmt.Printf("\n[-] An Error Occured while Fetching Data from Crt.sh : %d\n", resp.StatusCode)
		os.Exit(0)
	}
	if err != nil {
		fmt.Printf("In2")
		return nil, err
	}
	defer resp.Body.Close()

	var certificates []Certificate
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&certificates); err != nil {
		return nil, err
	}

	filteredCommonNames := make([]string, 0)
	for _, cert := range certificates {
		parts := strings.Split(cert.CommonName, ".")
		for _, part := range parts {
			if strings.Contains(strings.ToLower(part), org) {
				filteredCommonNames = append(filteredCommonNames, cert.CommonName)
				break
			}
		}
	}

	return filteredCommonNames, nil
}

func deDupe(secondLevelDomains []string, data string, org string) []string {
	exists := false

	for _, existingData := range secondLevelDomains {
		if strings.EqualFold(existingData, data) {
			exists = true
			break
		}
	}

	if !exists {
		if strings.Contains(strings.ToLower(data), org) {
			if !strings.Contains(strings.ToLower(data), " ") {
				secondLevelDomains = append(secondLevelDomains, data)
			}
		}
	}

	return secondLevelDomains
}
