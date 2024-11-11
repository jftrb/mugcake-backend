package dbWrapper

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jftrb/mugacke-backend/src/api"
	"github.com/jftrb/mugacke-backend/src/api/models"
	"github.com/rs/zerolog/log"
)

type DbWrapper interface {
	Disconnect() error
	GetUsers() ([]models.User, error)
	GetRecipeSummaries(userID string, paginationToken api.RecipeSummaryPaginationRequest, searchParams api.RecipeSearchRequest) ([]models.RecipeSummary, error)
	GetRecipe(recipeID int) (models.Recipe, error)
	AddRecipe(userID string, recipe models.Recipe) (int, error)
	PutRecipe(recipeID int, recipe models.Recipe) (error)
	PatchRecipe(recipeID int, favorite bool) (error)
	DeleteRecipe(recipeID int) (error)
}

type dbwrapper struct {
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

	return &dbwrapper{client: client}
}

func (d *dbwrapper) Disconnect() error {
	d.client.Release()
	log.Info().Msg("Released a connection from Postgres.")
	d.client = nil
	return nil
}

func (d *dbwrapper) GetUsers() ([]models.User, error) {
	rows, err := d.client.Query(context.Background(), "SELECT (id, email) from users")
	if err != nil {
		return []models.User{}, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (models.User, error) {
		var user models.User
		err := row.Scan(&user)
		return user, err
	})
}

func (d *dbwrapper) GetRecipeSummaries(userID string, paginationToken api.RecipeSummaryPaginationRequest, searchParams api.RecipeSearchRequest) ([]models.RecipeSummary, error) {
	sortQuery := BuildSortQuery(searchParams.SortBy)
	searchQuery := BuildSearchQuery(searchParams)
	rows, err := d.client.Query(context.Background(), 
	`SELECT id as recipeId, favorite, title, image_source, prep_info->>'totalTime' as total_time
			, ARRAY (
						SELECT tags.name
						FROM   unnest(r.tags) WITH ORDINALITY AS a(tag_id, ord)
						JOIN   tags ON tags.id = a.tag_id
						ORDER  BY a.ord
						)  AS tags
	FROM recipes r 
	WHERE r.user_id = $1 AND ` + searchQuery + " " + sortQuery, userID)

	if err != nil {
		return []models.RecipeSummary{}, err
	}
	defer rows.Close()

	log.Debug().Msg("Successful GetRecipesSummaries Query")

	return pgx.CollectRows(rows, pgx.RowToStructByName[models.RecipeSummary])
}

func (d *dbwrapper) GetRecipe(recipeID int) (models.Recipe, error) {
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

	recipe, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.Recipe])
	
	d.client.Exec(context.Background(),
	`UPDATE recipes SET last_viewed = current_timestamp WHERE id = $1`, recipeID)
	
	return recipe, err
}

func (d *dbwrapper)	AddRecipe(userID string, recipe models.Recipe) (int, error) {
	prepInfo, err := json.Marshal(recipe.PrepInfo)
	if err != nil {
		return -1, err
	}
	
	ingredientSections, err := json.Marshal(recipe.IngredientSections)
	if err != nil {
		return -1, err
	}

	err = d.AddTags(recipe.Tags)
	if err != nil {
		return -1, err
	}
	
	timeoutCtx, cancelFunc := context.WithTimeout(context.Background(), time.Second * 2)
	rows, err := d.client.Query(timeoutCtx, 
	`INSERT INTO recipes (user_id, favorite, title, url, image_source, prep_info, tags, ingredientSections, directions, notes)
	VALUES (
		$1, $2, $3, $4, $5, $6,
		ARRAY (
			SELECT tags.id
			FROM	 unnest($7::text[]) WITH ORDINALITY AS a(tag_name, ord)
			JOIN   tags ON tags.name = a.tag_name
			ORDER  BY a.ord
		),
		$8, $9, $10)
		RETURNING id`, userID, recipe.Favorite, recipe.Title, recipe.URL, recipe.ImageSource, prepInfo, recipe.Tags, ingredientSections, recipe.Directions, recipe.Notes)
	cancelFunc()
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	type PostReturn struct {
		ID int
	}
	postReturn, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[PostReturn]);
	
	return postReturn.ID, err
}

// Doesn't update Favorite status. Branched out to separate call
func (d *dbwrapper)	PutRecipe(recipeID int, recipe models.Recipe) (error) {
	prepInfo, err := json.Marshal(recipe.PrepInfo)
	if err != nil {
		return err
	}
	
	ingredientSections, err := json.Marshal(recipe.IngredientSections)
	if err != nil {
		return err
	}

	err = d.AddTags(recipe.Tags)
	if err != nil {
		return err
	}

	timeoutCtx, cancelFunc := context.WithTimeout(context.Background(), time.Second * 2)
	_, err = d.client.Exec(timeoutCtx, 
	`UPDATE recipes SET (title, url, image_source, prep_info, tags, ingredientSections, directions, notes, modified) = 
	(
		$1, $2, $3, $4,
		ARRAY (
			SELECT tags.id
			FROM	 unnest($5::text[]) WITH ORDINALITY AS a(tag_name, ord)
			JOIN   tags ON tags.name = a.tag_name
			ORDER  BY a.ord
		),
		$6, $7, $8, current_timestamp)
		WHERE id = $9`, recipe.Title, recipe.URL, recipe.ImageSource, prepInfo, recipe.Tags, ingredientSections, recipe.Directions, recipe.Notes, recipeID)

	cancelFunc()

	return err
}

func (d *dbwrapper) PatchRecipe(recipeID int, favorite bool) (error) {

	timeoutCtx, cancelFunc := context.WithTimeout(context.Background(), time.Second * 2)
	_, err := d.client.Exec(timeoutCtx, 
	`UPDATE recipes SET favorite = $1 WHERE id = $2`, favorite, recipeID)

	cancelFunc()
	return err
}

func (d *dbwrapper)	DeleteRecipe(recipeID int) (error) {
	_, err := d.client.Exec(context.Background(), `DELETE FROM recipes WHERE id = $1`, recipeID)
	return err
}

func (d *dbwrapper)	AddTags(tags []string) (error) {
	timeoutCtx, cancelFunc := context.WithTimeout(context.Background(), time.Second * 2)
	_, err := d.client.Exec(timeoutCtx, 
	`INSERT INTO tags (name) VALUES (unnest($1::text[])) ON CONFLICT (name) DO NOTHING`, 
	tags)
	cancelFunc()

	return err
}
