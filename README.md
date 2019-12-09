# o365verify
[![](https://godoc.org/github.com/jakewarren/o365verify?status.svg)](http://godoc.org/github.com/jakewarren/o365verify)
[![CircleCI](https://circleci.com/gh/jakewarren/o365verify.svg?style=shield)](https://circleci.com/gh/jakewarren/o365verify)
[![GitHub release](http://img.shields.io/github/release/jakewarren/o365verify.svg?style=flat-square)](https://github.com/jakewarren/o365verify/releases])
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](https://github.com/jakewarren/o365verify/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/jakewarren/o365verify)](https://goreportcard.com/report/github.com/jakewarren/o365verify)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=shields)](http://makeapullrequest.com)

This program uses the autodiscover JSON API of Office 365 to enumerate valid email addresses.

This is a Golang port of [Raikia/UhOh365](https://github.com/Raikia/UhOh365) with a few added features.

## Install
### Option 1: Binary

Download the latest release from [https://github.com/jakewarren/o365verify/releases/latest](https://github.com/jakewarren/o365verify/releases/latest)

### Option 2: From source

```
go get github.com/jakewarren/o365verify/...
```

## Usage

```
❯ o365verify --help
Usage: o365verify [flags] <email address...>

Flags:
  -h, --help          display help
  -t, --threads int   number of threads to run with (default 10)
  -V, --version       display version information
```

### Example
```
❯ o365verify test@example.com
[
  {
    "EmailAddress": "test@example.com",
    "CalculatedBETarget": "DM6PR06MB5690.NAMPRD06.PROD.OUTLOOK.COM",
    "MailboxGUID": "00037ffe-80ef-9bbe-0000-000000000000",
    "ValidAddress": true,
    "DomainIsO365": false
  }
]
```

## Acknowledgments 

[Raikia/UhOh365](https://github.com/Raikia/UhOh365) - for the discovery and initial work

## Changes

All notable changes to this project will be documented in the [changelog].

The format is based on [Keep a Changelog](http://keepachangelog.com/) and this project adheres to [Semantic Versioning](http://semver.org/).

## License

MIT © 2019 Jake Warren

[changelog]: https://github.com/jakewarren/o365verify/blob/master/CHANGELOG.md
