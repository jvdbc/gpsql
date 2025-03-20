package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "smtphost",
				Value: "email-smtp.eu-west-3.amazonaws.com",
				Usage: "AWS SES SMTP host",
			},
			&cli.StringFlag{
				Name:  "from",
				Value: "mail-sender@suite.woonoz.com",
				Usage: "origin of email",
			},
			&cli.StringFlag{
				Name:  "username",
				Value: "akia",
				Usage: "name of user",
			},
			&cli.StringFlag{
				Name:  "password",
				Value: "poney",
				Usage: "password",
			},
			&cli.UintFlag{
				Name:  "port",
				Value: 587,
				Usage: "smtp port",
			},
			&cli.StringFlag{
				Name:  "to",
				Value: "jvasdeboncoeur@ippon.fr",
				Usage: "destination of email",
			},
		},
		Name:  "gpsql",
		Usage: "Test a postgresql connection",
		Action: func(ctx *cli.Context) error {
			smtpHost := ctx.String("smtphost")
			smtpPort := ctx.Uint("port")
			from := ctx.String("from")
			username := ctx.String("username")
			password := ctx.String("password")
			to := ctx.String("to")

			// Contenu de l'e-mail
			subject := "Test d'envoi d'e-mail avec AWS SES"
			body := "Bonjour,\n\nVoici un test d'email envoyé via AWS SES avec Go !"

			// Créer le message

			// Envoi de l'e-mail

			log.Println("e-mail envoyé avec succès !")

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
