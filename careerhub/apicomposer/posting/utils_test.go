package posting

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeQuery(t *testing.T) {

	_, err := GetQuery("eyJwYWdlIjowLCJzaXplIjoxNiwiY2F0ZWdvcmllcyI6W3sic2l0ZSI6Imp1bXBpdCIsImNhdGVnb3J5TmFtZSI6IuybuSDtkoDsiqTtg50g6rCc67Cc7J6QIn0seyJzaXRlIjoianVtcGl0IiwiY2F0ZWdvcnlOYW1lIjoi7J246rO17KeA64qlL+uouOyLoOufrOuLnSJ9XSwic2tpbGxOYW1lcyI6W10sInRhZ0lkcyI6W10sIm1pbkNhcmVlciI6bnVsbCwibWF4Q2FyZWVyIjpudWxsfQ==")
	require.NoError(t, err)
}
