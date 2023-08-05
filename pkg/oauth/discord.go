package oauth

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/memnix/memnix-rest/config"
	"github.com/memnix/memnix-rest/infrastructures"
	"github.com/memnix/memnix-rest/views"
	"github.com/pkg/errors"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// GetDiscordAccessToken gets the access token from Discord
func GetDiscordAccessToken(ctx context.Context, code string) (string, error) {
	_, span := infrastructures.GetFiberTracer().Start(ctx, "GetDiscordAccessToken")
	defer span.End()
	reqBody := bytes.NewBuffer([]byte(fmt.Sprintf(
		"client_id=%s&client_secret=%s&grant_type=authorization_code&redirect_uri=%s&code=%s&scope=identify,email",
		infrastructures.AppConfig.DiscordConfig.ClientID,
		infrastructures.AppConfig.DiscordConfig.ClientSecret,
		config.GetCurrentURL()+"/v2/security/discord_callback",
		code,
	)))

	// POST request to set URL
	req, reqerr := http.NewRequestWithContext(ctx,
		http.MethodPost,
		"https://discord.com/api/oauth2/token",
		reqBody,
	)
	if reqerr != nil {
		otelzap.Ctx(ctx).Error("Failed to get Discord access token", zap.Error(reqerr))
		return "", errors.Wrap(reqerr, views.RequestFailed)
	}

	if req == nil || req.Body == nil || req.Header == nil {
		otelzap.Ctx(ctx).Error("Failed to get Discord access token", zap.Error(reqerr))
		return "", errors.New(views.RequestFailed)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	// Get the response
	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		otelzap.Ctx(ctx).Error("Failed to get Discord access token", zap.Error(resperr))
		return "", errors.Wrap(resperr, views.ResponseFailed)
	}

	if resp == nil || resp.Body == nil {
		otelzap.Ctx(ctx).Error("resp is empty", zap.Error(resperr))
		return "", errors.New(views.ResponseFailed)
	}

	// Response body converted to stringified JSON
	respbody, _ := io.ReadAll(resp.Body)

	otelzap.Ctx(ctx).Debug("Discord access token response", zap.String("response", string(respbody)))

	// Represents the response received from Github
	type discordAccessTokenResponse struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		Scope        string `json:"scope"`
		Expires      int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
	}

	// Convert stringified JSON to a struct object of type githubAccessTokenResponse
	var ghresp discordAccessTokenResponse
	err := config.JSONHelper.Unmarshal(respbody, &ghresp)
	if err != nil {
		return "", err
	}

	span.AddEvent("Discord access token received", trace.WithAttributes(attribute.String("access_token", ghresp.AccessToken)))

	// Return the access token (as the rest of the
	// details are relatively unnecessary for us)
	return ghresp.AccessToken, nil
}

// GetDiscordData gets the user data from Discord
func GetDiscordData(ctx context.Context, accessToken string) (string, error) {
	_, span := infrastructures.GetFiberTracer().Start(ctx, "GetDiscordData")
	defer span.End()

	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet,
		"https://discord.com/api/users/@me",
		nil,
	)
	if err != nil {
		return "", errors.Wrap(err, views.RequestFailed)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	// Get the response
	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		return "", errors.Wrap(resperr, views.ResponseFailed)
	}

	// Response body converted to stringified JSON
	respbody, _ := io.ReadAll(resp.Body)

	span.AddEvent("Discord user data received", trace.WithAttributes(attribute.String("user_data", string(respbody))))

	return string(respbody), nil
}
