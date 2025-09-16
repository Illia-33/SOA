package service

import (
	"context"
	"soa-socialnetwork/services/stats/internal/kafka"
	"soa-socialnetwork/services/stats/internal/repo"
	kafkajobs "soa-socialnetwork/services/stats/internal/service/jobs/kafka"
	"soa-socialnetwork/services/stats/internal/storage/clickhouse"
	"soa-socialnetwork/services/stats/pkg/models"
	pb "soa-socialnetwork/services/stats/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type StatsService struct {
	pb.UnimplementedStatsServiceServer

	Db             repo.Database
	KafkaProcessor kafkajobs.Processor
}

func New(cfg Config) (StatsService, error) {
	db, err := clickhouse.NewDB(clickhouse.Config{
		Hostname:     cfg.DbHost,
		Port:         cfg.DbPort,
		Database:     "default",
		Username:     cfg.DbUser,
		Password:     cfg.DbPassword,
		MaxIdleConns: 10,
		MaxOpenConns: 5,
	})
	if err != nil {
		return StatsService{}, err
	}

	kafkaProcessor, err := kafkajobs.NewProcessor(kafka.ConnectionConfig{
		Host: cfg.KafkaHost,
		Port: cfg.KafkaPort,
	}, &db)
	if err != nil {
		return StatsService{}, nil
	}

	return StatsService{
		Db:             &db,
		KafkaProcessor: kafkaProcessor,
	}, nil
}

func (s *StatsService) Start() {
	s.KafkaProcessor.Start(context.Background())
}

func (s *StatsService) GetPostMetric(ctx context.Context, req *pb.GetPostMetricRequest) (*pb.GetPostMetricResponse, error) {
	if req.Metric == pb.Metric_METRIC_UNSPECIFIED {
		return nil, status.Error(codes.InvalidArgument, "metric is not specified")
	}

	count, err := func() (uint64, error) {
		postId := models.PostId(req.PostId)
		switch req.Metric {
		case pb.Metric_METRIC_VIEW_COUNT:
			{
				return s.Db.PostsViews(ctx).GetCountForPost(postId)
			}

		case pb.Metric_METRIC_LIKE_COUNT:
			{
				return s.Db.PostsLikes(ctx).GetCountForPost(postId)
			}

		case pb.Metric_METRIC_COMMENT_COUNT:
			{
				return s.Db.PostsComments(ctx).GetCountForPost(postId)
			}
		}

		return 0, status.Error(codes.Internal, "unknown metric")
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

func (s *StatsService) GetTop10Posts(ctx context.Context, req *pb.GetTop10PostsRequest) (*pb.GetTop10PostsResponse, error) {
	if req.Metric == pb.Metric_METRIC_UNSPECIFIED {
		return nil, status.Error(codes.InvalidArgument, "metric is not specified")
	}

	metric := metricFromProto(req.Metric)
	postAggs, err := s.Db.Aggregation(ctx).GetTop10PostsByMetric(metric)
	if err != nil {
		return nil, err
	}

	var postStats []*pb.PostStats
	for i := range postAggs {
		postStats = append(postStats, postStatsFromAgg(postAggs[i], metric))
	}

	return &pb.GetTop10PostsResponse{
		Posts: postStats,
	}, nil
}

func (s *StatsService) GetTop10Users(ctx context.Context, req *pb.GetTop10UsersRequest) (*pb.GetTop10UsersResponse, error) {
	if req.Metric == pb.Metric_METRIC_UNSPECIFIED {
		return nil, status.Error(codes.InvalidArgument, "metric is not specified")
	}

	metric := metricFromProto(req.Metric)
	userAggs, err := s.Db.Aggregation(ctx).GetTop10UsersByMetric(metric)
	if err != nil {
		return nil, err
	}

	var userStats []*pb.UserStats
	for i := range userAggs {
		userStats = append(userStats, userStatsFromAgg(userAggs[i], metric))
	}

	return &pb.GetTop10UsersResponse{
		Users: userStats[:],
	}, nil
}
