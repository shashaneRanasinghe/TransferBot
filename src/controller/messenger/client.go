package messenger

import (
	"Pelmenner/TransferBot/orm"
	"context"
	msg "github.com/Pelmenner/TransferBot/proto/messenger"
	"google.golang.org/grpc"
)

type MessengerClient struct {
	msg.ChatServiceClient
}

func NewMessengerClient(cc grpc.ClientConnInterface) *MessengerClient {
	internalClient := msg.NewChatServiceClient(cc)
	return &MessengerClient{
		internalClient,
	}
}

func (m *MessengerClient) SendMessage(message *orm.Message, chat *orm.Chat) error {
	pbMessage := messageToProto(message)
	pbChat := chatToProto(chat)
	_, err := m.ChatServiceClient.SendMessage(context.TODO(), &msg.SendMessageRequest{
		Message: pbMessage,
		Chat:    pbChat,
	})
	return err
}

func messageToProto(message *orm.Message) *msg.Message {
	pbSender := senderToProto(&message.Sender)
	pbMessage := msg.Message{
		Text:   message.Text,
		Sender: pbSender,
	}
	for _, attachment := range message.Attachments {
		pbMessage.Attachments = append(pbMessage.Attachments, attachmentToProto(attachment))
	}
	return &pbMessage
}

func chatToProto(chat *orm.Chat) *msg.Chat {
	if chat == nil {
		return nil
	}
	return &msg.Chat{
		Id:    chat.ID,
		RowID: chat.RowID,
		Name:  chat.Type, // TODO: pass an actual name
		Type:  chat.Type,
		Token: chat.Token,
	}
}

func attachmentToProto(attachment *orm.Attachment) *msg.Attachment {
	return &msg.Attachment{
		Type: attachment.Type,
		Url:  attachment.URL,
	}
}
func senderToProto(sender *orm.Sender) *msg.Sender {
	return &msg.Sender{
		Name: sender.Name,
		Chat: &msg.Chat{Name: sender.Chat}, // TODO: make sender.Chat an orm.Chat instead of a string
	}
}
