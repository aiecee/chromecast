package protobuf

var (
	// Known Payload headers
	Connect   = PayloadHeader{Type: "CONNECT"}
	Close     = PayloadHeader{Type: "CLOSE"}
	GetStatus = PayloadHeader{Type: "GET_STATUS"}
	Pong      = PayloadHeader{Type: "PONG"}       // Response to PING payload
	Launch    = PayloadHeader{Type: "LAUNCH"}     // Launches a new chromecast app
	Play      = PayloadHeader{Type: "PLAY"}       // Plays / unpauses the running app
	Pause     = PayloadHeader{Type: "PAUSE"}      // Pauses the running app
	Seek      = PayloadHeader{Type: "SEEK"}       // Pauses the running app
	Volume    = PayloadHeader{Type: "SET_VOLUME"} // Sets the volume
	Load      = PayloadHeader{Type: "LOAD"}       // Loads an application onto the chromecast
)

type Payload interface {
	SetRequestID(id int)
}

type PayloadHeader struct {
	Type      string `json:"type"`
	RequestID int    `json:"requestId,omitempty"`
}

func (p *PayloadHeader) SetRequestID(id int) {
	p.RequestID = id
}

type MediaHeader struct {
	PayloadHeader
	MediaSessionID int     `json:"mediaSessionId"`
	CurrentTime    float32 `json:"currentTime"`
	ResumeState    string  `json:"resumeState"`
}

type VolumeConfig struct {
	Level float32 `json:"level"`
	Muted bool    `json:"muted"`
}

type ReceiverStatusResponse struct {
	PayloadHeader
	Status struct {
		Applications  []Application `json:"applications,omitempty"`
		IsStandBy     bool          `json:"isStandBy"`
		IsActiveInput bool          `json:"isActiveInput"`
		Volume        VolumeConfig  `json:"volume"`
	} `json:"status"`
}

type Application struct {
	AppID       string `json:"appId"`
	DisplayName string `json:"displayName"`
	SessionID   string `json:"sessionId"`
	StatusText  string `json:"statusText"`
	TransportID string `json:"transportId"`
}

type ReceiverStatusRequest struct {
	PayloadHeader
	Applications []Application `json:"applications"`

	Volume VolumeConfig `json:"volume"`
}

type LaunchRequest struct {
	PayloadHeader
	AppID string `json:"appId"`
}

type LoadMediaCommand struct {
	PayloadHeader
	Media       MediaItem   `json:"media"`
	CurrentTime int         `json:"currentTime"`
	Autoplay    bool        `json:"autoplay"`
	CustomData  interface{} `json:"customData"`
}

type MediaItem struct {
	ContentID   string  `json:"contentId"`
	ContentType string  `json:"contentType"`
	StreamType  string  `json:"streamType"`
	Duration    float32 `json:"duration"`
	Metadata    struct {
		MetadataType int    `json:"metadataType"`
		Title        string `json:"title"`
		SongName     string `json:"songName"`
		Artist       string `json:"artist"`
	} `json:"metadata"`
}

type Media struct {
	MediaSessionID int          `json:"mediaSessionId"`
	PlayerState    string       `json:"playerState"`
	CurrentTime    float32      `json:"currentTime"`
	IdleReason     string       `json:"idleReason"`
	Volume         VolumeConfig `json:"volume"`

	Media MediaItem `json:"media"`
}

type MediaStatusResponse struct {
	PayloadHeader
	Status []Media `json:"status"`
}
