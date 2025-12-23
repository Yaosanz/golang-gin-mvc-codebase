package google

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"fmt"
	"google.golang.org/api/option"
	"log"
	"path/filepath"
)

// newFirebaseApp initializes and returns the Firebase app instance.
func newFirebaseApp(serviceAccountPath string) (*firebase.App, error) {
	// path to serviceAccountKey.json in the root directory
	serviceAccount, err := filepath.Abs(serviceAccountPath)
	log.Print("[FCM-CRED-PATH]: ", serviceAccount)
	if err != nil {
		return nil, fmt.Errorf("could not get service account path: %v", err)
	}

	opt := option.WithCredentialsFile(serviceAccount)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("could not initialize firebase app: %v", err)
	}
	return app, nil
}
