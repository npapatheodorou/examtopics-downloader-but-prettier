# Exam Topics Downloader

This repo aims to make it possible for you to obtain all the exam questions from the examtopics website (which is paywalled).

## Setting it Up

### Running with Go

1. First, you must install [Golang >= 1.24](https://go.dev/doc/install) from the offical website.
2. Then, run `git clone https://github.com/thatonecodes/examtopics-downloader` in your terminal to clone the repo.
3. `cd` into the directory: `cd examtopics-downloader`
4. You can now run: `go run ./cmd/main.go`

(there will be compiled binaries in the future)

## Interactive Flow

There are no required startup parameters.
Optional:
- `-debug` enables verbose debug/error logs during fetch and extraction.

When the app starts it will:
1. scrape and display all available providers
2. ask you to select a provider
3. scrape and display all available exams for that provider
4. ask you to select an exam
5. run extraction with a single progress bar for the whole extraction workflow

Output is always HTML and defaults to `provider_exam-code.html` (for example, `cisco_200-301.html`).

## [For outputted file examples, see the examples folder](examples/google_devops.md)

## Demo

So, you have installed `go` on your system, and you're inside of the working directory. Let's say you would like the questions for the cisco exam 200-301.

Open your terminal and run:

```bash
go run ./cmd/main.go
```

Then choose `Cisco` from the provider list and `200-301` from the exam list when prompted.

After waiting a few moments, you would see the output end with:

```bash
Successfully saved output: {OUTPUT_LOCATION}.
```

If so, hooray, you have successfully saved all/most of the questions in an `.html` file!
The format would be such as (older, only scraping format):

```
----------------------------------------

## Exam 200-301 topic 1 question 532 discussion

Actual exam question from

Cisco's
200-301

Question #: 532
Topic #: 1

[All 200-301 Questions]

Refer to the exhibit. An engineer configured NAT translations and has verified that the configuration is correct. Which IP address is the source IP after the NAT has taken place?
Suggested Answer: D 

A. 10.4.4.4

B. 10.4.4.5

C. 172.23.103.10

D. 172.23.104.4

**Answer: D**

**Timestamp: Jan. 5, 2021, 9:48 p.m.**

[View on ExamTopics](https://www.examtopics.com/discussions/cisco/view/41599-exam-200-301-topic-1-question-532-discussion/)

----------------------------------------
```
