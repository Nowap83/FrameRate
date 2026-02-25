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

export const getMovieVideos = async (id) => {
    const response = await apiClient.get(`/tmdb/movie/${id}/videos`);
    return response.data;
};

// USER INTERACTIONS
export const getMovieInteraction = async (id) => {
    const response = await apiClient.get(`/movies/${id}/interaction`);
    return response.data;
};

export const trackMovie = async (id, data) => {
    const response = await apiClient.post(`/movies/${id}/track`, data);
    return response.data;
};

export const rateMovie = async (id, data) => {
    const response = await apiClient.post(`/movies/${id}/rate`, data);
    return response.data;
};


