package persist

import (
	"context"
	"fmt"

	"github.com/mikeydub/go-gallery/runtime"
	"go.mongodb.org/mongo-driver/bson"
)

const usersCollName = "users"

type User struct {
	VersionInt    int64   `bson:"version"` // schema version for this model
	IDstr         DbId    `bson:"_id"           json:"id"`
	CreationTimeF float64 `bson:"creation_time" json:"creation_time"`
	DeletedBool   bool    `bson:"deleted"`

	UserNameStr  string   `bson:"name"         json:"name"`       // mutable
	AddressesLst []string `bson:"addresses"     json:"addresses"` // IMPORTANT!! - users can have multiple addresses associated with their account
	BioStr       string   `bson:"bio"  json:"bio"`
}

//-------------------------------------------------------------
// USER
//-------------------------------------------------------------
// UPDATE
func UserUpdate(pUser *User,
	pCtx context.Context,
	pRuntime *runtime.Runtime) error {

	mp := NewMongoStorage(0, usersCollName, pRuntime)
	err := mp.Update(pCtx, bson.M{
		"_id": pUser.IDstr,
	}, pUser)
	if err != nil {
		return err
	}

	return nil
}

//-------------------------------------------------------------
// EXISTS BY ADDRESS
func UserExistsByAddress(pAddress string,
	pCtx context.Context,
	pRuntime *runtime.Runtime) (bool, error) {

	mp := NewMongoStorage(0, usersCollName, pRuntime)

	countInt, err := mp.Count(pCtx, bson.M{"addresses": bson.M{"$in": []string{pAddress}}})

	if err != nil {
		return false, err
	}

	return countInt > 0, nil
}

//-------------------------------------------------------------
// CREATE
func UserCreate(pUser *User,
	pCtx context.Context,
	pRuntime *runtime.Runtime) (DbId, error) {

	mp := NewMongoStorage(0, usersCollName, pRuntime)

	return mp.Insert(pCtx, pUser)

}

//-------------------------------------------------------------
// DELETE
func UserDelete(pUserID DbId,
	pCtx context.Context,
	pRuntime *runtime.Runtime) error {

	mp := NewMongoStorage(0, usersCollName, pRuntime)

	return mp.Update(pCtx, bson.M{"_id": pUserID}, bson.M{"$set": bson.M{"deleted": true}})

}

//-------------------------------------------------------------
// GET BY ID
func UserGetById(userId DbId,
	pCtx context.Context,
	pRuntime *runtime.Runtime) (*User, error) {

	mp := NewMongoStorage(0, usersCollName, pRuntime)

	result := []*User{}
	err := mp.Find(pCtx, bson.M{"_id": userId}, result)

	if err != nil {
		return nil, err
	}

	if len(result) == 0 || len(result) > 1 {
		return nil, fmt.Errorf("invalid amount of returned users: %d", len(result))
	}

	return result[0], nil
}

//-------------------------------------------------------------
// GET BY ADDRESS
func UserGetByAddress(pAddress string,
	pCtx context.Context,
	pRuntime *runtime.Runtime) (*User, error) {

	mp := NewMongoStorage(0, usersCollName, pRuntime)

	result := []*User{}
	err := mp.Find(pCtx, bson.M{"addresses": bson.M{"$in": []string{pAddress}}}, result)

	if err != nil {
		return nil, err
	}

	if len(result) == 0 || len(result) > 1 {
		return nil, fmt.Errorf("invalid amount of returned users: %d", len(result))
	}

	return result[0], nil
}
