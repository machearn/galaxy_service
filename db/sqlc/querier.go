// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2

package db

import (
	"context"
)

type Querier interface {
	CreateEntry(ctx context.Context, arg CreateEntryParams) (Entry, error)
	CreateItem(ctx context.Context, arg CreateItemParams) (Item, error)
	CreateMember(ctx context.Context, arg CreateMemberParams) (Member, error)
	DeleteEntry(ctx context.Context, id int32) error
	DeleteItem(ctx context.Context, id int32) error
	DeleteMember(ctx context.Context, id int32) error
	GetEntry(ctx context.Context, id int32) (Entry, error)
	GetItem(ctx context.Context, id int32) (Item, error)
	GetMember(ctx context.Context, id int32) (Member, error)
	GetMemberByName(ctx context.Context, username string) (Member, error)
	ListEntries(ctx context.Context, arg ListEntriesParams) ([]Entry, error)
	ListEntriesByItem(ctx context.Context, arg ListEntriesByItemParams) ([]Entry, error)
	ListEntriesByMember(ctx context.Context, arg ListEntriesByMemberParams) ([]Entry, error)
	ListItems(ctx context.Context, arg ListItemsParams) ([]Item, error)
	ListMembers(ctx context.Context, arg ListMembersParams) ([]Member, error)
	UpdateEntry(ctx context.Context, arg UpdateEntryParams) (Entry, error)
	UpdateItem(ctx context.Context, arg UpdateItemParams) (Item, error)
	UpdateMember(ctx context.Context, arg UpdateMemberParams) (Member, error)
}

var _ Querier = (*Queries)(nil)