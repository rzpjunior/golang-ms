package mongox

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client interface {
	// Insert : for insert one data
	// 	var user *model.user
	// 	user.ID = 1
	// 	user.Name = "Adan Aidan Teras"
	// 	ret, err := mongox.Client.Insert(ctx.TODO(), "users", user)
	Insert(ctx context.Context, colName string, payload interface{}) (ret interface{}, err error)
	// InsertBulk : for insert multiple data
	// 	var users []*model.user
	// 	users = append(users, &model.user{ID:1, Name:"Adan Aidan Teras"})
	// 	ret, err := mongox.Client.InsertBulk(ctx.TODO(), "users", users)
	InsertBulk(ctx context.Context, colName string, payload []interface{}) (ret interface{}, err error)

	// Update : for update one data
	Update(ctx context.Context, colName string, filter interface{}, payload interface{}) (err error)
	// UpdateBulk : for update one data
	UpdateBulk(ctx context.Context, colName string, filter interface{}, payload interface{}) (err error)

	// Get : for get data list
	// 	ret, err := mongox.Client.Get(ctx.TODO(), "users", &options.FindOptions{Skip:0,Limit:100})
	Get(ctx context.Context, colName string, opts ...*options.FindOptions) (ret []byte, err error)
	// GetByFilter : for get data list by filter
	// 	ret, err := mongox.Client.Get(ctx.TODO(), "users", &model.user{ID:1}, &options.FindOptions{Skip:0,Limit:100})
	GetByFilter(ctx context.Context, colName string, filter interface{}, opts ...*options.FindOptions) (ret []byte, err error)

	// Find : for get one data
	Find(ctx context.Context, colName string, opts ...*options.FindOneOptions) (ret []byte, err error)
	// FindByFilter : for get one data by filter
	FindByFilter(ctx context.Context, colName string, filter interface{}, opts ...*options.FindOneOptions) (ret []byte, err error)

	// GetCount : for get counted data by filter
	GetCount(ctx context.Context, colName string, opts ...*options.CountOptions) (ret int64, err error)
	// GetCountByFilter : for get counted data by filter
	GetCountByFilter(ctx context.Context, colName string, filter interface{}, opts ...*options.CountOptions) (ret int64, err error)

	// CreateIndex : for indexing single field
	// err = mongox.Client.CreateIndex(ctx.TODO(),"reference_id",true)
	// set unique = true if the value cannot duplicate, example: data document_code cannot duplicate
	// set unique = false if the value can be duplicate, example: data delivery_date can be duplicate
	CreateIndex(ctx context.Context, colName string, field string, unique bool) (err error)

	DeleteMany(ctx context.Context, colName string, filter interface{}) (int64, error)

	FindAggregate(ctx context.Context, colName string, pipeline interface{}, opts ...*options.AggregateOptions) (ret []byte, err error)

	// RemoveCollection removes a MongoDB collection from the specified database.
	//
	// Parameters:
	// - ctx: The context for the MongoDB operation.
	// - colName: The name of the collection to drop.
	//
	// Example:
	//   err := RemoveCollection(ctx, "mycollection")
	//   if err != nil {
	//       log.Fatal(err)
	//   }
	RemoveCollection(ctx context.Context, colName string) (err error)
}

type Mongox struct {
	Database string
	Client   *mongo.Client
	Log      *logrus.Logger
}

func NewMongox(database string, client *mongo.Client, lgr *logrus.Logger) Client {
	return &Mongox{
		Database: database,
		Client:   client,
		Log:      lgr,
	}
}

func (m *Mongox) InsertBulk(ctx context.Context, colName string, payload []interface{}) (ret interface{}, err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	opts := options.InsertMany().SetOrdered(false)
	collection := m.Client.Database(m.Database).Collection(colName)

	fmt.Println(opts, collection, ">>>>>>InsertBulk<<<<<<<", payload)

	var res *mongo.InsertManyResult
	if res, err = collection.InsertMany(ctx, payload, opts); err != nil {
		fmt.Println(">>>>>>InsertBulk Error:<<<<<<<", err)
		m.ErrorLog("InsertBulk", colName, nil, payload, err)
		return
	}

	ret = res.InsertedIDs

	m.DebugLog("InsertBulk", colName, nil, payload)

	return
}

func (m *Mongox) Insert(ctx context.Context, colName string, payload interface{}) (ret interface{}, err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	collection := m.Client.Database(m.Database).Collection(colName)
	rawData, err := bson.Marshal(&payload)
	if err != nil {
		fmt.Println(rawData, collection, ">>>>>>InsertBulk<<<<<<<", payload)
		fmt.Println(">>>>>>Insert Error 1:<<<<<<<", err)
		m.ErrorLog("Insert", colName, nil, payload, err)
		return
	}

	var res *mongo.InsertOneResult
	res, err = collection.InsertOne(ctx, rawData)
	if err != nil {
		fmt.Println(">>>>>>Insert Error:<<<<<<<", err)
		m.ErrorLog("Insert", colName, nil, payload, err)
		return
	}

	ret = res.InsertedID

	m.DebugLog("Insert", colName, nil, payload)

	return
}

func (m *Mongox) Get(ctx context.Context, colName string, opts ...*options.FindOptions) (ret []byte, err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var res []bson.M
	collection := m.Client.Database(m.Database).Collection(colName)
	if len(opts) == 0 {
		options2 := &options.FindOptions{}
		opts = append(opts, options2)
	}

	opts[0].SetProjection(bson.M{"_id": 0})

	cursor, err := collection.Find(ctx, bson.M{}, opts[0])
	if err != nil {
		m.ErrorLog("Get", colName, nil, nil, err)
		return
	}

	defer cursor.Close(ctx)
	cursor.All(ctx, &res)

	ret, err = json.Marshal(res)
	if err != nil {
		m.ErrorLog("Get", colName, nil, nil, err)
		return
	}

	m.DebugLog("Get", colName, nil, nil)

	return
}

func (m *Mongox) GetByFilter(ctx context.Context, colName string, filter interface{}, opts ...*options.FindOptions) (ret []byte, err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var res []bson.M
	data, err := bson.Marshal(filter)
	if err != nil {
		m.ErrorLog("GetByFilter", colName, filter, nil, err)
		return
	}

	collection := m.Client.Database(m.Database).Collection(colName)
	if len(opts) == 0 {
		options2 := &options.FindOptions{}
		opts = append(opts, options2)
	}

	opts[0].SetProjection(bson.M{"_id": 0})

	cursor, err := collection.Find(ctx, data, opts[0])
	if err != nil {
		m.ErrorLog("GetByFilter", colName, filter, nil, err)
		return
	}
	defer cursor.Close(ctx)
	cursor.All(ctx, &res)

	ret, err = json.Marshal(res)
	if err != nil {
		m.ErrorLog("GetByFilter", colName, filter, nil, err)
		return
	}

	m.DebugLog("GetByFilter", colName, filter, nil)

	return
}

func (m *Mongox) Find(ctx context.Context, colName string, opts ...*options.FindOneOptions) (ret []byte, err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var dataCollection bson.M
	collection := m.Client.Database(m.Database).Collection(colName)
	//to remove object id from mongo
	if len(opts) == 0 {
		options2 := &options.FindOneOptions{}
		opts = append(opts, options2)
	}

	opts[0].SetProjection(bson.M{"_id": 0})

	err = collection.FindOne(ctx, bson.M{}, opts[0]).Decode(&dataCollection)
	if err != nil {
		m.ErrorLog("Find", colName, nil, nil, err)
		return
	}

	ret, err = json.Marshal(dataCollection)
	if err != nil {
		m.ErrorLog("Find", colName, nil, nil, err)
		return
	}

	m.DebugLog("Find", colName, nil, nil)

	return
}

func (m *Mongox) FindByFilter(ctx context.Context, colName string, filter interface{}, opts ...*options.FindOneOptions) (ret []byte, err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var dataCollection bson.M
	data, err := bson.Marshal(filter)
	if err != nil {
		m.ErrorLog("FindByFilter", colName, filter, nil, err)
		return
	}
	collection := m.Client.Database(m.Database).Collection(colName)

	if len(opts) == 0 {
		options2 := &options.FindOneOptions{}
		opts = append(opts, options2)
	}

	opts[0].SetProjection(bson.M{"_id": 0})

	err = collection.FindOne(ctx, data, opts[0]).Decode(&dataCollection)
	if err != nil {
		m.ErrorLog("FindByFilter", colName, filter, nil, err)
		return
	}
	ret, err = json.Marshal(dataCollection)
	if err != nil {
		m.ErrorLog("FindByFilter", colName, filter, nil, err)
		return
	}

	m.DebugLog("FindByFilter", colName, filter, nil)
	return
}

func (m *Mongox) GetCountByFilter(ctx context.Context, colName string, filter interface{}, opts ...*options.CountOptions) (ret int64, err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	data, err := bson.Marshal(filter)
	if err != nil {
		m.ErrorLog("GetCountByFilter", colName, filter, nil, err)
		return 0, err
	}
	collection := m.Client.Database(m.Database).Collection(colName)
	if len(opts) == 0 {
		options2 := &options.CountOptions{}
		opts = append(opts, options2)
	}

	ret, err = collection.CountDocuments(ctx, data, opts[0])
	if err != nil {
		m.ErrorLog("GetCountByFilter", colName, filter, nil, err)
		return 0, err
	}

	m.DebugLog("GetCountByFilter", colName, filter, nil)
	return ret, err
}

func (m *Mongox) GetCount(ctx context.Context, colName string, opts ...*options.CountOptions) (ret int64, err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	collection := m.Client.Database(m.Database).Collection(colName)
	if len(opts) == 0 {
		options2 := &options.CountOptions{}
		opts = append(opts, options2)
	}

	ret, err = collection.CountDocuments(ctx, bson.M{}, opts[0])
	if err != nil {
		m.ErrorLog("GetCount", colName, nil, nil, err)
		return 0, err
	}

	m.DebugLog("GetCount", colName, nil, nil)
	return ret, err
}

func (m *Mongox) Update(ctx context.Context, colName string, filter interface{}, payload interface{}) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	data, err := bson.Marshal(filter)
	if err != nil {
		m.ErrorLog("Update", colName, filter, payload, err)
		return
	}

	collection := m.Client.Database(m.Database).Collection(colName)
	update := bson.M{
		"$set": payload,
	}
	_, err = collection.UpdateOne(ctx, data, update)
	if err != nil {
		m.ErrorLog("Update", colName, filter, payload, err)
		return
	}

	m.DebugLog("Update", colName, filter, payload)
	return
}

func (m *Mongox) UpdateBulk(ctx context.Context, colName string, filter interface{}, payload interface{}) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	data, err := bson.Marshal(filter)
	if err != nil {
		m.ErrorLog("UpdateBulk", colName, filter, payload, err)
		return
	}
	collection := m.Client.Database(m.Database).Collection(colName)
	update := bson.M{
		"$set": payload,
	}

	_, err = collection.UpdateMany(ctx, data, update)
	if err != nil {
		m.ErrorLog("UpdateBulk", colName, filter, payload, err)
		return
	}

	m.DebugLog("UpdateBulk", colName, filter, payload)
	return
}

func (m *Mongox) DebugLog(funcName string, colName string, filter interface{}, payload interface{}) {
	field := logrus.Fields{
		"name":       funcName,
		"time":       time.Now(),
		"collection": colName,
	}

	if filter != nil {
		field["filter"] = filter
	}

	if payload != nil {
		field["payload"] = payload
	}

	m.Log.WithFields(field).Info("Mongox")
}

func (m *Mongox) ErrorLog(funcName string, colName string, filter interface{}, payload interface{}, err error) {
	field := logrus.Fields{
		"name":       funcName,
		"time":       time.Now(),
		"collection": colName,
		"errors":     err,
	}

	if filter != nil {
		field["filter"] = filter
	}

	if payload != nil {
		field["payload"] = payload
	}

	m.Log.WithFields(field).Error("Mongox")
}

func (m *Mongox) CreateIndex(ctx context.Context, colName string, field string, unique bool) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	mod := mongo.IndexModel{
		Keys:    bson.M{field: 1}, // index in ascending order or -1 for descending order
		Options: options.Index().SetUnique(unique),
	}
	collection := m.Client.Database(m.Database).Collection(colName)

	_, err = collection.Indexes().CreateOne(ctx, mod)
	if err != nil {
		m.ErrorLog("CreateIndex", colName, nil, field, err)
		return err
	}
	m.DebugLog("CreateIndex", colName, nil, field)

	return err
}

func (m *Mongox) DeleteMany(ctx context.Context, colName string, filter interface{}) (int64, error) {
	collection := m.Client.Database(m.Database).Collection(colName)

	result, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		m.ErrorLog("DeleteMany", colName, filter, nil, err)
		return 0, err
	}

	m.DebugLog("DeleteMany", colName, filter, nil)
	return result.DeletedCount, nil
}

func (m *Mongox) FindAggregate(ctx context.Context, colName string, pipeline interface{}, opts ...*options.AggregateOptions) (ret []byte, err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	collection := m.Client.Database(m.Database).Collection(colName)

	if len(opts) == 0 {
		options := &options.AggregateOptions{}
		opts = append(opts, options)
	}

	cursor, err := collection.Aggregate(ctx, pipeline, opts[0])
	if err != nil {
		m.ErrorLog("FindAggregate", colName, pipeline, nil, err)
		return nil, err // Return nil for ret and propagate the error
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		m.ErrorLog("FindAggregate", colName, pipeline, nil, err)
		return nil, err // Return nil for ret and propagate the error
	}

	ret, err = json.Marshal(results)
	if err != nil {
		m.ErrorLog("FindAggregate", colName, pipeline, nil, err)
		return nil, err // Return nil for ret and propagate the error
	}

	m.DebugLog("FindAggregate", colName, pipeline, nil)
	return ret, nil // Return the result and nil for err
}

func (m *Mongox) RemoveCollection(ctx context.Context, colName string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	collection := m.Client.Database(m.Database).Collection(colName)

	err = collection.Drop(ctx)
	if err != nil {
		m.ErrorLog("CreateIndex", colName, nil, nil, err)
		return err
	}

	return err
}
