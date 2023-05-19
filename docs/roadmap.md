# Roadmap

Matcha is actively maintained, and we're always looking for feature requests, bug reports, and general improvements. Below is what we have planned for the near future. Please keep in mind that this is not an exhaustive list!

## Version 1.1.1

Planned for release by May 15, 2023

- Specified behavior for duplicate routes in preparation for 1.2.0 route validation
- Performance improvements

## Version 1.2.0

Planned for release by July 1, 2023

- New format for regex; use as part or whole token for improved performance with known formats
- Additional route validation: match against host, origin, etc.
- Error handling: shorthand access to a generic error response and the [RFC 7801](https://datatracker.ietf.org/doc/draft-ietf-httpapi-rfc7807bis/) problem+json proposal
- Integrate with common `http.Handler` chain middleware pattern
