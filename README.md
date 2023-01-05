# clYR - Command-line YR
### My favorite weather service as a binary

This project uses the free and public API's from the Norwegian Meteorological Institute and Norwegian Broadcasting Corporation (NRK/YR)  

As of now the CLI uses the "Compact" endpoint to display the current air temperature in the supplied city. Future development will add flags to print more of the forecast from the same endpoint.

Their API Documentation can be found here - [developer.yr.no](https://developer.yr.no/)  

Please use this tool with caution and read their [TOS](https://developer.yr.no/doc/TermsOfService/) before use. The current implementation of their API in this tool will not protect against abusive behaviour.

## Documentation:

Getting started
```bash
git clone git@github.com:WilliamDahlen/clyr.git
cd clyr
mkdir cache
go build
./clyr
```

Use
```bash
#If no city is supplied the geographical location of your IP address is used by default
./clyr
This application uses the met API to display weather information in the terminal.
No local cache exist, GET Oslo
Current Air Temp @ Thu, 05 Jan 2023 23:00:00 GMT: -2.9

#Valid local cache will be displayed until the Expires header timestamp in the local filename is past current system time.
./clyr 
This application uses the met API to display weather information in the terminal.
Valid Cache 
Current Air Temp @ Thu, 05 Jan 2023 23:00:00 GMT: -2

# City names can be supplied with '-c'
./clyr -c Trondheim
This application uses the met API to display weather information in the terminal.
No local cache exist, GET Trondheim
Current Air Temp @ Thu, 05 Jan 2023 23:00:00 GMT: -2

```
No proper error handling have been implemented yet, so it pretty much always failes if you do anything unexpected. I'll fix it later.

### Quirks:

The city to longitude and latitude mapping uses a free and downloadable CSV file found [here](https://simplemaps.com/data/world-cities). This must be placed in the local directory of the binary for now.

~~Caching is not pretty. It does not clean up from panics and the proper way to handle this according to the API [TOS](https://developer.yr.no/doc/TermsOfService/) is to use the `If-Modified-Since` header instead. I have not done that. Might fix it in the future.~~  

Hi! Some real cache is now implemented. All local files are used as cache until their 'Exires' timestamp is reached. Next on the journey to caching is to implement a HEAD check before pulling down a new copy when the local file expire timestamp is reached. Instead it will update the existing filename with a new expire timestamp from the HEAD response.

### Resources used:

To lookup the city where your IP is located i used [ip-api.com](https://ip-api.com/docs/api). Only the City is collected and used by the binary.


### GDPR Stuff:

The 'User-Agent' is currently hardcoded to identify this application. It contains my e-mail for now. This is not considered private information.

The Meteorological Institute stores the IP from where the HTTP request originates. Please keep this in mind.
