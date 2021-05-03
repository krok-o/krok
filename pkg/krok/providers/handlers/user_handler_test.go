package handlers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/krok-o/krok/pkg/krok/providers/mocks"
	"github.com/krok-o/krok/pkg/models"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserHandler_CreateUser(t *testing.T) {
	// setup handler
	mus := &mocks.UserStorer{}
	log := zerolog.New(os.Stderr)
	uh := NewUserHandler(UserHandlerDependencies{
		Logger:    log,
		UserStore: mus,
	})

	// setup expected mock calls
	mus.On("Create", mock.Anything, &models.User{
		DisplayName: "Gergely",
		Email:       "bla@bla.com",
	}).Return(&models.User{
		DisplayName: "Gergely",
		Email:       "bla@bla.com",
		ID:          0,
		LastLogin:   time.Date(1981, 1, 1, 1, 1, 1, 1, time.UTC),
		APIKeys:     nil,
	}, nil)

	token, err := generateTestToken("test@email.com")
	assert.NoError(t, err)
	userPost := `{"display_name":"Gergely","email":"bla@bla.com"}`
	userExpected := `{"display_name":"Gergely","email":"bla@bla.com","id":0,"last_login":"1981-01-01T01:01:01.000000001Z"}
`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/user/create", strings.NewReader(userPost))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err = uh.CreateUser()(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	content, err := ioutil.ReadAll(rec.Body)
	assert.NoError(t, err)
	assert.Equal(t, string(content), userExpected)
}
