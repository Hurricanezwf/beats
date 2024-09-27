// Note: BY ZWF;
package cls

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pierrec/lz4"
)

const logUri = "/structuredlog"

type CLSHTTPClient struct {
	endpoint  string
	host      string
	accessKey string
	secretKey string
	client    *http.Client
}

func NewCLSHTTPClient(endpoint, accessKey, secretKey string) (*CLSHTTPClient, error) {
	if endpoint == "" {
		return nil, errors.New("cls endpoint cannot be empty")
	}
	urlobj, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid cls endpoint, %w", err)
	}
	if urlobj.Scheme != "http" && urlobj.Scheme != "https" {
		return nil, errors.New("cls endpoint must be prefixed with http or https")
	}
	if accessKey == "" {
		return nil, errors.New("accessKey is required")
	}
	if secretKey == "" {
		return nil, errors.New("secretKey is required")
	}

	return &CLSHTTPClient{
		endpoint:  endpoint,
		host:      urlobj.Host,
		accessKey: accessKey,
		secretKey: secretKey,
		client: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   5 * time.Second,
					KeepAlive: 300 * time.Second,
				}).DialContext,
				MaxIdleConns:        50,
				MaxIdleConnsPerHost: 50,
				MaxConnsPerHost:     50,
				IdleConnTimeout:     300 * time.Second,
			},
			Timeout: 10 * time.Second, // 默认请求10秒超时
		},
	}, nil
}

func (client *CLSHTTPClient) lz4Compress(body []byte, params url.Values, urlReport string) (*http.Request, *CLSError) {
	out := make([]byte, lz4.CompressBlockBound(len(body)))
	var hashTable [1 << 16]int
	n, err := lz4.CompressBlock(body, out, hashTable[:])
	if err != nil {
		return nil, NewError(-1, "", BAD_REQUEST, err)
	}
	// copy incompressible data as lz4 format
	if n == 0 {
		n, _ = copyIncompressible(body, out)
	}
	req, err := http.NewRequest(http.MethodPost, urlReport, bytes.NewBuffer(out[:n]))
	if err != nil {
		return nil, NewError(-1, "", BAD_REQUEST, err)
	}
	req.URL.RawQuery = params.Encode()
	req.Header.Add("x-cls-compress-type", "lz4")
	return req, nil
}

type ErrorMessage struct {
	Code    string `json:"errorcode"`
	Message string `json:"errormessage"`
}

// Send cls实际发送接口
func (client *CLSHTTPClient) Send(ctx context.Context, topicId string, logGroupList LogGroupList) *CLSError {
	params := url.Values{"topic_id": []string{topicId}}
	headers := url.Values{"Host": {client.host}, "Content-Type": {"application/x-protobuf"}}
	authorization := signature(client.accessKey, client.secretKey, http.MethodPost, logUri, params, headers, 300)

	urlReport := strings.TrimSuffix(client.endpoint, "/") + logUri

	body, err := logGroupList.Marshal()
	if err != nil {
		return NewError(-1, "", INTERNAL_ERROR, fmt.Errorf("failed to encode logs, %v", err))
	}

	var req *http.Request
	var clsErr *CLSError

	if req, clsErr = client.lz4Compress(body, params, urlReport); clsErr != nil {
		return clsErr
	}

	req.Header.Add("Host", client.host)
	req.Header.Add("Content-Type", "application/x-protobuf")
	req.Header.Add("Authorization", authorization)
	req.Header.Add("User-Agent", "cls-go-sdk-1.0.7")

	req = req.WithContext(ctx)
	resp, err := client.client.Do(req)
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Temporary() {
				return NewError(-1, "--NetError--", TEMPORARY_ERROR, err)
		}
		return NewError(-1, "--UnknownError--", UNKNOWN_ERROR, err)
	}
	defer resp.Body.Close()

	// 401, 403, 404, 413 直接返回错误
	if resp.StatusCode == 401 || resp.StatusCode == 403 || resp.StatusCode == 404 || resp.StatusCode == 413 {
		v, err := io.ReadAll(resp.Body)
		if err != nil {
			return NewError(int32(resp.StatusCode), resp.Header.Get("X-Cls-Requestid"), BAD_REQUEST, errors.New("bad request"))
		}
		var message ErrorMessage
		if err := json.Unmarshal(v, &message); err != nil {
			return NewError(int32(resp.StatusCode), resp.Header.Get("X-Cls-Requestid"), BAD_REQUEST, errors.New("bad request"))
		}
		return NewError(int32(resp.StatusCode), resp.Header.Get("X-Cls-Requestid"), message.Code, errors.New(message.Message))
	}
	// 200 直接返回
	if resp.StatusCode == 200 {
		return nil
	}

	// 如果被服务端写入限速
	if resp.StatusCode == 429 {
		return NewError(int32(resp.StatusCode), resp.Header.Get("X-Cls-Requestid"), WRITE_QUOTA_EXCEED, errors.New("write quota exceed"))
	}
	// 如果是服务端错误
	if resp.StatusCode >= 500 {
		return NewError(int32(resp.StatusCode), resp.Header.Get("X-Cls-Requestid"), INTERNAL_SERVER_ERROR, errors.New("server internal error"))
	}
	return NewError(int32(resp.StatusCode), resp.Header.Get("X-Cls-Requestid"), UNKNOWN_ERROR, errors.New("unknown error"))
}

func copyIncompressible(src, dst []byte) (int, error) {
	lLen, dn := len(src), len(dst)
	di := 0
	if lLen < 0xF {
		dst[di] = byte(lLen << 4)
	} else {
		dst[di] = 0xF0
		if di++; di == dn {
			return di, nil
		}
		lLen -= 0xF
		for ; lLen >= 0xFF; lLen -= 0xFF {
			dst[di] = 0xFF
			if di++; di == dn {
				return di, nil
			}
		}
		dst[di] = byte(lLen)
	}
	if di++; di+len(src) > dn {
		return di, nil
	}
	di += copy(dst[di:], src)
	return di, nil
}
