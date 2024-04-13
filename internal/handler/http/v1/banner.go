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
	Delete(ctx context.Context, id int) error
}

type Router struct {
	l *zap.Logger
	s BannerService
}

func New(s BannerService) *Router {
	return &Router{
		s: s,
	}
}

func (r *Router) GetBanner(c *gin.Context, params api.GetBannerParams) {
	if us := c.GetString(domain.UserStatusHeader); us != domain.Admin {
		c.Status(http.StatusForbidden)
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
	if us := c.GetString(domain.UserStatusHeader); us != domain.Admin {
		c.Status(http.StatusForbidden)
		return
	}

	var requestBody api.PostBannerRequestBody
	if err := c.Bind(&requestBody); err != nil {
		msg := err.Error()
		c.JSON(http.StatusBadRequest, api.ErrorResponse{Error: &msg})
		return
	}

	if requestBody.TagIds == nil || requestBody.FeatureId == nil {
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
	if us := c.GetString(domain.UserStatusHeader); us != domain.Admin {
		c.Status(http.StatusForbidden)
		return
	}

	err := r.s.Delete(c, id)

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
	if us := c.GetString(domain.UserStatusHeader); us != domain.Admin {
		c.Status(http.StatusForbidden)
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
	if us := c.GetString(domain.UserStatusHeader); us != domain.User && us != domain.Admin {
		c.Status(http.StatusForbidden)
		return
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
