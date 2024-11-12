package middleware_test

import (
	"testing"

	"github.com/jftrb/mugacke-backend/internal/encoders"
	"github.com/jftrb/mugacke-backend/internal/middleware"
	"github.com/stretchr/testify/assert"
)

type DummyCursor struct {
	Offset int
	Limit int
	Query string
}

func Test_DecodeCursor(t *testing.T) {
	type args struct {
		cursor  string
	}
	tests := []struct {
		name          string
		args          args
		expectedCursor DummyCursor
	}{
		{name: "Empty cursor", args: args{cursor: ""}, expectedCursor: DummyCursor{Offset: 0, Limit: 0, Query: ""}},
		{name: "Only Offset cursor", args: args{cursor: "offset:10"}, expectedCursor: DummyCursor{Offset: 10, Limit: 0, Query: ""}},
		{name: "Only Limit cursor", args: args{cursor: "limit:10"}, expectedCursor: DummyCursor{Offset: 0, Limit: 10, Query: ""}},
		{name: "Only Query cursor", args: args{cursor: "query:search"}, expectedCursor: DummyCursor{Offset: 0, Limit: 0, Query: "search"}},
		{name: "Full cursor", args: args{cursor: "offset:10,limit:9,query:search"}, expectedCursor: DummyCursor{Offset: 10, Limit: 9, Query: "search"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encodedCursor := encoders.EncodeToBase64(tt.args.cursor)
			actualCursor, err := middleware.DecodeCursorParams[DummyCursor](encodedCursor)
			assert.Nil(t, err)
			assert.Equal(t, tt.expectedCursor, *actualCursor)
		})
	}

}