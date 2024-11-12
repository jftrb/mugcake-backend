package middleware

import (
	"context"
	"net/http"

	"github.com/jftrb/mugacke-backend/src/api"
	"github.com/rs/zerolog/log"
)

func ParseSummariesSearchParams(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		searchParams, err := DecodeQueryParams[api.RecipeSearchRequest](r.URL.Query())
		if err != nil {
			api.RequestErrorHandler(w, err)
			return
		}

		log.Debug().Msg(r.URL.Query().Get("tags"))
		log.Debug().Int("tags length", len(searchParams.Tags)).Send()
		if searchParams.Tags == nil {
			log.Debug().Msg("nil")
			searchParams.Tags = []string{}
		}

		if searchParams.SortBy == nil {
			searchParams.SortBy = []api.Sort{}
		}

		ctx := context.WithValue(r.Context(), ContextKeySearchParams, searchParams)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}