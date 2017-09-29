package cmd

import (
	"fmt"
	"log"

	"github.com/jckimble/releasetool/update"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Check for update",
	Run:   startUpdate,
}

func init() {
	RootCmd.AddCommand(updateCmd)
}

func startUpdate(cmd *cobra.Command, args []string) {
	var version string
	fmt.Sscanf(Version, "v%s", &version)
	u := &update.GitHub{
		CurrentVersion: version,
		GithubOwner:    "jckimble",
		GithubRepo:     "lighttower",
		Certificate: `-----BEGIN CERTIFICATE-----
MIIFcDCCA1igAwIBAgIJAOMTW/DwB0O9MA0GCSqGSIb3DQEBCwUAME0xCzAJBgNV
BAYTAlVTMRQwEgYDVQQIDAtNaXNzaXNzaXBwaTETMBEGA1UECgwKTGlnaHRUb3dl
cjETMBEGA1UEAwwKTGlnaHRUb3dlcjAeFw0xNzA5MTcwOTU3MTBaFw0xOTA5MTcw
OTU3MTBaME0xCzAJBgNVBAYTAlVTMRQwEgYDVQQIDAtNaXNzaXNzaXBwaTETMBEG
A1UECgwKTGlnaHRUb3dlcjETMBEGA1UEAwwKTGlnaHRUb3dlcjCCAiIwDQYJKoZI
hvcNAQEBBQADggIPADCCAgoCggIBAMSJZXQUolPc7xJDcBuWKJYd1NyGCZUbRA8/
zxCk5wh80/+vsFfxLNMbykeGSTOCcoOWieavMPRzpzD6Vk3euZH1wcg8ZQE4Em6P
e4eXQ7zgrsqdBjI9M29Z6gWnHiyT69iGSB/85KuuLIiZeoWMLpzhIvt7y57rviRj
AFx4giV1kTuU0kwv0MWDBNiWNRTk2pUUxdBvjUNb9FW5qpi6d9CtoIQBVxMAVwwS
6Ky05dA++sKWVJYOFGb7sJ+EkdI7wXpeDYel8ucK/x3uTrgrM5NYDeeE6rCGaNHU
Q/5qgnfpinQAUbozxf1l5BbTASX9AsWFZJ2A+5KKCdw0DUf7gw7+NXOIfyJDeVeX
PgDSwIFhTngqfLHVl6FX6J/7TpoC/NE7QGWF99SlV61rH5XkPSkR0I9Ec+ktBxIi
HC8QxSC7GW0RU+zUsHsbaKXgK/U2UBskcRniJxd3mymKNx5XQ3uqHoNncb6ad6fS
QOYBUdCXk8nNNu6A6BWHTXroxj1l8HtNQqXBMoMJyL3loW7VOh1Jkcw5etSath35
m04kCvMUMZiJN4LvV8Y6Wpq+Y8eZZeS0xXtia5HfSLo5zEFZkMGNmCM8cVawVoHj
DR9D0R6q1xtxnAnO9Fyl33wQF56daf3tFkqoZiS9CbHzTsOc4DdWqfh8lBXCeyQe
DdGnQrcPAgMBAAGjUzBRMB0GA1UdDgQWBBSgKT0W5IwVd+phq2uvL1Pn3bTOnjAf
BgNVHSMEGDAWgBSgKT0W5IwVd+phq2uvL1Pn3bTOnjAPBgNVHRMBAf8EBTADAQH/
MA0GCSqGSIb3DQEBCwUAA4ICAQCLcjiruok5v5Qol5TtPqFJuinp/LlxAoi+Muuh
pR7LySb6c60iKfNLU/lHmAd8gLPO4TEVH7UfPwJxlT3sMilg7FuYm1Vnk6J5KZ9k
WJ2sRPOl2EX++/DGv4h7fdSPOPtSwbH+ZPutBvUtrhIPjfAOOzDsrLXk0mVIzUmG
aBommAPtpDnAINqWs0/O+ZJHlYHXmSrd86XsWeazvcQmk/ZKlhSdQ4Hqp7w4khS3
x13M5LuWmSmhvn+Toqr8M5QCFVeMJuXljUUhqQeyL8GwO9w4jM1d+a8bvt6znIf0
KQaGbxAP0ag6QKx+y1GFOed7v194EjOXt5fRpiksre4MElOt/l10aH9/3Syr1qLb
wh6N9cD1USZQfEdEQD7pSyv/vqelfs6BIxSyg5yS1snihKvtt7Y9njgn1KdYnbTW
p+6RBRMg1zQ+HDCBCKuKLqFYqMOemlYhBvgvUJnwRQIPWdzYzBH3xZJxPNT3GZ5W
/9tyOF1cDcC5SlxAby53QVX7ScdRymCoZpo3D0jRpGfStk9HtQPGkI0jOPCs8eCe
9Ffsk140t0nQloa2qaDAt6NMaRaTjU4M6QVu4MdnA7Z8F3W161yWPpG/AGCLLY2H
MgVN36DSDsVbECniqiRqHYV6B5I1OY6eoQL4K5Ny2z3mDiJ79+T1A8NTfZuktfLt
8391DQ==
-----END CERTIFICATE-----`,
	}
	available, err := u.CheckUpdateAvailable()
	if err != nil {
		log.Printf("Unable to check Update: %s\n", err)
	}

	if available != "" {
		log.Printf("Version %s available\n", available)
		err := update.Update(u)
		if err != nil {
			log.Printf("Unable to Update: %s\n", err)
		}
	} else {
		log.Println("LightTower is current Version")
	}
}
