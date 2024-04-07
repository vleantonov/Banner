package banner

import (
	api "banner/internal/api/gen"
	"banner/internal/models"
	"banner/internal/repository"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type DefaultResponse struct {
	Message string `json:"message"`
}

type Repository interface {
	GetBanner(ctx context.Context, tagId, featureId int) (*models.Banner, error)
	GetByFilterTagFeatureId(ctx context.Context, f models.FilterBanner) (*[]models.Banner, error)
	Insert(ctx context.Context, b models.Banner) (int, error)
	Update(ctx context.Context, b models.UpdBanner) error
	DeleteById(ctx context.Context, id int) (int64, error)
}

type Router struct {
	r Repository
}

func New(r Repository) *Router {
	return &Router{
		r: r,
	}
}

func (r *Router) GetBanner(c *gin.Context, params api.GetBannerParams) {
	b, err := r.r.GetByFilterTagFeatureId(c, models.FilterBanner{
		FeatureID: params.FeatureId,
		TagID:     params.TagId,
		Limit:     params.Limit,
		Offset:    params.Offset,
	})

	if err != nil {
		if errors.Is(err, repository.ErrBannerNotExists) {
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
	var requestBody api.PostBannerRequestBody

	if err := c.Bind(&requestBody); err != nil {
		msg := err.Error()
		c.JSON(http.StatusBadRequest, api.ErrorResponse{Error: &msg})
		return
	}

	// TODO: Check If tags and feature exists or both unavailable is service
	if requestBody.TagIds == nil || requestBody.FeatureId == nil {
		msg := "tag_ids and feature_id are required"
		c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Error: &msg,
		})
	}

	idxResp, err := r.r.Insert(c, models.Banner{
		Content: requestBody.Content,
		Tags:    *requestBody.TagIds,
		Feature: *requestBody.FeatureId,
	})

	if err != nil {
		msg := err.Error()
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Error: &msg,
		})
		return
	}
	c.JSON(http.StatusCreated, api.PostBannerResponse{BannerId: &idxResp})
}

func (r *Router) DeleteBannerId(c *gin.Context, id int, params api.DeleteBannerIdParams) {
	_, err := r.r.DeleteById(c, id)

	if err != nil {
		if errors.Is(err, repository.ErrBannerNotExists) {
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
	var requestBody api.PatchBannerRequestBody
	if err := c.Bind(&requestBody); err != nil {
		msg := err.Error()
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Error: &msg,
		})
		return
	}

	var cnt interface{}
	if requestBody.Content != nil {
		cnt = requestBody.Content
	}
	err := r.r.Update(c, models.UpdBanner{
		ID:      id,
		Tags:    requestBody.TagIds,
		Feature: requestBody.FeatureId,
		Content: cnt,
	})

	if err != nil {
		if errors.Is(err, repository.ErrBannerNotExists) {
			c.Status(http.StatusNotFound)
			return
		}
		msg := err.Error()
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Error: &msg,
		})
		return
	}
	c.Status(http.StatusOK)
}

func (r *Router) GetUserBanner(c *gin.Context, params api.GetUserBannerParams) {
	b, err := r.r.GetBanner(c, params.TagId, params.FeatureId)
	if err != nil {
		if errors.Is(err, repository.ErrBannerNotExists) {
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
