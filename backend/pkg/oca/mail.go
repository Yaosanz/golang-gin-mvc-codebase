package oca

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// SendEmailResponse Response untuk Send Email
type SendEmailResponse struct {
	Success bool   `json:"success"`
	MsgID   string `json:"msg_id"`
}

// EmailHistoryResponse Response untuk Get History
type EmailHistoryResponse struct {
	Page        int `json:"page"`
	Pages       int `json:"pages"`
	TotalInPage int `json:"total_inpage"`
	Limit       int `json:"limit"`
	TotalAll    int `json:"total_all"`
	Data        []struct {
		History   []interface{} `json:"history"`
		Email     string        `json:"email"`
		Status    string        `json:"status"`
		CreatedAt time.Time     `json:"created_at"`
		Source    string        `json:"source"`
		MsgID     string        `json:"msgid"`
	} `json:"data"`
}

// EmailDetailResponse Response untuk Email Detail
type EmailDetailResponse struct {
	Status bool `json:"status"`
	Data   struct {
		MsgID       string        `json:"msgid"`
		Attachment  bool          `json:"attachment"`
		History     []interface{} `json:"history"`
		Email       string        `json:"email"`
		Status      string        `json:"status"`
		CreatedAt   time.Time     `json:"created_at"`
		Source      string        `json:"source"`
		MessageBody string        `json:"message_body"`
	} `json:"data"`
}

type email struct {
	baseUrl     string
	urlPrefix   string
	authToken   string
	senderEmail string
	senderName  string
}

type Sender struct {
	Email string
	Name  string
}

// SendSingleEmail send single email (attachment optional)
func (c *Client) SendSingleEmail(to, subject, htmlBody string, attachmentPath string, sender *Sender) (*SendEmailResponse, error) {
	url := c.email.baseUrl + c.email.urlPrefix + "/send-single"

	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	// form fields
	_ = writer.WriteField("email", to)
	_ = writer.WriteField("subject", subject)

	// encode message body ke base64
	encodedMsg := base64.StdEncoding.EncodeToString([]byte(htmlBody))
	_ = writer.WriteField("message", encodedMsg)

	// default sender
	fromEmail := c.email.senderEmail
	fromName := c.email.senderName
	if sender != nil {
		if sender.Email != "" {
			fromEmail = sender.Email
		}
		if sender.Name != "" {
			fromName = sender.Name
		}
	}

	_ = writer.WriteField("sender_email", fromEmail)
	_ = writer.WriteField("sender_name", fromName)

	// attachment (optional)
	if attachmentPath != "" {
		file, err := os.Open(attachmentPath)
		if err != nil {
			return nil, err
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				fmt.Println("Error closing file:", err)
			}
		}(file)

		part, err := writer.CreateFormFile("attachment", filepath.Base(attachmentPath))
		if err != nil {
			return nil, err
		}
		if _, err = io.Copy(part, file); err != nil {
			return nil, err
		}
	}

	_ = writer.Close()

	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.email.authToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	var result SendEmailResponse
	if err := c.doRequest(req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetEmailHistory get email history
func (c *Client) GetEmailHistory(page, limit int) (*EmailHistoryResponse, error) {
	url := fmt.Sprintf("%s/email/history?page=%d&limit=%d", c.email.baseUrl+c.email.urlPrefix, page, limit)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.email.authToken)

	var result EmailHistoryResponse
	if err := c.doRequest(req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetEmailDetail get email detail by msgID
func (c *Client) GetEmailDetail(msgID string) (*EmailDetailResponse, error) {
	url := fmt.Sprintf("%s/email/status/%s", c.email.baseUrl+c.email.urlPrefix, msgID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.email.authToken)

	var result EmailDetailResponse
	if err := c.doRequest(req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
