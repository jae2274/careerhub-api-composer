package posting

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/common/query"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/posting/restapi_grpc"
)

func ExtractJobPostingsRequest(r *http.Request, initPage int) (*JobPostingsRequest, error) {
	queryValues := r.URL.Query()

	// "page" 값 추출
	pageStr := queryValues.Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return nil, fmt.Errorf("invalid page value. %s", pageStr)
	} else if page < initPage {
		return nil, fmt.Errorf("invalid page value. %d", page)
	}

	// "size" 값 추출
	sizeStr := queryValues.Get("size")
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return nil, fmt.Errorf("invalid size value. %s", sizeStr)
	} else if size < 1 || size > 100 {
		return nil, fmt.Errorf("invalid size value. %d", size)
	}

	queryReq, err := ExtractJobPostingQuery(r)
	if err != nil {
		return nil, err
	}

	return &JobPostingsRequest{
		Page:     int32(page),
		Size:     int32(size),
		QueryReq: queryReq,
	}, nil
}

func ExtractJobPostingQuery(r *http.Request) (*query.Query, error) {
	queryValues := r.URL.Query()
	queryReq, err := GetQuery(queryValues.Get("encoded_query"))
	if err != nil {
		return nil, err
	}

	return queryReq, nil
}

func GetQuery(encodedQuery string) (*query.Query, error) {

	bytes, err := base64.StdEncoding.DecodeString(encodedQuery)
	if err != nil {
		query := string(bytes)
		log.Println(query)
		return nil, fmt.Errorf("invalid encoded_query value. failed to decode. %s", encodedQuery)
	}

	var queryReq query.Query
	err = json.Unmarshal(bytes, &queryReq)
	if err != nil {
		return nil, fmt.Errorf("invalid encoded_query value. failed to unmarshal. %s", string(bytes))
	}

	return &queryReq, nil // TODO
}

func ExtractJobPostingDetailRequest(r *http.Request) (*restapi_grpc.JobPostingDetailRequest, error) {
	vars := mux.Vars(r)
	site, ok := vars["site"]
	if !ok {
		return nil, fmt.Errorf("invalid site value. %s", site)
	}

	postingId, ok := vars["postingId"]
	if !ok {
		return nil, fmt.Errorf("invalid postingId value. %s", postingId)
	}

	return &restapi_grpc.JobPostingDetailRequest{
		Site:      site,
		PostingId: postingId,
	}, nil
}

func IsValid(queryReq *query.Query) error {
	if queryReq.MinCareer != nil && *queryReq.MinCareer < 0 {
		return fmt.Errorf("invalid minCareer value. %d", *queryReq.MinCareer)
	}

	if queryReq.MaxCareer != nil && *queryReq.MaxCareer < 0 {
		return fmt.Errorf("invalid maxCareer value. %d", *queryReq.MaxCareer)
	}

	if queryReq.MinCareer != nil && queryReq.MaxCareer != nil && *queryReq.MinCareer > *queryReq.MaxCareer {
		return fmt.Errorf("invalid minCareer and maxCareer value. minCareer(%d) > maxCareer(%d)", *queryReq.MinCareer, *queryReq.MaxCareer)
	}

	return nil
}
