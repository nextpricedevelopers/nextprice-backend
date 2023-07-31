package payment

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/nextpricedevelopers/go-next/internal/config"
	"github.com/nextpricedevelopers/go-next/internal/config/logger"
	"github.com/nextpricedevelopers/go-next/pkg/adapter/mongodb"
	"github.com/nextpricedevelopers/go-next/pkg/model"
	"github.com/nextpricedevelopers/go-next/pkg/service/validation"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PaymentServiceInterface interface {
	Create(ctx context.Context, payment *model.PaymentOptions) (*model.PaymentOptions, error)
	Update(ctx context.Context, ID string, payment model.PaymentOptions) (bool, error)
	GetByID(ctx context.Context, ID string) (*model.PaymentOptions, error)
	GetByName(ctx context.Context, paymentName string) (payment *model.PaymentOptions, err error)
	ListPayament(ctx context.Context, paymentName string, enabled ...bool) ([]*model.PaymentOptions, error)
}

var conf = config.NewConfig()

type PaymentDataService struct {
	dbp mongodb.MongoDBInterface
}

func NewPaymentService(db_pool mongodb.MongoDBInterface) *PaymentDataService {
	collection := db_pool.GetCollection(conf.MongoDBConfig.MDB_COLLECTION)
	
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"payment_name": 1},
		Options: options.Index().SetUnique(true),
	}

	// Crie o índice na coleção.
	if _, err := collection.Indexes().CreateOne(context.Background(), indexModel); err != nil {
		log.Fatalf("Erro ao criar índice único: %v", err)
	}
	return &PaymentDataService{
		dbp: db_pool,
	}
}

func (pds *PaymentDataService) Create(ctx context.Context, pay *model.PaymentOptions) (*model.PaymentOptions, error) {
	collection := pds.dbp.GetCollection(conf.MongoDBConfig.MDB_COLLECTION)

	dt := time.Now().Format(time.RFC3339)
	pay.PaymentName = validation.CareString(pay.PaymentName)
	pay.Enabled = true
	pay.CreatedAt = dt
	pay.UpdatedAt = dt

	log.Println(pay.PaymentName)

	result, err := collection.InsertOne(ctx, pay)
	if err != nil {
		log.Println(err.Error())
		logger.Error("Erro to create:"+err.Error(), err)
		return pay, err
	}

	pay.ID = result.InsertedID.(primitive.ObjectID)
	logger.Info("Payment created successfully")

	return pay, nil
}

func (cds *PaymentDataService) Update(ctx context.Context, ID string, payment model.PaymentOptions) (bool, error) {
	collection := cds.dbp.GetCollection(conf.MongoDBConfig.MDB_COLLECTION)

	opts := options.Update().SetUpsert(true)

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		log.Println("Error to parse ObjectIDFromHex")
		return false, err
	}

	filter := bson.D{

		{Key: "_id", Value: objectID},
	}
	payment.PaymentName = validation.CareString(payment.PaymentName)
	update := bson.D{{Key: "$set",
		Value: bson.D{
			{Key: "payment_name", Value: payment.PaymentName},

			{Key: "updated_at", Value: time.Now().Format(time.RFC3339)},
		},
	}}

	_, err = collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Println("Error while updating data")
		return false, err
	}

	return true, nil
}

func (cds *PaymentDataService) GetByID(ctx context.Context, ID string) (*model.PaymentOptions, error) {

	collection := cds.dbp.GetCollection(conf.MongoDBConfig.MDB_COLLECTION)

	payment := &model.PaymentOptions{}

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	filter := bson.D{
		{Key: "data_type", Value: "customer"},
		{Key: "_id", Value: objectID},
	}

	err = collection.FindOne(ctx, filter).Decode(payment)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return payment, nil
}

func (cds *PaymentDataService) GetByName(ctx context.Context, name string) (payment *model.PaymentOptions, err error) {

	collection := cds.dbp.GetCollection(conf.MongoDBConfig.MDB_COLLECTION)

	filter := bson.D{
		{Key: "email", Value: name},
	}

	err = collection.FindOne(ctx, filter).Decode(&payment)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return payment, nil
}

func (cds *PaymentDataService) ListPayament(ctx context.Context, pymenteName string, enabled ...bool) ([]*model.PaymentOptions, error) {
	collection := cds.dbp.GetCollection(conf.MongoDBConfig.MDB_COLLECTION)

	where := bson.M{}

	if pymenteName != "" {
		where["pymenteName"] = pymenteName
	}

	if len(enabled) > 0 {
		where["enabled"] = enabled[0]
	}

	cursor, err := collection.Find(ctx, where)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*model.PaymentOptions
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
