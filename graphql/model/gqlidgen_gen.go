// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/mikeydub/go-gallery/service/persist"
)

func (r *Collection) ID() GqlID {
	return GqlID(fmt.Sprintf("Collection:%s", r.Dbid))
}

func (r *CollectionNft) ID() GqlID {
	//-----------------------------------------------------------------------------------------------
	//-----------------------------------------------------------------------------------------------
	// Some fields specified by @goGqlId require manual binding because one of the following is true:
	// (a) the field does not exist on the CollectionNft type, or
	// (b) the field exists but is not a string type
	//-----------------------------------------------------------------------------------------------
	// Please create binding methods on the CollectionNft type with the following signatures:
	// func (r *CollectionNft) GetGqlIDField_NftID() string
	// func (r *CollectionNft) GetGqlIDField_CollectionID() string
	//-----------------------------------------------------------------------------------------------
	return GqlID(fmt.Sprintf("CollectionNft:%s:%s", r.GetGqlIDField_NftID(), r.GetGqlIDField_CollectionID()))
}

func (r *Gallery) ID() GqlID {
	return GqlID(fmt.Sprintf("Gallery:%s", r.Dbid))
}

func (r *GalleryUser) ID() GqlID {
	return GqlID(fmt.Sprintf("GalleryUser:%s", r.Dbid))
}

func (r *MembershipTier) ID() GqlID {
	return GqlID(fmt.Sprintf("MembershipTier:%s", r.Dbid))
}

func (r *Nft) ID() GqlID {
	return GqlID(fmt.Sprintf("Nft:%s", r.Dbid))
}

func (r *Wallet) ID() GqlID {
	return GqlID(fmt.Sprintf("Wallet:%s", *r.Address))
}

type NodeFetcher struct {
	OnCollection     func(ctx context.Context, dbid persist.DBID) (*Collection, error)
	OnCollectionNft  func(ctx context.Context, nftId string, collectionId string) (*CollectionNft, error)
	OnGallery        func(ctx context.Context, dbid persist.DBID) (*Gallery, error)
	OnGalleryUser    func(ctx context.Context, dbid persist.DBID) (*GalleryUser, error)
	OnMembershipTier func(ctx context.Context, dbid persist.DBID) (*MembershipTier, error)
	OnNft            func(ctx context.Context, dbid persist.DBID) (*Nft, error)
	OnWallet         func(ctx context.Context, address persist.Address) (*Wallet, error)
}

func (n *NodeFetcher) GetNodeByGqlID(ctx context.Context, id GqlID) (Node, error) {
	parts := strings.Split(string(id), ":")
	if len(parts) == 1 {
		return nil, ErrInvalidIDFormat{message: "no ID components specified after type name"}
	}

	typeName := parts[0]
	ids := parts[1:]

	switch typeName {
	case "Collection":
		if len(ids) != 1 {
			return nil, ErrInvalidIDFormat{message: fmt.Sprintf("'Collection' type requires 1 ID component(s) (%d component(s) supplied)", len(ids))}
		}
		return n.OnCollection(ctx, persist.DBID(ids[0]))
	case "CollectionNft":
		if len(ids) != 2 {
			return nil, ErrInvalidIDFormat{message: fmt.Sprintf("'CollectionNft' type requires 2 ID component(s) (%d component(s) supplied)", len(ids))}
		}
		return n.OnCollectionNft(ctx, string(ids[0]), string(ids[1]))
	case "Gallery":
		if len(ids) != 1 {
			return nil, ErrInvalidIDFormat{message: fmt.Sprintf("'Gallery' type requires 1 ID component(s) (%d component(s) supplied)", len(ids))}
		}
		return n.OnGallery(ctx, persist.DBID(ids[0]))
	case "GalleryUser":
		if len(ids) != 1 {
			return nil, ErrInvalidIDFormat{message: fmt.Sprintf("'GalleryUser' type requires 1 ID component(s) (%d component(s) supplied)", len(ids))}
		}
		return n.OnGalleryUser(ctx, persist.DBID(ids[0]))
	case "MembershipTier":
		if len(ids) != 1 {
			return nil, ErrInvalidIDFormat{message: fmt.Sprintf("'MembershipTier' type requires 1 ID component(s) (%d component(s) supplied)", len(ids))}
		}
		return n.OnMembershipTier(ctx, persist.DBID(ids[0]))
	case "Nft":
		if len(ids) != 1 {
			return nil, ErrInvalidIDFormat{message: fmt.Sprintf("'Nft' type requires 1 ID component(s) (%d component(s) supplied)", len(ids))}
		}
		return n.OnNft(ctx, persist.DBID(ids[0]))
	case "Wallet":
		if len(ids) != 1 {
			return nil, ErrInvalidIDFormat{message: fmt.Sprintf("'Wallet' type requires 1 ID component(s) (%d component(s) supplied)", len(ids))}
		}
		return n.OnWallet(ctx, persist.Address(ids[0]))
	}

	return nil, ErrInvalidIDFormat{typeName}
}

func (n *NodeFetcher) ValidateHandlers() {
	switch {
	case n.OnCollection == nil:
		panic("NodeFetcher handler validation failed: no handler set for NodeFetcher.OnCollection")
	case n.OnCollectionNft == nil:
		panic("NodeFetcher handler validation failed: no handler set for NodeFetcher.OnCollectionNft")
	case n.OnGallery == nil:
		panic("NodeFetcher handler validation failed: no handler set for NodeFetcher.OnGallery")
	case n.OnGalleryUser == nil:
		panic("NodeFetcher handler validation failed: no handler set for NodeFetcher.OnGalleryUser")
	case n.OnMembershipTier == nil:
		panic("NodeFetcher handler validation failed: no handler set for NodeFetcher.OnMembershipTier")
	case n.OnNft == nil:
		panic("NodeFetcher handler validation failed: no handler set for NodeFetcher.OnNft")
	case n.OnWallet == nil:
		panic("NodeFetcher handler validation failed: no handler set for NodeFetcher.OnWallet")
	}
}
