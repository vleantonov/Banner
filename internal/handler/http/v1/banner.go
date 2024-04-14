package v1

import (
	"banner/internal/domain"
	"banner/internal/handler/http/v1/gen"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type DefaultResponse struct {
	Message string `json:"message"`
}

type BannerService interface {
	GetActiveContentByTagFeatureID(ctx context.Context, tagID, FeatureID int, useLast bool) (*map[string]interface{}, error)
	GetByFilter(ctx context.Context, f domain.FilterBanner) (*[]domain.Banner, error)
	Update(ctx context.Context, banner domain.UpdBanner) error
	Create(ctx context.Context, banner domain.Banner) (int, error)
	DeleteById(ctx context.Context, id int) error
	Delete(ctx context.Context, tagID, featureID *int) error
}

type AuthService interface {
	Login(ctx context.Context, login string, password string) (string, error)
	RegisterNewUser(ctx context.Context, login string, password string) error
	IsAdmin(ctx context.Context, token string) (bool, error)
}

type Router struct {
	l *zap.Logger
	s BannerService
	a AuthService
}

func New(s BannerService, a AuthService) *Router {
	return &Router{
		s: s,
		a: a,
	}
}

func (r *Router) GetBanner(c *gin.Context, params api.GetBannerParams) {
	if !r.adminVerify(c, params.Token) {
		return
	}

	b, err := r.s.GetByFilter(c, domain.FilterBanner{
		FeatureID: params.FeatureId,
		TagID:     params.TagId,
		Limit:     params.Limit,
		Offset:    params.Offset,
	})

	if err != nil {
		if errors.Is(err, domain.ErrBannerNotFound) {
			c.Status(http.StatusNotFound)
			return
		}
		msg := err.Error()
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Error: &msg,
		})
		return
	}
	c.JSON(http.StatusOK, b)
}

func (r *Router) PostBanner(c *gin.Context, params api.PostBannerParams) {
	if !r.adminVerify(c, params.Token) {
		return
	}

	var requestBody api.PostBannerRequestBody
	if err := c.Bind(&requestBody); err != nil {
		msg := err.Error()
		c.JSON(http.StatusBadRequest, api.ErrorResponse{Error: &msg})
		return
	}

	if requestBody.TagIds == nil || len(*requestBody.TagIds) == 0 || requestBody.FeatureId == nil {
		msg := domain.ErrTagFeatureRequired.Error()
		c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Error: &msg,
		})
		return
	}

	if requestBody.Content == nil {
		requestBody.Content = &map[string]interface{}{}
	}
	if requestBody.IsActive == nil {
		b := false
		requestBody.IsActive = &b
	}

	idxResp, err := r.s.Create(c, domain.Banner{
		Content: *requestBody.Content,
		Tags:    *requestBody.TagIds,
		Feature: *requestBody.FeatureId,
		Active:  *requestBody.IsActive,
	})

	if err != nil {
		msg := err.Error()
		if errors.Is(err, domain.ErrTagFeatureAlreadyExists) {
			c.JSON(http.StatusBadRequest, api.ErrorResponse{
				Error: &msg,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Error: &msg,
		})
		return
	}
	c.JSON(http.StatusCreated, api.PostBannerResponse{BannerId: &idxResp})
}

func (r *Router) DeleteBannerId(c *gin.Context, id int, params api.DeleteBannerIdParams) {
	if !r.adminVerify(c, params.Token) {
		return
	}

	err := r.s.DeleteById(c, id)

	if err != nil {
		if errors.Is(err, domain.ErrBannerNotFound) {
			c.Status(http.StatusNotFound)
			return
		}
		msg := err.Error()
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Error: &msg,
		})
		return
	}
	c.Status(http.StatusNoContent)
}

func (r *Router) PatchBannerId(c *gin.Context, id int, params api.PatchBannerIdParams) {
	if !r.adminVerify(c, params.Token) {
		return
	}

	var requestBody api.PatchBannerRequestBody
	var err error

	if err := c.Bind(&requestBody); err != nil {
		msg := err.Error()
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Error: &msg,
		})
		return
	}

	err = r.s.Update(c, domain.UpdBanner{
		ID:      id,
		Tags:    requestBody.TagIds,
		Feature: requestBody.FeatureId,
		Content: requestBody.Content,
		Active:  requestBody.IsActive,
	})

	if err != nil {
		if errors.Is(err, domain.ErrBannerNotFound) {
			c.Status(http.StatusNotFound)
			return
		}
		msg := err.Error()
		if errors.Is(err, domain.ErrTagFeatureAlreadyExists) {
			c.JSON(http.StatusBadRequest, api.ErrorResponse{
				Error: &msg,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Error: &msg,
		})
		return
	}
	c.Status(http.StatusOK)
}

func (r *Router) GetUserBanner(c *gin.Context, params api.GetUserBannerParams) {
	if params.Token == nil || *params.Token == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	if _, err := r.a.IsAdmin(c, *params.Token); err != nil {
		c.AbortWithStatus(http.StatusForbidden)
	}

	if params.UseLastRevision == nil {
		u := false
		params.UseLastRevision = &u
	}
	content, err := r.s.GetActiveContentByTagFeatureID(c, params.TagId, params.FeatureId, *params.UseLastRevision)
	if err != nil {
		if errors.Is(err, domain.ErrBannerNotFound) {
			c.Status(http.StatusNotFound)
			return
		}
		r.l.Error("can't process user banner query", zap.Error(err))
		msg := domain.ErrInternalServerError.Error()
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Error: &msg,
		})
		return
	}
	c.JSON(http.StatusOK, content)
}

func (r *Router) DeleteBanner(c *gin.Context, params api.DeleteBannerParams) {
	if !r.adminVerify(c, params.Token) {
		return
	}

	if params.TagId == nil && params.FeatureId == nil {
		msg := domain.ErrTagOrFeatureRequired.Error()
		c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Error: &msg,
		})
		return
	}

	err := r.s.Delete(c, params.TagId, params.FeatureId)
	if err != nil {
		msg := err.Error()
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Error: &msg,
		})
		return
	}
	c.Status(http.StatusAccepted)
}

func (r *Router) PostAuthLogin(c *gin.Context) {
	var requestBody api.UserRequestBody

	if err := c.Bind(&requestBody); err != nil {
		r.l.Error("can't bind request body", zap.Error(err))
		msg := domain.ErrInternalServerError.Error()
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Error: &msg,
		})
		return
	}

	token, err := r.a.Login(c, requestBody.Login, requestBody.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			msg := err.Error()
			c.JSON(http.StatusBadRequest, api.ErrorResponse{
				Error: &msg,
			})
			return
		}
		msg := domain.ErrInternalServerError.Error()
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Error: &msg,
		})
		return
	}

	c.JSON(http.StatusOK, api.PostAuthLoginResponse{
		Token: token,
	})
}

func (r *Router) PostAuthRegister(c *gin.Context) {
	var requestBody api.UserRequestBody
	if err := c.Bind(&requestBody); err != nil {
		r.l.Error("can't bind request body", zap.Error(err))
		msg := domain.ErrInternalServerError.Error()
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Error: &msg,
		})
		return
	}
	if err := validateUserRequestBody(requestBody); err != nil {
		msg := err.Error()
		c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Error: &msg,
		})
		return
	}

	err := r.a.RegisterNewUser(c, requestBody.Login, requestBody.Password)
	if err != nil {
		if errors.Is(err, domain.ErrUserExists) {
			msg := err.Error()
			c.JSON(http.StatusConflict, api.ErrorResponse{
				Error: &msg,
			})
			return
		}
		msg := domain.ErrInternalServerError.Error()
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Error: &msg,
		})
		return
	}
	c.Status(http.StatusCreated)
}

func (r *Router) adminVerify(c *gin.Context, token *string) bool {
	if token == nil || *token == "" {
		c.Status(http.StatusUnauthorized)
		return false
	}
	if admin, err := r.a.IsAdmin(c, *token); err != nil || !admin {
		c.Status(http.StatusForbidden)
		return false
	}
	return true
}

func validateUserRequestBody(body api.UserRequestBody) error {
	if body.Login == "" || body.Password == "" {
		return domain.ErrUserLoginPassword
	}
	return nil
}
