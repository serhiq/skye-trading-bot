package performer

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/serhiq/skye-trading-bot/internal/app"
	b "github.com/serhiq/skye-trading-bot/internal/bot"
	"github.com/serhiq/skye-trading-bot/internal/bot/commands"
	"github.com/serhiq/skye-trading-bot/internal/contorller"
	"github.com/serhiq/skye-trading-bot/internal/logger"
	"github.com/serhiq/skye-trading-bot/internal/repository"
	"github.com/serhiq/skye-trading-bot/internal/repository/chat"
	"strings"
)

type Options struct {
	Token    string
	TimeZone string
}

type Performer struct {
	App              *app.App
	menuHandlers     []menuHandler
	commandHandler   []commandHandler
	freeInputHandler []freeInputHandler
	callbackHandlers []callbackHandler
	contactHandler   func(app *app.App, message *tgbotapi.Message) error
}

type menuHandler struct {
	name    string
	handler func(app *app.App, message *tgbotapi.Message) error
}

type commandHandler struct {
	name    string
	handler func(app *app.App, message *tgbotapi.Message) error
}

type freeInputHandler struct {
	state   string
	handler func(app *app.App, input string, session *chat.Chat) error
}

type callbackHandler struct {
	name    string
	handler func(app *app.App, callback *tgbotapi.CallbackQuery) error
}

func (p *Performer) Dispatch(update *tgbotapi.Update) {

	if err := p.process(update); err != nil {
		p.processError(err, *update)
	}
}

func (p *Performer) isCommandButton(text string) bool {
	for _, hanler := range p.commandHandler {
		if hanler.name == text {
			return true
		}
	}
	return false
}

func (p *Performer) AddCallbackHandler(name string, handl func(app *app.App, callback *tgbotapi.CallbackQuery) error) {
	ch := callbackHandler{
		name:    name,
		handler: handl,
	}
	callbackHanlers := append(p.callbackHandlers, ch)
	p.callbackHandlers = callbackHanlers
}

func (p *Performer) AddMenuCommandHandler(name string, handl func(app *app.App, message *tgbotapi.Message) error) {
	ch := menuHandler{
		name:    name,
		handler: handl,
	}
	commandHanlers := append(p.menuHandlers, ch)
	p.menuHandlers = commandHanlers
}

func (p *Performer) AddCommandHandler(name string, handl func(app *app.App, message *tgbotapi.Message) error) {
	ch := commandHandler{
		name:    name,
		handler: handl,
	}
	commandHanlers := append(p.commandHandler, ch)
	p.commandHandler = commandHanlers
}

func (p *Performer) AddInputHandler(name string, handl func(app *app.App, input string, session *chat.Chat) error) {
	ch := freeInputHandler{
		state:   name,
		handler: handl,
	}
	handlers := append(p.freeInputHandler, ch)
	p.freeInputHandler = handlers
}

func (p *Performer) RegisterContactHandler(handl func(app *app.App, message *tgbotapi.Message) error) {
	p.contactHandler = handl
}

func (p *Performer) processMenuCommand(message *tgbotapi.Message) error {
	var menuCommand = message.Command()

	for _, command := range p.menuHandlers {
		if command.name == menuCommand {
			return command.handler(p.App, message)
		}
	}
	return b.NewCommandNotFound(message.Command())
}

func (p *Performer) processCommand(message *tgbotapi.Message) error {

	for _, command := range p.commandHandler {
		if command.name == message.Text {
			return command.handler(p.App, message)
		}
	}
	return b.NewCommandNotFound(message.Text)
}

func (p *Performer) processCallback(callback *tgbotapi.CallbackQuery) error {
	c := commands.New(callback.Data)

	for _, command := range p.callbackHandlers {
		if command.name == c.Command {
			return command.handler(p.App, callback)
		}
	}
	return b.NewCommandNotFound(callback.Data)
}

func (p *Performer) process(update *tgbotapi.Update) error {

	if update.Message != nil && update.Message.IsCommand() {
		return p.processMenuCommand(update.Message)
	}

	if update.Message != nil && StartsWith("🛒", update.Message.Text) {
		return b.DisplayOrderHandler(p.App, update.Message)
	}

	if update.Message != nil && p.isCommandButton(update.Message.Text) {
		return p.processCommand(update.Message)
	}

	if update.CallbackQuery != nil {
		return p.processCallback(update.CallbackQuery)
	}

	if update.Message != nil && update.Message.Contact != nil {
		return p.contactHandler(p.App, update.Message)
	}

	// check free input

	if update.Message != nil {
		session, err := p.App.RepoChat.GetOrCreateChat(update.Message.Chat.ID)
		if err != nil {
			return fmt.Errorf("Failed to get chat  %s", err.Error())
		}

		return p.processFreeInput(update.Message.Text, session)
	}

	if update.Message.Text != "" {
		return b.NewCommandNotFound(update.Message.Text)
	}

	return nil
}

func (p *Performer) processFreeInput(input string, session *chat.Chat) error {
	for _, command := range p.freeInputHandler {
		if command.state == session.ChatState {
			return command.handler(p.App, input, session)
		}
	}
	return b.NewCommandNotFound("free input for state " + session.ChatState)
}

func (p *Performer) processError(err error, update tgbotapi.Update) {

	if err != nil {
		if b.IsCommandNotFoundError(err) {
			logger.SugaredLogger.Errorw("bot_command_not_found", "update", update,
				"chatId", update.FromChat(), "err", err)
			return
		}

		logger.SugaredLogger.Errorw("bot_update",
			"chatId", update.FromChat(), "err", err)
		return
	}
}

func New(options Options, productController contorller.ProductController, repoChat repository.ChatRepository, orderController contorller.OrderController) (*Performer, error) {
	bot, err := tgbotapi.NewBotAPI(options.Token)
	if err != nil {
		return nil, err
	}

	var p = Performer{
		App: &app.App{
			ProductController: productController,
			RepoChat:          repoChat,
			OrderController:   orderController,
			Bot:               app.NewTelegramBot(bot),
			Cfg:               &app.AppConfig{TimeZone: options.TimeZone},
		}}

	p.AddMenuCommandHandler("start", b.StartCommand)

	p.RegisterContactHandler(b.ContactHandler)
	p.AddInputHandler(chat.STATE_INPUT_PHONE, b.ContactInputHandler)
	p.AddInputHandler(chat.STATE_INPUT_NAME, b.NameInputHandler)

	p.AddCommandHandler(b.DISPLAY_MENU_BUTTON, b.DisplayMenuHandler)
	p.AddCommandHandler(b.DISPLAY_ORDER_BUTTON, b.DisplayOrderHandler)
	p.AddCommandHandler(b.HELP_BUTTON, b.AboutCommand)

	p.AddCommandHandler(b.CLEAR_ORDER_BUTTON, b.ClearOrderHandler)
	p.AddCommandHandler(b.BACK_ORDER_BUTTON, b.BackToCatalogFromOrder)

	//  меню
	p.AddCallbackHandler(b.CLICK_ON_FOLDER, b.ClickOnFolderCallbackHandler)
	p.AddCallbackHandler(b.CLICK_ON_PRODUCT_ITEM, b.ClickOnItemCallbackHandler)
	p.AddCallbackHandler(b.CLICK_ON_BACK, b.ClickOnBackInFolderCallbackHandler)

	// редактирование заказа через меню
	p.AddCallbackHandler(b.COMMAND_DECREASE_POSITION, b.ClickOnDecreasePositionCallbackHandler)
	p.AddCallbackHandler(b.COMMAND_ADD_POSITION, b.ClickOnIncreasePositionCallbackHandler)
	p.AddCallbackHandler(b.CLICK_ON_EDIT_POSITION, b.ClickOnEditPositionCallbackHandler)

	// редактирование заказа при просмотре корзины
	p.AddCallbackHandler(b.CLICK_ON_BACK_EDIT_ORDER, b.ClickDisplayEditOrder)
	p.AddCallbackHandler(b.CLICK_ON_ADD_POSITION, b.ClickOnIncreasePositionEditOrderCallbackHandler)
	p.AddCallbackHandler(b.CLICK_ON_DECREASE_POSITION, b.ClickOnDecreasePositionEditOrderCallbackHandler)
	/////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//оформление заказа
	p.AddCommandHandler(b.ORDER_BUTTON, b.CreateOrderHandler)

	//доставка
	p.AddCallbackHandler(b.DELIVERY_COMMAND, b.ClickOnSetDeliveryCallbackHandler)
	//самовывоз
	p.AddCallbackHandler(b.SELF_PICKUP_COMMAND, b.ClickOnSetDeliveryCallbackHandler)

	p.AddInputHandler(chat.INPUT_DELIVERY_LOCATION, b.InputLocationHandler)

	p.AddCallbackHandler(b.TIME_COMMAND_40M, b.ClickOnSetTimeCallback)
	p.AddCallbackHandler(b.TIME_COMMAND_120M, b.ClickOnSetTimeCallback)
	p.AddCallbackHandler(b.TIME_COMMAND_SOON, b.ClickOnSetTimeCallback)

	p.AddInputHandler(chat.INPUT_COMMENT, b.InputCommentHandler)
	p.AddCallbackHandler(b.NO_COMMENT_COMMAND, b.ClickOnNoCommentCallback)

	p.AddCallbackHandler(b.CONFIRM_COMMAND, b.ClickOnConfirm)
	p.AddCallbackHandler(b.CANCEL_COMMAND, b.ClickOnCancel)

	//профиль
	p.AddCommandHandler(b.DISPLAY_PROFILE_BUTTON, b.ProfileMenuHandler)

	p.AddCallbackHandler(b.DISPLAY_HISTORY_COMMAND, b.ProfileDisplayHistory)
	p.AddCallbackHandler(b.CLICK_ON_REPEAT_ORDER, b.ProfileRepeatOrder)

	p.AddCallbackHandler(b.CANCEL_FROM_PROFILE, b.ProfileOnCancel)

	p.AddCallbackHandler(b.CHANGE_PHONE_COMMAND, b.ClickOnChangePhone)
	p.AddInputHandler(chat.STATE_CHANGE_PHONE, b.ContactInputHandler)

	p.AddCallbackHandler(b.CHANGE_NAME_COMMAND, b.ClickOnChangeName)
	p.AddInputHandler(chat.STATE_CHANGE_NAME, b.NameInputHandler)

	bot.Debug = false

	return &p, nil
}

func StartsWith(prefix string, content string) bool {
	return (strings.Split(content, " ")[0] == prefix)
}
