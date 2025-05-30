// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: collection_items.sql

package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	gofeed "github.com/mmcdole/gofeed"
	typeext "github.com/rhajizada/gazette/internal/typeext"
)

const addItemToCollection = `-- name: AddItemToCollection :one
INSERT INTO collection_items (collection_id, item_id)
VALUES ($1, $2)
RETURNING collection_id, item_id, added_at
`

type AddItemToCollectionParams struct {
	CollectionID uuid.UUID `json:"collectionId"`
	ItemID       uuid.UUID `json:"itemId"`
}

func (q *Queries) AddItemToCollection(ctx context.Context, arg AddItemToCollectionParams) (CollectionItem, error) {
	row := q.db.QueryRow(ctx, addItemToCollection, arg.CollectionID, arg.ItemID)
	var i CollectionItem
	err := row.Scan(&i.CollectionID, &i.ItemID, &i.AddedAt)
	return i, err
}

const countItemsInCollection = `-- name: CountItemsInCollection :one
SELECT COUNT(*) AS count
FROM collection_items
WHERE collection_id = $1
`

func (q *Queries) CountItemsInCollection(ctx context.Context, collectionID uuid.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, countItemsInCollection, collectionID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getCollectionItem = `-- name: GetCollectionItem :one
SELECT collection_id, item_id, added_at
FROM collection_items
WHERE collection_id = $1
  AND item_id       = $2
`

type GetCollectionItemParams struct {
	CollectionID uuid.UUID `json:"collectionId"`
	ItemID       uuid.UUID `json:"itemId"`
}

func (q *Queries) GetCollectionItem(ctx context.Context, arg GetCollectionItemParams) (CollectionItem, error) {
	row := q.db.QueryRow(ctx, getCollectionItem, arg.CollectionID, arg.ItemID)
	var i CollectionItem
	err := row.Scan(&i.CollectionID, &i.ItemID, &i.AddedAt)
	return i, err
}

const listItemsInCollection = `-- name: ListItemsInCollection :many
SELECT ci.collection_id, ci.item_id, ci.added_at, i.id, i.feed_id, i.title, i.description, i.content, i.link, i.links, i.updated_parsed, i.published_parsed, i.authors, i.guid, i.image, i.categories, i.enclosures, i.created_at, i.updated_at
FROM collection_items ci
JOIN items i ON i.id = ci.item_id
WHERE ci.collection_id = $1
ORDER BY ci.added_at DESC
LIMIT  $2
OFFSET $3
`

type ListItemsInCollectionParams struct {
	CollectionID uuid.UUID `json:"collectionId"`
	Limit        int32     `json:"limit"`
	Offset       int32     `json:"offset"`
}

type ListItemsInCollectionRow struct {
	CollectionID    uuid.UUID          `json:"collectionId"`
	ItemID          uuid.UUID          `json:"itemId"`
	AddedAt         time.Time          `json:"addedAt"`
	ID              uuid.UUID          `json:"id"`
	FeedID          uuid.UUID          `json:"feedId"`
	Title           *string            `json:"title"`
	Description     *string            `json:"description"`
	Content         *string            `json:"content"`
	Link            string             `json:"link"`
	Links           []string           `json:"links"`
	UpdatedParsed   *time.Time         `json:"updatedParsed"`
	PublishedParsed *time.Time         `json:"publishedParsed"`
	Authors         typeext.Authors    `json:"authors"`
	Guid            *string            `json:"guid"`
	Image           *gofeed.Image      `json:"image"`
	Categories      []string           `json:"categories"`
	Enclosures      typeext.Enclosures `json:"enclosures"`
	CreatedAt       time.Time          `json:"createdAt"`
	UpdatedAt       time.Time          `json:"updatedAt"`
}

func (q *Queries) ListItemsInCollection(ctx context.Context, arg ListItemsInCollectionParams) ([]ListItemsInCollectionRow, error) {
	rows, err := q.db.Query(ctx, listItemsInCollection, arg.CollectionID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListItemsInCollectionRow
	for rows.Next() {
		var i ListItemsInCollectionRow
		if err := rows.Scan(
			&i.CollectionID,
			&i.ItemID,
			&i.AddedAt,
			&i.ID,
			&i.FeedID,
			&i.Title,
			&i.Description,
			&i.Content,
			&i.Link,
			&i.Links,
			&i.UpdatedParsed,
			&i.PublishedParsed,
			&i.Authors,
			&i.Guid,
			&i.Image,
			&i.Categories,
			&i.Enclosures,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const removeItemFromCollection = `-- name: RemoveItemFromCollection :exec
DELETE FROM collection_items
WHERE collection_id = $1
  AND item_id       = $2
`

type RemoveItemFromCollectionParams struct {
	CollectionID uuid.UUID `json:"collectionId"`
	ItemID       uuid.UUID `json:"itemId"`
}

func (q *Queries) RemoveItemFromCollection(ctx context.Context, arg RemoveItemFromCollectionParams) error {
	_, err := q.db.Exec(ctx, removeItemFromCollection, arg.CollectionID, arg.ItemID)
	return err
}
