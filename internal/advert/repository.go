package advert

import (
	"context"
	"github.com/tvitcom/local-adverts/internal/entity"
	"github.com/tvitcom/local-adverts/internal/dto"
	"github.com/tvitcom/local-adverts/pkg/dbcontext"
	"github.com/tvitcom/local-adverts/pkg/log"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

// repository persists adverts in database
type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

// Repository encapsulates the logic to access adverts from the data source.
type Repository interface {
	// Get returns the users with follow by specified limit-offset params.
	GetUsersWithLimitOffset(ctx context.Context, limit, offset int64) ([]entity.User, error)
	// Get returns the advert with the specified advert ID.
	GetAdvertById(ctx context.Context, id int64) (entity.Advert, error)
	// Get returns the advert with the specified advert ID.
	GetAdvertsLast(ctx context.Context) ([]entity.Advert, error)
	// Get returns the Messages by specified User_ID.
	GetMessagesSendersByUserId(ctx context.Context, uid int64) ([]dto.MessageSender, error)
	// Get returns the advert by word filtered like function.
	GetAdvertsSearch(ctx context.Context, word string) ([]entity.Advert, error)
	// Get returns the adverts by user_id.
	GetAdvertsByUserId(ctx context.Context, uid int64) ([]entity.Advert, error)
	// Get returns the user with the specified user_id.
	GetUserById(ctx context.Context, id int64) (entity.User, error)
	// Get returns the user with the specified advert_id.
	GetUserByAdvertId(ctx context.Context, aid int64) (entity.User, error)
	// Get returns the advert with the specified email.
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	// Count returns the number of adverts.
	Count(ctx context.Context) (int, error)
	// Query returns the list of adverts with the given offset and limit.
	QueryAdverts(ctx context.Context, offset, limit int) ([]entity.Advert, error)
	// Create saves a new user in the storage.
	CreateUser(ctx context.Context, user entity.User) (int64, error)
	// Update user in the storage.
	UpdateUser(ctx context.Context, user entity.User, uid int64) error
	// Update User lastlogin field
	UpdateUserLastlogin(ctx context.Context, uid int64, dtstring string) error
	// Update User Adverts Picture
	UpdateAdvertsPicture(ctx context.Context, aid int64, field, fname string) error
	// Create saves a new advert in the storage.
	CreateAdvert(ctx context.Context, advert entity.Advert) (int64, error)
	// Getting advert by user_id.
	GetAdvertsDisplayByUserId(ctx context.Context, id int64) ([]dto.AdvertDisplay, error)
	// Create saves a new advert in the storage.
	CreateMessage(ctx context.Context, message entity.Message) error
	// Update updates the advert with given ID in the storage.
	Update(ctx context.Context, advert entity.Advert) error
	// Delete removes the advert with given ID from the storage.
	DeleteAdvertById(ctx context.Context, id int64) error
	// GetCategory(ctx) returns the list of adverts with the given offset and limit.
	GetCategory(ctx context.Context) ([]entity.Category, error)
	// GetCategoryPath(ctx) returns the list categories by special tree pathes query.
	GetCategoryPath(ctx context.Context) ([]entity.CategoryPath, error)
}

// NewRepository creates a new advert repository
func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

// Get the advert with the specified ID from the database.
func (r repository) GetAdvertById(ctx context.Context, id int64) (entity.Advert, error) {
	var advert entity.Advert
	err := r.db.With(ctx).Select().Model(id, &advert)
	return advert, err
}

// Get the user with the specified ID from the database.
func (r repository) GetUserById(ctx context.Context, id int64) (entity.User, error) {
	var user entity.User
	err := r.db.With(ctx).Select().From("user").Where(dbx.HashExp{"user_id": id}).One(&user)
	return user, err
}

// Get the user with the specified Advert_ID from the database.
func (r repository) GetUserByAdvertId(ctx context.Context, aid int64) (entity.User, error) {
	var user entity.User
	sql := `select u.* from user u left join advert a on u.user_id=a.user_id where a.advert_id={:aid}`
	err := r.db.With(ctx).NewQuery(sql).Bind(dbx.Params{"aid": aid}).One(&user)
	return user, err
}

// Get the user with the specified ID from the database.
func (r repository) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	var user entity.User
	err := r.db.With(ctx).Select().From("user").Where(dbx.HashExp{"email": email}).One(&user)
	return user, err
}

// Create saves a new advert record in the database.
// It returns the ID of the newly inserted advert record.
func (r repository) CreateUser(ctx context.Context, user entity.User) (int64, error) {
	err := r.db.With(ctx).Model(&user).Insert()
	if err != nil {
		return 0, err
	}
	return user.UserId, nil
}

// Create saves a new advert record in the database.
// It returns the ID of the newly inserted advert record.
func (r repository) CreateAdvert(ctx context.Context, advert entity.Advert) (int64, error) {
	err := r.db.With(ctx).Model(&advert).Insert()
	if err != nil {
		return 0, err
	}
	return advert.AdvertId, nil
}

// return advert records in the database.
func (r repository) GetAdvertsDisplayByUserId(ctx context.Context, uid int64) ([]dto.AdvertDisplay, error) {
	var adverts []dto.AdvertDisplay
	sql := `
	SELECT a.advert_id, c.name category_name, a.title, a.price, a.currency, a.moderator_id, a.created, a.active 
	FROM advert a, category c 
	WHERE a.category_id = c.category_id AND a.user_id = {:uid}
	GROUP BY a.category_id 
	ORDER BY a.created ASC 
	LIMIT 100;`
	err := r.db.With(ctx).NewQuery(sql).Bind(dbx.Params{"uid": uid}).All(&adverts)
	return adverts, err
}

// returns recently added advert records
func (r repository) GetUsersWithLimitOffset(ctx context.Context, limit, offset int64) ([]entity.User, error) {
	var items []entity.User
	err := r.db.With(ctx).
		Select().
		From("user").
		// Where(dbx.HashExp{"active": 0}).
		Offset(offset).
		Limit(limit).
		OrderBy("lastlogin").
		All(&items)
	return items, err
}

// Create saves a new advert record in the database.
// It returns the ID of the newly inserted advert record.
// 
func (r repository) CreateMessage(ctx context.Context, message entity.Message) error {
	return r.db.With(ctx).Model(&message).Insert()
}

// Update saves the changes to an user in the database.
func (r repository) UpdateUser(ctx context.Context, user entity.User, uid int64) error {
	dbxvar := dbx.Params{
			"name": user.Name,
			"tel": user.Tel,
		}
	if user.Passhash != "" {
		dbxvar["passhash"]= user.Passhash
	}
	if user.Picture != "" {
		dbxvar["picture"] = user.Picture
	}
	// UPDATE `users` SET `status`={:p0} WHERE `id`={:p1}
	_, err := r.db.With(ctx).Update("user", dbxvar, dbx.HashExp{
		"user_id": uid,
	}).Execute()
	return err
}

// Update saves the changes to an user in the database.
func (r repository) UpdateUserLastlogin(ctx context.Context, uid int64, dtstring string) error {
	dbxvar := dbx.Params{
			"lastlogin": dtstring,
		}
	// UPDATE `users` SET `status`={:p0} WHERE `id`={:p1}
	_, err := r.db.With(ctx).Update("user", dbxvar, dbx.HashExp{
		"user_id": uid,
	}).Execute()
	return err
}
// UpdateAdvertsPicture(ctx context.Context, aid int64, field, fname string) error
func (r repository) UpdateAdvertsPicture(ctx context.Context, aid int64, field, fname string) error {
	dbxvar := dbx.Params{
			field: fname,
		}
	_, err := r.db.With(ctx).Update("advert", dbxvar, dbx.HashExp{
		"advert_id": aid,
	}).Execute()
	return err
}

// Update saves the changes to an advert in the database.
func (r repository) Update(ctx context.Context, advert entity.Advert) error {
	return r.db.With(ctx).Model(&advert).Update()
}

// Delete deletes an advert with the specified ID from the database.
func (r repository) DeleteAdvertById(ctx context.Context, id int64) error {
	_, err := r.db.With(ctx).Delete("advert", dbx.HashExp{"advert_id": id}).Execute()
	return err
}

// Count returns the number of the advert records in the database.
func (r repository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.With(ctx).Select("COUNT(*)").From("advert").Row(&count)
	return count, err
}

// Query retrieves the advert records with the specified offset and limit from the database.
func (r repository) QueryAdverts(ctx context.Context, offset, limit int) ([]entity.Advert, error) {
	var adverts []entity.Advert
	err := r.db.With(ctx).
		Select().
		OrderBy("id").
		Offset(int64(offset)).
		Limit(int64(limit)).
		All(&adverts)
	return adverts, err
}

// Count returns the number of the advert records in the database.
func (r repository) GetCategory(ctx context.Context) ([]entity.Category, error) {
	var categories []entity.Category
	err := r.db.With(ctx).
		Select().
		From("category").
		OrderBy("category_id").
		All(&categories)
	return categories, err
}

// Count returns the number of the advert records in the database.
func (r repository) GetCategoryPath(ctx context.Context) ([]entity.CategoryPath, error) {
	var items []entity.CategoryPath
	sql := `
	select c.category_id, concat_ws("-",( 
		select parc.name  
		from category parc  
		where c.parent_category_id = parc.category_id 
	),c.name) as path  
	from category c  where c.parent_category_id > 0 
	UNION 
	select c.category_id, c.name 
	from category c 
	where c.parent_category_id = 0 and c.category_id not in (
		select cc.parent_category_id 
		from category cc
	) order by 1 asc;`
	err := r.db.With(ctx).NewQuery(sql).All(&items)
	return items, err
}

// returns recently added advert records
func (r repository) GetAdvertsLast(ctx context.Context) ([]entity.Advert, error) {
	var limit int = 99
	var items []entity.Advert
	err := r.db.With(ctx).
		Select().
		From("advert").
		Where(dbx.HashExp{"active": 1}).
		// Offset(int64(0)).
		Limit(int64(limit)).
		OrderBy("created").
		All(&items)
	return items, err
}

func (r repository) GetMessagesSendersByUserId(ctx context.Context, uid int64) ([]dto.MessageSender, error) {
	var items []dto.MessageSender
	sql := `
	select m.sender_id, u.name, u.email, u.tel, m.sended, m.moderator_id 
	from message m left join user u on m.sender_id = u.user_id 
	where m.receiver_id = {:uid}
	order by m.sended ASC
	`
	err := r.db.With(ctx).NewQuery(sql).Bind(dbx.Params{"uid": uid}).All(&items)
	return items, err
}

func (r repository) GetAdvertsSearch(ctx context.Context, word string) ([]entity.Advert, error) {
	var limit int = 200
	var items []entity.Advert
	err := r.db.With(ctx).
		Select().
		From("advert").
		Where(dbx.And(dbx.HashExp{"active": 1}, dbx.Like("nanopost", word))).
		// Offset(int64(0)).
		Limit(int64(limit)).
		OrderBy("created").
		All(&items)
	return items, err
}

func (r repository) GetAdvertsByUserId(ctx context.Context, uid int64) ([]entity.Advert, error) {
	var items []entity.Advert
	err := r.db.With(ctx).
		Select().
		From("advert").
		Where(dbx.HashExp{"user_id": uid}).
		OrderBy("created").
		All(&items)
	return items, err
}