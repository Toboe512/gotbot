package telegram

const msgHelp = `Я могу сохранять и хранить ваши страницы. Также я могу предложить вам их для чтения.

	Чтобы сохранить страницу, просто отправьте мне ссылку на нее.

	Чтобы получить случайную страницу из вашего списка, отправьте мне команду /rnd.
	Внимание! После этого эта страница будет удалена из вашего списка`

const msgHello = "Привет! 👾\n\n" + msgHelp

const (
	msgUnknownCommand = "Неизвестная команда 🤔"
	msgNoSavedPages   = "У вас нет сохраненных страниц 🙊"
	msgSaved          = "Сохранено! 👌"
	msgAlreadyExists  = "Эта страница уже есть в вашем списке 🤗"
)
