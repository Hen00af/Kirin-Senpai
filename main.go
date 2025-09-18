package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "time"
)

type Contest struct {
    ID        string    `json:"id"`
    StartTime time.Time `json:"start"`
    Title     string    `json:"title"`
    URL       string    `json:"url"`
}

const savedFile = "saved_contests.json"
const apiURL = "https://api.example.com/atcoder/upcoming"  // glokta1 の API 等を使う

func fetchUpcoming() ([]Contest, error) {
    resp, err := http.Get(apiURL)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("bad status: %d", resp.StatusCode)
    }

    var data struct {
        Contests []Contest `json:"contests"`
    }
    body, _ := ioutil.ReadAll(resp.Body)
    err = json.Unmarshal(body, &data)
    if err != nil {
        return nil, err
    }
    return data.Contests, nil
}

func loadSaved() (map[string]Contest, error) {
    m := map[string]Contest{}
    if _, err := os.Stat(savedFile); os.IsNotExist(err) {
        return m, nil
    }
    b, err := ioutil.ReadFile(savedFile)
    if err != nil {
        return nil, err
    }
    var arr []Contest
    err = json.Unmarshal(b, &arr)
    if err != nil {
        return nil, err
    }
    for _, c := range arr {
        m[c.ID] = c
    }
    return m, nil
}

func save(contests []Contest) error {
    b, err := json.MarshalIndent(contests, "", "  ")
    if err != nil {
        return err
    }
    return ioutil.WriteFile(savedFile, b, 0644)
}

func sendDiscordWebhook(contest Contest) error {
    webhookUrl := os.Getenv("DISCORD_WEBHOOK")
    if webhookUrl == "" {
        return fmt.Errorf("missing webhook URL")
    }
    content := fmt.Sprintf("🕒 新しいコンテスト: **%s** が %s に始まります！\n%s", contest.Title, contest.StartTime.Format(time.RFC3339), contest.URL)
    payload := map[string]string{"content": content}
    pd, err := json.Marshal(payload)
    if err != nil {
        return err
    }
    resp, err := http.Post(webhookUrl, "application/json", bytes.NewBuffer(pd))
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    if resp.StatusCode != 204 && resp.StatusCode != 200 {
        return fmt.Errorf("webhook error status: %d", resp.StatusCode)
    }
    return nil
}

func main() {
    upcoming, err := fetchUpcoming()
    if err != nil {
        fmt.Println("fetch error:", err)
        return
    }
    saved, err := loadSaved()
    if err != nil {
        fmt.Println("load error:", err)
        return
    }

    var toNotify []Contest
    for _, c := range upcoming {
        if _, exists := saved[c.ID]; !exists {
            toNotify = append(toNotify, c)
        }
    }

    for _, c := range toNotify {
        err := sendDiscordWebhook(c)
        if err != nil {
            fmt.Println("notify error:", err)
        }
    }

    if len(toNotify) > 0 {
        err := save(upcoming)
        if err != nil {
            fmt.Println("save error:", err)
        }
    }
}
