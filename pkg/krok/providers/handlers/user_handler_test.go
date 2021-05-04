package handlers

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	kerr "github.com/krok-o/krok/errors"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/krok-o/krok/pkg/krok/providers/mocks"
	"github.com/krok-o/krok/pkg/models"
)

func TestUserHandler_CreateUser(t *testing.T) {
	t.Run("normal flow", func(tt *testing.T) {
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
		req := httptest.NewRequest(http.MethodPost, "/user", strings.NewReader(userPost))
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
	})
	t.Run("when there is an error creating a user", func(tt *testing.T) {
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
		}).Return(nil, errors.New("nope"))

		token, err := generateTestToken("test@email.com")
		assert.NoError(t, err)
		userPost := `{"display_name":"Gergely","email":"bla@bla.com"}`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/user", strings.NewReader(userPost))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = uh.CreateUser()(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
	t.Run("when there is an invalid body which can't be parsed", func(tt *testing.T) {
		// setup handler
		mus := &mocks.UserStorer{}
		log := zerolog.New(os.Stderr)
		uh := NewUserHandler(UserHandlerDependencies{
			Logger:    log,
			UserStore: mus,
		})

		token, err := generateTestToken("test@email.com")
		assert.NoError(t, err)
		userPost := `<xml>`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/user", strings.NewReader(userPost))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = uh.CreateUser()(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mus.AssertNumberOfCalls(tt, "Create", 0)
	})
}

func TestUserHandler_GetUser(t *testing.T) {
	// setup handler
	mus := &mocks.UserStorer{}
	log := zerolog.New(os.Stderr)
	uh := NewUserHandler(UserHandlerDependencies{
		Logger:    log,
		UserStore: mus,
	})

	// setup expected mock calls
	mus.On("Get", mock.Anything, 0).Return(&models.User{
		DisplayName: "Gergely",
		Email:       "bla@bla.com",
		ID:          0,
		LastLogin:   time.Date(1981, 1, 1, 1, 1, 1, 1, time.UTC),
		APIKeys:     nil,
	}, nil)

	token, err := generateTestToken("test@email.com")
	assert.NoError(t, err)
	userExpected := `{"display_name":"Gergely","email":"bla@bla.com","id":0,"last_login":"1981-01-01T01:01:01.000000001Z"}
`
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/user/:id")
	c.SetParamNames("id")
	c.SetParamValues("0")
	err = uh.GetUser()(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	content, err := ioutil.ReadAll(rec.Body)
	assert.NoError(t, err)
	assert.Equal(t, string(content), userExpected)
}

func TestUserHandler_ListUsers(t *testing.T) {
	t.Run("normal flow", func(tt *testing.T) {
		// setup handler
		mus := &mocks.UserStorer{}
		log := zerolog.New(os.Stderr)
		uh := NewUserHandler(UserHandlerDependencies{
			Logger:    log,
			UserStore: mus,
		})

		// setup expected mock calls
		mus.On("List", mock.Anything).Return([]*models.User{
			{
				DisplayName: "Gergely",
				Email:       "bla@bla.com",
				ID:          0,
				LastLogin:   time.Date(1981, 1, 1, 1, 1, 1, 1, time.UTC),
				APIKeys:     nil,
			},
		}, nil)

		token, err := generateTestToken("test@email.com")
		assert.NoError(t, err)
		userExpected := `[{"display_name":"Gergely","email":"bla@bla.com","id":0,"last_login":"1981-01-01T01:01:01.000000001Z"}]
`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/users", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = uh.ListUsers()(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		content, err := ioutil.ReadAll(rec.Body)
		assert.NoError(t, err)
		assert.Equal(t, string(content), userExpected)
	})
	t.Run("when there is an error listing users", func(tt *testing.T) {
		// setup handler
		mus := &mocks.UserStorer{}
		log := zerolog.New(os.Stderr)
		uh := NewUserHandler(UserHandlerDependencies{
			Logger:    log,
			UserStore: mus,
		})

		// setup expected mock calls
		mus.On("List", mock.Anything).Return(nil, errors.New("nope"))

		token, err := generateTestToken("test@email.com")
		assert.NoError(t, err)
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/users", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = uh.ListUsers()(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestUserHandler_DeleteUser(t *testing.T) {
	t.Run("normal flow", func(tt *testing.T) {
		// setup handler
		mus := &mocks.UserStorer{}
		log := zerolog.New(os.Stderr)
		uh := NewUserHandler(UserHandlerDependencies{
			Logger:    log,
			UserStore: mus,
		})

		// setup expected mock calls
		mus.On("Delete", mock.Anything, 0).Return(nil)

		token, err := generateTestToken("test@email.com")
		assert.NoError(t, err)

		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/user/:id")
		c.SetParamNames("id")
		c.SetParamValues("0")
		err = uh.DeleteUser()(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
	t.Run("when there is an error deleting the user", func(tt *testing.T) {
		// setup handler
		mus := &mocks.UserStorer{}
		log := zerolog.New(os.Stderr)
		uh := NewUserHandler(UserHandlerDependencies{
			Logger:    log,
			UserStore: mus,
		})

		// setup expected mock calls
		mus.On("Delete", mock.Anything, 0).Return(errors.New("nope"))

		token, err := generateTestToken("test@email.com")
		assert.NoError(t, err)

		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/user/:id")
		c.SetParamNames("id")
		c.SetParamValues("0")
		err = uh.DeleteUser()(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
	t.Run("when the user id is invalid", func(tt *testing.T) {
		// setup handler
		mus := &mocks.UserStorer{}
		log := zerolog.New(os.Stderr)
		uh := NewUserHandler(UserHandlerDependencies{
			Logger:    log,
			UserStore: mus,
		})

		// setup expected mock calls
		mus.On("Delete", mock.Anything, 0).Return(nil)

		token, err := generateTestToken("test@email.com")
		assert.NoError(t, err)

		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/user/:id")
		c.SetParamNames("id")
		c.SetParamValues("asdf")
		err = uh.DeleteUser()(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
	t.Run("when the use doesn't exist", func(tt *testing.T) {
		// setup handler
		mus := &mocks.UserStorer{}
		log := zerolog.New(os.Stderr)
		uh := NewUserHandler(UserHandlerDependencies{
			Logger:    log,
			UserStore: mus,
		})

		// setup expected mock calls
		mus.On("Delete", mock.Anything, 0).Return(kerr.ErrNotFound)

		token, err := generateTestToken("test@email.com")
		assert.NoError(t, err)

		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/user/:id")
		c.SetParamNames("id")
		c.SetParamValues("0")
		err = uh.DeleteUser()(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestUserHandler_UpdateUser(t *testing.T) {
	t.Run("normal flow", func(tt *testing.T) {
		// setup handler
		mus := &mocks.UserStorer{}
		log := zerolog.New(os.Stderr)
		uh := NewUserHandler(UserHandlerDependencies{
			Logger:    log,
			UserStore: mus,
		})

		// setup expected mock calls
		mus.On("Update", mock.Anything, &models.User{
			DisplayName: "NewName",
			Email:       "bla@bla.com",
		}).Return(&models.User{
			DisplayName: "NewName",
			Email:       "bla@bla.com",
			ID:          0,
			LastLogin:   time.Date(1981, 1, 1, 1, 1, 1, 1, time.UTC),
			APIKeys:     nil,
		}, nil)

		token, err := generateTestToken("test@email.com")
		assert.NoError(t, err)
		userPost := `{"display_name":"NewName","email":"bla@bla.com"}`
		userExpected := `{"display_name":"NewName","email":"bla@bla.com","id":0,"last_login":"1981-01-01T01:01:01.000000001Z"}
`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/user/update", strings.NewReader(userPost))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = uh.UpdateUser()(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		content, err := ioutil.ReadAll(rec.Body)
		assert.NoError(t, err)
		assert.Equal(t, string(content), userExpected)
	})
	t.Run("when the user does not exist", func(tt *testing.T) {
		// setup handler
		mus := &mocks.UserStorer{}
		log := zerolog.New(os.Stderr)
		uh := NewUserHandler(UserHandlerDependencies{
			Logger:    log,
			UserStore: mus,
		})

		// setup expected mock calls
		mus.On("Update", mock.Anything, &models.User{
			DisplayName: "NewName",
			Email:       "bla@bla.com",
		}).Return(nil, kerr.ErrNotFound)

		token, err := generateTestToken("test@email.com")
		assert.NoError(t, err)
		userPost := `{"display_name":"NewName","email":"bla@bla.com"}`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/user/update", strings.NewReader(userPost))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = uh.UpdateUser()(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("when there is an error parsing the body", func(tt *testing.T) {
		// setup handler
		mus := &mocks.UserStorer{}
		log := zerolog.New(os.Stderr)
		uh := NewUserHandler(UserHandlerDependencies{
			Logger:    log,
			UserStore: mus,
		})

		token, err := generateTestToken("test@email.com")
		assert.NoError(t, err)
		userPost := `<xml>`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/user/update", strings.NewReader(userPost))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = uh.UpdateUser()(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mus.AssertNumberOfCalls(tt, "Update", 0)
	})

	t.Run("when update fails", func(tt *testing.T) {
		// setup handler
		mus := &mocks.UserStorer{}
		log := zerolog.New(os.Stderr)
		uh := NewUserHandler(UserHandlerDependencies{
			Logger:    log,
			UserStore: mus,
		})

		// setup expected mock calls
		mus.On("Update", mock.Anything, &models.User{
			DisplayName: "Gergely",
			Email:       "bla@bla.com",
		}).Return(nil, errors.New("nope"))

		token, err := generateTestToken("test@email.com")
		assert.NoError(t, err)
		userPost := `{"display_name":"Gergely","email":"bla@bla.com"}`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/user/update", strings.NewReader(userPost))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = uh.UpdateUser()(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
