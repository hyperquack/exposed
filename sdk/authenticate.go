package sdk

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"os"
	"path/filepath"
)

type CognitoClient struct {
	ClientID string `json:"client_id"`
	API      string `json:"api"`
	Username string `json:"username"`
	Password string `json:"password"`
	// temp attributes
	AccessToken  string `json:"access_token"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

func NewCognitoClientFromJSON(filePath string) (*CognitoClient, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("valid keychain is required: %v", err)
	}
	defer file.Close()

	var client CognitoClient
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&client); err != nil {
		return nil, fmt.Errorf("failed to parse keychain as JSON: %v", err)
	}

	return &client, nil
}

func Authenticate() (*CognitoClient, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	keychain := filepath.Join(home, ".pdc", "keychain.json")

	c, err := NewCognitoClientFromJSON(keychain)
	if err != nil {
		return nil, err
	}

	err = c.Refresh()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *CognitoClient) Refresh() error {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-2"),
	}))

	svc := cognitoidentityprovider.New(sess)

	authInput := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(c.Username),
			"PASSWORD": aws.String(c.Password),
		},
		ClientId: aws.String(c.ClientID),
	}

	authOutput, err := svc.InitiateAuth(authInput)
	if err != nil {
		return fmt.Errorf("failed to authenticate: %v", err)
	}

	c.AccessToken = *authOutput.AuthenticationResult.AccessToken
	c.IDToken = *authOutput.AuthenticationResult.IdToken
	c.RefreshToken = *authOutput.AuthenticationResult.RefreshToken

	return nil
}
