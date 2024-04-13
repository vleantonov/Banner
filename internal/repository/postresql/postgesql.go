package postresql

import (
	"banner/internal/domain"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/pgtype"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"strings"
)

type Banner struct {
	ID        pgtype.Int8      `db:"id"`
	TagIDs    pgtype.Int8Array `db:"tag_ids"`
	FeatureID pgtype.Int8      `db:"feature_id"`
	Content   pgtype.JSONB     `db:"content"`
	Active    pgtype.Bool      `db:"is_active"`

	CreatedAT pgtype.Timestamp `db:"created_at"`
	UpdatedAT pgtype.Timestamp `db:"updated_at"`
}

func (b *Banner) mustConvertToStd() domain.Banner {

	tags := make([]int, 0)
	for _, tag := range b.TagIDs.Elements {
		tags = append(tags, int(tag.Int))
	}

	return domain.Banner{
		ID:        int(b.ID.Int),
		Tags:      tags,
		Feature:   int(b.FeatureID.Int),
		Content:   b.Content.Get().(map[string]interface{}),
		Active:    b.Active.Bool,
		CreatedAT: b.CreatedAT.Time.String(),
		UpdatedAT: b.UpdatedAT.Time.String(),
	}
}

type PostgresRepo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *PostgresRepo {
	return &PostgresRepo{
		db: db,
	}
}

func (p *PostgresRepo) GetActiveContentByTagFeatureID(
	ctx context.Context, tagId, featureId int,
) (*map[string]interface{}, error) {

	const getBannerQuery = `
		SELECT id, content FROM Banners JOIN tag_feature_banners
		ON Banners.id = tag_feature_banners.banner_id
		WHERE tag_feature_banners.tag_id = $1
		AND tag_feature_banners.feature_id = $2
		AND is_active;
	`

	var banner Banner
	if err := p.db.GetContext(ctx, &banner, getBannerQuery, tagId, featureId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrBannerNotFound
		}
		return nil, fmt.Errorf("can't get banner: %w", err)
	}

	res := banner.mustConvertToStd()

	return &res.Content, nil
}

func (p *PostgresRepo) GetByFilter(ctx context.Context, t domain.FilterBanner) (*[]domain.Banner, error) {

	filters := make([]string, 0)
	args := make([]interface{}, 0)

	if t.TagID != nil {
		args = append(args, *t.TagID)
		filters = append(filters, fmt.Sprintf("AND tag_id=$%d", len(args)))
	}
	if t.FeatureID != nil {
		args = append(args, *t.FeatureID)
		filters = append(filters, fmt.Sprintf("AND feature_id=$%d", len(args)))
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
		SELECT id, content, is_active,
		       tfb.feature_id, 
		       array_agg(tfb.tag_id) as tag_ids,
		       created_at, updated_at
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

	banners := make([]Banner, 0)
	if err := p.db.SelectContext(ctx, &banners, getBannersQuery, args...); err != nil {
		return nil, err
	}

	res := make([]domain.Banner, 0)
	for _, b := range banners {
		cnv := b.mustConvertToStd()
		res = append(res, cnv)
	}

	if len(banners) == 0 {
		return nil, domain.ErrBannerNotFound
	}

	return &res, nil
}

func (p *PostgresRepo) Insert(ctx context.Context, b domain.Banner) (int, error) {

	const insertBannerQuery = `INSERT INTO banners (content, is_active) VALUES ($1, $2) RETURNING id`

	tx, err := p.db.Beginx()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	cJson, err := json.Marshal(b.Content)
	if err != nil {
		return 0, nil
	}

	row, err := p.db.Queryx(insertBannerQuery, cJson, b.Active)
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

	tagFeatureBanners := make([]domain.TagFeatureBanner, 0)
	for _, tag := range b.Tags {
		tagFeatureBanners = append(tagFeatureBanners, domain.TagFeatureBanner{
			TagID:     tag,
			FeatureID: b.Feature,
			BannerID:  insId,
		})
	}

	_, err = tx.NamedExecContext(ctx, insertBannerTagFeatureQuery, tagFeatureBanners)

	var e *pgconn.PgError
	if err != nil {
		if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
			return 0, domain.ErrTagFeatureAlreadyExists
		}
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return insId, nil
}

func (p *PostgresRepo) Update(ctx context.Context, u domain.UpdBanner) error {

	tx := p.db.MustBegin()
	defer tx.Rollback()

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

	updateBannerInfoMap := make(map[string]interface{})
	if u.Content != nil {
		cJson, err := json.Marshal(u.Content)
		if err != nil {
			return err
		}
		updateBannerInfoMap["content"] = cJson
	}
	if u.Active != nil {
		updateBannerInfoMap["is_active"] = *u.Active
	}
	if err := p.updateBannerInfoByMap(ctx, tx, u.ID, updateBannerInfoMap); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (p *PostgresRepo) updateBannerInfoByMap(
	ctx context.Context, tx *sqlx.Tx,
	id int,
	m map[string]interface{},
) error {
	const updateBannerQueryTmpl = `
		UPDATE banners 
		SET %s
		WHERE id=:id
	`

	delete(m, "id")
	setVal := make([]string, 0)
	for key := range m {
		setVal = append(setVal, fmt.Sprintf("%s=:%s", key, key))
	}
	m["id"] = id

	query := fmt.Sprintf(updateBannerQueryTmpl, strings.Join(setVal, ", "))

	res, err := tx.NamedExecContext(ctx, query, m)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return domain.ErrBannerNotFound
	}

	return nil
}

func (p *PostgresRepo) updateFeatureID(ctx context.Context, tx *sqlx.Tx, bannerID, featureID int) error {
	const updateFeatureIDQuery = `
		UPDATE tag_feature_banners
		SET feature_id=$1
		WHERE banner_id=$2
	`

	res, err := tx.ExecContext(ctx, updateFeatureIDQuery, featureID, bannerID)
	var e *pgconn.PgError
	if err != nil {
		if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
			return domain.ErrTagFeatureAlreadyExists
		}
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return domain.ErrBannerNotFound
	}
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
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrBannerNotFound
		}
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

	res, err := tx.NamedExecContext(ctx, insertTagsQuery, newRows)

	var e *pgconn.PgError
	if err != nil {
		if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
			return domain.ErrTagFeatureAlreadyExists
		}
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return domain.ErrBannerNotFound
	}

	return nil
}

func (p *PostgresRepo) DeleteByID(ctx context.Context, id int) error {

	const deleteBannerQuery = `DELETE FROM banners WHERE id=$1`

	res, err := p.db.ExecContext(ctx, deleteBannerQuery, id)
	n, err := res.RowsAffected()
	if n == 0 {
		return domain.ErrBannerNotFound
	}
	if err != nil {
		return err
	}
	return nil
}
