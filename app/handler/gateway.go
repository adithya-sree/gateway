package handler

import (
	"crypto/tls"
	"fmt"
	"github.com/adithya-sree/commons"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

// Gateway Request
func (h Handler) GatewayRouterV1() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// Get Route DTO from Context
		dto := request.Context().Value("session").(*RouteDTO)
		// Create Request Endpoint
		endpoint := dto.Definition.Config.BackendFQDN + dto.QualifyingPath
		out.Infof("[%s] - GatewayRouterV1 routing request to [%s]", dto.SessionId, endpoint)
		start := time.Now()
		req, err := http.NewRequest(request.Method, endpoint, request.Body)
		if err != nil {
			// Return 500 if error building request
			msg := fmt.Sprintf("[%s] - Unable to build request [%v]", dto.SessionId, err)
			out.Error(msg)
			_, _ = commons.RespondError(writer, http.StatusInternalServerError, msg)
			return
		}
		req.Header = request.Header
		// Create Request Client & Executed Request
		out.Infof("[%s] - Executing request to backend service", dto.SessionId)
		mutableTime := time.Now()
		client := client(dto.Definition.Config.BackendTimeout)
		resp, err := client.Do(req)
		if err != nil {
			// Return 502 if error executing request
			msg := fmt.Sprintf("[%s] - error connecting to backend service [%v], failed in [%d] ms",
				dto.SessionId, err, time.Since(mutableTime).Milliseconds())
			out.Info(msg)
			_, _ = commons.RespondError(writer, http.StatusBadGateway, msg)
			return
		}
		// Defer closing request body
		defer closeQuietly(dto.SessionId, resp.Body)
		// Return 200 with content if successfully
		out.Infof("[%s] - Reached backend service [%s] in [%d] ms", dto.SessionId, endpoint, time.Since(mutableTime).Milliseconds())
		body, _ := ioutil.ReadAll(resp.Body)
		headers := headers(resp.Header)
		_, _ = commons.Respond(writer, resp.StatusCode, body, headers)
		out.Infof("[%s] - Routed request, total latency [%d] ms", dto.SessionId, time.Since(start).Milliseconds())
		return
	}
}

func closeQuietly(session string, streams ...io.ReadCloser) {
	for _, stream := range streams {
		if stream != nil {
			err := stream.Close()
			if err != nil {
				out.Errorf("[%s] - error closing input stream, potential memory leak [%v]", session, err)
			}
		}
	}
}

func headers(resp http.Header) map[string]string {
	headers := make(map[string]string)
	for k, v := range resp {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}
	return headers
}

func client(tmOut int) *http.Client {
	to := time.Duration(tmOut) * time.Second
	tr := &http.Transport{
		DialContext: (&net.Dialer{
			KeepAlive: to,
		}).DialContext,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	return &http.Client{
		Transport: tr,
		Timeout:   to,
	}
}