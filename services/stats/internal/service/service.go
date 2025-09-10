package service

import (
	"context"
	"soa-socialnetwork/services/stats/internal/repo"
	"soa-socialnetwork/services/stats/internal/storage/clickhouse"
	"soa-socialnetwork/services/stats/pkg/models"
	pb "soa-socialnetwork/services/stats/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type StatsService struct {
	pb.UnimplementedStatsServiceServer

	Db repo.Database
}

func New(cfg Config) (StatsService, error) {
	db, err := clickhouse.NewDB(clickhouse.Config{
		Hostname:     cfg.DbHost,
		Port:         cfg.DbPort,
		Database:     "stats_clickhouse",
		Username:     cfg.DbUser,
		Password:     cfg.DbPassword,
		MaxIdleConns: 10,
		MaxOpenConns: 5,
	})
	if err != nil {
		return StatsService{}, err
	}

	return StatsService{
		Db: &db,
	}, nil
}

func (s *StatsService) GetPostMetric(ctx context.Context, req *pb.GetPostMetricRequest) (*pb.GetPostMetricResponse, error) {
	if req.Metric == pb.Metric_METRIC_UNSPECIFIED {
		return nil, status.Error(codes.InvalidArgument, "metric is not specified")
	}

	count, err := func() (int64, error) {
		postId := models.PostId(req.PostId)
		switch req.Metric {
		case pb.Metric_METRIC_VIEW_COUNT:
			{
				count, err := s.Db.PostsViews(ctx).GetCountForPost(postId)
				if err != nil {
					return -1, err
				}

				return count, nil
			}

		case pb.Metric_METRIC_LIKE_COUNT:
			{
				count, err := s.Db.PostsLikes(ctx).GetCountForPost(postId)
				if err != nil {
					return -1, err
				}

				return count, nil
			}

		case pb.Metric_METRIC_COMMENT_COUNT:
			{
				count, err := s.Db.PostsComments(ctx).GetCountForPost(postId)
				if err != nil {
					return -1, err
				}

				return count, nil
			}
		}

		return -1, status.Error(codes.Internal, "unknown metric")
	}()

	if err != nil {
		return nil, err
	}

	return &pb.GetPostMetricResponse{
		Count: count,
	}, nil
}

func (s *StatsService) GetPostMetricDynamics(ctx context.Context, req *pb.GetPostMetricDynamicsRequest) (*pb.GetPostMetricDynamicsResponse, error) {
	if req.Metric == pb.Metric_METRIC_UNSPECIFIED {
		return nil, status.Error(codes.InvalidArgument, "metric is not specified")
	}

	postId := models.PostId(req.PostId)
	switch req.Metric {
	case pb.Metric_METRIC_VIEW_COUNT:
		{
			dynamics, err := s.Db.PostsViews(ctx).GetDynamicsForPost(postId)
			if err != nil {
				return nil, err
			}

			pbDynamics := make([]*pb.DayDynamics, len(dynamics))
			for i, dayStat := range dynamics {
				pbDynamics[i] = &pb.DayDynamics{
					Date:  timestamppb.New(dayStat.Date),
					Count: int32(dayStat.Count),
				}
			}

			return &pb.GetPostMetricDynamicsResponse{
				Dynamics: pbDynamics,
			}, nil
		}

	case pb.Metric_METRIC_LIKE_COUNT:
		{
			dynamics, err := s.Db.PostsLikes(ctx).GetDynamicsForPost(postId)
			if err != nil {
				return nil, err
			}

			pbDynamics := make([]*pb.DayDynamics, len(dynamics))
			for i, dayStat := range dynamics {
				pbDynamics[i] = &pb.DayDynamics{
					Date:  timestamppb.New(dayStat.Date),
					Count: int32(dayStat.Count),
				}
			}

			return &pb.GetPostMetricDynamicsResponse{
				Dynamics: pbDynamics,
			}, nil
		}

	case pb.Metric_METRIC_COMMENT_COUNT:
		{
			dynamics, err := s.Db.PostsComments(ctx).GetDynamicsForPost(postId)
			if err != nil {
				return nil, err
			}

			pbDynamics := make([]*pb.DayDynamics, len(dynamics))
			for i, dayStat := range dynamics {
				pbDynamics[i] = &pb.DayDynamics{
					Date:  timestamppb.New(dayStat.Date),
					Count: int32(dayStat.Count),
				}
			}

			return &pb.GetPostMetricDynamicsResponse{
				Dynamics: pbDynamics,
			}, nil
		}

	default:
		return nil, status.Error(codes.Internal, "unknown metric")
	}
}
