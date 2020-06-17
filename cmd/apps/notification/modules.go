package main

import (
	"context"
	"github.com/Sirupsen/logrus"
	"github.com/gidyon/micros"

	channel_app "github.com/gidyon/umrs/internal/apps/notification/channel"
	notification_app "github.com/gidyon/umrs/internal/apps/notification/notification"
	subscriber_app "github.com/gidyon/umrs/internal/apps/notification/subscriber"
	"github.com/gidyon/umrs/pkg/api/channel"
	"github.com/gidyon/umrs/pkg/api/notification"
	"github.com/gidyon/umrs/pkg/api/subscriber"
)

func registerModules(ctx context.Context, app *micros.Service) {
	// email dialer
	dialer, err := emailDialer()
	handleErr(err)

	opt := &notification_app.Options{
		SQLDB:       app.GormDB(),
		RedisClient: app.RedisClient(),
		SMTPDialer:  dialer,
	}
	// notification module
	notificationModule, err := notification_app.NewNotificationServiceServer(ctx, opt)
	handleErr(err)

	// channel module
	channelModule, err := channel_app.NewChannelAPIServerServer(ctx, app.GormDB())
	handleErr(err)

	// subscriber module
	subscriberModule, err := subscriber_app.NewSubscriberAPIServer(ctx, app.GormDB())
	handleErr(err)

	// Register app modules to the gRPC server
	notification.RegisterNotificationServiceServer(app.GRPCServer(), notificationModule)
	channel.RegisterChannelAPIServer(app.GRPCServer(), channelModule)
	subscriber.RegisterSubscriberAPIServer(app.GRPCServer(), subscriberModule)

	// Register client(s) for reverse proxy server
	handleErr(notification.RegisterNotificationServiceHandler(ctx, app.RuntimeMux(), app.ClientConn()))
	handleErr(channel.RegisterChannelAPIHandler(ctx, app.RuntimeMux(), app.ClientConn()))
	handleErr(subscriber.RegisterSubscriberAPIHandler(ctx, app.RuntimeMux(), app.ClientConn()))
}

func handleErr(err error) {
	if err != nil {
		logrus.Fatalln(err)
	}
}
