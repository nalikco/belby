class Message:
    id = 0
    user_id = 0
    message = ""
    created_at = 0


def fetch_message_to_object(result):
    message = Message()
    message.id = result[0]
    message.user_id = result[1]
    message.message = result[2]
    message.created_at = result[3]

    return message
