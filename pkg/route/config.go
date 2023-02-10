package route

// RouteConfigFuncs can be applied to a Route at creation.
type ConfigFunc func(Route) error
