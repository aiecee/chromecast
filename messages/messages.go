package messages

var (
	// Connection Payloads
	ConnectPayload = PayloadHeader{Type: "CONNECT"}
	ClosePayload   = PayloadHeader{Type: "CLOSE"}
	// Heartbeat Payloads
	PingPayload = PayloadHeader{Type: "PING"}
	PongPayload = PayloadHeader{Type: "PONG"}
	//Media Payloads
	MediaStatusPayload = PayloadHeader{Type: "GET_STATUS"}
	PlayMediaPayload   = PayloadHeader{Type: "PLAY"}
	PauseMediaPayload  = PayloadHeader{Type: "PAUSE"}
	StopMediaPayload   = PayloadHeader{Type: "STOP"}
	LoadMediaPayload   = PayloadHeader{Type: "LOAD"}
	// Reciever Payloads
	RecieverStatusPayload = PayloadHeader{Type: "GET_STATUS"}
	LaunchRecieverPayload = PayloadHeader{Type: "LAUNCH"}
	StopRecieverPayload   = PayloadHeader{Type: "STOP"}
	SetVolumePayload      = PayloadHeader{Type: "SET_VOLUME"}
	// URL Payloads
	URLStatusPayload = PayloadHeader{Type: "GET_STATUS"}
	LoadURLPayload   = PayloadHeader{Type: "LOAD"}
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

type MediaCommand struct {
	PayloadHeader
	MediaSessionID int `json:"mediaSessionId"`
}

type MediaItem struct {
	ContentID   string `json:"contentId"`
	StreamType  string `json:"streamType"`
	ContentType string `json:"contentType"`
}

type LoadMediaCommand struct {
	PayloadHeader
	Media       MediaItem   `json:"media"`
	CurrentTime int         `json:"currentTime"`
	Autoplay    bool        `json:"autoplay"`
	CustomData  interface{} `json:"customData"`
}

type MediaStatusMedia struct {
	ContentID   string  `json:"contentId"`
	StreamType  string  `json:"streamType"`
	ContentType string  `json:"contentType"`
	Duration    float64 `json:"duration"`
}

type Volume struct {
	Level *float64 `json:"level,omitempty"`
	Muted *bool    `json:"muted,omitempty"`
}

type MediaStatus struct {
	PayloadHeader
	MediaSessionID         int                    `json:"mediaSessionId"`
	PlaybackRate           float64                `json:"playbackRate"`
	PlayerState            string                 `json:"playerState"`
	CurrentTime            float64                `json:"currentTime"`
	SupportedMediaCommands int                    `json:"supportedMediaCommands"`
	Volume                 *Volume                `json:"volume,omitempty"`
	Media                  *MediaStatusMedia      `json:"media"`
	CustomData             map[string]interface{} `json:"customData"`
	RepeatMode             string                 `json:"repeatMode"`
	IdleReason             string                 `json:"idleReason"`
}

type MediaStatusResponse struct {
	PayloadHeader
	Status []*MediaStatus `json:"status,omitempty"`
}

type StatusResponse struct {
	PayloadHeader
	Status *ReceiverStatus `json:"status,omitempty"`
}

type ApplicationSession struct {
	AppID       *string      `json:"appId,omitempty"`
	DisplayName *string      `json:"displayName,omitempty"`
	Namespaces  []*Namespace `json:"namespaces"`
	SessionID   *string      `json:"sessionId,omitempty"`
	StatusText  *string      `json:"statusText,omitempty"`
	TransportID *string      `json:"transportId,omitempty"`
}

type Namespace struct {
	Name string `json:"name"`
}
type ReceiverStatus struct {
	PayloadHeader
	Applications []*ApplicationSession `json:"applications"`
	Volume       *Volume               `json:"volume,omitempty"`
}

func (s *ReceiverStatus) GetSessionByAppId(appId string) *ApplicationSession {
	for _, app := range s.Applications {
		if *app.AppID == appId {
			return app
		}
	}
	return nil
}

type LaunchRequest struct {
	PayloadHeader
	AppID string `json:"appId"`
}

type LoadURLCommand struct {
	PayloadHeader
	URL  string `json:"url"`
	Type string `json:"type"`
}

type URLStatusURL struct {
	ContentID   string  `json:"contentId"`
	StreamType  string  `json:"streamType"`
	ContentType string  `json:"contentType"`
	Duration    float64 `json:"duration"`
}

type URLStatusResponse struct {
	PayloadHeader
	Status []*URLStatus `json:"status,omitempty"`
}

type URLStatus struct {
	PayloadHeader
	URLSessionID         int                    `json:"mediaSessionId"`
	PlaybackRate         float64                `json:"playbackRate"`
	PlayerState          string                 `json:"playerState"`
	CurrentTime          float64                `json:"currentTime"`
	SupportedURLCommands int                    `json:"supportedURLCommands"`
	Volume               *Volume                `json:"volume,omitempty"`
	URL                  *URLStatusURL          `json:"media"`
	CustomData           map[string]interface{} `json:"customData"`
	RepeatMode           string                 `json:"repeatMode"`
	IdleReason           string                 `json:"idleReason"`
}
