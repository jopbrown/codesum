package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/jopbrown/gobase/errors"
	"github.com/jopbrown/gobase/log"
	"github.com/sashabaranov/go-openai"
)

func UpdateApiServerAccessToken(endpoint, token string) error {
	url, err := url.ParseRequestURI(endpoint)
	if err != nil {
		return errors.ErrorAt(err)
	}

	reqUrl := fmt.Sprintf("%s://%s/admin/tokens", url.Scheme, url.Host)
	tokens := []string{token}

	log.Debugf("send update token request to `%s`", reqUrl)

	payload, err := json.Marshal(tokens)
	if err != nil {
		return errors.ErrorAt(err)
	}

	req, err := http.NewRequest("PATCH", reqUrl, bytes.NewBuffer(payload))
	if err != nil {
		return errors.ErrorAt(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "TotallySecurePassword")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		if err != nil {
			return errors.ErrorAt(err)
		}

	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("update token request got non-200 status: %d", resp.StatusCode)
	}

	return nil
}

func TrimMDTableHeader(summary string) string {
	sb := &strings.Builder{}
	s := bufio.NewScanner(strings.NewReader(summary))
	lineNO := 0
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" {
			continue
		}

		lineNO++
		if lineNO > 2 {
			sb.WriteString(s.Text())
			sb.WriteByte('\n')
		}
	}

	return sb.String()
}

func IsErrHTTP413(err error) bool {
	if reqErr, ok := errors.AsIs[*openai.RequestError](err); ok {
		if reqErr.HTTPStatusCode == 413 {
			return true
		}
	}

	if errRes, ok := errors.AsIs[*openai.APIError](err); ok {
		if errRes.HTTPStatusCode == 413 {
			return true
		}
	}

	return false
}

func IsErrUnexpectedEOF(err error) bool {
	return errors.Is(err, io.ErrUnexpectedEOF)
}
