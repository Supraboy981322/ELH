package elh

// render src with specific registery
func RenderWithRegistry(src string, registry map[string]Runner, r *http.Request) (string, error) {
	return parseAndRun(src, registry, r)
}

// wrapper that uses the DefaultRegistry.
func Render(src string, r *http.Request) (string, error) {
	return RenderWithRegistry(src, DefaultRegistry(), r)
}
