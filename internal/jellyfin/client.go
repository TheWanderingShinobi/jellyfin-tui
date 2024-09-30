package jellyfin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

type MediaItem struct {
	ID              string  `json:"Id"`
	Name            string  `json:"Name"`
	Type            string  `json:"Type"`
	Path            string  `json:"Path"`
	Overview        string  `json:"Overview"`
	CommunityRating float64 `json:"CommunityRating"`
}

type Playlist struct {
	ID   string `json:"Id"`
	Name string `json:"Name"`
}

type User struct {
	ID   string `json:"Id"`
	Name string `json:"Name"`
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
}

func (c *Client) Login(username, password string) error {
	data := url.Values{}
	data.Set("Username", username)
	data.Set("Pw", password)

	resp, err := c.HTTPClient.PostForm(c.BaseURL+"/Users/AuthenticateByName", data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("authentication failed: %s", resp.Status)
	}

	var result struct {
		AccessToken string `json:"AccessToken"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	c.Token = result.AccessToken
	return nil
}

func (c *Client) GetMediaItems(page, itemsPerPage int, filter string) ([]MediaItem, int, error) {
	req, err := http.NewRequest("GET", c.BaseURL+"/Items", nil)
	if err != nil {
		return nil, 0, err
	}

	q := req.URL.Query()
	q.Add("StartIndex", fmt.Sprintf("%d", (page-1)*itemsPerPage))
	q.Add("Limit", fmt.Sprintf("%d", itemsPerPage))
	if filter != "" {
		q.Add("IncludeItemTypes", filter)
	}
	req.URL.RawQuery = q.Encode()

	req.Header.Set("X-Emby-Authorization", fmt.Sprintf("MediaBrowser Token=\"%s\"", c.Token))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	var result struct {
		Items            []MediaItem `json:"Items"`
		TotalRecordCount int         `json:"TotalRecordCount"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, 0, err
	}

	return result.Items, result.TotalRecordCount, nil
}

func (c *Client) GetItemDetails(itemID string) (*MediaItem, error) {
	req, err := http.NewRequest("GET", c.BaseURL+"/Items/"+itemID, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Emby-Authorization", fmt.Sprintf("MediaBrowser Token=\"%s\"", c.Token))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var item MediaItem
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		return nil, err
	}

	return &item, nil
}

func (c *Client) Search(query string) ([]MediaItem, error) {
	req, err := http.NewRequest("GET", c.BaseURL+"/Items", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("SearchTerm", query)
	req.URL.RawQuery = q.Encode()

	req.Header.Set("X-Emby-Authorization", fmt.Sprintf("MediaBrowser Token=\"%s\"", c.Token))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Items []MediaItem `json:"Items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Items, nil
}

func (c *Client) GetPlaylists() ([]Playlist, error) {
	req, err := http.NewRequest("GET", c.BaseURL+"/Playlists", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Emby-Authorization", fmt.Sprintf("MediaBrowser Token=\"%s\"", c.Token))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Items []Playlist `json:"Items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Items, nil
}

func (c *Client) AddToPlaylist(playlistID, itemID string) error {
	url := fmt.Sprintf("%s/Playlists/%s/Items?Ids=%s", c.BaseURL, playlistID, itemID)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-Emby-Authorization", fmt.Sprintf("MediaBrowser Token=\"%s\"", c.Token))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to add item to playlist: %s", resp.Status)
	}

	return nil
}

func (c *Client) GetStreamURL(itemID string) string {
	return fmt.Sprintf("%s/Videos/%s/stream?static=true&api_key=%s", c.BaseURL, itemID, c.Token)
}

func (c *Client) CreatePlaylist(name string) error {
	data := url.Values{}
	data.Set("Name", name)
	data.Set("MediaType", "Video")

	req, err := http.NewRequest("POST", c.BaseURL+"/Playlists", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Emby-Authorization", fmt.Sprintf("MediaBrowser Token=\"%s\"", c.Token))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create playlist: %s", resp.Status)
	}

	return nil
}

func (c *Client) GetUsers() ([]User, error) {
	req, err := http.NewRequest("GET", c.BaseURL+"/Users", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Emby-Authorization", fmt.Sprintf("MediaBrowser Token=\"%s\"", c.Token))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var users []User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, err
	}

	return users, nil
}

func (c *Client) SwitchUser(userID string) error {
	// This is a placeholder. The actual implementation would depend on how
	// Jellyfin handles user switching. You might need to re-authenticate
	// as the new user, or there might be a specific API endpoint for switching users.
	return fmt.Errorf("SwitchUser not implemented")
}
