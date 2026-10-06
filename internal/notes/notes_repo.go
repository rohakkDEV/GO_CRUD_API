package notes

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var ErrNotFound = errors.New("note not found")

// repo -> data access layer
type Repo struct {
	coll *mongo.Collection
}

func NewRepo(db *mongo.Database) *Repo {
	return &Repo{
		coll: db.Collection("notes"),
	}
}

func (r *Repo) Create(ctx context.Context, note Note) (Note, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// your model's _id has no omitempty, so the ID must be set here
	now := time.Now().UTC()
	note.ID = primitive.NewObjectID()
	note.CreatedAt = now
	note.UpdatedAt = now

	if _, err := r.coll.InsertOne(opCtx, note); err != nil {
		return Note{}, fmt.Errorf("insert note failed: %w", err)
	}
	return note, nil
}

func (r *Repo) List(ctx context.Context) ([]Note, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// find returns a cursor (an iterator over matching docs)
	cursor, err := r.coll.Find(opCtx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("find notes failed: %w", err)
	}
	// cursor must be closed after use to avoid leaks
	defer cursor.Close(opCtx)

	notes := []Note{} // empty slice so JSON is [] instead of null
	if err := cursor.All(opCtx, &notes); err != nil {
		return nil, fmt.Errorf("decode notes failed: %w", err)
	}
	return notes, nil
}

func (r *Repo) GetByID(ctx context.Context, id primitive.ObjectID) (Note, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var note Note
	err := r.coll.FindOne(opCtx, bson.M{"_id": id}).Decode(&note)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return Note{}, ErrNotFound
	}
	if err != nil {
		return Note{}, fmt.Errorf("find note by id failed: %w", err)
	}
	return note, nil
}

func (r *Repo) UpdateByID(ctx context.Context, id primitive.ObjectID, req UpdateNoteRequest) (Note, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"title":     req.Title,
			"content":   req.Content,
			"pinned":    req.Pinned,
			"updatedAt": time.Now().UTC(), 
		},
	}

	
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updated Note
	err := r.coll.FindOneAndUpdate(opCtx, bson.M{"_id": id}, update, opts).Decode(&updated)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return Note{}, ErrNotFound
	}
	if err != nil {
		return Note{}, fmt.Errorf("update note failed: %w", err)
	}
	return updated, nil
}

func ( r *Repo) DeleteByID(ctx context.Context, id primitive.ObjectID)(bool,error){
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"_id" : id}

	res, err := r.coll.DeleteOne(opCtx, filter)
	if err != nil{
		return false, fmt.Errorf("failed to delete the given note: %w", err)
	}
	if(res.DeletedCount == 0){
		return false,nil
	}
	return true,nil
}