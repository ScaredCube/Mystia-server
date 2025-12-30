package livekit

import (
	"time"

	"github.com/livekit/protocol/auth"
)

type LiveKitProvider struct {
	ApiKey    string
	ApiSecret string
	Host      string
}

func NewLiveKitProvider(apiKey, apiSecret, host string) *LiveKitProvider {
	return &LiveKitProvider{
		ApiKey:    apiKey,
		ApiSecret: apiSecret,
		Host:      host,
	}
}

func (p *LiveKitProvider) GenerateToken(roomName, identity, nickname string) (string, error) {
	at := auth.NewAccessToken(p.ApiKey, p.ApiSecret)
	at.SetIdentity(identity)
	at.SetName(nickname)
	at.SetValidFor(time.Hour)

	grant := &auth.VideoGrant{
		RoomJoin: true,
		Room:     roomName,
	}
	at.AddGrant(grant)

	return at.ToJWT()
}
