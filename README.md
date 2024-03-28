# weather-exercise

## Request
Write an http server that uses the Open Weather API that exposes an endpoint that takes in lat/long coordinates. This endpoint should return what the weather condition is outside in that area (snow, rain, etc), whether it’s hot, cold, or moderate outside (use your own discretion on what temperature equates to each type).

The API can be found here:https://openweathermap.org/api. Even though most of the API calls found on OpenWeather aren’t free, you should be able to use the free “current weather data” API call for this project.  First, sign-up for an account, which shouldn’t require credit card or payment information.  Once you’ve created an account, use https://openweathermap.org/faq to get your API Key setup to start using the API.

## Implementation
This is broken into three basic parts: http server (server/), domain layer (domain/), and weather service (repo).

### HTTP Server
Very basic setup.  It only has the single route.  If this service was meant to be RESTful, obviously we would organize the single route into an appropriate path.
There's two example middleware: logging and authentication.  Authentication just passes through at the moment, but would be simple to implement.  The logging middleware assigns a logger to the request context and ties a request ID to it.  This helps with monitoring, and debugging.  It could easily be extended to contain much more information on incoming and outgoing requests.

### Domain
The domain service simply remaps the weather service data into the out going data.  Obviously if we had business logic, this is where we would do that.

### Weather Service
Basic client for interacting with the Open Weather service.  Again very simple handling here.


## Configuration
Configuration is handled purely with environment variables:

| Env Var | Required | Description | Default |
| ------- | -------- | ----------- | ------- |
| WEATHER_ADDRESS | No | Address for the http server to listen on | 0.0.0.0 |
| WEATHER_PORT | No | Port for the http server to listen on | 80 |
| WEATHER_READWRITETIMEOUT | No | Read and Write timeout for the server | 20s |
| WEATHER_IDLETIMEOUT | No | Idle timeout for the server | 75s |
| WEATHER_SHUTDOWNTIMEOUT | No | Graceful shutdown time out | 20s |
| WEATHER_LOGLEVEL | No | Zerolog log level | info |
| WEATHER_OPENWEATHER_APIID | Yes | Open Weather API ID | |
| WEATHER_OPENWEATHER_BASEURL | Yes | Base URL for Open Weather API | |
| WEATHER_OPENWEATHER_TIMEOUT | No | Client timeout for Open Weather connections | 5s |
| WEATHER_AUTHSERVICE_URL | No | Auth service URL | http://some.auth.com |

The required 

## Build and Run
Use standard go build and run commands with `cmd/server/main.go`

## NOTES

* We are making the assumption that responses are in JSON.

* There's nothing specifying float precision in the Open weather API, so I used 6 digits as it should get you around the millimeter precision.

* I am assuming this is a service that will be extended.  If this was meant to be stand alone, it should be much smaller and more streamlined.