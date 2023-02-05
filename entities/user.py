class User:
    id = 0
    vk_id = 0
    created_at = 0


def fetch_user_to_object(result):
    user = User()
    user.id = result[0]
    user.vk_id = result[1]
    user.created_at = result[2]

    return user
