package server_test

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/tui95/go-url-shortener/internal/database"
	"github.com/tui95/go-url-shortener/internal/server"
)

type RouterTestSuite struct {
	suite.Suite
	db     *sql.DB
	router *server.Router
}

func (suite *RouterTestSuite) SetupTest() {
	db := database.NewDB(":memory:")
	database.CreateTableIfNotExists(db)
	suite.db = db
	suite.router = server.NewRouter(db)
}

func (suite *RouterTestSuite) TearDownTest() {
	suite.db.Close()
}

func TestRouterTestSuite(t *testing.T) {
	suite.Run(t, new(RouterTestSuite))
}

func (suite *RouterTestSuite) TestCreateShortUrlHandler() {
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"url": "https://example.com"}`))
	// if not set, will default to example.com
	request.Host = "localhost:8080"
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	suite.Equal(http.StatusCreated, response.Code)
	assertResponseJsonEqual(suite.T(), response, map[string]string{"url": "http://localhost:8080/b"})
}

func (suite *RouterTestSuite) TestCreateShortUrlHandlerRequiredBody() {
	request := httptest.NewRequest(http.MethodPost, "/", nil)
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	suite.Equal(http.StatusBadRequest, response.Code)
	assertResponseJsonEqual(suite.T(), response, map[string]string{"detail": "Required body."})
}

func (suite *RouterTestSuite) TestRedirectToOriginalUrlHandler() {
	_, err := database.CreateUrlMapping(suite.db, "https://example.com")
	suite.Nil(err)
	request := httptest.NewRequest(http.MethodGet, "/b", nil)
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	suite.Equal(http.StatusSeeOther, response.Code)
	suite.Equal("https://example.com", response.Result().Header.Get("Location"))
}

func (suite *RouterTestSuite) TestRedirectToOriginalUrlHandlerInvalidBase64Id() {
	request := httptest.NewRequest(http.MethodGet, "/b", nil)
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	suite.Equal(http.StatusNotFound, response.Code)
	assertResponseJsonEqual(suite.T(), response, map[string]string{"detail": "Not found."})
}

func assertResponseJsonEqual(t *testing.T, response *httptest.ResponseRecorder, expected interface{}) {
	t.Helper()
	assert.Equal(t, "application/json", response.Header().Get("content-type"))
	assertJsonEqual(t, expected, response.Body.Bytes())
}

func assertJsonEqual(t *testing.T, expected interface{}, actual []byte) {
	t.Helper()
	expectedJsonString := toJson(t, expected)
	assert.JSONEq(t, expectedJsonString, string(actual))
}

func toJson(t *testing.T, v interface{}) string {
	t.Helper()
	jsonBytes, err := json.Marshal(v)
	assert.Nil(t, err)
	return string(jsonBytes)
}
