package infrastructure

import (
	"common"
	"common/errors"
	"common/logging"
	"context"
	"errors"
	"github.com/apex/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"order/domain"
)

type OrderRepositoryImpl struct {
	mongoCollection *mongo.Collection
}

func NewOrderRepository(mongo *mongo.Client, database string) *OrderRepositoryImpl {
	collection := mongo.Database(database).Collection("order")
	return &OrderRepositoryImpl{
		mongoCollection: collection,
	}
}

func (r *OrderRepositoryImpl) GetById(ctx context.Context, id string) (*domain.Order, error) {
	logger := r.getLogger(ctx)
	logger.Infof("GetById id: %s", id)

	filter := bson.M{"_id": id}
	var result domain.Order

	err := r.mongoCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperrors.EntityNotFound("EncryptedData not found", "id", id, err)
		} else {
			return nil, apperrors.InternalServerError("Unexpected error when querying EncryptedData", err)
		}
	}
	return &result, nil
}

func (r *OrderRepositoryImpl) GetAll(ctx context.Context, merchantFilter *domain.OrderFilter, pageFilter *common.PageFilter) (*common.Paginated[domain.Order], error) {
	logger := r.getLogger(ctx)
	logger.Infof("GetAll")

	filter := bson.M{}

	if len(merchantFilter.Id) > 0 {
		filter["_id"] = primitive.Regex{
			Pattern: merchantFilter.Id,
		}
	}
	if len(merchantFilter.Name) > 0 {
		filter["name"] = primitive.Regex{
			Pattern: merchantFilter.Name,
		}
	}
	if merchantFilter.CreatedFrom != nil && merchantFilter.CreatedTo != nil {
		filter["created"] = bson.M{
			"$gte": primitive.NewDateTimeFromTime(*merchantFilter.CreatedFrom),
			"$lte": primitive.NewDateTimeFromTime(*merchantFilter.CreatedTo),
		}
	} else {
		if merchantFilter.CreatedFrom != nil {
			filter["created"] = bson.M{
				"$gte": primitive.NewDateTimeFromTime(*merchantFilter.CreatedFrom),
			}
		}
		if merchantFilter.CreatedTo != nil {
			filter["created"] = bson.M{
				"$lte": primitive.NewDateTimeFromTime(*merchantFilter.CreatedTo),
			}
		}
	}

	pageSize := pageFilter.PageSize
	skip := pageFilter.GetSkip()

	opt := &options.FindOptions{
		Limit: &pageSize,
		Skip:  &skip,
		Sort:  bson.D{{pageFilter.SortField, pageFilter.GetSortTypeInt()}},
	}

	documentCount, err := r.mongoCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, apperrors.InternalServerError("Failed to get document count", err)
	}

	cursor, err := r.mongoCollection.Find(ctx, filter, opt)
	if err != nil {
		return nil, apperrors.InternalServerError("Failed to get all merchants", err)
	}
	defer cursor.Close(ctx)
	var encryptedDatas []*domain.Order
	for cursor.Next(ctx) {
		var encryptedData domain.Order
		if err := cursor.Decode(&encryptedData); err == nil {
			encryptedDatas = append(encryptedDatas, &encryptedData)
		}
	}

	return common.NewPaginated[domain.Order](encryptedDatas, documentCount, pageFilter.PageSize, pageFilter.Page), nil
}

func (r *OrderRepositoryImpl) Save(ctx context.Context, order *domain.Order) error {
	logger := r.getLogger(ctx)
	logger.Infof("Save %s", order.Id)

	order.Version = order.Version + 1

	_, err := r.mongoCollection.ReplaceOne(ctx, bson.M{"_id": order.Id}, order, options.Replace().SetUpsert(true))
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepositoryImpl) getLogger(ctx context.Context) *log.Entry {
	return logging.Log(ctx, "OrderRepository")
}
