# clYR - Command-line YR
### My favorite weather service as a binary

This project uses the free and public API's from the Norwegian Meteorological Institute and Norwegian Broadcasting Corporation (NRK/YR)

Their API Documentation can be found here - [developer.yr.no](https://developer.yr.no/)  

Please use this tool with caution and read their [TOS](https://developer.yr.no/doc/TermsOfService/) before use. The current implementation of their API in this tool will not protect against abusive behaviour.

## Documentation:

Install
```bash
git clone git@github.com:WilliamDahlen/clyr.git
cd clyr
mkdir cache
go build
./clyr
```

Use
```bash
#Checks the graphical location of your IP address
./clyr
This application uses the met API to display weather information in the terminal.
2023/01/02 23:26:57 No local cache exist, GET Oslo
Current Air Temp: -5.2

#Valid local cache will be displayed until 2 hours after last 'GET' or if the 'Last Modified is older then the timeseries data'
./clyr 
This application uses the met API to display weather information in the terminal.
2023/01/02 23:24:59 Valid Cache 
Current Air Temp: -4.7

# City names can be supplied with '-c'
./clyr -c Trondheim
This application uses the met API to display weather information in the terminal.
2023/01/02 23:31:40 No local cache exist, GET Trondheim
Current Air Temp: -4.4

```
No proper error handling have been implemented yet, so it pretty much always failes if you do anything unexpected. I'll fix it later.

### Quirks:

The city to longitude and latitude mapping uses a free and downloadable CSV file found [here](https://simplemaps.com/data/world-cities). This must be placed in the local directory of the binary for now.

Caching is not pretty. It does not clean up from panics and the proper way to handle this according to the API [TOS](https://developer.yr.no/doc/TermsOfService/) is to use the `If-Modified-Since` header instead. I have not done that. Might fix it in the future.

### Resources used:

To lookup the city where your IP is located i used [ip-api.com](https://ip-api.com/docs/api). Only the City is collected and used by the binary.


### GDPR Stuff:

The 'User-Agent' is currently hardcoded to identify this application. It contains my e-mail for now. This is not considered private information.

The Meteorological Institute stores the IP from where the HTTP request originates. Please keep this in mind.
