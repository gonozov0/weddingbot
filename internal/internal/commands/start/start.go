package start

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gonozov0/weddingtgbot/internal/internal/commands/shared"
	"github.com/gonozov0/weddingtgbot/internal/internal/commands/shared/owner_chat"
	"github.com/gonozov0/weddingtgbot/pkg/logger"
)

type DTO struct {
	ChatID int64
	Login  string
}

func Do(bot *tgbotapi.BotAPI, dto DTO) *logger.SlogError {
	if dto.Login == "" {
		return requestPhoneNumber(bot, dto.ChatID)
	}

	if !isLoginInvited(dto.Login) {
		return shared.SendNotInvitedInfo(bot, dto.ChatID)
	}

	personInfo := shared.GetPersonInfo(dto.Login)
	if err := owner_chat.SendStart(bot, personInfo.GetFullName()); err != nil {
		return err
	}

	return sendInvitation(bot, dto.ChatID, personInfo.Name)
}

func requestPhoneNumber(bot *tgbotapi.BotAPI, chatID int64) *logger.SlogError {
	photoGroup := tgbotapi.NewPhoto(chatID, tgbotapi.FileID(photoFileID))
	if _, err := bot.Send(photoGroup); err != nil {
		return logger.NewSlogError(err, "error sending photo")
	}

	msg := tgbotapi.NewMessage(
		chatID,
		invitationAnonText,
	)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButtonContact("Отправить мой номер боту для идентификации"),
		),
	)
	if _, err := bot.Send(msg); err != nil {
		return logger.NewSlogError(err, "error sending message")
	}

	return nil
}

func sendInvitation(bot *tgbotapi.BotAPI, chatID int64, name string) *logger.SlogError {
	photoGroup := tgbotapi.NewPhoto(chatID, tgbotapi.FileID(photoFileID))
	if _, err := bot.Send(photoGroup); err != nil {
		return logger.NewSlogError(err, "error sending photo")
	}

	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(invitationGuestText, name))
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = shared.GetStartReplyKeyboard()
	if _, err := bot.Send(msg); err != nil {
		return logger.NewSlogError(err, "error sending message")
	}

	return nil
}
