// github.com/rudSarkar
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/fatih/color"
)

type TopHolder struct {
	Owner string  `json:"owner"`
	Pct   float64 `json:"pct"`
}

type TokenData struct {
	Mint       string      `json:"mint"`
	Token      TokenInfo   `json:"token"`
	TokenMeta  TokenMeta   `json:"tokenMeta"`
	TopHolders []TopHolder `json:"topHolders"`
	Risks      []Risk      `json:"risks"`
	FileMeta   FileMeta    `json:"fileMeta"`
	Rugged     bool        `json:"rugged"`
	Markets    []Market    `json:"markets"`
}

type TokenInfo struct {
	MintAuthority   string  `json:"mintAuthority"`
	Supply          float64 `json:"supply"` // Change type to float64
	Decimals        uint8   `json:"decimals"`
	FreezeAuthority string  `json:"freezeAuthority"`
}

type TokenMeta struct {
	Name            string `json:"name"`
	Symbol          string `json:"symbol"`
	Mutable         bool   `json:"mutable"`
	UpdateAuthority string `json:"updateAuthority"`
}

type Risk struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Score       int    `json:"score"`
	Level       string `json:"level"`
}

type FileMeta struct {
	Image string `json:"image"`
}

type Market struct {
	MarketType string `json:"marketType"`
	Lp         LP     `json:"lp"`
}

type LP struct {
	LpLockedPct float64 `json:"lpLockedPct"`
}

func rugcheck(token string) (TokenData, error) {
	url := fmt.Sprintf("https://api.rugcheck.xyz/v1/tokens/%s/report", token)

	resp, err := http.Get(url)
	if err != nil {
		return TokenData{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return TokenData{}, fmt.Errorf("HTTP status error: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return TokenData{}, err
	}

	var data TokenData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return TokenData{}, err
	}

	return data, nil
}

func main() {
	var token string
	flag.StringVar(&token, "token", "", "Token contract address")
	flag.Parse()

	if token == "" {
		fmt.Println("Please provide the token contract address using the -token flag")
		return
	}

	data, err := rugcheck(token)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Access data fields from 'data' struct here
	fmt.Println("Token Name:", data.TokenMeta.Name)
	fmt.Println("Token Symbol:", data.TokenMeta.Symbol)
	fmt.Println("Token Supply:", data.Token.Supply)
	fmt.Println("Rug Risk Name:", data.Risks[0].Name)

	var ScoreAnalysis string

	switch {
	case data.Risks[0].Score < 1000:
		ScoreAnalysis = color.GreenString("Good")
	case data.Risks[0].Score < 5000:
		ScoreAnalysis = color.YellowString("Warning")
	default:
		ScoreAnalysis = color.RedString("Danger")
	}

	fmt.Println("Rug Risk Score:", ScoreAnalysis)

	fmt.Println("Top Holder 1:", data.TopHolders[0].Owner)
	fmt.Println("LP Locked Percentage:", data.Markets[0].Lp.LpLockedPct)
}
