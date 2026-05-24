import { produce } from "immer";
import { create } from "zustand";

const useAuthStore = create((set) => ({
	token: null,
	user: null,
	isAuthed: false,
	isLoading: true,

	login: (token, user) => {
		set(
			produce((state) => {
				state.token = token;
				state.user = user;
				state.isAuthed = true;
			})
		);
	},

	logout: () => {
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
