package domain

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	//"time"

	"github.com/alphamystic/odin/lib/utils"
	dfn "github.com/alphamystic/odin/lib/definers"
)

const (
	createBlogsStmt = `
    INSERT INTO blogs
    (uuid, ownerid, title, author, content, maintag, types, categories, subcats, tags, public, archived, created_at, updated_at)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW());
    `

	getRecentBlogsStmt = `
    SELECT id, uuid, ownerid, title, author, content, maintag, types, categories, subcats, tags, public, archived, created_at, updated_at
    FROM blogs WHERE ownerid = ? ORDER BY created_at DESC LIMIT 10;`

	getBlogByUUIDStmt = `
    SELECT id, uuid, ownerid, title, author, content, maintag, types, categories, subcats, tags, public, archived, created_at, updated_at
    FROM blogs
    WHERE uuid = ?;
    `

	listBlogsByAuthorStmt = `
    SELECT id, uuid, ownerid, title, author, content, maintag, types, categories, subcats, tags, public, archived, created_at, updated_at
    FROM blogs
    WHERE author = ?
    ORDER BY updated_at DESC LIMIT ? OFFSET ?;
    `

	updateBlogStmt = `
    UPDATE blogs
    SET title = ?, content = ?, maintag = ?, types = ?, categories = ?, subcats = ?, tags = ?, public = ?, archived = ?, updated_at = NOW()
    WHERE uuid = ?;
    `

	listBlogByMainTagStmt = `
    SELECT id, uuid, ownerid, title, author, content, maintag, types, categories, subcats, tags, public, archived, created_at, updated_at
    FROM blogs
    WHERE maintag = ?
    ORDER BY updated_at DESC LIMIT ? OFFSET ?;
    `

	// Query blogs that have a particular type in JSON array (MySQL JSON_CONTAINS)
	listByTypeStmt = `
    SELECT id, uuid, ownerid, title, author, content, maintag, types, categories, subcats, tags, public, archived, created_at, updated_at
    FROM blogs
    WHERE JSON_CONTAINS(types, CAST(? AS JSON)) LIMIT ? OFFSET ?;
    `

	// Query blogs that include category id in categories JSON array
	listByCategoryStmt = `
    SELECT id, uuid, ownerid, title, author, content, maintag, types, categories, subcats, tags, public, archived, created_at, updated_at
    FROM blogs
    WHERE JSON_CONTAINS(categories, CAST(? AS JSON)) LIMIT ? OFFSET ?;
    `

	archiveOrNoneArchiveStmt = `UPDATE blogs SET archived = ?, updated_at = NOW() WHERE uuid = ?;`
	listArchivedOrNoneArchivedStmt = `
    SELECT id, uuid, ownerid, title, author, content, maintag, types, categories, subcats, tags, public, archived, created_at, updated_at
    FROM blogs
    WHERE archived = ? LIMIT ? OFFSET ?;
    `

	createCommentStmt = `
    INSERT INTO comments (comment_uuid, blog_uuid, comment, commentor, created_at, updated_at)
    VALUES (?, ?, ?, ?, NOW(), NOW());
    `

	listCommentsByBlogUUIDStmt = `
    SELECT comment_uuid, blog_uuid, comment, commentor, created_at, updated_at
    FROM comments
    WHERE blog_uuid = ?
    ORDER BY created_at ASC LIMIT ? OFFSET ?;
    `
)

// helper: scan a row into dfn.Blog
func ScanBlogRow(rows *sql.Rows, b *dfn.Blog) error {
	var (
		typesRaw      sql.NullString
		categoriesRaw sql.NullString
		subcatsRaw    sql.NullString
		tagsRaw       sql.NullString
		createdAt     sql.RawBytes
		updatedAt     sql.RawBytes
		public        sql.NullBool
		archived      sql.NullBool
		ownerID       sql.NullString
		uuidStr       sql.NullString
	)

	err := rows.Scan(
		&b.ID,
		&uuidStr,
		&ownerID,
		&b.Title,
		&b.Author,
		&b.Content,
		&b.MainTag,
		&typesRaw,
		&categoriesRaw,
		&subcatsRaw,
		&tagsRaw,
		&public,
		&archived,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return err
	}

	if uuidStr.Valid {
		b.UUID = uuidStr.String
	}
	if ownerID.Valid {
		b.OwnerID = ownerID.String
	}
	if public.Valid {
		b.Public = public.Bool
	}
	if archived.Valid {
		b.Archived = archived.Bool
	}

	// Parse JSON fields into slices
	if typesRaw.Valid && typesRaw.String != "" {
		var types []string
		if err := json.Unmarshal([]byte(typesRaw.String), &types); err == nil {
			b.Types = types
		} else {
			// fallback: try to treat as comma list
			_ = err
		}
	}
	if categoriesRaw.Valid && categoriesRaw.String != "" {
		var cats []int
		if err := json.Unmarshal([]byte(categoriesRaw.String), &cats); err == nil {
			b.Categories = cats
		}
	}
	if subcatsRaw.Valid && subcatsRaw.String != "" {
		var sc []int
		if err := json.Unmarshal([]byte(subcatsRaw.String), &sc); err == nil {
			b.SubCats = sc
		}
	}
	if tagsRaw.Valid && tagsRaw.String != "" {
		var tags []string
		if err := json.Unmarshal([]byte(tagsRaw.String), &tags); err == nil {
			b.Tags = tags
		}
	}

	// timestamps
	if err := utils.ScanTimeStamps(&b.TimeStamps, createdAt, updatedAt); err != nil {
		// not fatal; just log inside utils if needed
		_ = err
	}

	return nil
}

// CreateBlog inserts a new blog. If blog.UUID is empty, a new UUID is generated.
func (d *Domain) CreateBlog(ctx context.Context, blog dfn.Blog) error {
	conn, err := d.GetConnection(ctx)
	if err != nil {
		return fmt.Errorf("error getting db connection: %w", err)
	}
	defer conn.Close()

	// ensure UUID
	if blog.UUID == "" {
		blog.UUID = utils.GenerateUUID()
	}

	// marshal json fields
	typesJSON, _ := json.Marshal(blog.Types)
	categoriesJSON, _ := json.Marshal(blog.Categories)
	subcatsJSON, _ := json.Marshal(blog.SubCats)
	tagsJSON, _ := json.Marshal(blog.Tags)

	stmt, err := conn.PrepareContext(ctx, createBlogsStmt)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "blog_sql", Text: fmt.Sprintf("Error preparing blog insert: %s", err)})
		return errors.New("error preparing to create blog")
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx,
		blog.UUID,
		blog.OwnerID,
		blog.Title,
		blog.Author,
		blog.Content,
		blog.MainTag,
		string(typesJSON),
		string(categoriesJSON),
		string(subcatsJSON),
		string(tagsJSON),
		blog.Public,
		blog.Archived,
	)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "blog_sql", Text: fmt.Sprintf("Error executing blog insert: %s", err)})
		return errors.New("error creating blog entry")
	}

	rowsAff, err := res.RowsAffected()
	if rowsAff != 1 || err != nil{
		d.LogToFile(utils.Logger{Name: "blog_sql", Text: fmt.Sprintf("Create blog: unexpected rows affected = %d. ERROR: %s", rowsAff, err)})
		return errors.New("server encountered an error while creating blog")
	}

	return nil
}

// GetRecentBlogs returns the most recent blogs (limit 10).
func (d *Domain) GetRecentBlogs(ctx context.Context, ownerid string) ([]dfn.Blog, error) {
	conn, err := d.GetConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting db connection: %w", err)
	}
	defer conn.Close()

	rows, err := conn.QueryContext(ctx, getRecentBlogsStmt, ownerid)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "blog_sql", Text: fmt.Sprintf("Error retrieving recent blogs: %s", err)})
		return nil, errors.New("error retrieving recent blogs")
	}
	defer rows.Close()

	var blogs []dfn.Blog
	for rows.Next() {
		var b dfn.Blog
		if err := ScanBlogRow(rows, &b); err != nil {
			d.LogToFile(utils.Logger{Name: "blog_sql", Text: fmt.Sprintf("Error scanning recent blog rows: %s", err)})
			continue
		}
		blogs = append(blogs, b)
	}
	return blogs, nil
}

// UpdateBlog updates a blog identified by UUID.
func (d *Domain) UpdateBlog(ctx context.Context, blog dfn.Blog) error {
	conn, err := d.GetConnection(ctx)
	if err != nil {
		return fmt.Errorf("error getting db connection: %w", err)
	}
	defer conn.Close()

	typesJSON, _ := json.Marshal(blog.Types)
	categoriesJSON, _ := json.Marshal(blog.Categories)
	subcatsJSON, _ := json.Marshal(blog.SubCats)
	tagsJSON, _ := json.Marshal(blog.Tags)

	stmt, err := conn.PrepareContext(ctx, updateBlogStmt)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "blog_sql", Text: fmt.Sprintf("Error preparing update blog: %s", err)})
		return errors.New("error preparing to update blog")
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx,
		blog.Title,
		blog.Content,
		blog.MainTag,
		string(typesJSON),
		string(categoriesJSON),
		string(subcatsJSON),
		string(tagsJSON),
		blog.Public,
		blog.Archived,
		blog.UUID,
	)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "blog_sql", Text: fmt.Sprintf("Error executing blog update: %s", err)})
		return errors.New("error updating blog entry")
	}

	rowsAff, err := res.RowsAffected()
	if rowsAff != 1 || err != nil{
		d.LogToFile(utils.Logger{Name: "blog_sql", Text: fmt.Sprintf("Update blog: unexpected rows affected = %d", rowsAff)})
		// Not necessarily an error if no fields changed; we treat as success though you can enforce otherwise
	}

	return nil
}

// ViewBlog returns a blog by UUID.
func (d *Domain) ViewBlog(ctx context.Context, blogUUID string) (*dfn.Blog, error) {
	conn, err := d.GetConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting db connection: %w", err)
	}
	defer conn.Close()

	// We will scan row using a temporary rows-like approach by using QueryContext
	// Simpler: use QueryContext with the same query and call ScanBlogRow over rows.Next()
	rows, err := conn.QueryContext(ctx, getBlogByUUIDStmt, blogUUID)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "blog_sql", Text: fmt.Sprintf("Error querying blog uuid %s: %s", blogUUID, err)})
		return nil, errors.New("error retrieving blog")
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, errors.New("blog not found")
	}
	var b dfn.Blog
	if err := ScanBlogRow(rows, &b); err != nil {
		d.LogToFile(utils.Logger{Name: "blog_sql", Text: fmt.Sprintf("Error scanning blog row for uuid %s: %s", blogUUID, err)})
		return nil, errors.New("error parsing blog")
	}
	return &b, nil
}

// ListBlogByAuthor returns blogs authored by a given author.
func (d *Domain) ListBlogByAuthor(ctx context.Context, authorName string, limit, offset int) ([]dfn.Blog, error) {
	conn, err := d.GetConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting db connection: %w", err)
	}
	defer conn.Close()

	rows, err := conn.QueryContext(ctx, listBlogsByAuthorStmt, authorName, limit, offset)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "blog_sql", Text: fmt.Sprintf("Error listing blogs by author %s: %s", authorName, err)})
		return nil, errors.New("error retrieving blogs by author")
	}
	defer rows.Close()

	var blogs []dfn.Blog
	for rows.Next() {
		var b dfn.Blog
		if err := ScanBlogRow(rows, &b); err != nil {
			d.LogToFile(utils.Logger{Name: "blog_sql", Text: fmt.Sprintf("Error scanning blog by author: %s", err)})
			continue
		}
		blogs = append(blogs, b)
	}
	return blogs, nil
}

// ListBlogByMainTag returns blogs matching maintag.
func (d *Domain) ListBlogByMainTag(ctx context.Context, tag string, limit, offset int) ([]dfn.Blog, error) {
	conn, err := d.GetConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting db connection: %w", err)
	}
	defer conn.Close()

	rows, err := conn.QueryContext(ctx, listBlogByMainTagStmt, tag, limit, offset)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "blog_sql", Text: fmt.Sprintf("Error listing blogs by maintag %s: %s", tag, err)})
		return nil, errors.New("error retrieving blogs by tag")
	}
	defer rows.Close()

	var blogs []dfn.Blog
	for rows.Next() {
		var b dfn.Blog
		if err := ScanBlogRow(rows, &b); err != nil {
			d.LogToFile(utils.Logger{Name: "blog_sql", Text: fmt.Sprintf("Error scanning blog by tag: %s", err)})
			continue
		}
		blogs = append(blogs, b)
	}
	return blogs, nil
}

// ListByType returns blogs that include the requested type in their types JSON array.
// example typeArg should be a JSON string like: "\"framework\"" OR simply "framework" - we'll marshal.
func (d *Domain) ListByType(ctx context.Context, typeArg string, limit, offset int) ([]dfn.Blog, error) {
	conn, err := d.GetConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting db connection: %w", err)
	}
	defer conn.Close()

	// ensure a JSON string is passed to JSON_CONTAINS:
	typeJSON, _ := json.Marshal(typeArg)

	rows, err := conn.QueryContext(ctx, listByTypeStmt, string(typeJSON), limit, offset)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "blog_sql", Text: fmt.Sprintf("Error listing blogs by type %s: %s", typeArg, err)})
		return nil, errors.New("error retrieving blogs by type")
	}
	defer rows.Close()

	var blogs []dfn.Blog
	for rows.Next() {
		var b dfn.Blog
		if err := ScanBlogRow(rows, &b); err != nil {
			d.LogToFile(utils.Logger{Name: "blog_sql", Text: fmt.Sprintf("Error scanning blog by type: %s", err)})
			continue
		}
		blogs = append(blogs, b)
	}
	return blogs, nil
}

// ListByCategory returns blogs that include the requested category ID in their categories JSON array.
func (d *Domain) ListByCategory(ctx context.Context, categoryID, limit, offset int) ([]dfn.Blog, error) {
	conn, err := d.GetConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting db connection: %w", err)
	}
	defer conn.Close()

	catJSON, _ := json.Marshal(categoryID)

	rows, err := conn.QueryContext(ctx, listByCategoryStmt, string(catJSON), limit, offset)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "blog_sql", Text: fmt.Sprintf("Error listing blogs by category %d: %s", categoryID, err)})
		return nil, errors.New("error retrieving blogs by category")
	}
	defer rows.Close()

	var blogs []dfn.Blog
	for rows.Next() {
		var b dfn.Blog
		if err := ScanBlogRow(rows, &b); err != nil {
			d.LogToFile(utils.Logger{Name: "blog_sql", Text: fmt.Sprintf("Error scanning blog by category: %s", err)})
			continue
		}
		blogs = append(blogs, b)
	}
	return blogs, nil
}

// ListArchivedOrNoneArchivedBlog returns blogs filtered by archived flag.
func (d *Domain) ListArchivedOrNoneArchivedBlog(ctx context.Context, archived bool, limit, offset int) ([]dfn.Blog, error) {
	conn, err := d.GetConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting db connection: %w", err)
	}
	defer conn.Close()

	rows, err := conn.QueryContext(ctx, listArchivedOrNoneArchivedStmt, archived, limit, offset)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "blog_sql", Text: fmt.Sprintf("Error listing archived blogs: %s", err)})
		return nil, errors.New("error retrieving archived/non-archived blogs")
	}
	defer rows.Close()

	var blogs []dfn.Blog
	for rows.Next() {
		var b dfn.Blog
		if err := ScanBlogRow(rows, &b); err != nil {
			d.LogToFile(utils.Logger{Name: "blog_sql", Text: fmt.Sprintf("Error scanning archived blog: %s", err)})
			continue
		}
		blogs = append(blogs, b)
	}
	return blogs, nil
}

// ArchiveOrRemoveFromArchive sets archived flag for a blog UUID.
func (d *Domain) ArchiveOrRemoveFromArchive(ctx context.Context, blogUUID string, archive bool) error {
	conn, err := d.GetConnection(ctx)
	if err != nil {
		return fmt.Errorf("error getting db connection: %w", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx, archiveOrNoneArchiveStmt)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "blog_sql", Text: fmt.Sprintf("Error preparing archive update for blog %s: %s", blogUUID, err)})
		return errors.New("error preparing to update blog archive status")
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, archive, blogUUID)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "blog_sql", Text: fmt.Sprintf("Error updating archive status for blog %s: %s", blogUUID, err)})
		return errors.New("error updating blog archive status")
	}
	rowsAff, _ := res.RowsAffected()
	if rowsAff != 1 {
		// If 0 rows affected, maybe blog not found. Return meaningful error.
		return errors.New("no blog updated (uuid not found or no change)")
	}
	return nil
}

// CreateComment inserts a comment. If comment.CommentUUID empty, generate one.
func (d *Domain) CreateComment(ctx context.Context, comment dfn.Comment) error {
	conn, err := d.GetConnection(ctx)
	if err != nil {
		return fmt.Errorf("error getting db connection: %w", err)
	}
	defer conn.Close()

	if comment.CommentUUID == "" {
		comment.CommentUUID = utils.GenerateUUID()
	}

	stmt, err := conn.PrepareContext(ctx, createCommentStmt)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "comment_sql", Text: fmt.Sprintf("Error preparing create comment: %s", err)})
		return errors.New("error preparing comment insert")
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, comment.CommentUUID, comment.BlogUUID, comment.Comment, comment.Commentor)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "comment_sql", Text: fmt.Sprintf("Error executing create comment for blog %s: %s", comment.BlogUUID, err)})
		return errors.New("error creating comment")
	}
	rowsAff, err := res.RowsAffected()
	if rowsAff != 1 || err != nil {
		d.LogToFile(utils.Logger{Name: "comment_sql", Text: fmt.Sprintf("Create comment: unexpected rows affected = %d. ERROR: %s", rowsAff, err)})
		return errors.New("server encountered an error while creating comment")
	}
	return nil
}

// ListComments returns comments for a blog UUID.
func (d *Domain) ListComments(ctx context.Context, blogUUID string, limit, offset int) ([]dfn.Comment, error) {
	conn, err := d.GetConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting db connection: %w", err)
	}
	defer conn.Close()

	rows, err := conn.QueryContext(ctx, listCommentsByBlogUUIDStmt, blogUUID, limit, offset)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "comment_sql", Text: fmt.Sprintf("Error listing comments for blog %s: %s", blogUUID, err)})
		return nil, errors.New("error retrieving comments")
	}
	defer rows.Close()

	var comments []dfn.Comment
	for rows.Next() {
		var c dfn.Comment
		var createdAt, updatedAt sql.NullTime
		if err := rows.Scan(&c.CommentUUID, &c.BlogUUID, &c.Comment, &c.Commentor, &createdAt, &updatedAt); err != nil {
			d.LogToFile(utils.Logger{Name: "comment_sql", Text: fmt.Sprintf("Error scanning comment: %s", err)})
			continue
		}

		// convert to utils.TimeStamps if needed; assuming dfn.Comment embeds TimeStamps
		if createdAt.Valid {
			c.CreatedAt = createdAt.Time
		}
		if updatedAt.Valid {
			c.UpdatedAt = updatedAt.Time
		}
		comments = append(comments, c)
	}
	return comments, nil
}
