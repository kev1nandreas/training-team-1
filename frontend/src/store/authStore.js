import { produce } from "immer";
import { create } from "zustand";
import { removeToken, setToken } from "../utils/cookies";

const useAuthStore = create((set) => ({
	token: null,
	user: null,
	isAuthed: false,
	isLoading: true,

	login: async (token, user) => {
		await setToken(token);
		set(
			produce((state) => {
				state.token = token;
				state.user = user;
				state.isAuthed = true;
			})
		);
	},

	logout: () => {
		removeToken();
		set(
			produce((state) => {
				state.token = null;
				state.user = null;
				state.isAuthed = false;
			})
		);
	},

	stopLoading: () => {
		set(
			produce((state) => {
				state.isLoading = false;
			})
		);
	},
}));

export default useAuthStore;
