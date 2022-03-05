package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"
	"guarantorplace.com/internal/data"
)



func (app *application) handleApproveChannel() error {

	var (
		selector = &tele.ReplyMarkup{}
	
		btnApprove = selector.Data(app.config.Messages.ApprovePublicPageBtn, "approve")
		btnDecline = selector.Data(app.config.Messages.DeclinePublicPageBtn, "decline")
	
	)	



	
	adminBoard, err := app.bot.ChatByID(app.config.AdminChannel)


	selector.Inline(
		selector.Row(btnApprove, btnDecline),
	)		
	

	
	app.bot.Handle(tele.OnText, func(c tele.Context) error {
		if err != nil {
			return err
		}

		if c.Message().Sender.Username  == ""	{
		_, err = app.bot.Send(adminBoard, fmt.Sprintf("%d\n%s", c.Message().Sender.ID, c.Text()), selector)

		if err != nil {
			return err
		}

		}else {
			_, err = app.bot.Send(adminBoard, fmt.Sprintf("%d\n%s\n\n%s", c.Message().Sender.ID, "@" + c.Message().Sender.Username, c.Text()), selector)
			if err != nil {
				return err
			}
	
		}

	
		_, err := app.bot.Send(c.Sender(), app.config.Messages.AfterSubmittingPublicPage, &tele.SendOptions{ParseMode:"MarkdownV2"})

		if err != nil {
			return err
		}

		app.bot.Handle(tele.OnText, func(c tele.Context) error {
			return nil
		})
		return nil
	})

	
	app.bot.Handle(&btnApprove, func(c tele.Context) error {
		var (
			btnPaid = selector.Data(app.config.Messages.PaidBtn, "paid")
			btnDeclinePaid = selector.Data(app.config.Messages.DeclinePaidBtn, "decline_paid")

		)
		selector.Inline(
			selector.Row(btnPaid),
			selector.Row(btnDeclinePaid),
		)

		arrayOfString := strings.Split(c.Text(), "\n")
		

		id, err := strconv.ParseInt(arrayOfString[0], 10, 64)

		if err != nil {
			return err
		}


		chat, err := app.bot.ChatByID(id)


		if err != nil {
			return err
		}

		message, err := generateUniqueMessage()

		if err != nil {
			return err
		}


		ad := &data.Ad{
				UserId:   id,
				Link:    arrayOfString[len(arrayOfString)-1],
				Msg: message,
		}


		err = app.models.Ads.Insert(ad)
		if err != nil {
			return err
		}



		paymentMessage := fmt.Sprintf(app.config.Messages.PaymentMessage, "0\\.\\5", app.config.Messages.PaidBtn, app.config.Wallet, message,)
		
		_, err = app.bot.Send(chat, paymentMessage, &tele.SendOptions{ParseMode: "MarkdownV2", ReplyMarkup: selector})
		
		if err != nil {
			return err
		}

		done := make(chan bool) 


		app.bot.Handle(&btnPaid, func(c tele.Context) error {


			ad, err := app.models.Ads.Get(c.Chat().ID)

			if err != nil {
				return err
			}

			errs := make(chan error, 1)

			go app.checkTransaction(done, app.bot, ad.Link, c.Chat(), ad.Msg, errs)

			if err := <-errs; err != nil {
				return err
			}
			
			time.Sleep(time.Minute * 5)
			close(done)


			return nil

		})


		

		app.bot.Handle(&btnDeclinePaid, func(c tele.Context) error {	
			close(done)

			return nil
		})



		

		return nil
	})

	app.bot.Handle(&btnDecline, func(c tele.Context) error {

		id, err := strconv.ParseInt(strings.Split(c.Text(), "\n")[0], 10, 64)

		if err != nil {
			return err
		}

		chat, err := app.bot.ChatByID(id)


		if err != nil {
			return err
		}
  
		_, err = app.bot.Send(chat, app.config.Messages.RejectPublicPage, &tele.SendOptions{
			ParseMode: "MarkdownV2",
		})

		if err != nil {
			return err
		}

		return nil
	})

	return nil

}

func (app *application) handleStartCommand() error {

	var (
		mainMenu = &tele.ReplyMarkup{}
	
		btnChoosePlace = mainMenu.Data(app.config.Messages.ChoosePublicPageBtn, "place")
		btnCreateRequest = mainMenu.Data(app.config.Messages.CreateRequestBtn, "create_request")

		photo = &tele.Photo{File: tele.FromURL("https://ibb.co/G5mHG0w")}
	)


	
	mainMenu.Inline(
		mainMenu.Row(btnChoosePlace),
		mainMenu.Row(btnCreateRequest),
	)
	
	
	app.bot.Handle("/start", func(c tele.Context) error {
		return c.Send(photo, mainMenu)
	})

	app.bot.Handle(&btnChoosePlace, func(c tele.Context) error {
		return c.Respond(&tele.CallbackResponse{Text: "В разработке!"})
	})

	app.bot.Handle(&btnCreateRequest, func(c tele.Context) error {

		var backToMenu = &tele.ReplyMarkup{}

		var btnBack = backToMenu.Data(app.config.Messages.BackBtn, "back")

		backToMenu.Inline(
			backToMenu.Row(btnBack),
		)

		c.Delete()

		_, err := app.bot.Send(c.Sender(), app.config.Messages.PutPublicPage, &tele.SendOptions{ParseMode:"MarkdownV2", ReplyMarkup: backToMenu})
		if err != nil {
			return err
		}
		
		app.bot.Handle(&btnBack, func(c tele.Context) error {
			c.Delete()
			app.bot.Handle(tele.OnText, func(c tele.Context) error {
				return nil;
			})		
			return c.Send(photo, mainMenu)
		})

			
		err = app.handleApproveChannel()

		if err != nil {
			return err
		}

		return nil

	})
	return nil


}


func (app *application) handleUpdates() error {
	err := app.handleStartCommand()

	if err != nil {
		return err
	}
	return nil
}