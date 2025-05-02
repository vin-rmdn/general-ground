# User
id - string - uuid
name - string
## Index
id

# Messages
chat_id - string - uuid
from_user_id - string - fk
to_user_id - string - fk
message - string
created_at - timestamp
## Index
created_at