package main

import (
	"code.google.com/p/gcfg"
	"flag"
	"fmt"
	"github.com/sprintly/sprintly-go/sprintly"
	"github.com/sprintly/uservoice-go/uservoice"
	"log"
	"os"
)

var ticket_number = flag.Int("uservoice_ticket_number", 0, "Ticket number to import from uservoice")
var link_with = flag.Int("link_to", 0, "Sprint.ly ticket number to link it with")
var config_location = flag.String("config_location", "${HOME}/.config/sprintly/uservoice_import.ini", "Location of config file. See README for format.")

type ConfigFile struct {
	Uservoice uservoice.UservoiceConfig
	Sprintly  struct {
		BaseUrl   string
		ProductId int
		Email     string
		ApiKey    string
	}
}

func (cf ConfigFile) validate() error {
	if cf.Uservoice.Subdomain == "" {
		return fmt.Errorf("Invalid Uservoice Subdomain")
	}
	if cf.Uservoice.ApiKey == "" {
		return fmt.Errorf("Invalid Uservoice API Key")
	}
	if cf.Uservoice.ApiSecret == "" {
		return fmt.Errorf("Invalid Uservoice API Secret")
	}
	if cf.Sprintly.BaseUrl == "" {
		return fmt.Errorf("Invalid Sprintly Base URL")
	}
	if cf.Sprintly.ProductId == 0 {
		return fmt.Errorf("Invalid Sprintly Product ID")
	}
	if cf.Sprintly.Email == "" {
		return fmt.Errorf("Invalid Sprintly Email")
	}
	if cf.Sprintly.ApiKey == "" {
		return fmt.Errorf("Invalid Sprintly API Key")
	}
	return nil
}

func main() {
	flag.Parse()
	if *ticket_number == 0 {
		log.Fatal("You must specify a ticket number. View usage for details.")
	}
	if *link_with == 0 {
		log.Print("No linking ticket specified. Will create Sprint.ly ticket.")
	}

	config := ConfigFile{}
	err := gcfg.ReadFileInto(&config, os.ExpandEnv(*config_location))
	if err != nil {
		log.Fatal(err)
	}
	err = config.validate()
	if err != nil {
		log.Fatalf("Invalid config file: %s", err)
	}

	uv_client := uservoice.NewUservoiceClient(config.Uservoice)

	// get ticket from uservoice
	ticket, err := uv_client.GetTicketByNumber(*ticket_number)
	if err != nil {
		log.Fatalf("Error fetching ticket by number: %s\n", err)
	}

	s_client := sprintly.NewSprintlyClient(
		config.Sprintly.Email,
		config.Sprintly.ApiKey,
		config.Sprintly.ProductId,
	)

	var url string
	if *link_with == 0 {
		// post to sprint.ly
		if len(ticket.Messages) > 0 {
			message := fmt.Sprintf("%s\n\n(Link to Uservoice)[%s]",
				ticket.Messages[len(ticket.Messages)-1].PlaintextBody,
				uv_client.UrlForTicket(ticket.Number))
			url, err = s_client.CreateDefect(ticket.Subject, message)
			if err != nil {
				log.Fatalf("Couldn't create Sprint.ly defect: %s\n", err)
			}
		} else {
			log.Fatal("There were no messages in the provided ticket.")
		}
	} else {
		url = s_client.ItemLink(*link_with)
		s_client.AddAnnotation(*link_with,
			"uservoice", "received a message from a customer",
			fmt.Sprintf("[%d](%s)", ticket.Number,
				uv_client.UrlForTicket(ticket.Number)))
	}

	// post url of sprintly ticket as note to uservoice
	uv_client.PostNote(ticket.Id, url)

	fmt.Printf("Ticket ( %s ) posted from Uservoice #%v\n", url, *ticket_number)
}
