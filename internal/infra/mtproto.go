package infra

import (
	"context"
	"github.com/JscorpTech/paymento/internal/usecase"
	"github.com/gotd/td/examples"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/telegram/updates"
	updhook "github.com/gotd/td/telegram/updates/hook"
	"github.com/gotd/td/tg"

	"github.com/go-faster/errors"
	"go.uber.org/zap"
)

func mtproto(ctx context.Context, log *zap.Logger, watch_id int64) error {

	d := tg.NewUpdateDispatcher()
	gaps := updates.New(updates.Config{
		Handler: d,
		Logger:  log.Named("gaps"),
	})

	flow := auth.NewFlow(examples.Terminal{}, auth.SendCodeOptions{})

	client, err := telegram.ClientFromEnvironment(telegram.Options{
		Logger:        log,
		UpdateHandler: gaps,
		Middlewares: []telegram.Middleware{
			updhook.UpdateHook(gaps.Handle),
		},
	})
	if err != nil {
		return err
	}

	d.OnNewMessage(func(ctx context.Context, e tg.Entities, update *tg.UpdateNewMessage) error {
		msg, ok := update.Message.(*tg.Message)
		if !ok {
			return nil
		}

		text := msg.Message
		if text == "" {
			return nil
		}

		var senderUsername string
		var senderID int64

		switch from := msg.FromID.(type) {
		case *tg.PeerUser:
			senderID = from.UserID
		case *tg.PeerChannel:
			return nil
		case *tg.PeerChat:
			return nil
		default:
			return nil
		}

		if !isWatched(senderID, senderUsername, watch_id) {
			return nil
		}

		if res := usecase.ParseTopUp(text, log); res != nil {
			log.Info("To'ldirish aniqlandi",
				zap.String("raw", res.AmountRaw),
				zap.Int64("raw", res.AmountInt),
			)
		} else {
			log.Debug("Xabar top_up emas", zap.String("text", limit(text, 120)))
		}

		return nil
	})

	return client.Run(ctx, func(ctx context.Context) error {
		if err := client.Auth().IfNecessary(ctx, flow); err != nil {
			return errors.Wrap(err, "auth")
		}
		user, err := client.Self(ctx)
		if err != nil {
			return errors.Wrap(err, "call self")
		}
		return gaps.Run(ctx, client.API(), user.ID, updates.AuthOptions{
			OnStart: func(ctx context.Context) {
				log.Info("Gaps started")
			},
		})
	})
}

func isWatched(id int64, username string, watch_id int64) bool {
	if watch_id != 0 && id == watch_id {
		return true
	}
	return false
}

func limit(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
