package tasker

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	gocache "github.com/patrickmn/go-cache"
	"golang.org/x/oauth2/google"
)

var (
	cache               = gocache.New(15*time.Minute, 30*time.Minute)
	accessTokenCacheKey = "access_token"
)

func getAccessToken(ctx context.Context) (string, error) {
	// try to retrieve access token from cache
	accessToken, foundToken := cache.Get(accessTokenCacheKey)
	if foundToken {
		// return access token from cache
		return accessToken.(string), nil
	} else {
		// request new access token and cache it
		serviceAccountJson, err := os.ReadFile("service-account.json")
		if err != nil {
			log.Println(err)
			return "", err
		}

		jwtConfig, err := google.JWTConfigFromJSON(
			serviceAccountJson,
			"https://www.googleapis.com/auth/firebase.messaging",
		)
		if err != nil {
			return "", fmt.Errorf("error creating JWT config: %w", err)
		}

		tokenSource := jwtConfig.TokenSource(ctx)
		token, err := tokenSource.Token()
		if err != nil {
			return "", fmt.Errorf("error retrieving access token: %w", err)
		}

		// get ttl for cache
		ttl := time.Until(token.Expiry)
		if ttl <= 0 {
			//default ttl if there is any issue with expiration time
			ttl = 30 * time.Minute
		}

		// cache access token
		cache.Set(accessTokenCacheKey, token.AccessToken, ttl)

		return token.AccessToken, nil
	}
}
