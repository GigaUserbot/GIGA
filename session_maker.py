# GIGA Userbot Session String Generator
import sys
from pyrogram.client import Client
from telethon.sync import TelegramClient
from telethon.sessions import StringSession

print("Welcome to GIGA Userbot Session Generator\n")

APP_ID = input("Enter your App ID: ")
API_HASH = input("Enter your API Hash: ")

SESSION_MSG = """
**GIGA Userbot Session Generator**

Hello! You have successfully generated the {type} session string for GIGA Userbot as follows:

`{session}`

**Note**: **DO NOT SHARE** this session string with anyone as it may cause hijacking of your account.
"""

SESSION_TYPE = input(
"""Please specificy session type,
Enter 0 for telethon
Enter 1 for pyrogram
Your answer: """)

if SESSION_TYPE == "0":
    client = TelegramClient(StringSession(), APP_ID, API_HASH)
    # client.session.set_dc(2, '149.154.167.40', 80)
    client.start()
    client.send_message("me", SESSION_MSG.format(type="telethon",session=client.session.save()))
elif SESSION_TYPE == "1":
    client = Client("", APP_ID, API_HASH, in_memory=True)
    client.start()
    client.send_message("me", SESSION_MSG.format(type="pyrogram", session=client.export_session_string()))
else:
    print("Your input was invalid, expected 0 or 1.")
    sys.exit(1)

print("\nYou successfully generated the session string for GIGA Userbot,")
print("Please check your saved messages on telegram to get it noted.")
