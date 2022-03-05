package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	tele "gopkg.in/telebot.v3"
	"guarantorplace.com/internal/data"
)

func(app *application) requestForTransaction(ch chan bool, msg string, errs chan error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://toncenter.com/api/v2/getTransactions?address=%s&limit=100&to_lt=0&archival=false&api_key=%s", app.config.Wallet, app.config.ApiKey), nil)
	if err != nil {
		errs <- err
	}

	resp, err := client.Do(req)

	fmt.Println(resp.StatusCode, resp.Body)

	if err != nil {
		fmt.Println("Errored when sending request to the server")
		errs <- err
	}

	var res data.Response

	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&res)

	for _, item := range res.Result {
		for _, outMsgs := range item.OutMessage {

			value, err := strconv.Atoi(outMsgs.Value)

			if err != nil {
				errs <- err
			}

			ch <- false

			fmt.Println(false)

			if outMsgs.Message == msg && math.Floor(float64(value)*100)/100000000000 == float64(app.config.Fee) {
				ch <- true
			}

		}
	}

}


func generateUniqueMessage() (string, error) {
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))	  
	const chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	bytes := make([]byte, 16)
 
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	for i := range bytes {
		bytes[i] = chars[seededRand.Intn(len(chars))]
	}
 
	return string(bytes), nil
}


func(app *application) checkTransaction(done chan bool, b *tele.Bot, link string, chat *tele.Chat, message string, errs chan error) {

	_, err := b.Send(chat, app.config.Messages.AfterPaymentResponse, &tele.SendOptions{ParseMode: "MarkdownV2"})

	if err != nil {
		errs <- err
	}

	ch := make(chan bool)

	ticker := time.NewTicker(3 * time.Second)

	out:for range ticker.C {
		select {
			case <- done:
				_, err := b.Send(chat, app.config.Messages.FailedPaymentResponse, &tele.SendOptions{ParseMode: "MarkdownV2"})

				if err != nil {
					errs <- err
				}

				ticker.Stop()
				break out
			case status := <-ch: 
				if status {
				target, err := b.ChatByID(app.config.ExchangeChannel)

				if err != nil {
					errs <- err
				}

				_, err = b.Send(target, link)
				if err != nil {
					errs <- err
				}
				
				ad, err := app.models.Ads.GetByMessage(message)

				if err != nil {
					switch {
					case errors.Is(err, data.ErrRecordNotFound):
						fmt.Println("No record in DB")
						errs <- err
					default:
						errs <- err
					}
				}

				input := &data.Ad{
				   	UserId:	ad.UserId,
					Msg: ad.Msg,
					Link: ad.Link,
					Paid: true,
					CreatedAt: ad.CreatedAt,
					ID: ad.ID,
					Version: ad.Version,
				}
			
				err = app.models.Ads.Update(input)
				
				if err != nil {
					switch {
					case errors.Is(err, data.ErrEditConflict):
						errs <- err
					default:
						errs <- err
					}
				}

				_, err = b.Send(chat, app.config.Messages.SuccessPaymentResponse, &tele.SendOptions{ParseMode:"MarkdownV2"})

				if err != nil {
					errs <- err
				}
		
			
				ticker.Stop()
			}else {
				continue
			}
		case <-ticker.C:
				go app.requestForTransaction(ch, message, errs)
		}
	
	}

}
