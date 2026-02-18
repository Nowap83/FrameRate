import apiClient from "./apiClient";

export const getPopularMovies = async () => {
    const response = await apiClient.get("/tmdb/popular");
    return response.data;
};

export const getMovieDetails = async (id) => {
    const response = await apiClient.get(`/tmdb/movie/${id}`);
    return response.data;
};

export const getMovieCredits = async (id) => {
    const response = await apiClient.get(`/tmdb/movie/${id}/credits`);
    return response.data;
};


