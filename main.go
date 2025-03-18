package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

type Member struct {
	ID          string `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Avatar      string `json:"avatar"`
	IsBot       bool   `json:"is_bot"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
}

type User struct {
	ID     string `json:"id"`
	Member Member `json:"member"`
}

type UsersData struct {
	Users []User `json:"users"`
}

type Issue struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	StateID         string   `json:"state_id"`
	SortOrder       int      `json:"sort_order"`
	CompletedAt     *string  `json:"completed_at"`
	EstimatePoint   int      `json:"estimate_point"`
	Priority        string   `json:"priority"`
	StartDate       *string  `json:"start_date"`
	TargetDate      *string  `json:"target_date"`
	SequenceID      int      `json:"sequence_id"`
	ProjectID       string   `json:"project_id"`
	ParentID        *string  `json:"parent_id"`
	CreatedAt       string   `json:"created_at"`
	UpdatedAt       string   `json:"updated_at"`
	CreatedBy       string   `json:"created_by"`
	UpdatedBy       string   `json:"updated_by"`
	IsDraft         bool     `json:"is_draft"`
	ArchivedAt      *string  `json:"archived_at"`
	CycleID         string   `json:"cycle_id"`
	LinkCount       int      `json:"link_count"`
	AttachmentCount int      `json:"attachment_count"`
	SubIssuesCount  int      `json:"sub_issues_count"`
	LabelIDs        []string `json:"label_ids"`
	AssigneeIDs     []string `json:"assignee_ids"`
	ModuleIDs       []string `json:"module_ids"`
	Assignees       Member   `json:"assignees"`
}

type IssuesData struct {
	Issues Issue `json:"issues"`
}

func main() {
	// Load users data
	userFile, err := os.Open("member.json")
	if err != nil {
		fmt.Println("Error opening users file:", err)
		return
	}
	defer userFile.Close()

	userBytes, err := ioutil.ReadAll(userFile)
	if err != nil {
		fmt.Println("Error reading users file:", err)
		return
	}

	var usersData []User // FIX: JSON is an array, so we use a slice
	err = json.Unmarshal(userBytes, &usersData)
	if err != nil {
		fmt.Println("Error unmarshalling users JSON:", err)
		return
	}

	// Create a map of users by their Member ID
	userMap := make(map[string]Member)
	for _, user := range usersData {
		userMap[user.Member.ID] = user.Member
	}

	// Load issues data
	issueFile, err := os.Open("issues.json")
	if err != nil {
		fmt.Println("Error opening issues file:", err)
		return
	}
	defer issueFile.Close()

	issueBytes, err := ioutil.ReadAll(issueFile)
	if err != nil {
		fmt.Println("Error reading issues file:", err)
		return
	}

	var issuesData []Issue
	err = json.Unmarshal(issueBytes, &issuesData)
	if err != nil {
		fmt.Println("Error unmarshalling issues JSON:", err)
		return
	}

	// Assign members to issues based on AssigneeIDs
	for i, issue := range issuesData {
		for _, assigneeID := range issue.AssigneeIDs {
			if member, exists := userMap[assigneeID]; exists {
				issuesData[i].Assignees = member
				break
			}
		}
	}

	sort.SliceStable(issuesData, func(i, j int) bool {
		return issuesData[i].CreatedBy < issuesData[j].CreatedBy
	})

	// 	for _, issue := range issuesData {
	// 		fmt.Printf("%s by %s\n", issue.Name, issue.Assignees.DisplayName)
	// 	}

	file, err := os.Create("tasks.csv")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, issue := range issuesData {
		title := fmt.Sprintf("%s by %s", strings.ReplaceAll(issue.Name, ",", " &"), issue.Assignees.DisplayName)
		record := []string{"Task", title, "", "", "", ""}
		writer.Write(record)
	}

	fmt.Println("CSV file created successfully.")
}
