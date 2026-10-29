package proxy_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blazejsewera/go-test-proxy/config"
	"github.com/blazejsewera/go-test-proxy/proxy"
	"github.com/blazejsewera/go-test-proxy/test/assert"
)

func TestBuilder(t *testing.T) {
	cfg := config.Configuration{Target: "target", Port: 1}
	endpointPattern := "/endpoint"

	builderPrototype := func() *proxy.Builder {
		return proxy.NewBuilder().
			WithProxyTarget(cfg.Target).
			WithPort(cfg.Port).
			WithMockGroup("mockGroup1", proxy.Mock{RoutePattern: endpointPattern, HandlerFunc: func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(201) }}).
			WithMockGroup("mockGroup2", proxy.Mock{RoutePattern: endpointPattern, HandlerFunc: func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(202) }}).
			WithMockGroup("mockGroup3", proxy.Mock{RoutePattern: endpointPattern, HandlerFunc: func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(203) }})
	}

	assertServerRespondsWithTheFirstMock := func(t testing.TB, builder *proxy.Builder) {
		t.Helper()
		server := buildTestServer(builder)
		defer server.Close()

		client := server.Client()
		res, err := client.Get(server.URL + endpointPattern)
		assert.Nil(t, err)
		assert.Equal(t, 201, res.StatusCode)
	}

	t.Run("enable only one mock when there are multiple mocks for one endpoint", func(t *testing.T) {
		t.Run("when all mocks are enabled", func(t *testing.T) {
			builder := builderPrototype().WithEnabledAllMocks()
			assertServerRespondsWithTheFirstMock(t, builder)
		})

		t.Run("when mocks are enabled explicitly", func(t *testing.T) {
			builder := builderPrototype().WithEnabledMockGroups([]string{"mockGroup1", "mockGroup2"})
			assertServerRespondsWithTheFirstMock(t, builder)
		})
	})

	t.Run("deduplicate explicitly enabled mock groups", func(t *testing.T) {
		builder := builderPrototype().WithEnabledMockGroups([]string{"mockGroup1", "mockGroup1", "mockGroup2"})
		assertServerRespondsWithTheFirstMock(t, builder)
	})
}

func buildTestServer(builder *proxy.Builder) *httptest.Server {
	return httptest.NewUnstartedServer(builder.Router)
}
