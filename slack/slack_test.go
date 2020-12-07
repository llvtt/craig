package slack

import (
	"encoding/json"
	"fmt"
	"github.com/slack-go/slack"
	"testing"
)

const webhook = ""
const craigAccessToken = ""
const clTestChannelId = "C01BDF1CCP3"

const workingImageUrl = "https://addisonchoate.com/wp-content/uploads/2020/01/Antiques.jpg"
const brokenImageUrl = "https://images.craigslist.org/00H0H_3uR69KiZUnz_0CI0t2_1200x900.jpg"

func TestWebhook(t *testing.T) {
	attachment := slack.Attachment{
		Fallback:      "You successfully posted by Incoming Webhook URL!",
		Text:          "Test image",
		//ImageURL:      "https://images.craigslist.org/00H0H_3uR69KiZUnz_0CI0t2_1200x900.jpg",
		ImageURL:      "https://addisonchoate.com/wp-content/uploads/2020/01/Antiques.jpg",
	}
	msg := slack.WebhookMessage{
		Text: "Item from craig",
		Attachments: []slack.Attachment{attachment},
	}
	err := slack.PostWebhook(webhook, &msg)
	if err != nil {
		fmt.Println(err)
	}
}

func TestSlack(t *testing.T) {
	api := slack.New(craigAccessToken, slack.OptionDebug(true))
	b := buildCraigBlocks()

	//jsonBlocks := marshalBlocks(b)
	//fmt.Println(jsonBlocks)

	//msg := slack.WebhookMessage{Text: "P!"}

	//err := slack.PostWebhook(webhook, &msg)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//slack.MsgOptionBlocks(b...)
	//_, _, err := api.PostMessage(clTestChannelId, slack.MsgOptionText("Some text", false))

	_, _, err := api.PostMessage(clTestChannelId, slack.MsgOptionBlocks(b...))
	//_, _, err := api.PostMessage(clTestChannelId, slack.MsgOptionBlocks(b...))
	if err != nil {
		fmt.Println(err)
	}
}

func marshalBlocks(msg slack.Message) string {
	b, err := json.MarshalIndent(msg, "", "    ")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(b)
}

func buildCraigBlocks() []slack.Block {
	//titleBlock := slack.NewTextBlockObject("plain_text", "Image from craigslist", true, false)
	//slack.NewAccessory(slack.NewImageBlockElement("https://images.craigslist.org/00H0H_3uR69KiZUnz_0CI0t2_1200x900.jpg", "Derp", "imgblock", titleBlock))
	//imgBlock := slack.NewImageBlock("https://addisonchoate.com/wp-content/uploads/2020/01/Antiques.jpg", "Derp", "imgblock", titleBlock)
	//txtBlockObj := slack.NewTextBlockObject("plain_text", "fuck you slack", false, true)
	//txtBlock := slack.NewSectionBlock(txtBlockObj, nil, nil)

	headerText := slack.NewTextBlockObject("plain_text", workingImageUrl, false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)
	return []slack.Block{headerSection}
}

func buildBlocks1() slack.Message {
	// Header Section
	headerText := slack.NewTextBlockObject("mrkdwn", "You have a new request:\n*<fakeLink.toEmployeeProfile.com|Fred Enriquez - New device request>*", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	// Fields
	typeField := slack.NewTextBlockObject("mrkdwn", "*Type:*\nComputer (laptop)", false, false)
	whenField := slack.NewTextBlockObject("mrkdwn", "*When:*\nSubmitted Aut 10", false, false)
	lastUpdateField := slack.NewTextBlockObject("mrkdwn", "*Last Update:*\nMar 10, 2015 (3 years, 5 months)", false, false)
	reasonField := slack.NewTextBlockObject("mrkdwn", "*Reason:*\nAll vowel keys aren't working.", false, false)
	specsField := slack.NewTextBlockObject("mrkdwn", "*Specs:*\n\"Cheetah Pro 15\" - Fast, really fast\"", false, false)

	fieldSlice := make([]*slack.TextBlockObject, 0)
	fieldSlice = append(fieldSlice, typeField)
	fieldSlice = append(fieldSlice, whenField)
	fieldSlice = append(fieldSlice, lastUpdateField)
	fieldSlice = append(fieldSlice, reasonField)
	fieldSlice = append(fieldSlice, specsField)

	fieldsSection := slack.NewSectionBlock(nil, fieldSlice, nil)

	// Approve and Deny Buttons
	approveBtnTxt := slack.NewTextBlockObject("plain_text", "Approve", false, false)
	approveBtn := slack.NewButtonBlockElement("", "click_me_123", approveBtnTxt)

	denyBtnTxt := slack.NewTextBlockObject("plain_text", "Deny", false, false)
	denyBtn := slack.NewButtonBlockElement("", "click_me_123", denyBtnTxt)

	actionBlock := slack.NewActionBlock("", approveBtn, denyBtn)

	// Build Message with blocks created above

	msg := slack.NewBlockMessage(
		headerSection,
		fieldsSection,
		actionBlock,
	)

	return msg
}

func buildBlocks2() []slack.Block {
	// Header Section
	headerText := slack.NewTextBlockObject("mrkdwn", "You have a new request:\n*<google.com|Fred Enriquez - Time Off request>*", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	approvalText := slack.NewTextBlockObject("mrkdwn", "*Type:*\nPaid time off\n*When:*\nAug 10-Aug 13\n*Hours:* 16.0 (2 days)\n*Remaining balance:* 32.0 hours (4 days)\n*Comments:* \"Family in town, going camping!\"", false, false)
	approvalImage := slack.NewImageBlockElement("https://api.slack.com/img/blocks/bkb_template_images/approvalsNewDevice.png", "computer thumbnail")

	fieldsSection := slack.NewSectionBlock(approvalText, nil, slack.NewAccessory(approvalImage))

	// Approve and Deny Buttons
	approveBtnTxt := slack.NewTextBlockObject("plain_text", "Approve", false, false)
	approveBtn := slack.NewButtonBlockElement("", "click_me_123", approveBtnTxt)

	denyBtnTxt := slack.NewTextBlockObject("plain_text", "Deny", false, false)
	denyBtn := slack.NewButtonBlockElement("", "click_me_123", denyBtnTxt)

	actionBlock := slack.NewActionBlock("", approveBtn, denyBtn)

	// Build Message with blocks created above
	//msg := slack.NewBlockMessage(
	//	headerSection,
	//	fieldsSection,
	//	actionBlock,
	//)

    return []slack.Block{headerSection, fieldsSection, actionBlock}
}

