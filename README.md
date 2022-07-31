<p align="center">
<img src="./logo.png">
<b>A Powerful Telegram Userbot</b>
<br>
<a href="https://telegram.me/GIGAupdates"><img src="https://img.shields.io/badge/Updates%20Channel-blue?logo=telegram"></a>
<a href="https://telegram.me/GIGAsupport"><img src="https://img.shields.io/badge/Support%20Group-blue?logo=telegram"></a>
</p>

GIGA is a powerful telegram userbot written in Go with the help of [gotd](https://github.com/gotd/td) and [gotgproto](https://github.com/anonyindian/gotgproto).

## Deployment
The userbot is still under development but you can deploy it through your console by following a few steps:
- **Console based deployment**
    - Create config file
        -> `cp sample_config.json config.json`
    
        After copying the sample config to build config, just fill up the required fields in config file. 
    - Build the project
        -> `go build . -o giga`
    - Run the binary built 
        -> `./giga`
- **Deploying on Heroku**
    - You can quickly deploy GIGA on Heroku using the following deploy button:
    
        [![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy?template=https://github.com/GigaUserbot/GIGA)

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License
[![GNU Affero General Public License v3.0](https://www.gnu.org/graphics/agplv3-155x51.png)](https://www.gnu.org/licenses/agpl-3.0.en.html#header)    
[Licensed under GNU AGPL v3.0.](https://www.gnu.org/licenses/agpl-3.0.en.html#header)   
**Note**: Selling the codes to other people for money is *strictly prohibited*.