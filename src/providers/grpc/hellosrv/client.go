package hellosrv

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"

	// TODO: should be replace to an appropriate api
	extapi "github.com/akhripko/dummy/api"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	Target string
}

var ErrInit = errors.New("initialisation error")
var ErrConn = errors.New("connection error")

func New(ctx context.Context, cfg Config) (*Client, error) {
	s := Client{
		ctx: ctx,
	}

	// init connection
	var err error
	s.conn, err = grpc.DialContext(
		ctx,
		cfg.Target,
		grpc.WithUnaryInterceptor(grpc_prometheus.UnaryClientInterceptor),
		grpc.WithStreamInterceptor(grpc_prometheus.StreamClientInterceptor),
		grpc.WithInsecure(),
		//grpc.WithBlock(),
	)
	if err != nil {
		log.Info("providers.grpc.hellosrv init err:", err)
		return nil, errors.Wrap(err, "providers.grpc.hellosrv: failed to init connection")
	}

	// handle shutdown
	go func() {
		<-ctx.Done()
		if s.conn != nil {
			log.Info("providers.grpc.hellosrv: close connection")
			s.conn.Close()
		}
	}()

	s.client = extapi.NewDummyServiceClient(s.conn)

	return &s, nil
}

type Client struct {
	ctx    context.Context
	conn   *grpc.ClientConn
	client extapi.DummyServiceClient
}

func (s *Client) Check() error {
	if s.conn == nil {
		return ErrInit
	}

	if s.conn.GetState() != connectivity.Ready {
		return ErrConn
	}

	return nil
}
