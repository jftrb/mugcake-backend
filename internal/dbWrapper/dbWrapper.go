package dbwrapper

import (
	"context"
	"encoding/json"
	"os"
	"time"

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
	AddRecipe(userID string, recipe models.Recipe) (error)
	UpdateRecipe(recipeID int, recipe models.Recipe) (error)
	DeleteRecipe(recipeID int) (error)
}

type dbWrapper struct {
	client *pgxpool.Conn
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
	`SELECT id as recipeId, favorite, title, image_source, prep_info->>'totalTime' as total_time
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
	`SELECT favorite, title, url, image_source, prep_info, ingredientSections, directions, notes
			, ARRAY (
						SELECT tags.name
						FROM   unnest(r.tags) WITH ORDINALITY AS a(tag_id, ord)
						JOIN   tags ON tags.id = a.tag_id
						ORDER  BY a.ord
						)  AS tags
	FROM recipes r WHERE r.id = $1`, recipeID)

	if err != nil {
		return models.Recipe{}, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.Recipe])
}

func (d *dbWrapper)	AddRecipe(userID string, recipe models.Recipe) (error) {
	prepInfo, err := json.Marshal(recipe.PrepInfo)
	if err != nil {
		return err
	}
	
	ingredientSections, err := json.Marshal(recipe.IngredientSections)
	if err != nil {
		return err
	}

	timeout, _ := context.WithTimeout(context.Background(), time.Second * 2)
	_, err = d.client.Query(timeout, 
	`INSERT INTO recipes (user_id, favorite, title, url, image_source, prep_info, tags, ingredientSections, directions, notes)
	VALUES (
		$1, $2, $3, $4, $5, $6,
		ARRAY (
			SELECT tags.id
			FROM	 unnest($7::text[]) WITH ORDINALITY AS a(tag_name, ord)
			JOIN   tags ON tags.name = a.tag_name
			ORDER  BY a.ord
		),
		$8, $9, $10)`, userID, recipe.Favorite, recipe.Title, recipe.URL, recipe.ImageSource, prepInfo, recipe.Tags, ingredientSections, recipe.Directions, recipe.Notes)

	return err
}

func (d *dbWrapper)	UpdateRecipe(recipeID int, recipe models.Recipe) (error) {
	prepInfo, err := json.Marshal(recipe.PrepInfo)
	if err != nil {
		return err
	}
	
	ingredientSections, err := json.Marshal(recipe.IngredientSections)
	if err != nil {
		return err
	}

	timeout, _ := context.WithTimeout(context.Background(), time.Second * 2)
	_, err = d.client.Query(timeout, 
	`UPDATE recipes SET (favorite, title, url, image_source, prep_info, tags, ingredientSections, directions, notes) = 
	(
		$1, $2, $3, $4, $5,
		ARRAY (
			SELECT tags.id
			FROM	 unnest($6::text[]) WITH ORDINALITY AS a(tag_name, ord)
			JOIN   tags ON tags.name = a.tag_name
			ORDER  BY a.ord
		),
		$7, $8, $9)
		WHERE id = $10`, recipe.Favorite, recipe.Title, recipe.URL, recipe.ImageSource, prepInfo, recipe.Tags, ingredientSections, recipe.Directions, recipe.Notes, recipeID)

	return err
}

func (d *dbWrapper)	DeleteRecipe(recipeID int) (error) {
	_, err := d.client.Query(context.Background(), `DELETE FROM recipes WHERE id = $1`, recipeID)
	return err
}
