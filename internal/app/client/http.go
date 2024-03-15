package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	log "github.com/sirupsen/logrus"

	"github.com/gsk148/gophkeeper/internal/app/client/config"
	"github.com/gsk148/gophkeeper/internal/app/models"
)

type HTTPKeeperClient struct {
	http   *http.Client
	apiURL *url.URL
}

const (
	SBinary   string = "/storage/binary/"
	SCard     string = "/storage/card/"
	SPassword string = "/storage/password/"
	SText     string = "/storage/text/"
)

var ErrUnauthorized = errors.New("incorrect username or password")

func NewHTTPClient(cfg *config.ClientConfig) (HTTPKeeperClient, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return HTTPKeeperClient{}, err
	}

	uri, err := url.Parse(cfg.GetAPIAddress())
	if err != nil {
		return HTTPKeeperClient{}, err
	}

	client := HTTPKeeperClient{
		http:   &http.Client{Jar: jar},
		apiURL: uri,
	}

	caCertPool, caErr := cfg.GetCACertPool()
	c, cErr := cfg.GetCertificate()
	if caErr == nil && cErr == nil {
		client.http.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{c},
				MinVersion:   tls.VersionTLS12,
			},
		}
	}

	return client, nil
}

func (c HTTPKeeperClient) Login(ctx context.Context, user, password string) error {
	res, err := c.makeRequest(ctx, http.MethodPost, "/auth/login", models.UserRequest{
		Name:     user,
		Password: password,
	})
	if err != nil {
		if res.StatusCode == http.StatusUnauthorized {
			return ErrUnauthorized
		}
		return err
	}
	defer closeResponseBody(res.Body)

	c.http.Jar.SetCookies(c.apiURL, res.Cookies())
	return nil
}

func (c HTTPKeeperClient) Logout(ctx context.Context) error {
	res, err := c.makeRequest(ctx, http.MethodPost, "/auth/logout", nil)
	if err != nil {
		return err
	}
	defer closeResponseBody(res.Body)

	c.http.Jar.SetCookies(c.apiURL, nil)
	return nil
}

func (c HTTPKeeperClient) Register(ctx context.Context, user, password string) error {
	res, err := c.makeRequest(ctx, http.MethodPost, "/auth/register", models.UserRequest{
		Name:     user,
		Password: password,
	})
	if err != nil {
		return err
	}
	defer closeResponseBody(res.Body)
	return nil
}

func (c HTTPKeeperClient) DeleteBinary(ctx context.Context, id string) error {
	return c.deleteData(ctx, SBinary, id)
}

func (c HTTPKeeperClient) GetAllBinaries(ctx context.Context) ([]models.BinaryResponse, error) {
	body, err := c.getAllData(ctx, SBinary)
	if err != nil {
		return nil, err
	}
	defer closeResponseBody(body)

	var bins []models.BinaryResponse
	err = json.NewDecoder(body).Decode(&bins)
	return bins, err
}

func (c HTTPKeeperClient) GetBinaryByID(ctx context.Context, id string) (models.BinaryResponse, error) {
	var data models.BinaryResponse
	body, err := c.getDataByID(ctx, SBinary, id)
	if err != nil {
		return data, err
	}
	defer closeResponseBody(body)

	err = json.NewDecoder(body).Decode(&data)
	return data, err
}

func (c HTTPKeeperClient) StoreBinary(ctx context.Context, name string,
	data []byte, note string,
) (string, error) {
	return c.storeData(ctx, SBinary, models.BinaryRequest{
		Name: name,
		Data: data,
		Note: note,
	})
}

func (c HTTPKeeperClient) DeleteText(ctx context.Context, id string) error {
	return c.deleteData(ctx, SText, id)
}

func (c HTTPKeeperClient) GetAllTexts(ctx context.Context) ([]models.TextResponse, error) {
	body, err := c.getAllData(ctx, SText)
	if err != nil {
		return nil, err
	}
	defer closeResponseBody(body)

	var data []models.TextResponse
	err = json.NewDecoder(body).Decode(&data)
	return data, err
}

func (c HTTPKeeperClient) GetTextByID(ctx context.Context, id string) (models.TextResponse, error) {
	var data models.TextResponse
	body, err := c.getDataByID(ctx, SText, id)
	if err != nil {
		return data, err
	}
	defer closeResponseBody(body)

	err = json.NewDecoder(body).Decode(&data)
	return data, err
}

func (c HTTPKeeperClient) StoreText(ctx context.Context, name, data, note string) (string, error) {
	return c.storeData(ctx, SText, models.TextRequest{
		Name: name,
		Data: data,
		Note: note,
	})
}

func (c HTTPKeeperClient) DeleteCard(ctx context.Context, id string) error {
	return c.deleteData(ctx, SCard, id)
}

func (c HTTPKeeperClient) GetAllCards(ctx context.Context) ([]models.CardResponse, error) {
	body, err := c.getAllData(ctx, SCard)
	if err != nil {
		return nil, err
	}
	defer closeResponseBody(body)

	var data []models.CardResponse
	err = json.NewDecoder(body).Decode(&data)
	return data, err
}

func (c HTTPKeeperClient) GetCardByID(ctx context.Context, id string) (models.CardResponse, error) {
	var data models.CardResponse
	body, err := c.getDataByID(ctx, SCard, id)
	if err != nil {
		return data, err
	}
	defer closeResponseBody(body)

	err = json.NewDecoder(body).Decode(&data)
	return data, err
}

func (c HTTPKeeperClient) StoreCard(ctx context.Context,
	name, number, holder, expDate, cvv, note string,
) (string, error) {
	return c.storeData(ctx, SCard, models.CardRequest{
		Name:    name,
		Number:  number,
		Holder:  holder,
		ExpDate: expDate,
		CVV:     cvv,
		Note:    note,
	})
}

func (c HTTPKeeperClient) DeletePassword(ctx context.Context, id string) error {
	return c.deleteData(ctx, SPassword, id)
}

func (c HTTPKeeperClient) GetAllPasswords(ctx context.Context) ([]models.PasswordResponse, error) {
	body, err := c.getAllData(ctx, SPassword)
	if err != nil {
		return nil, err
	}
	defer closeResponseBody(body)

	var data []models.PasswordResponse
	err = json.NewDecoder(body).Decode(&data)
	return data, err
}

func (c HTTPKeeperClient) GetPasswordByID(ctx context.Context, id string) (models.PasswordResponse, error) {
	var data models.PasswordResponse
	body, err := c.getDataByID(ctx, SPassword, id)
	if err != nil {
		return data, err
	}
	defer closeResponseBody(body)

	err = json.NewDecoder(body).Decode(&data)
	return data, err
}

func (c HTTPKeeperClient) StorePassword(ctx context.Context, name, user, password, note string) (string, error) {
	return c.storeData(ctx, SPassword, models.PasswordRequest{
		Name:     name,
		User:     user,
		Password: password,
		Note:     note,
	})
}

func (c HTTPKeeperClient) deleteData(ctx context.Context, url, id string) error {
	_, err := c.makeRequest(ctx, http.MethodDelete, url+id, nil)
	return err
}

func (c HTTPKeeperClient) getAllData(ctx context.Context, url string) (io.ReadCloser, error) {
	res, err := c.makeRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return res.Body, nil
}

func (c HTTPKeeperClient) getDataByID(ctx context.Context, url, id string) (io.ReadCloser, error) {
	res, err := c.makeRequest(ctx, http.MethodGet, url+id, nil)
	if err != nil {
		return nil, err
	}
	return res.Body, nil
}

func (c HTTPKeeperClient) storeData(ctx context.Context, url string, data any) (string, error) {
	res, err := c.makeRequest(ctx, http.MethodPost, url, data)
	if err != nil {
		return "", err
	}
	defer closeResponseBody(res.Body)

	id, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(id), err
}

func (c HTTPKeeperClient) makeRequest(ctx context.Context, method, url string, data any) (*http.Response, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var reader io.Reader
	if len(body) > 0 {
		reader = bytes.NewBuffer(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.apiURL.String()+url, reader)
	if err != nil {
		return nil, err
	}

	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return res, errors.New(res.Status)
	}
	return res, err
}

func closeResponseBody(b io.Closer) {
	if err := b.Close(); err != nil {
		log.Error(err)
	}
}
