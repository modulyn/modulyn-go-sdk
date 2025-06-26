package modulyn

import "slices"

// this function will evaluate whether the propvided feature name is enabled or not
// irrespective of the possibility of having fine grained control
func IsEnabled(featureName string) bool {
	datastore.mu.RLock()
	defer datastore.mu.RUnlock()

	feature, ok := datastore.features[featureName]
	if !ok {
		return false
	}

	return feature.Enabled
}

// this function will evaluate whether the provided feature name is enabled or not
// only if the provided key is configured and the provided value is a part of the
// values configured for fine grained control
func IsEnabledForKeyValue(featureName, key, value string) bool {
	datastore.mu.RLock()
	defer datastore.mu.RUnlock()

	feature, ok := datastore.features[featureName]
	if !ok {
		return false
	}

	if feature.JsonValue.Key != key || !slices.Contains(feature.JsonValue.Values, value) {
		return false
	}

	return feature.JsonValue.Enabled
}

// this function will evaluate whether the provided feature name is enabled or not
// only if the provided key is configured and the provided values exactly match the values
// configured for fine grained control
func IsEnabledForKeyValues(featureName, key string, values []string) bool {
	datastore.mu.RLock()
	defer datastore.mu.RUnlock()

	feature, ok := datastore.features[featureName]
	if !ok {
		return false
	}

	if feature.JsonValue.Key != key || !slices.Equal(feature.JsonValue.Values, values) {
		return false
	}

	return feature.JsonValue.Enabled
}
