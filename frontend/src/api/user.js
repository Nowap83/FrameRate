import apiClient from "./apiClient";

export const userService = {
    getUserFilms: async (page = 1, limit = 20) => {
        const response = await apiClient.get(`/users/me/films?page=${page}&limit=${limit}`);
        return response.data;
    },
    getUserReviews: async (page = 1, limit = 20) => {
        const response = await apiClient.get(`/users/me/reviews?page=${page}&limit=${limit}`);
        return response.data;
    },
    getUserWatchlist: async (page = 1, limit = 20) => {
        try {
            const response = await apiClient.get(`/users/me/watchlist?page=${page}&limit=${limit}`);
            return response.data;
        } catch (error) {
            console.error("Error fetching user watchlist:", error);
            throw error;
        }
    },
};
