package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/google/go-github/v62/github"
	"golang.org/x/oauth2"
)

var format = flag.String("format", "text", "Format (json|text)")
var authToken = flag.String("auth_token", "", "Auth token for better rate limits")

var urlRE = regexp.MustCompile(`\Ahttps://github.com/([^/]+)/([^/]+)/issues/(\d+)\z`)

func main() {
	flag.Parse()
	ctx := context.Background()
	args := flag.Args()

	if len(args) == 1 && urlRE.MatchString(args[0]) {
		match := urlRE.FindStringSubmatch(args[0])
		args = match[1:]
	} else if len(args) < 3 {
		log.Fatalf("Usage: %s [<url>|<owner> <repo> <issue_number>]", os.Args[0])
	}

	owner := args[0]
	repo := args[1]
	issueNumber, err := strconv.Atoi(args[2])
	if err != nil {
		log.Fatalf("Invalid issue number: %s", err)
	}

	if *format != "json" && *format != "text" {
		log.Fatalf("unsupported format. Value values are json and text")
	}

	var hc *http.Client
	if *authToken != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: *authToken},
		)
		hc = oauth2.NewClient(ctx, ts)
	}

	client := github.NewClient(hc)

	issue, _, err := client.Issues.Get(ctx, owner, repo, issueNumber)
	if err != nil {
		log.Fatalf("Failed to fetch issue %d for %s/%s: %s", issueNumber, owner, repo, err)
	}

	if *format == "text" {
		fmt.Printf("[%s] %s:\n%s\n================================================================================\n", issue.CreatedAt, *issue.User.Login, *issue.Body)
	}

	var commentOps github.IssueListCommentsOptions
	var allComments []github.IssueComment
	for {
		comments, commentsResp, err := client.Issues.ListComments(ctx, owner, repo, issueNumber, &commentOps)
		if err != nil {
			log.Fatalf("Failed to fetch comments for issue %d in %s/%s: %s", issueNumber, owner, repo, err)
		}
		for _, comment := range comments {
			if *format == "text" {
				fmt.Printf("[%s] %s:\n%s\n================================================================================\n", comment.CreatedAt, *comment.User.Login, *comment.Body)
			} else {
				allComments = append(allComments, *comment)
			}
		}
		if commentsResp.NextPage == 0 {
			break
		}
		commentOps.Page = commentsResp.NextPage
	}

	if *format == "json" {
		result := Issue{
			Issue:    *issue,
			Comments: allComments,
		}

		enc := json.NewEncoder(os.Stdout)
		err = enc.Encode(result)
		if err != nil {
			log.Fatalf("Failed to encode JSON: %s", err)
		}
	}
}

type Issue struct {
	Issue    github.Issue          `json:"issue"`
	Comments []github.IssueComment `json:"comments"`
}
