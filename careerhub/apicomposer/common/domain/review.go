package domain

import (
	"time"

	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/review/restapi_grpc"
)

type Review struct {
	Score            int32     `json:"score"`
	Summary          string    `json:"summary"`
	EmploymentStatus bool      `json:"employmentStatus"`
	ReviewUserId     string    `json:"reviewUserId"`
	JobType          string    `json:"jobType"`
	Date             time.Time `json:"date"`
}

func ConvertGrpcToReviews(reviews []*restapi_grpc.Review) []*Review {
	convertedReviews := make([]*Review, 0, len(reviews))
	for _, review := range reviews {
		convertedReviews = append(convertedReviews, ConvertGrpcToReview(review))
	}
	return convertedReviews
}

func ConvertGrpcToReview(review *restapi_grpc.Review) *Review {
	return &Review{
		Score:            review.Score,
		Summary:          review.Summary,
		EmploymentStatus: review.EmploymentStatus,
		ReviewUserId:     review.ReviewUserId,
		JobType:          review.JobType,
		Date:             time.UnixMilli(review.UnixMilli),
	}
}
