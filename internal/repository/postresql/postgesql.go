package postresql

import (
	"banner/internal/models"
	repo "banner/internal/repository"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/pgtype"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"strings"
)

const (
	driverName = "pgx"
)

type Banner struct {
	ID        pgtype.Int8      `db:"id"`
	TagIDs    pgtype.Int8Array `db:"tag_ids"`
	FeatureID pgtype.Int8      `db:"feature_id"`
	Content   pgtype.JSONB     `db:"content"`
}

func (b *Banner) mustConvertToStd() models.Banner {

	tags := make([]int, 0)
	for _, tag := range b.TagIDs.Elements {
		tags = append(tags, int(tag.Int))
	}

	return models.Banner{
		ID:      int(b.ID.Int),
		Tags:    tags,
		Feature: int(b.FeatureID.Int),
		Content: b.Content.Get(),
	}
}

type PostgresRepo struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func New(host, port, username, password, dbName, sslMode string, l *zap.Logger) (*PostgresRepo, error) {

	// TODO: Use pgpool with wrapped logger
	db, err := sqlx.Connect(driverName, fetchSourcePath(host, port, username, password, dbName, sslMode))
	if err != nil {
		return nil, fmt.Errorf("can't connect to postgres database: %w", err)
	}

	return &PostgresRepo{
		db:     db,
		logger: l,
	}, nil
}

func fetchSourcePath(host, port, username, password, dbName, sslMode string) string {
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		host, port, dbName, username, password, sslMode,
	)
}

func (p *PostgresRepo) GetBanner(ctx context.Context, tagId, featureId int) (*models.Banner, error) {

	const getBannerQuery = `
		SELECT id, content FROM Banners JOIN tag_feature_banners
		ON Banners.id = tag_feature_banners.banner_id
		WHERE tag_feature_banners.tag_id = $1
		AND tag_feature_banners.feature_id = $2;
	`

	log := p.logger.With(
		zap.Int("tag_id", tagId),
		zap.Int("feature_id", featureId),
	)
	log.Info(
		"trying to get banner from database",
	)

	var banner Banner
	if err := p.db.GetContext(ctx, &banner, getBannerQuery, tagId, featureId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repo.ErrBannerNotExists
		}
		return nil, fmt.Errorf("can't get banner: %w", err)
	}
	log.Info("successfully get banner from database")

	res := banner.mustConvertToStd()
	return &res, nil
}

func (p *PostgresRepo) GetByFilterTagFeatureId(ctx context.Context, t models.FilterBanner) (*[]models.Banner, error) {

	log := p.logger
	log.Info("building get banners with filters query")

	filters := make([]string, 0)
	args := make([]interface{}, 0)

	if t.TagID != nil {
		args = append(args, *t.TagID)
		filters = append(filters, fmt.Sprintf("AND tag_id=$%d", len(args)))
		log = log.With(zap.Int("tag_id", *t.TagID))
	}
	if t.FeatureID != nil {
		args = append(args, *t.FeatureID)
		filters = append(filters, fmt.Sprintf("AND feature_id=$%d", len(args)))
		log = log.With(zap.Int("feature_id", *t.FeatureID))
	}
	queryFilters := strings.Join(filters, " ")

	paginationFilter := make([]string, 0)
	if t.Limit != nil {
		paginationFilter = append(paginationFilter, fmt.Sprintf("LIMIT %d", *t.Limit))
	}
	if t.Offset != nil {
		paginationFilter = append(paginationFilter, fmt.Sprintf("OFFSET %d", *t.Offset))
	}

	const getBannersQueryTemplate = `
		SELECT id, content, tfb.feature_id, array_agg(tfb.tag_id) as tag_ids
		FROM banners as b 
		JOIN tag_feature_banners as tfb 
		ON b.id = tfb.banner_id 
		WHERE TRUE %s 
		GROUP BY b.id, tfb.feature_id 
		ORDER BY b.id %s;
	`

	getBannersQuery := fmt.Sprintf(
		getBannersQueryTemplate,
		queryFilters,
		strings.Join(paginationFilter, " "),
	)

	log.Info(
		"trying to get banners from database with filters",
		zap.String("query", getBannersQuery),
	)

	banners := []Banner{}
	if err := p.db.SelectContext(ctx, &banners, getBannersQuery, args...); err != nil {
		return nil, err
	}

	log.Info("successfully get banners")
	res := make([]models.Banner, 0)
	for _, b := range banners {
		cnv := b.mustConvertToStd()
		res = append(res, cnv)
	}

	if len(banners) == 0 {
		return nil, repo.ErrBannerNotExists
	}

	return &res, nil
}

func (p *PostgresRepo) Insert(ctx context.Context, b models.Banner) (int, error) {
	// TODO: make duplicate tags processing with equal tag values in one field
	log := p.logger.With(zap.Int("feature_id", b.Feature), zap.Ints("tag_ids", b.Tags))
	p.logger.Info("trying to create banner")

	const insertBannerQuery = `INSERT INTO banners (content) VALUES ($1) RETURNING id`

	log.Info("insert banner into banners")

	tx, err := p.db.Beginx()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	cJson, err := json.Marshal(b.Content)
	if err != nil {
		return 0, nil
	}

	row, err := p.db.Queryx(insertBannerQuery, cJson)
	if err != nil {
		return 0, err
	}

	var insId int
	if !row.Next() {
		return 0, fmt.Errorf("TODO ERROR WITH NULL ROW")
	}

	err = row.Scan(&insId)
	if err != nil {
		return 0, err
	}
	if err := row.Close(); err != nil {
		return 0, err
	}

	const insertBannerTagFeatureQuery = `
		INSERT INTO tag_feature_banners (tag_id, feature_id, banner_id) 
		VALUES (:tag_id, :feature_id, :banner_id)
	`

	log.Info("insert tag and feature into tag_feature_banners")
	tagFeatureBanners := make([]models.TagFeatureBanner, 0)
	for _, tag := range b.Tags {
		tagFeatureBanners = append(tagFeatureBanners, models.TagFeatureBanner{
			TagID:     tag,
			FeatureID: b.Feature,
			BannerID:  insId,
		})
	}

	// TODO: Check with pointer
	_, err = tx.NamedExecContext(ctx, insertBannerTagFeatureQuery, tagFeatureBanners)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	log.Info("banner has been successfully created", zap.Int("id", insId))
	return insId, nil
}

// TODO: Get models.Banner
func (p *PostgresRepo) Update(ctx context.Context, u models.UpdBanner) error {

	log := p.logger.With(zap.Int("id", u.ID))
	log.Info("trying to update banner")

	tx := p.db.MustBegin()
	defer tx.Rollback()

	if u.Content != nil {
		if err := p.updateContent(ctx, tx, u.ID, u.Content); err != nil {
			return err
		}
	}

	if u.Tags != nil {
		if err := p.updateTagIDs(ctx, tx, u.ID, *u.Tags); err != nil {
			return err
		}
	}
	if u.Feature != nil {
		if err := p.updateFeatureID(ctx, tx, u.ID, *u.Feature); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (p *PostgresRepo) updateContent(ctx context.Context, tx *sqlx.Tx, bannerID int, content interface{}) error {
	const updateContentQuery = `
		UPDATE banners 
		SET content=$1
		WHERE id=$2
	`

	cJson, err := json.Marshal(content)
	if err != nil {
		return err
	}

	_ = tx.MustExecContext(ctx, updateContentQuery, cJson, bannerID)
	return nil
}

func (p *PostgresRepo) updateFeatureID(ctx context.Context, tx *sqlx.Tx, bannerID, featureID int) error {
	const updateFeatureIDQuery = `
		UPDATE tag_feature_banners
		SET feature_id=$1
		WHERE banner_id=$2
	`

	_ = tx.MustExecContext(ctx, updateFeatureIDQuery, featureID, bannerID)
	return nil
}

func (p *PostgresRepo) updateTagIDs(ctx context.Context, tx *sqlx.Tx, bannerID int, tagIDs []int) error {

	const deleteCurTagIDsQuery = `
		WITH f_del AS (
			DELETE FROM tag_feature_banners
			WHERE tag_feature_banners.banner_id=$1
			RETURNING tag_feature_banners.feature_id as id
		)
		SELECT DISTINCT id FROM f_del;
	`

	var featureID int
	err := tx.GetContext(ctx, &featureID, deleteCurTagIDsQuery, bannerID)
	if err != nil {
		return err
	}

	const insertTagsQuery = `INSERT INTO tag_feature_banners VALUES (:tag_id,:feature_id,:banner_id)`

	var newRows []map[string]interface{}
	for _, tagID := range tagIDs {
		row := map[string]interface{}{
			"tag_id":     tagID,
			"feature_id": featureID,
			"banner_id":  bannerID,
		}
		newRows = append(newRows, row)
	}

	_, err = tx.NamedExecContext(ctx, insertTagsQuery, newRows)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresRepo) DeleteById(ctx context.Context, id int) (int64, error) {

	p.logger.Info("trying to delete banner", zap.Int("id", id))

	const deleteBannerQuery = `DELETE FROM banners WHERE id=$1`

	res, err := p.db.ExecContext(ctx, deleteBannerQuery, id)
	if err != nil {
		return 0, err
	}

	// TODO: add error or return only integer
	return res.RowsAffected()
}
