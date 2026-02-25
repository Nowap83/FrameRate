import apiClient from "./apiClient";

export const userService = {
    getUserFilms: async (page = 1, limit = 20) => {
        const response = await apiClient.get(`/users/me/films?page=${page}&limit=${limit}`);
        return response.data;
    },
};
