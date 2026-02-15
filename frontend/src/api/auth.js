import apiClient from "./apiClient";

/**
 * Service gérant les appels API liés à l'authentification
 */
export const authService = {
    /**
     * Connecte un utilisateur
     * @param {Object} credentials - { email, password }
     */
    login: async (credentials) => {
        const response = await apiClient.post("/auth/login", credentials);
        return response.data;
    },

    /**
     * Inscrit un nouvel utilisateur
     * @param {Object} userData - { email, username, password }
     */
    register: async (userData) => {
        const response = await apiClient.post("/auth/register", userData);
        return response.data;
    },

    /**
     * Déconnecte l'utilisateur (optionnel selon ton implémentation backend)
     */
    logout: async () => {
        // Si tu as un endpoint de logout ou si tu nettoies juste le localStorage
        localStorage.removeItem("token");
    },
};
