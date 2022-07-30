# GIGA Userbot Session String Generator

from telethon.sync import TelegramClient
from telethon.sessions import StringSession

print("Welcome to GIGA Userbot Session Generator\n")

APP_ID = input("Enter your App ID: ")
API_HASH = input("Enter your API Hash: ")

SESSION_MSG = """
**GIGA Userbot Session Generator**

Hello! You have successfully generated the session string for GIGA Userbot as follows:

`{session}`

**Note**: **DO NOT SHARE** this session string with anyone as it may cause hijacking of your account.
"""

client = TelegramClient(StringSession(), APP_ID, API_HASH)
# client.session.set_dc(2, '149.154.167.40', 80)
client.start()
client.send_message("me", SESSION_MSG.format(session=client.session.save()))

print("\nYou successfully generated the session string for GIGA Userbot,")
print("Please check your saved messages on telegram to get it noted.")
