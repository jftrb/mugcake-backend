package dbwrapper

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jftrb/mugacke-backend/src/api/models"
	"github.com/rs/zerolog/log"
)

type DbWrapper interface {
	Disconnect() error
	GetUsers() ([]models.User, error)
	GetRecipeSummaries(userID string) ([]models.RecipeSummary, error)
	GetRecipe(recipeID int) (models.Recipe, error)
}

type dbWrapper struct {
	client *pgxpool.Conn
	// movies *mongo.Collection
}

var pool *pgxpool.Pool = nil

func ConnectPool() {
	if pool == nil {
		dbPool, err := pgxpool.New(context.Background(), os.Getenv("POSTGRES_URL"))
		if err != nil {
			log.Error().AnErr("Error", err).Msg("Unable to connect to database.")
			panic(err)
		}
		pool = dbPool
	}
}

func DisconnectAll() {
	pool.Close()
	log.Info().Msg("Disconnected all connection from Postgres.")
	pool = nil
}


func NewDbWrapper() DbWrapper {
	ConnectPool()

	client, err := pool.Acquire(context.Background())
	if err != nil {
		log.Error().AnErr("Error", err).Msg("Unable to acquire a new connection from connection pool.")
		panic(err)
	}

	err = client.Ping(context.Background())
	if err != nil {
		log.Error().AnErr("Error", err).Msg("Connected to Postgres but unable to communicate: Initial Ping failed.")
		panic(err)
	}
	log.Info().Msg("Pinged your deployment. You successfully connected to Postgres!")

	return &dbWrapper{client: client}
}

func (d *dbWrapper) Disconnect() error {
	d.client.Release()
	log.Info().Msg("Released a connection from Postgres.")
	d.client = nil
	return nil
}

func (d *dbWrapper) GetUsers() ([]models.User, error) {
	rows, err := d.client.Query(context.Background(), "SELECT (id, email) from users")
	if err != nil {
		return []models.User{}, err
	}

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (models.User, error) {
		var user models.User
		err := row.Scan(&user)
		return user, err
	})
}

func (d *dbWrapper) GetRecipeSummaries(userID string) ([]models.RecipeSummary, error) {
	rows, err := d.client.Query(context.Background(), 
	`SELECT id, favorite, title, image_source, prep_info->>'totalTime' as total_time
			, ARRAY (
						SELECT tags.name
						FROM   unnest(r.tags) WITH ORDINALITY AS a(tag_id, ord)
						JOIN   tags ON tags.id = a.tag_id
						ORDER  BY a.ord
						)  AS tags
	FROM recipes r WHERE r.user_id = $1`, userID)

	if err != nil {
		return []models.RecipeSummary{}, err
	}

	log.Debug().Msg("Successful GetRecipesSummaries Query")

	return pgx.CollectRows(rows, pgx.RowToStructByName[models.RecipeSummary])
}

func (d *dbWrapper) GetRecipe(recipeID int) (models.Recipe, error) {
	rows, err := d.client.Query(context.Background(), 
	`SELECT id, favorite, title, url, image_source, prep_info, ingredients, directions, notes
			, ARRAY (
						SELECT tags.name
						FROM   unnest(r.tags) WITH ORDINALITY AS a(tag_id, ord)
						JOIN   tags ON tags.id = a.tag_id
						ORDER  BY a.ord
						)  AS tags
	FROM recipes r WHERE r.id = ` + fmt.Sprint(recipeID))

	if err != nil {
		return models.Recipe{}, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.Recipe])
}