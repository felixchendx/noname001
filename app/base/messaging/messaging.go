package messaging

import (
	"fmt"
	"slices"
)

// === bundle ===
type MessageTemplateBundle struct {
	Registry map[string]*MessageTemplate
}

func NewMessageTemplateBundle() (*MessageTemplateBundle) {
	bundle := &MessageTemplateBundle{}
	bundle.Registry = make(map[string]*MessageTemplate)
	bundle.Registry["MSG.NOT.FOUND"] = &MessageTemplate{"MSG.NOT.FOUND", "Message Template with code '%s' not found."}

	return bundle
}

func (bundle *MessageTemplateBundle) AddTemplate(code string, template string) {
	bundle.Registry[code] = &MessageTemplate{code, template}
}
func (bundle *MessageTemplateBundle) AddTemplate2(code string, template string) {
	bundle.Registry[code] = &MessageTemplate{code, template}
}

func (bundle *MessageTemplateBundle) NewMessage(code string) (*Message) {
	tmpl, ok := bundle.Registry[code]
	if !ok {
		tmpl = bundle.Registry["MSG.NOT.FOUND"]
	}

	return &Message{tmpl.Code, tmpl.Template, nil}
}

// === template ===
type MessageTemplate struct {
	Code     string
	Template string
}

func NewMessageTemplate(code, template string) (*MessageTemplate) {
	return &MessageTemplate{code, template}
}

func (msgTemplate *MessageTemplate) NewMessage(args ...any) (*Message) {
	return &Message{msgTemplate.Code, msgTemplate.Template, args}
}

// === Message Without s ===
type Message struct {
	Code string

	DescTemplate string
	Args []any

	// Hint
}

func (msg *Message) String() (string) {
	desc := fmt.Sprintf(msg.DescTemplate, msg.Args...)
	return fmt.Sprintf("[%s] %s", msg.Code, desc)
}

func (msg *Message) Description() (string) {
	return fmt.Sprintf(msg.DescTemplate, msg.Args...)
}

// === OrderedMessages ===
type OrderedMessages struct {
	messages []*Message
}

// === Messages ===
type Messages struct {
	Notices  []*Message
	Warnings []*Message
	Errors   []*Message
}

func NewMessages() (*Messages) {
	return &Messages{
		Notices: make([]*Message, 0),
		Warnings: make([]*Message, 0),
		Errors: make([]*Message, 0),
	}
}

func OneLinerNotice(msg *Message) (*Messages) {
	ms := NewMessages()
	ms.AddNotice(msg)
	return ms
}
func OneLinerWarning(msg *Message) (*Messages) {
	ms := NewMessages()
	ms.AddWarning(msg)
	return ms
}
func OneLinerError(msg *Message) (*Messages) {
	ms := NewMessages()
	ms.AddError(msg)
	return ms
}

func (ms *Messages) Append(otherMessages *Messages) {
	ms.Notices = slices.Concat(ms.Notices, otherMessages.Notices)
	ms.Warnings = slices.Concat(ms.Warnings, otherMessages.Warnings)
	ms.Errors = slices.Concat(ms.Errors, otherMessages.Errors)
}

func (ms *Messages) AddNotice(msg *Message) { ms.Notices = append(ms.Notices, msg) }
func (ms *Messages) AddWarning(msg *Message) { ms.Warnings = append(ms.Warnings, msg) }
func (ms *Messages) AddError(msg *Message) { ms.Errors = append(ms.Errors, msg) }

func (ms *Messages) HasError() (bool) {
	return len(ms.Errors) > 0
}

func (ms *Messages) HasWarning() (bool) {
	return len(ms.Warnings) > 0
}

func (ms *Messages) HasMessage() (bool) {
	if len(ms.Errors) > 0 { return true }
	if len(ms.Notices) > 0 { return true }
	if len(ms.Warnings) > 0 { return true }

	return false
}

func (ms *Messages) FirstErrorMessage() (*Message) {
	if len(ms.Errors) > 0 {
		return ms.Errors[0]
	}
	return nil
}
func (ms *Messages) FirstErrorMessageString() (string) {
	msg := ms.FirstErrorMessage()
	if msg == nil { return "" }
	return msg.String()
}
func (ms *Messages) FirstErrorMessageDescription() (string) {
	msg := ms.FirstErrorMessage()
	if msg == nil { return "" }
	return msg.Description()
}

func (ms *Messages) LastErrorMessage() (*Message) {
	if len(ms.Errors) > 0 {
		return ms.Errors[len(ms.Errors) - 1]
	}
	return nil
}
func (ms *Messages) LastErrorMessageString() (string) {
	msg := ms.LastErrorMessage()
	if msg == nil { return "" }
	return msg.String()
}
func (ms *Messages) LastErrorMessageDescription() (string) {
	msg := ms.LastErrorMessage()
	if msg == nil { return "" }
	return msg.Description()
}

func (ms *Messages) FirstNoticeMessage() (*Message) {
	if len(ms.Notices) > 0 {
		return ms.Notices[0]
	}
	return nil
}
func (ms *Messages) FirstNoticeMessageString() (string) {
	msg := ms.FirstNoticeMessage()
	if msg == nil { return "" }
	return msg.String()
}
func (ms *Messages) FirstNoticeMessageDescription() (string) {
	msg := ms.FirstNoticeMessage()
	if msg == nil { return "" }
	return msg.Description()
}

func (ms *Messages) LastNoticeMessage() (*Message) {
	if len(ms.Notices) > 0 {
		return ms.Notices[len(ms.Notices) - 1]
	}
	return nil
}
func (ms *Messages) LastNoticeMessageString() (string) {
	msg := ms.LastNoticeMessage()
	if msg == nil { return "" }
	return msg.String()
}
func (ms *Messages) LastNoticeMessageDescription() (string) {
	msg := ms.LastNoticeMessage()
	if msg == nil { return "" }
	return msg.Description()
}

func (ms *Messages) Dump() (string) {
	output := "\n"

	if len(ms.Notices) > 0 {
		output += "=== Notices ===\n"
		for idx, item := range ms.Notices {
			output += fmt.Sprintf("%v. %s\n", idx, item.String())
		}
	}
	if len(ms.Warnings) > 0 {
		output += "=== Warnings ===\n"
		for idx, item := range ms.Warnings {
			output += fmt.Sprintf("%v. %s\n", idx, item.String())
		}
	}
	if len(ms.Errors) > 0 {
		output += "=== Errors ===\n"
		for idx, item := range ms.Errors {
			output += fmt.Sprintf("%v. %s\n", idx, item.String())
		}
	}

	output += "\n"

	return output
}
