const store = window.localStorage;

function get(key) {
	return store.getItem(key);
}

function set(key, val) {
	store.setItem(key, val);
	store.getItem(key);
}

export default {
	get,
	set
}