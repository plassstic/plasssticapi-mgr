package api

type ArtistSI struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Type         string            `json:"type"`
	URI          string            `json:"uri"`
	ExternalURLs map[string]string `json:"externalUrls"`
}

type ImageSI struct {
	Height int    `json:"height"`
	Width  int    `json:"width"`
	URL    string `json:"url"`
}

type ContextSI struct {
	Type         string            `json:"type"`
	URI          string            `json:"uri"`
	ExternalURLs map[string]string `json:"externalUrls"`
}

type AlbumSI struct {
	ID                   string            `json:"id"`
	Name                 string            `json:"name"`
	Type                 string            `json:"type"`
	URI                  string            `json:"uri"`
	ExternalURLs         map[string]string `json:"externalUrls"`
	ReleaseDate          string            `json:"releaseDate"`
	Artists              []ArtistSI        `json:"artists"`
	ReleaseDatePrecision string            `json:"releaseDatePrecision"`
	TotalTracks          int               `json:"totalTracks"`
}

type DeviceSI struct {
	ID            string `json:"id"`
	Type          string `json:"type"`
	IsActive      bool   `json:"isActive"`
	Name          string `json:"name"`
	VolumePercent int    `json:"volumePercent"`
}

type ItemSI struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Type         string            `json:"type"`
	URI          string            `json:"uri"`
	ExternalURLs map[string]string `json:"externalUrls"`
	Album        AlbumSI           `json:"album"`
	Artists      []ArtistSI        `json:"artists"`
	DurationMs   int               `json:"durationMs"`
	ExternalIDs  map[string]string `json:"externalIds"`
	Explicit     bool              `json:"explicit"`
	IsLocal      bool              `json:"isLocal"`
	TrackNumber  int               `json:"trackNumber"`
}

type PlayerSI struct {
	Device               DeviceSI   `json:"device"`
	IsPlaying            bool       `json:"isPlaying"`
	Timestamp            int        `json:"timestamp"`
	Context              *ContextSI `json:"context,omitempty"`
	ProgressMs           int        `json:"progressMs"`
	Item                 ItemSI     `json:"item"`
	CurrentlyPlayingType string     `json:"currentlyPlayingType"`
}
