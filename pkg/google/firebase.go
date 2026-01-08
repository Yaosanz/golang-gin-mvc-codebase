package google

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"fmt"
	"google.golang.org/api/option"
	"log"
	"os"
	"path/filepath"
)

// newFirebaseApp initializes and returns the Firebase app instance.
func newFirebaseApp(serviceAccountPath string) (*firebase.App, error) {
	serviceAccount, err := filepath.Abs(serviceAccountPath)
	if err != nil {
		return nil, fmt.Errorf("could not get service account path: %v", err)
	}

	log.Println("[FCM-CRED-PATH]:", serviceAccount)

	projectID := os.Getenv("FCM_PROJECT_ID")
	log.Println("[FCM-PROJECT-ID]:", projectID)

	if projectID == "" {
		return nil, fmt.Errorf("FCM_PROJECT_ID is required but not set")
	}

	conf := &firebase.Config{
		ProjectID: projectID,
	}

	opt := option.WithCredentialsFile(serviceAccount)

	app, err := firebase.NewApp(
		context.Background(),
		conf,
		opt,
	)
	if err != nil {
		return nil, fmt.Errorf("could not initialize firebase app: %v", err)
	}

	log.Println("✅ Firebase app initialized successfully")
	return app, nil
}
