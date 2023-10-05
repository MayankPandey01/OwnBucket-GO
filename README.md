<p align="center"><img width="750" height="350" src="https://user-images.githubusercontent.com/29165227/210492763-1f3e82ba-9a77-470a-be57-e02e12017335.jpg"></p>

<p align="center">
<a href="https://go.dev/"><img src="https://img.shields.io/badge/Made%20with-Go-1f425f.svg"></a>
<a href="https://github.com/MayankPandey01/Jira-Lens/issues"><img src="https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square"></a>
<a href="https://twitter.com/mayank_pandey01"><img src="https://img.shields.io/twitter/follow/mayank_pandey01?style=social"></a>
<a href="https://github.com/ellerbrock/open-source-badges/"><img src="https://badges.frapsoft.com/os/v1/open-source.svg?v=103"></a>  
</p>

# 🤔 What's OwnBucket?

OwnBucket is a Fast GO-based Recon tool for Storage Buckets. It scans for AWS S3 Bucket, Azure Storage Blob, and GCP Buckets by brute-forcing using different permutations. It also enumerates different organization names based on the target name using certificate names from  [Crt.sh](https://crt.sh )


# 🚀 Usage
OwnBucket can be easily  used from the command line 
- `go run main.go -t {COMPANY}`

![Screenshot 2023-01-04 150553](https://github.com/faressoft/terminalizer/assets/29165227/399b6487-006c-4815-b483-5160dbd079e5)

 Additional Arguments can be passed to use the tool in different ways:
 
 - `-t` : To Provide a Company Name for Scanning
 - `--aws` : Only Check for AWS S3 Buckets (Default)
 - `--gcp` : Only Check for GCP Buckets 
 - `--azure` : Only Check for Azure Storage Blob 
 - `--all` : Check for both AWS S3 and GCP Buckets
 - `--enumerate` : This Flag will enumerate Different Organizations Names Based on the Target Name


- `go run main.go -t {COMPANY} --enumerate ` 

# 🔧Installation

## 🔨 Using Git
- ` git clone https://github.com/MayankPandey01/OwnBucket-GO`
- `go run main.go -h` 

## 🧪 Recommended GO Version:
- This tool was made on GO 1.18.1 (go1.18.1 linux/amd64)

## ⛳ Dependencies:

This Tool Uses an External Package [DomianParser](https://github.com/Cgboal/DomainParser) by [CGlobal](https://github.com/Cgboal/)


## 🐞 Bug Bounties

This tool is focused mainly on `Bug Bounty Hunters` and `Security Professionals`. You Can Use OwnBucket to Scan For Different Storage Buckets of the Target Company. 


## 🎯 Contribution ![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)
We Love to Get Contribution from the Open Source Community💙. You are Welcome to Provide your Important Suggestions to make this tool more Awesome. Open a PR  and we will See to it ASAP.

**Ways to contribute**
- Suggest a feature
- Report a bug
- Fix something and open a pull request
- Spread the word

## 📚 DISCLAIMER

This project is a [personal development](https://en.wikipedia.org/wiki/Personal_development). Please respect its philosophy and don't use it for evil purposes. By using OwnBucket, you agree to the MIT license included in the repository. For more details at [The MIT License &mdash; OpenSource](https://opensource.org/licenses/MIT).

Happy Hacking ✨✨

> This Tool is Highly Motivated by [LazyS3](https://github.com/nahamsec/lazys3)

## 📃 Licensing

This project is licensed under the MIT license.
