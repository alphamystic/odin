package domain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"


	"github.com/alphamystic/odin/lib/utils"
	dfn "github.com/alphamystic/odin/lib/definers"
)


const (
	// ENUM TYPES
	createEnumTypeStmt = `
		INSERT INTO blog_enum_types (name, description, created_at)
		VALUES (?, ?, NOW());
		`

	listEnumTypesStmt = `
		SELECT id, name, description, created_at
		FROM blog_enum_types
		ORDER BY name ASC
		LIMIT ? OFFSET ?;
		`

	// CATEGORIES
	createCategoryStmt = `
		INSERT INTO blog_categories (category_name, created_at)
		VALUES (?, NOW());
		`

	listCategoriesStmt = `
		SELECT id, category_name, created_at
		FROM blog_categories
		ORDER BY category_name ASC
		LIMIT ? OFFSET ?;
		`

	// Search blogs by tag using JSON_CONTAINS
	searchByTagStmt = `
		SELECT id, uuid, ownerid, title, author, content, maintag, types, categories, subcats, tags,
		       public, archived, created_at, updated_at
		FROM blogs
		WHERE JSON_CONTAINS(tags, CAST(? AS JSON))
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?;
		`
)


// CreateEnumType inserts an enum type (e.g. "technology", "security", "tutorial").
func (d *Domain) CreateEnumType(ctx context.Context, name, description string) error {
	conn, err := d.GetConnection(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	res, err := conn.ExecContext(ctx, createEnumTypeStmt, name, description)
	if err != nil {
		return fmt.Errorf("error creating enum type: %w", err)
	}
	rowsAff, err := res.RowsAffected()
	if rowsAff != 1 || err != nil{
		d.LogToFile(utils.Logger{Name: "blog_sql", Text: fmt.Sprintf("Create EnumType: unexpected rows affected = %d. ERROR: %s", rowsAff, err)})
		return errors.New("server encountered an error while creating EnumType")
	}
	return nil
}

// ListEnumTypes returns paginated enum types.
func (d *Domain) ListEnumTypes(ctx context.Context, limit, offset int) ([]dfn.EnumType, error) {
	conn, err := d.GetConnection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rows, err := conn.QueryContext(ctx, listEnumTypesStmt, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error listing enum types: %w", err)
	}
	defer rows.Close()

	var enums []dfn.EnumType
	for rows.Next() {
		var e dfn.EnumType
		if err := rows.Scan(&e.ID, &e.Type, &e.Description); err != nil {
			continue
		}
		enums = append(enums, e)
	}

	return enums, nil
}

// CreateCategory creates a blog category (e.g. "Cybersecurity", "DevOps")
func (d *Domain) CreateCategory(ctx context.Context, categoryName string) error {
	conn, err := d.GetConnection(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	res, err := conn.ExecContext(ctx, createCategoryStmt, categoryName)
	if err != nil {
		return fmt.Errorf("error creating category: %w", err)
	}
	rowsAff, err := res.RowsAffected()
	if rowsAff != 1 || err != nil{
		d.LogToFile(utils.Logger{Name: "blog_sql", Text: fmt.Sprintf("Create category: unexpected rows affected = %d. ERROR: %s", rowsAff, err)})
		return errors.New("server encountered an error while creating category")
	}
	return nil
}

// ListCategories lists categories with pagination.
func (d *Domain) ListCategories(ctx context.Context, limit, offset int) ([]dfn.Category, error) {
	conn, err := d.GetConnection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rows, err := conn.QueryContext(ctx, listCategoriesStmt, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error listing categories: %w", err)
	}
	defer rows.Close()

	var cats []dfn.Category
	for rows.Next() {
		var c dfn.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.ParentID); err != nil {
			continue
		}
		cats = append(cats, c)
	}

	return cats, nil
}


// SearchBlogsByTag uses JSON_CONTAINS to find blogs whose tags array includes the given tag.
func (d *Domain) SearchBlogsByTag(ctx context.Context, tag string, limit, offset int) ([]dfn.Blog, error) {
	conn, err := d.GetConnection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// tag must be JSON encoded: e.g. `"security"`
	tagJSON, _ := json.Marshal(tag)

	rows, err := conn.QueryContext(ctx, searchByTagStmt, string(tagJSON), limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error searching blogs by tag: %w", err)
	}
	defer rows.Close()

	var blogs []dfn.Blog
	for rows.Next() {
		var b dfn.Blog
		if err := ScanBlogRow(rows, &b); err != nil {
			continue
		}
		blogs = append(blogs, b)
	}

	return blogs, nil
}
