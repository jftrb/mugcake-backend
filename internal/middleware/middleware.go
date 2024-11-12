package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/schema"
	"github.com/jftrb/mugacke-backend/internal/encoders"
	"github.com/jftrb/mugacke-backend/src/api"
	"github.com/rs/zerolog/log"
)

type ContextKey string

const (
	ContextKeyRecipeId ContextKey = "pageToken"
	ContextKeyPagination ContextKey = "pagination"
	ContextKeyCursorParams ContextKey = "cursorParams"
	ContextKeySearchParams ContextKey = "searchParams"
)

var validMethods []string = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
var validHeaders []string = []string{"content-type", "accept"}

func EncodeResponse[T any](w http.ResponseWriter, response T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Err(err).Msg("Error during response JSON encoding.")
		api.InternalErrorHandler(w)
	}
}

func DecodeQueryParams[T any](query map[string][]string) (T, error) {
	var params T
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	decoder.ZeroEmpty(true)
	err := decoder.Decode(&params, query); 
	if err != nil {
		log.Err(err).Msg("Error while decoding query params")
	}

	return params, err
}

func Paginate[T any](defaultLimit int) func (http.Handler) http.Handler {
	return func (next http.Handler) http.Handler {
		return paginate[T](next, defaultLimit)
	}
}

// Parses query for valid pagination parameters and passes them to context.
func paginate[T any](next http.Handler, defaultLimit int) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		pageRequest, err := DecodeQueryParams[api.PaginationRequest](r.URL.Query())
		if err != nil {
			log.Err(err).Msg("Error while decoding Pagination Request from url.")
			return
		}

		if pageRequest.Limit == 0 {
			pageRequest.Limit = defaultLimit
		}

		queryParams, err := DecodeCursorParams[T](pageRequest.Cursor)
		if err != nil {
			api.RequestErrorHandler(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), ContextKeyPagination, pageRequest)
		ctx = context.WithValue(ctx, ContextKeyCursorParams, *queryParams)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

func DecodeCursorParams[T any](cursor string) (*T, error) {
	nextCursor, err := encoders.DecodeBase64(cursor)
	if err != nil {
		log.Err(err).Str("Cursor", cursor).Msg("Error while decoding Cursor from base64.")
		return nil, err
	}

	cursorParams := map[string][]string{}
	
	params := []string{}
	if len(cursor) > 0 {
		params = strings.Split(nextCursor, ",")
	}

	for _, param := range params {
		keyValuePair := strings.Split(param, ":")
		if len(keyValuePair) != 2 {
			err := api.ErrBadQuery
			log.Err(err).Str("KeyValue Pair", param).Msg("Error while trying to parse KeyValue Pair.")
			return nil, err
		}

		key := keyValuePair[0]
		value := keyValuePair[1]
		cursorParams[key] = []string{value}
	}

	out, err := DecodeQueryParams[T](cursorParams)
	if err != nil {
		log.Err(err).Msg("Error while decoding Cursor query params")
		return nil, err
	}
	
	return &out, nil
}

func CorsAllowOrigin(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// TODO : Later, set address to what would be the Mugcake website that would host the web version of the app
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081") 
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// See : https://www.html5rocks.com/static/images/cors_server_flowchart.png
func CorsPreflight(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin"); 
	log.Debug().Str("Origin", origin).Send()
	if origin == "" { // If no origin header
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	accessControlRequestMethod := r.Header.Get("Access-Control-Request-Method"); 
	if accessControlRequestMethod != "" {
		log.Debug().Str("Access-Control-Request-Method", accessControlRequestMethod).Send()
		// Refuse access on invalid method otherwise keep going. 
		if !validateRequestMethod(accessControlRequestMethod) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		accessControlRequestHeaders := r.Header.Get("Access-Control-Request-Headers"); 
		if accessControlRequestHeaders != "" {
			log.Debug().Str("Access-Control-Request-Headers", accessControlRequestHeaders).Send()
			// Refuse access on invalid header.
			if !validateRequestHeaders(accessControlRequestHeaders) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		} 
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(validMethods, ","))
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(validHeaders, ",")) 

		} else {
			// TODO : verify why headers should be exposed
			// Set Access-Control-Expose-Headers if headers should be exposed
		}

		w.WriteHeader(http.StatusOK)
	}

func validateRequestMethod(requestMethod string) bool {
	for _, val := range validMethods {
		if val == requestMethod {
			return true
		}
	}
	return false
}

func validateRequestHeaders(requestHeaders string) bool {
	headers := strings.Split(requestHeaders, ",")
	valid := true
	for _, header := range headers {
		valid = valid && validateRequestHeader(header)
	}
	return valid
}

func validateRequestHeader(requestHeader string) bool {
	for _, val := range validHeaders {
		if val == requestHeader {
			return true
		}
	}
	return false
}
